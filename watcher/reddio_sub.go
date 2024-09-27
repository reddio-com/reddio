package watcher

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/contract"
	"github.com/yu-org/yu/core/kernel"
)

type ReddioSubscriber struct {
	ethClient                 *ethclient.Client
	client                    *rpc.Client
	filterer                  *contract.ChildBridgeCoreFacetFilterer
	processedEvents           map[string]bool
	chain                     *kernel.Kernel
	chainConfig               *params.ChainConfig
	childLayerContractAddress common.Address
}

func NewReddioSubscriber(chain *kernel.Kernel, cfg *evm.GethConfig) (*ReddioSubscriber, error) {
	clientAddress := cfg.L2ClientAddress
	childLayerContractAddress := common.HexToAddress(cfg.ChildLayerContractAddress)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// TODO replace with our own client once we have one.
	// Geth pulls in a lot of dependencies that we don't use.
	client, err := rpc.DialContext(ctx, clientAddress)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)
	filterer, err := contract.NewChildBridgeCoreFacetFilterer(childLayerContractAddress, ethClient)
	if err != nil {
		return nil, err
	}
	return &ReddioSubscriber{
		ethClient:                 ethClient,
		client:                    client,
		filterer:                  filterer,
		processedEvents:           make(map[string]bool),
		chain:                     chain,
		chainConfig:               cfg.ChainConfig,
		childLayerContractAddress: childLayerContractAddress,
	}, nil
}

// Note: this function only works when reddio client is wss mode
func (s *ReddioSubscriber) WatchUpwardMessageWss(
	ctx context.Context,
	sink chan<- *contract.ChildBridgeCoreFacetUpwardMessage,
	sequence []*big.Int,
) (event.Subscription, error) {
	return s.filterer.WatchUpwardMessage(&bind.WatchOpts{Context: ctx}, sink, sequence)
}

func (s *ReddioSubscriber) ChainID(ctx context.Context) (*big.Int, error) {
	return s.ethClient.ChainID(ctx)
}

func (s *ReddioSubscriber) Close() {
	s.ethClient.Close()
}

func (s *ReddioSubscriber) WatchUpwardMessageHttp(ctx context.Context,
	sink chan<- *contract.ChildBridgeCoreFacetUpwardMessage,
	sequence []*big.Int) error {
	// Parse contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(contract.ChildBridgeCoreFacetABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Get the event signature hash
	eventSignatureHash := contractAbi.Events["UpwardMessage"].ID
	fmt.Printf("eventSignatureHash: %s\n", eventSignatureHash.Hex())

	// Use a goroutine to handle event logs
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		var lastBlock uint64 = 0 // Starting block
		const bufferBlocks = 6   // Buffer blocks

		for {
			select {
			case <-ticker.C:
				// Get the latest block number
				header, err := s.ethClient.HeaderByNumber(context.Background(), nil)
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

				query := ethereum.FilterQuery{
					FromBlock: big.NewInt(int64(fromBlock)),
					ToBlock:   big.NewInt(int64(latestBlock)),
					Addresses: []common.Address{
						s.childLayerContractAddress,
					},
					Topics: [][]common.Hash{
						{eventSignatureHash},
					},
				}

				// Get event logs
				logs, err := s.ethClient.FilterLogs(context.Background(), query)
				if err != nil {
					log.Fatalf("Failed to filter logs: %v", err)
				}

				// Handle event logs
				for _, vLog := range logs {
					eventID := vLog.TxHash.Hex() + fmt.Sprintf("%d", vLog.Index)
					if s.processedEvents[eventID] {
						// If the event has already been processed, skip it
						continue
					}

					fmt.Printf("Log block number: %d\n", vLog.BlockNumber)

					// Parse log data
					var upwardMessageEvent struct {
						Sequence    *big.Int
						PayloadType uint32
						Payload     []byte
					}
					err := contractAbi.UnpackIntoInterface(&upwardMessageEvent, "UpwardMessage", vLog.Data)
					if err != nil {
						log.Fatalf("Failed to unpack log data: %v", err)
					}

					fmt.Printf("UpwardMessage event: Sequence=%s, PayloadType=%d, Payload=%x\n", upwardMessageEvent.Sequence.String(), upwardMessageEvent.PayloadType, upwardMessageEvent.Payload)

					// Create a new ChildBridgeCoreFacetUpwardMessage instance
					upwardMessage := &contract.ChildBridgeCoreFacetUpwardMessage{
						Sequence:    upwardMessageEvent.Sequence,
						PayloadType: upwardMessageEvent.PayloadType,
						Payload:     upwardMessageEvent.Payload,
						Raw:         vLog,
					}

					// Send the event to the sink channel
					select {
					case sink <- upwardMessage:
						fmt.Println("Event sent to sink channel")
					case <-ctx.Done():
						fmt.Println("Context done, stopping event processing")
						return
					}

					// Record the processed event
					s.processedEvents[eventID] = true
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
