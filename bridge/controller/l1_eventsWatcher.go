package controller

import (
	"context"
	"math/big"
	"sync"

	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// L1ReorgSafeDepth represents the number of block confirmations considered safe against L1 chain reorganizations.
// Reorganizations at this depth under normal cases are extremely unlikely.
const L1ReorgSafeDepth = 64

type L1FilterResult struct {
	DepositMessages []*orm.RawBridgeEvent
	RelayedMessages []*orm.RawBridgeEvent
}
type L1EventsWatcher struct {
	ctx                 context.Context
	cfg                 *evm.GethConfig
	l1Client            *ethclient.Client
	l1EventParser       *logic.L1EventParser
	mu                  sync.Mutex
	l1SyncHeight        uint64
	l1LastSyncBlockHash common.Hash
	contractAddressList []common.Address

	rawBridgeEventsOrm *orm.RawBridgeEvent
}

func NewL1EventsWatcher(ctx context.Context, cfg *evm.GethConfig, ethClient *ethclient.Client, db *gorm.DB) (*L1EventsWatcher, error) {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ParentLayerContractAddress)}
	c := &L1EventsWatcher{
		ctx:                 ctx,
		cfg:                 cfg,
		l1Client:            ethClient,
		l1EventParser:       logic.NewL1EventParser(cfg),
		rawBridgeEventsOrm:  orm.NewRawBridgeEvent(db, cfg),
		contractAddressList: contractAddressList,
	}
	return c, nil
}

// Start starts the L1 message fetching process.
func (w *L1EventsWatcher) Start() {
	messageSyncedHeight, dbErr := w.GetL1SyncHeight(w.ctx)
	if dbErr != nil {
		logrus.Error("failed to get L1 cross message synced height", "error", dbErr)

	}
	l1SyncHeight := messageSyncedHeight

	// Sync from an older block to prevent reorg during restart.
	if l1SyncHeight < L1ReorgSafeDepth {
		l1SyncHeight = 0
	} else {
		l1SyncHeight -= L1ReorgSafeDepth
	}
	if w.cfg.L1WatcherConfig.StartHeight > l1SyncHeight {
		l1SyncHeight = w.cfg.L1WatcherConfig.StartHeight - 1
	}
	header, err := w.l1Client.HeaderByNumber(w.ctx, new(big.Int).SetUint64(l1SyncHeight))
	if err != nil {
		log.Crit("failed to get L1 header by number", "block number", l1SyncHeight, "err", err)
		return
	}

	w.updateL1SyncHeight(l1SyncHeight, header.Hash())

	logrus.Info("Start L1 message fetcher ",
		" message synced height", messageSyncedHeight,
		" config start height", w.cfg.L1WatcherConfig.StartHeight,
		" sync start height", w.l1SyncHeight+1,
	)

	tick := time.NewTicker(time.Duration(w.cfg.L1WatcherConfig.BlockTime) * time.Second)
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				tick.Stop()
				return
			case <-tick.C:
				w.fetchAndSaveEvents(w.cfg.L1WatcherConfig.Confirmation)
			}
		}
	}()
}

/*****************************
 *    [Functions:Handler]    *
 *****************************/
// func (w *L1EventsWatcher) handleDownwardMessage(
// 	msg *contract.ParentBridgeCoreFacetQueueTransaction,
// ) error {
// 	err := w.l1toL2Relayer.HandleDownwardMessageWithSystemCall(msg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (w *L1EventsWatcher) handleRelayerMessage(msg *contract.UpwardMessageDispatcherFacetRelayedMessage) error {
// 	err := w.l1toL2Relayer.HandleRelayerMessage(msg)
// 	if err != nil {
// 		logrus.Errorf("Failed to handle RelayerMessage: %v", err)
// 		return err
// 	}
// 	return nil
// }
//save downward message to db

func (w *L1EventsWatcher) ChainID(ctx context.Context) (*big.Int, error) {
	return w.l1Client.ChainID(ctx)
}

func (w *L1EventsWatcher) Close() {
	w.l1Client.Close()
}

func (w *L1EventsWatcher) fetchAndSaveEvents(confirmation uint64) {
	startHeight := w.l1SyncHeight + 1
	endHeight, rpcErr := utils.GetBlockNumber(w.ctx, w.l1Client, confirmation)
	if rpcErr != nil {
		logrus.Error("failed to get L1 block number", "confirmation", confirmation, "err", rpcErr)
		return
	}

	logrus.Info("fetch and save missing L1 events", "start height", startHeight, "end height", endHeight, "confirmation", confirmation)

	for from := startHeight; from <= endHeight; from += w.cfg.L1WatcherConfig.FetchLimit {
		to := from + w.cfg.L1WatcherConfig.FetchLimit - 1
		if to > endHeight {
			to = endHeight
		}

		isReorg, resyncHeight, lastBlockHash, l1FetcherResult, fetcherErr := w.L1Fetcher(w.ctx, from, to, w.l1LastSyncBlockHash)
		if fetcherErr != nil {
			log.Error("failed to fetch L1 events", "from", from, "to", to, "err", fetcherErr)
			return
		}

		if isReorg {
			log.Warn("L1 reorg happened, exit and re-enter fetchAndSaveEvents", "re-sync height", resyncHeight)
			w.updateL1SyncHeight(resyncHeight, lastBlockHash)
			return
		}

		if insertUpdateErr := w.L1InsertOrUpdate(w.ctx, l1FetcherResult); insertUpdateErr != nil {
			log.Error("failed to save L1 events", "from", from, "to", to, "err", insertUpdateErr)
			return
		}

		w.updateL1SyncHeight(to, lastBlockHash)
	}
}

