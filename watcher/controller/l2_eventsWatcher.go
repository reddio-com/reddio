package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

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

type L2EventsWatcher struct {
	ctx            context.Context
	cfg            *evm.GethConfig
	ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	bridgeRelayer  relayer.BridgeRelayerInterface
}

func NewL2EventsWatcher(ctx context.Context, cfg *evm.GethConfig, ethClient *ethclient.Client, bridgeRelayer relayer.BridgeRelayerInterface) (*L2EventsWatcher, error) {

	c := &L2EventsWatcher{
		ctx:            ctx,
		cfg:            cfg,
		ethClient:      ethClient,
		l2WatcherLogic: logic.NewL2WatcherLogic(cfg, ethClient),
		bridgeRelayer:  bridgeRelayer,
	}
	return c, nil
}

func (w *L2EventsWatcher) Run(cfg *evm.GethConfig, ctx context.Context) error {
	upwardMsgChan := make(chan *contract.ChildBridgeCoreFacetUpwardMessage)
	if w.ethClient.Client().SupportsSubscriptions() {
		sub, err := w.WatchUpwardMessageWss(ctx, upwardMsgChan, nil)
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case msg := <-upwardMsgChan:
					fmt.Println("WatchUpwardMessageWss, msgChan: ", msg)
					jsonData, err := json.Marshal(msg)
					if err != nil {
						logrus.Errorf("Error converting upwardMsgChan txn to JSON: %v", err)
						continue
					}
					fmt.Println("msg as JSON:", string(jsonData))
					fmt.Println("handleupwardMessage")
					fmt.Println("handleUpwardMessage end")
				case subErr := <-sub.Err():
					logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
					sub.Unsubscribe()
					sub, err = w.WatchUpwardMessageWss(ctx, upwardMsgChan, nil)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
					}
				case <-ctx.Done():
					sub.Unsubscribe()
					return
				}
			}
		}()
	} else {
		err := w.WatchUpwardMessageHttp(ctx, upwardMsgChan, nil)
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case msg := <-upwardMsgChan:
					fmt.Println(": ", msg)
					jsonData, err := json.Marshal(msg)
					if err != nil {
						logrus.Errorf("Error converting upwardMsgChan txn to JSON: %v", err)
						continue
					}
					fmt.Println("msg as JSON:", string(jsonData))
					fmt.Println("handleupwardMessage")
					fmt.Println("handleUpwardMessage end")
				case <-ctx.Done():
					fmt.Println("Context done, stopping event processing")
					return
				}
			}
		}()
	}
	return nil
}

func (w *L2EventsWatcher) Close() {
	w.ethClient.Close()
}

func (w *L2EventsWatcher) WatchUpwardMessageWss(
	ctx context.Context,
	sink chan<- *contract.ChildBridgeCoreFacetUpwardMessage,
	sequence []*big.Int,
) (event.Subscription, error) {
	filterer, err := contract.NewChildBridgeCoreFacetFilterer(common.HexToAddress(w.cfg.ChildLayerContractAddress), w.ethClient)
	if err != nil {
		return nil, err
	}
	return filterer.WatchUpwardMessage(&bind.WatchOpts{Context: ctx}, sink, sequence)
}

func (w *L2EventsWatcher) WatchUpwardMessageHttp(ctx context.Context,
	sink chan<- *contract.ChildBridgeCoreFacetUpwardMessage,
	sequence []*big.Int) error {
	// Use a goroutine to handle event logs
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		var lastBlock uint64 = 0 // Starting block
		const bufferBlocks = 1   // Buffer blocks

		for {
			select {
			case <-ticker.C:
				// Get the latest block number
				header, err := w.ethClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					log.Fatalf("Failed to get latest block header: %v", err)
				}
				latestBlock := header.Number.Uint64()

				if lastBlock == 0 {
					lastBlock = latestBlock
				}
				// Set filter query
				fromBlock := lastBlock
				if fromBlock > bufferBlocks {
					fromBlock -= bufferBlocks
				} else {
					fromBlock = 0
				}

				upwardMessage, err := w.l2WatcherLogic.L2FetcherUpwardMessageFromLogs(ctx, fromBlock, latestBlock)
				if err != nil {
					log.Fatalf("Failed to fetch L2 event logs", "from", fromBlock, "to", latestBlock, "err", err)
					continue
				}
				if upwardMessage == nil {
					lastBlock = latestBlock + 1
					continue
				}
				select {
				case sink <- upwardMessage:
					w.bridgeRelayer.HandleUpwardMessage(upwardMessage)
					fmt.Println("Event sent to sink channel")
				case <-ctx.Done():
					fmt.Println("Context done, stopping event processing")
					return
				}

				// Update lastBlock to the next block after the latest block
				lastBlock = latestBlock + 1
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

/*****************************
 *     [Functions:Handler]   *
 *****************************/
