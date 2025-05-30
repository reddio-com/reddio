package controller

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	backendabi "github.com/reddio-com/reddio/bridge/abi"
	rdoclient "github.com/reddio-com/reddio/bridge/client"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

const L2ReorgSafeDepth = 64

type L2FilterResult struct {
	WithdrawMessages []*orm.RawBridgeEvent
	RelayedMessages  []*orm.RawBridgeEvent
}
type L2EventsWatcher struct {
	ctx                 context.Context
	cfg                 *evm.GethConfig
	l2Client            *rdoclient.Client
	l2WatcherLogic      *logic.L2WatcherLogic
	l2EventParser       *logic.L2EventParser
	mu                  sync.Mutex
	l2SyncHeight        uint64
	l2LastSyncBlockHash common.Hash
	contractAddressList []common.Address
	rawBridgeEventsOrm  *orm.RawBridgeEvent
}

func NewL2EventsWatcher(ctx context.Context, cfg *evm.GethConfig, rdoclient *rdoclient.Client, db *gorm.DB) (*L2EventsWatcher, error) {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ChildLayerContractAddress)}
	c := &L2EventsWatcher{
		ctx:                 ctx,
		cfg:                 cfg,
		l2Client:            rdoclient,
		l2EventParser:       logic.NewL2EventParser(cfg),
		rawBridgeEventsOrm:  orm.NewRawBridgeEvent(db, cfg),
		contractAddressList: contractAddressList,
	}
	return c, nil
}

// Start starts the L1 message fetching process.
func (w *L2EventsWatcher) Start() {
	messageSyncedHeight, dbErr := w.GetL2SyncHeight(w.ctx)
	if dbErr != nil {
		logrus.Error("failed to get L2 cross message synced height", "error", dbErr)

	}
	l2SyncHeight := messageSyncedHeight

	// Sync from an older block to prevent reorg during restart.
	if l2SyncHeight < L2ReorgSafeDepth {
		l2SyncHeight = 0
	} else {
		l2SyncHeight -= L2ReorgSafeDepth
	}
	if w.cfg.L2WatcherConfig.StartHeight > l2SyncHeight {
		l2SyncHeight = w.cfg.L2WatcherConfig.StartHeight - 1
	}
	header, err := w.l2Client.HeaderByNumberNoType(w.ctx, new(big.Int).SetUint64(l2SyncHeight))
	if err != nil {
		logrus.Warn("failed to get L1 header by number", "block number", l2SyncHeight, "err", err)
		return
	}
	//blockNumberHex := fmt.Sprintf("0x%x", l2SyncHeight)
	blockHash, _ := (*header)["hash"].(string)
	w.updateL2SyncHeight(l2SyncHeight, common.HexToHash(blockHash))

	logrus.Info("Start L2 message fetcher ",
		" message synced height", messageSyncedHeight,
		" config start height", w.cfg.L2WatcherConfig.StartHeight,
		" sync start height", w.l2SyncHeight+1,
	)

	tick := time.NewTicker(time.Duration(w.cfg.L2WatcherConfig.BlockTime) * time.Second)
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				tick.Stop()
				return
			case <-tick.C:
				w.fetchAndSaveEvents(w.cfg.L2WatcherConfig.Confirmation)
			}
		}
	}()
}

func (w *L2EventsWatcher) ChainID(ctx context.Context) (*big.Int, error) {
	return w.l2Client.ChainID(ctx)
}

func (w *L2EventsWatcher) Close() {
	w.l2Client.Close()
}

func (w *L2EventsWatcher) fetchAndSaveEvents(confirmation uint64) {
	startHeight := w.l2SyncHeight + 1
	endHeight, rpcErr := utils.GetRdoBlockNumber(w.ctx, w.l2Client, confirmation)
	if rpcErr != nil {
		logrus.Error("failed to get L2 block number", "confirmation", confirmation, "err", rpcErr)
		return
	}

	for from := startHeight; from <= endHeight; from += w.cfg.L2WatcherConfig.FetchLimit {
		to := from + w.cfg.L2WatcherConfig.FetchLimit - 1
		if to > endHeight {
			to = endHeight
		}

		isReorg, resyncHeight, lastBlockHash, l2FetcherResult, fetcherErr := w.L2Fetcher(w.ctx, from, to, w.l2LastSyncBlockHash)
		if fetcherErr != nil {
			log.Error("failed to fetch L2 events", "from", from, "to", to, "err", fetcherErr)
			return
		}

		if isReorg {
			log.Warn("L2 reorg happened, exit and re-enter fetchAndSaveEvents", "re-sync height", resyncHeight)
			w.updateL2SyncHeight(resyncHeight, lastBlockHash)
			return
		}

		if insertUpdateErr := w.L2InsertOrUpdate(w.ctx, l2FetcherResult); insertUpdateErr != nil {
			log.Error("failed to save L2 events", "from", from, "to", to, "err", insertUpdateErr)
			return
		}

		w.updateL2SyncHeight(to, lastBlockHash)
	}
}