func (w *L1EventsWatcher) updateL1SyncHeight(height uint64, blockHash common.Hash) {
	w.l1LastSyncBlockHash = blockHash
	w.l1SyncHeight = height
}

func (w *L1EventsWatcher) L1Fetcher(ctx context.Context, from, to uint64, lastBlockHash common.Hash) (bool, uint64, common.Hash, *L1FilterResult, error) {
	log.Info("fetch and save L1 events", "from", from, "to", to)

	isReorg, reorgHeight, blockHash, _, getErr := w.getBlocksAndDetectReorg(ctx, from, to, lastBlockHash)
	if getErr != nil {
		log.Error("L1Fetcher getBlocksAndDetectReorg failed", "from", from, "to", to, "error", getErr)
		return false, 0, common.Hash{}, nil, getErr
	}

	if isReorg {
		return isReorg, reorgHeight, blockHash, nil, nil
	}

	eventLogs, err := w.l1FetcherLogs(ctx, from, to)

	if err != nil {
		log.Error("L1Fetcher l1FetcherLogs failed", "from", from, "to", to, "error", err)
		return false, 0, common.Hash{}, nil, err
	}

	l1DepositMessages, l1RelayedMessages, err := w.l1EventParser.ParseL1EventToRawBridgeEvents(ctx, eventLogs)
	if err != nil {
		log.Error("failed to parse L1 cross chain event logs", "from", from, "to", to, "err", err)
		return false, 0, common.Hash{}, nil, err
	}

	res := L1FilterResult{
		DepositMessages: l1DepositMessages,
		RelayedMessages: l1RelayedMessages,
	}

	return false, 0, blockHash, &res, nil
}

func (w *L1EventsWatcher) getBlocksAndDetectReorg(ctx context.Context, from, to uint64, lastBlockHash common.Hash) (bool, uint64, common.Hash, []*types.Block, error) {
	blocks, err := utils.GetBlocksInRange(ctx, w.l1Client, from, to)

	if err != nil {
		logrus.Error("failed to get L1 blocks in range", "from", from, "to", to, "err", err)
		return false, 0, common.Hash{}, nil, err
	}

	for _, block := range blocks {
		if block.ParentHash() != lastBlockHash {
			logrus.Warn("L1 reorg detected", " reorg height", block.NumberU64()-1, "expected hash", block.ParentHash().String(), "local hash", lastBlockHash.String())
			var resyncHeight uint64
			if block.NumberU64() > L1ReorgSafeDepth+1 {
				resyncHeight = block.NumberU64() - L1ReorgSafeDepth - 1
			}
			header, err := w.l1Client.HeaderByNumber(ctx, new(big.Int).SetUint64(resyncHeight))
			if err != nil {
				log.Error("failed to get L1 header by number", "block number", resyncHeight, "err", err)
				return false, 0, common.Hash{}, nil, err
			}
			return true, resyncHeight, header.Hash(), nil, nil
		}
		lastBlockHash = block.Hash()
	}

	return false, 0, lastBlockHash, blocks, nil
}
func (w *L1EventsWatcher) GetL1SyncHeight(ctx context.Context) (uint64, error) {
	messageSyncedHeight, err := w.rawBridgeEventsOrm.GetMaxBlockNumber(ctx, w.cfg.L1_RawBridgeEventsTableName)
	if err != nil {
		log.Error("failed to get L1 cross message synced height", "error", err)
		return 0, err
	}

	return messageSyncedHeight, nil
}
func (w *L1EventsWatcher) L1InsertOrUpdate(ctx context.Context, l1FetcherResult *L1FilterResult) error {
	if err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), w.cfg.L1_RawBridgeEventsTableName, l1FetcherResult.DepositMessages); err != nil {
		logrus.Error("failed to insert L1 deposit messages", "err", err)
		return err
	}

	if err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), w.cfg.L1_RawBridgeEventsTableName, l1FetcherResult.RelayedMessages); err != nil {
		logrus.Error("failed to insert L1 relayed messages", "err", err)
		return err
	}

	return nil
}
func (w *L1EventsWatcher) l1FetcherLogs(ctx context.Context, from, to uint64) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), // inclusive
		ToBlock:   new(big.Int).SetUint64(to),   // inclusive
		Addresses: w.contractAddressList,
		Topics:    make([][]common.Hash, 1),
	}

	query.Topics[0] = make([]common.Hash, 2)
	query.Topics[0][0] = backendabi.L1QueueTransactionEventSig
	query.Topics[0][1] = backendabi.L1RelayedMessageEventSig

	eventLogs, err := w.l1Client.FilterLogs(ctx, query)
	if err != nil {
		log.Error("failed to filter L1 event logs", "from", from, "to", to, "err", err)
		return nil, err
	}
	return eventLogs, nil
}
