package controller

import (
	"context"
	"fmt"
	"math/big"

	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/tripod"
	yutypes "github.com/yu-org/yu/core/types"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
)

type L2EventsWatcher struct {
	//ctx context.Context
	cfg *evm.GethConfig
	//ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	*tripod.Tripod
	solidity           *evm.Solidity `tripod:"solidity"`
	rawBridgeEventsOrm *orm.RawBridgeEvent
	db                 *gorm.DB
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

func (w *L2EventsWatcher) WatchL2BridgeEvent(ctx context.Context, block *yutypes.Block, Solidity *evm.Solidity) error {
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

	for _, event := range l2WithdrawMessages {
		if event.EventType == int(btypes.SentMessage) {
			metrics.WithdrawMessageNonceGauge.WithLabelValues("withdrawMessageNonce").Set(float64(event.MessageNonce))
		}
	}
	return nil
}

func (w *L2EventsWatcher) InitChain(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		w.rawBridgeEventsOrm = orm.NewRawBridgeEvent(w.db)

		l2WatcherLogic, err := logic.NewL2WatcherLogic(w.cfg, w.solidity)
		if err != nil {
			logrus.Fatal("init l2WatcherLogic failed: ", err)
		}
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
			go func() {
				err := w.WatchL2BridgeEvent(context.Background(), block, w.solidity)
				if err != nil {
					logrus.Errorf("WatchUpwardMessage error: %v", err)
				}
			}()
		}
	}
}
func (w *L2EventsWatcher) savel2BridgeEvents(
	rawBridgeEvents []*orm.RawBridgeEvent,
) error {
	//fmt.Println("savel2BridgeEvents rawBridgeEvents: ", rawBridgeEvents)
	if len(rawBridgeEvents) == 0 {
		return nil
	}
	err := w.rawBridgeEventsOrm.InsertRawBridgeEvents(context.Background(), orm.TableRawBridgeEvents50341, rawBridgeEvents)
	if err != nil {
		return err
	}

	return nil
}
