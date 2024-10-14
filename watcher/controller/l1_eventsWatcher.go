package controller

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/relayer"
	"github.com/reddio-com/reddio/watcher/contract"
	"github.com/reddio-com/reddio/watcher/logic"
	"github.com/sirupsen/logrus"
)

type L1EventsWatcher struct {
	ctx            context.Context
	cfg            *evm.GethConfig
	ethClient      *ethclient.Client
	l1WatcherLogic *logic.L1WatcherLogic
	bridgeRelayer  relayer.BridgeRelayerInterface
}

func NewL1EventsWatcher(ctx context.Context, cfg *evm.GethConfig, ethClient *ethclient.Client, bridgeRelayer relayer.BridgeRelayerInterface) (*L1EventsWatcher, error) {

	c := &L1EventsWatcher{
		ctx:            ctx,
		cfg:            cfg,
		ethClient:      ethClient,
		l1WatcherLogic: logic.NewL1WatcherLogic(cfg, ethClient),
		bridgeRelayer:  bridgeRelayer,
	}
	return c, nil
}

func (w *L1EventsWatcher) Run(cfg *evm.GethConfig, ctx context.Context) error {
	downwardMsgChan := make(chan *contract.ParentBridgeCoreFacetDownwardMessage)
	if w.ethClient.Client().SupportsSubscriptions() {
		sub, err := w.watchDownwardMessage(ctx, downwardMsgChan, nil)
		if err != nil {
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
					//fmt.Println("Listen for msgChan", msg)
					// jsonData, err := json.Marshal(msg)
					// if err != nil {
					// 	logrus.Errorf("Error converting downwardMsgChan txn to JSON: %v", err)
					// 	continue
					// }
					// fmt.Println("msg as JSON:", string(jsonData))
					w.handleDownwardMessage(msg)
					//fmt.Println("handleDownwardMessage end")
				case subErr := <-sub.Err():
					logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
					sub, err = w.watchDownwardMessage(ctx, downwardMsgChan, nil)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
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
	sequence []*big.Int,
) (event.Subscription, error) {
	filterer, err := contract.NewParentBridgeCoreFacetFilterer(common.HexToAddress(w.cfg.ParentLayerContractAddress), w.ethClient)
	if err != nil {
		return nil, err
	}
	return filterer.WatchDownwardMessage(&bind.WatchOpts{Context: ctx}, sink, sequence)
}

/*****************************
 *    [Functions:Handler]    *
 *****************************/
func (w *L1EventsWatcher) handleDownwardMessage(
	msg *contract.ParentBridgeCoreFacetDownwardMessage,
) error {
	err := w.bridgeRelayer.HandleDownwardMessageWithSystemCall(msg)
	if err != nil {
		return err
	}
	return nil
}

func (w *L1EventsWatcher) ChainID(ctx context.Context) (*big.Int, error) {
	return w.ethClient.ChainID(ctx)
}

func (w *L1EventsWatcher) Close() {
	w.ethClient.Close()
}
