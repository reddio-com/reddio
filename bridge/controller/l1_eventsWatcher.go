package controller

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/sirupsen/logrus"
	
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/relayer"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
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
	downwardMsgChan := make(chan *contract.ParentBridgeCoreFacetQueueTransaction)
	relayerMsgChan := make(chan *contract.UpwardMessageDispatcherFacetRelayedMessage)

	if w.l1Client.Client().SupportsSubscriptions() {
		downwardSub, err := w.watchDownwardMessage(ctx, downwardMsgChan)
		if err != nil {
			metrics.L1EventWatcherFailureCounter.Inc()
			return err
		}
		relayerSub, err := w.watchRelayerMessage(ctx, relayerMsgChan)
		if err != nil {
			metrics.L1EventWatcherFailureCounter.Inc()
			return err
		}

		go func() {
			defer downwardSub.Unsubscribe()
			defer relayerSub.Unsubscribe()
			for {
				select {
				case msg := <-downwardMsgChan:
					if msg == nil {
						continue
					}
					w.handleDownwardMessage(msg)
				case msg := <-relayerMsgChan:
					if msg == nil {
						continue
					}
					w.handleRelayerMessage(msg)
				case subErr := <-downwardSub.Err():
					logrus.Errorf("L1 downward subscription failed: %v, Resubscribing...", subErr)
					metrics.L1EventWatcherFailureCounter.Inc()
					metrics.L1EventWatcherRetryCounter.Inc()
					downwardSub, err = w.watchDownwardMessage(ctx, downwardMsgChan)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
						metrics.L1EventWatcherFailureCounter.Inc()
						return
					}
				case subErr := <-relayerSub.Err():
					logrus.Errorf("L1 relayer subscription failed: %v, Resubscribing...", subErr)
					metrics.L1EventWatcherFailureCounter.Inc()
					metrics.L1EventWatcherRetryCounter.Inc()
					relayerSub, err = w.watchRelayerMessage(ctx, relayerMsgChan)
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

// func (w *L1EventsWatcher) watchDownwardMessage(
//
//	ctx context.Context,
//	sink chan<- *contract.ParentBridgeCoreFacetDownwardMessage,
//
//	) (event.Subscription, error) {
//		if !common.IsHexAddress(w.cfg.ParentLayerContractAddress) {
//			return nil, fmt.Errorf("invalid address: %s", w.cfg.ParentLayerContractAddress)
//		}
//		filterer, err := contract.NewParentBridgeCoreFacetFilterer(common.HexToAddress(w.cfg.ParentLayerContractAddress), w.l1Client)
//		if err != nil {
//			return nil, err
//		}
//		return filterer.WatchDownwardMessage(&bind.WatchOpts{Context: ctx}, sink)
//	}
func (w *L1EventsWatcher) watchDownwardMessage(
	ctx context.Context,
	sink chan<- *contract.ParentBridgeCoreFacetQueueTransaction,
) (event.Subscription, error) {
	if !common.IsHexAddress(w.cfg.ParentLayerContractAddress) {
		return nil, fmt.Errorf("invalid address: %s", w.cfg.ParentLayerContractAddress)
	}
	filterer, err := contract.NewParentBridgeCoreFacetFilterer(common.HexToAddress(w.cfg.ParentLayerContractAddress), w.l1Client)
	if err != nil {
		return nil, err
	}
	return filterer.WatchQueueTransaction(&bind.WatchOpts{Context: ctx}, sink, nil, nil)
}
func (w *L1EventsWatcher) watchRelayerMessage(
	ctx context.Context,
	sink chan<- *contract.UpwardMessageDispatcherFacetRelayedMessage,
) (event.Subscription, error) {
	if !common.IsHexAddress(w.cfg.ParentLayerContractAddress) {
		return nil, fmt.Errorf("invalid address: %s", w.cfg.ParentLayerContractAddress)
	}
	filterer, err := contract.NewUpwardMessageDispatcherFacetFilterer(common.HexToAddress(w.cfg.ParentLayerContractAddress), w.l1Client)
	if err != nil {
		return nil, err
	}
	return filterer.WatchRelayedMessage(&bind.WatchOpts{Context: ctx}, sink, nil)
}

/*****************************
 *    [Functions:Handler]    *
 *****************************/
func (w *L1EventsWatcher) handleDownwardMessage(
	msg *contract.ParentBridgeCoreFacetQueueTransaction,
) error {
	err := w.l1toL2Relayer.HandleDownwardMessageWithSystemCall(msg)
	if err != nil {
		return err
	}
	return nil
}

func (w *L1EventsWatcher) handleRelayerMessage(msg *contract.UpwardMessageDispatcherFacetRelayedMessage) error {
	err := w.l1toL2Relayer.HandleRelayerMessage(msg)
	if err != nil {
		logrus.Errorf("Failed to handle RelayerMessage: %v", err)
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
