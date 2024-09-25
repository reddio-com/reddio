package watcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/contract"
	"github.com/sirupsen/logrus"
)

type EventWatcher struct {
	l1Watcher *EthSubscriber
	l2Watcher *ReddioSubscriber
}

func NewEventWatcher(cfg *evm.GethConfig) (*EventWatcher, error) {

	l1Watcher, err := NewEthSubscriber(cfg.L1ClientAddress, common.HexToAddress(cfg.ParentLayerContractAddress))
	if err != nil {
		return nil, err
	}
	l2Watcher, err := NewReddioSubscriber(cfg.L2ClientAddress, common.HexToAddress(cfg.ParentLayerContractAddress))
	if err != nil {
		return nil, err
	}
	return &EventWatcher{
		l1Watcher: l1Watcher,
		l2Watcher: l2Watcher,
	}, nil
}

func StartupEventWatcher(cfg *evm.GethConfig) {
	if cfg.EnableEventWatcher {
		eventWatcher, err := NewEventWatcher(cfg)
		if err != nil {
			logrus.Fatal("init L1 client failed: ", err)
		}
		err = eventWatcher.Run(context.Background())
		if err != nil {
			logrus.Fatal("l1 client run failed: ", err)
		}

	}
}

func (w *EventWatcher) Run(ctx context.Context) error {
	msgChan := make(chan *contract.ParentBridgeCoreFacetDownwardMessage)
	sub, err := w.l1Watcher.WatchDownwardMessage(ctx, msgChan, nil)
	if err != nil {
		return err
	}

	// Monitor L1 chain
	go func() {
		for {
			select {
			case msg := <-msgChan:
				fmt.Println("Listen for msgChan", msg)
				jsonData, err := json.Marshal(msg)
				if err != nil {
					logrus.Errorf("Error converting L1 txn to JSON: %v", err)
					continue
				}
				fmt.Println("msg as JSON:", string(jsonData))
			case subErr := <-sub.Err():
				logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
				sub.Unsubscribe()

				sub, err = w.l1Watcher.WatchDownwardMessage(ctx, msgChan, nil)
				if err != nil {
					logrus.Errorf("Resubscribe failed: %v", err)
				}
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			}
		}
	}()

	// Monitor L2 chain
	// go func() {
	// 	msgChanL2 := make(chan *contract.ChildBridgeCoreFacetUpwardMessage)
	// 	subL2, err := w.ethL2.WatchUpwardMessage(ctx, msgChanL2, nil)
	// 	if err != nil {
	// 		logrus.Fatal("Failed to watch L2 upward message: ", err)
	// 	}
	// 	defer subL2.Unsubscribe()

	// 	for msg := range msgChanL2 {
	// 		l.handleUpwardMessage(msg)
	// 	}
	// }()
	return nil
}
