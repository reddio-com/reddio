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
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/reddio-com/reddio/watcher/contract"
)

type ReddioSubscriber struct {
	ethClient       *ethclient.Client
	client          *rpc.Client
	filterer        *contract.ChildBridgeCoreFacetFilterer
	processedEvents map[string]bool
}

func NewReddioSubscriber(clientAddress string, coreContractAddress common.Address) (*ReddioSubscriber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// TODO replace with our own client once we have one.
	// Geth pulls in a lot of dependencies that we don't use.
	client, err := rpc.DialContext(ctx, clientAddress)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)
	filterer, err := contract.NewChildBridgeCoreFacetFilterer(coreContractAddress, ethClient)
	if err != nil {
		return nil, err
	}
	return &ReddioSubscriber{
		ethClient: ethClient,
		client:    client,
		filterer:  filterer,
	}, nil
}

func (s *ReddioSubscriber) WatchUpwardMessage(
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

func (s *ReddioSubscriber) PollUpwardMessage(contractAddress common.Address, eventSignature string) {
	// Parse contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(contract.ChildBridgeCoreFacetABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Get the event signature hash
	eventSignatureHash := contractAbi.Events[eventSignature].ID
	fmt.Printf("eventSignatureHash: %s\n", eventSignatureHash.Hex())

	// Use a goroutine to handle event logs
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		var lastBlock uint64 = 0 // Starting block
		const bufferBlocks = 6   // Buffer blocks

		for range ticker.C {
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
					contractAddress,
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
					From    common.Address
					To      common.Address
					Message string
				}
				err := contractAbi.UnpackIntoInterface(&upwardMessageEvent, "UpwardMessage", vLog.Data)
				if err != nil {
					log.Fatalf("Failed to unpack log data: %v", err)
				}

				// Parse indexed fields
				upwardMessageEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
				upwardMessageEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
				fmt.Printf("UpwardMessage event: From=%s, To=%s, Message=%s\n", upwardMessageEvent.From.Hex(), upwardMessageEvent.To.Hex(), upwardMessageEvent.Message)

				// Record the processed event
				s.processedEvents[eventID] = true
			}

			// Update lastBlock to the next block after the latest block
			lastBlock = latestBlock + 1
		}
	}()
}