func (w *L2EventsWatcher) updateL2SyncHeight(height uint64, blockHash common.Hash) {
	w.l2LastSyncBlockHash = blockHash
	w.l2SyncHeight = height
}

func (w *L2EventsWatcher) L2Fetcher(ctx context.Context, from, to uint64, lastBlockHash common.Hash) (bool, uint64, common.Hash, *L2FilterResult, error) {
	log.Info("fetch and save L1 events", "from", from, "to", to)

	isReorg, reorgHeight, blockHash, _, getErr := w.getBlocksAndDetectReorg(ctx, from, to, lastBlockHash)
	if getErr != nil {
		log.Error("L1Fetcher getBlocksAndDetectReorg failed", "from", from, "to", to, "error", getErr)
		return false, 0, common.Hash{}, nil, getErr
	}

	if isReorg {
		return isReorg, reorgHeight, blockHash, nil, nil
	}

	eventLogs, err := w.l2FetcherLogs(ctx, from, to)

	if err != nil {
		log.Error("L1Fetcher l2FetcherLogs failed", "from", from, "to", to, "error", err)
		return false, 0, common.Hash{}, nil, err
	}

	l1WithdrawMessages, l2RelayedMessages, err := w.l2EventParser.ParseL2EventToRawBridgeEvents(ctx, eventLogs)
	if err != nil {
		log.Error("failed to parse L2 cross chain event logs", "from", from, "to", to, "err", err)
		return false, 0, common.Hash{}, nil, err
	}

	res := L2FilterResult{
		WithdrawMessages: l1WithdrawMessages,
		RelayedMessages:  l2RelayedMessages,
	}

	return false, 0, blockHash, &res, nil
}

func (w *L2EventsWatcher) getBlocksAndDetectReorg(ctx context.Context, from, to uint64, lastBlockHash common.Hash) (bool, uint64, common.Hash, []*rdoclient.RdoBlock, error) {
	blocks, err := utils.GetRdoBlocksInRange(ctx, w.l2Client, from, to)

	if err != nil {
		logrus.Error("failed to get L2 blocks in range", "from", from, "to", to, "err", err)
		return false, 0, common.Hash{}, nil, err
	}

	for _, block := range blocks {
		if block.ParentHash() != lastBlockHash {
			logrus.Warn("L2 reorg detected", " reorg height ", block.NumberU64()-1, " expected hash ", block.ParentHash().String(), "local hash", lastBlockHash.String(), "current block hash", block.Hash().String())
			var resyncHeight uint64
			if block.NumberU64() > L1ReorgSafeDepth+1 {
				resyncHeight = block.NumberU64() - L1ReorgSafeDepth - 1
			}
			header, err := w.l2Client.HeaderByNumberNoType(w.ctx, new(big.Int).SetUint64(resyncHeight))
			if err != nil {
				log.Error("failed to get L2 header by number", "block number", resyncHeight, "err", err)
				return false, 0, common.Hash{}, nil, err
			}
			blockHash, _ := (*header)["hash"].(string)
			return true, resyncHeight, common.HexToHash(blockHash), nil, nil
		}
		lastBlockHash = block.Hash()
	}

	return false, 0, lastBlockHash, blocks, nil
}

func (w *L2EventsWatcher) GetL2SyncHeight(ctx context.Context) (uint64, error) {
	messageSyncedHeight, err := w.rawBridgeEventsOrm.GetMaxBlockNumber(ctx, w.cfg.L2_RawBridgeEventsTableName)
	if err != nil {
		log.Error("failed to get L2 cross message synced height", "error", err)
		return 0, err
	}

	return messageSyncedHeight, nil
}

func (w *L2EventsWatcher) L2InsertOrUpdate(ctx context.Context, l2FetcherResult *L2FilterResult) error {
	if err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), w.cfg.L2_RawBridgeEventsTableName, l2FetcherResult.WithdrawMessages); err != nil {
		logrus.Error("failed to insert L2 deposit messages", "err", err)
		return err
	}

	if err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), w.cfg.L2_RawBridgeEventsTableName, l2FetcherResult.RelayedMessages); err != nil {
		logrus.Error("failed to insert L2 relayed messages", "err", err)
		return err
	}

	return nil
}

func (w *L2EventsWatcher) l2FetcherLogs(ctx context.Context, from, to uint64) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), // inclusive
		ToBlock:   new(big.Int).SetUint64(to),   // inclusive
		Addresses: w.contractAddressList,
		Topics:    make([][]common.Hash, 1),
	}
	query.Topics[0] = make([]common.Hash, 2)
	query.Topics[0][0] = backendabi.L2SentMessageEventSig
	query.Topics[0][1] = backendabi.L2RelayedMessageEventSig

	eventLogs, err := w.l2Client.FilterLogs(ctx, query)
	if err != nil {
		logrus.Error("failed to filter L2 event logs", "from", from, "to", to, "err", err)
		return nil, err
	}
	return eventLogs, nil
}
