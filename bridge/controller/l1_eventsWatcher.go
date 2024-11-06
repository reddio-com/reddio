package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/relayer"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
	"github.com/sirupsen/logrus"
)

type L1EventsWatcher struct {
	ctx            context.Context
	cfg            *evm.GethConfig
	l1Client       *ethclient.Client
	l1WatcherLogic *logic.L1WatcherLogic
	l1toL2Relayer  relayer.L1ToL2RelayerInterface
}

func NewL1EventsWatcher(ctx context.Context, cfg *evm.GethConfig, ethClient *ethclient.Client, l1toL2Relayer relayer.L1ToL2RelayerInterface) (*L1EventsWatcher, error) {

	c := &L1EventsWatcher{
		ctx:            ctx,
		cfg:            cfg,
		l1Client:       ethClient,
		l1WatcherLogic: logic.NewL1WatcherLogic(cfg, ethClient),
		l1toL2Relayer:  l1toL2Relayer,
	}
	return c, nil
}

func (w *L1EventsWatcher) Run(ctx context.Context) error {
	downwardMsgChan := make(chan *contract.ParentBridgeCoreFacetDownwardMessage)
	if w.l1Client.Client().SupportsSubscriptions() {
		sub, err := w.watchDownwardMessage(ctx, downwardMsgChan)
		if err != nil {
			metrics.L1EventWatcherFailureCounter.Inc()
			return err
		}
		go func() {
			defer sub.Unsubscribe()
			for {
				select {
				case msg := <-downwardMsgChan:
					if msg == nil {
						continue
					}
					// fmt.Println("Listen for msgChan", msg)
					// jsonData, err := json.Marshal(msg)
					// if err != nil {
					// 	logrus.Errorf("Error converting downwardMsgChan txn to JSON: %v", err)
					// 	continue
					// }
					// fmt.Println("msg as JSON:", string(jsonData))
					w.handleDownwardMessage(msg)
				case subErr := <-sub.Err():
					logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
					metrics.L1EventWatcherFailureCounter.Inc()
					metrics.L1EventWatcherRetryCounter.Inc()
					sub, err = w.watchDownwardMessage(ctx, downwardMsgChan)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
						metrics.L1EventWatcherFailureCounter.Inc()
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}

func (w *L1EventsWatcher) watchDownwardMessage(
	ctx context.Context,
	sink chan<- *contract.ParentBridgeCoreFacetDownwardMessage,
) (event.Subscription, error) {
	if !common.IsHexAddress(w.cfg.ParentLayerContractAddress) {
		return nil, fmt.Errorf("invalid address: %s", w.cfg.ParentLayerContractAddress)
	}
	filterer, err := contract.NewParentBridgeCoreFacetFilterer(common.HexToAddress(w.cfg.ParentLayerContractAddress), w.l1Client)
	if err != nil {
		return nil, err
	}
	return filterer.WatchDownwardMessage(&bind.WatchOpts{Context: ctx}, sink)
}

/*****************************
 *    [Functions:Handler]    *
 *****************************/
func (w *L1EventsWatcher) handleDownwardMessage(
	msg *contract.ParentBridgeCoreFacetDownwardMessage,
) error {
	err := w.l1toL2Relayer.HandleDownwardMessageWithSystemCall(msg)
	if err != nil {
		return err
	}
	return nil
}

func (w *L1EventsWatcher) ChainID(ctx context.Context) (*big.Int, error) {
	return w.l1Client.ChainID(ctx)
}

func (w *L1EventsWatcher) Close() {
	w.l1Client.Close()
}
