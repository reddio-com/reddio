package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/tripod"
	yutypes "github.com/yu-org/yu/core/types"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/evm"
)

type L2EventsWatcherTripod struct {
	//ctx context.Context
	cfg *evm.GethConfig
	//ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	*tripod.Tripod
	solidity           *evm.Solidity `tripod:"solidity"`
	rawBridgeEventsOrm *orm.RawBridgeEvent
	db                 *gorm.DB
}

func NewL2EventsWatcherTripod(cfg *evm.GethConfig, db *gorm.DB) *L2EventsWatcherTripod {
	tri := tripod.NewTripod()
	c := &L2EventsWatcherTripod{
		//ctx:            ctx,
		cfg:    cfg,
		Tripod: tri,
		db:     db,
	}
	return c
}

func (w *L2EventsWatcherTripod) WatchL2BridgeEvent(ctx context.Context, block *yutypes.Block, Solidity *evm.Solidity) error {
	l2WithdrawMessages, l2RelayedMessages, _, err := w.l2WatcherLogic.L2FetcherBridgeEventsFromLogs(ctx, block, w.cfg.L2BlockCollectionDepth)
	if err != nil {
		return fmt.Errorf("failed to fetch upward message from logs: %v", err)
	}
	if len(l2WithdrawMessages) != 0 {
		err = w.savel2BridgeEvents(l2WithdrawMessages)
		if err != nil {
			return fmt.Errorf("failed to save l2WithdrawMessages: %v", err)
		}
	}
	if len(l2RelayedMessages) != 0 {
		err = w.savel2BridgeEvents(l2RelayedMessages)
		if err != nil {
			return fmt.Errorf("failed to save l2RelayedMessages: %v", err)
		}
	}

	return nil
}

func (w *L2EventsWatcherTripod) InitChain(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		w.rawBridgeEventsOrm = orm.NewRawBridgeEvent(w.db, w.cfg)

		l2WatcherLogic, err := logic.NewL2WatcherLogic(w.cfg, w.solidity)
		if err != nil {
			logrus.Fatal("init l2WatcherLogic failed: ", err)
		}
		w.l2WatcherLogic = l2WatcherLogic
	}
}

func (w *L2EventsWatcherTripod) StartBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcherTripod) EndBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcherTripod) FinalizeBlock(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		//watch upward message
		blockHeightBigInt := big.NewInt(int64(block.Header.Height))
		if big.NewInt(0).Mod(blockHeightBigInt, w.cfg.L2BlockCollectionDepth).Cmp(big.NewInt(0)) == 0 {
			go func() {
				err := w.WatchL2BridgeEvent(context.Background(), block, w.solidity)
				if err != nil {
					logrus.Errorf("WatchUpwardMessage error: %v", err)
				}
			}()
		}
	}
}
func (w *L2EventsWatcherTripod) savel2BridgeEvents(
	rawBridgeEvents []*orm.RawBridgeEvent,
) error {
	//fmt.Println("savel2BridgeEvents rawBridgeEvents: ", rawBridgeEvents)
	if len(rawBridgeEvents) == 0 {
		return nil
	}
	err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), w.cfg.L2_RawBridgeEventsTableName, rawBridgeEvents)
	if err != nil {
		return err
	}

	return nil
}
