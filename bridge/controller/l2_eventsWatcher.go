package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/tripod"
	yutypes "github.com/yu-org/yu/core/types"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/relayer"
	"github.com/reddio-com/reddio/evm"
)

type L2EventsWatcher struct {
	//ctx context.Context
	cfg *evm.GethConfig
	//ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	l2toL1Relayer  relayer.L2ToL1RelayerInterface
	*tripod.Tripod
	solidity *evm.Solidity `tripod:"solidity"`
	db       *gorm.DB
}

func NewL2EventsWatcher(cfg *evm.GethConfig, db *gorm.DB) *L2EventsWatcher {
	tri := tripod.NewTripod()
	c := &L2EventsWatcher{
		//ctx:            ctx,
		cfg:    cfg,
		Tripod: tri,
		db:     db,
	}
	return c
}

func (w *L2EventsWatcher) WatchUpwardMessage(ctx context.Context, block *yutypes.Block, Solidity *evm.Solidity) error {
	upwardMessage, blockTimestampsMap, err := w.l2WatcherLogic.L2FetcherUpwardMessageFromLogs(ctx, block, w.cfg.L2BlockCollectionDepth)
	if err != nil {
		return fmt.Errorf("failed to fetch upward message from logs: %v", err)
	}

	if len(upwardMessage) == 0 {
		return nil
	}

	err = w.l2toL1Relayer.HandleUpwardMessage(upwardMessage, blockTimestampsMap)
	if err != nil {
		return fmt.Errorf("failed to handle upward message: %v", err)
	}
	return nil
}

func (w *L2EventsWatcher) InitChain(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		l1Client, err := ethclient.Dial(w.cfg.L1ClientAddress)
		if err != nil {
			logrus.Fatal("failed to connect to L1 geth", "endpoint", w.cfg.L1ClientAddress, "err", err)
		}

		l2toL1Relayer, err := relayer.NewL2ToL1Relayer(context.Background(), w.cfg, l1Client, w.db)
		if err != nil {
			logrus.Fatal("init bridge relayer failed: ", err)
		}
		l2WatcherLogic, err := logic.NewL2WatcherLogic(w.cfg, w.solidity)
		if err != nil {
			logrus.Fatal("init l2WatcherLogic failed: ", err)
		}

		w.l2toL1Relayer = l2toL1Relayer
		w.l2WatcherLogic = l2WatcherLogic
	}
}

func (w *L2EventsWatcher) StartBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcher) EndBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcher) FinalizeBlock(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		//watch upward message
		blockHeightBigInt := big.NewInt(int64(block.Header.Height))
		if big.NewInt(0).Mod(blockHeightBigInt, w.cfg.L2BlockCollectionDepth).Cmp(big.NewInt(0)) == 0 {
			err := w.WatchUpwardMessage(context.Background(), block, w.solidity)
			if err != nil {
				logrus.Errorf("WatchUpwardMessage error: %v", err)
			}
		}
	}
}
