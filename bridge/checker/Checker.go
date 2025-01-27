package checker

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/evm"
	"github.com/sirupsen/logrus"
)

type Checker struct {
	cfg               *evm.GethConfig
	l1Client          *ethclient.Client
	l1EventParser     *logic.L1EventParser
	rawBridgeEventOrm *orm.RawBridgeEvent
	crossMessageOrm   *orm.CrossMessage
	ctx               context.Context
	checkingSemaphore chan struct{}
}

// CalculateExpectedCount calculates the expected number of data entries between start and end (inclusive).
func (c *Checker) CalculateExpectedCount(start, end int) int {
	if start > end {
		return 0
	}
	return end - start + 1
}

// NewChecker creates a new Checker instance.
func NewChecker(ctx context.Context, cfg *evm.GethConfig, l1Client *ethclient.Client, rawBridgeEventOrm *orm.RawBridgeEvent, crossMessageOrm *orm.CrossMessage) *Checker {
	return &Checker{
		cfg:               cfg,
		l1Client:          l1Client,
		rawBridgeEventOrm: rawBridgeEventOrm,
		crossMessageOrm:   crossMessageOrm,
		ctx:               ctx,
		checkingSemaphore: make(chan struct{}, 1),
	}
}
func (c *Checker) StartChecking(rawBridgeEventTableName string, crossMessageTableName string, eventType int, clientAddress string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case c.checkingSemaphore <- struct{}{}:
				go func() {
					defer func() { <-c.checkingSemaphore }()
					if err := c.checkStep1(rawBridgeEventTableName, eventType, clientAddress); err != nil {
						logrus.Errorf("checkStep1 failed: %v", err)
					}
					if err := c.checkStep2(rawBridgeEventTableName, crossMessageTableName, eventType); err != nil {
						logrus.Errorf("checkStep2 failed: %v", err)
					}
				}()
			default:
				// skip this round if semaphore is full
			}
		case <-c.ctx.Done():
			return
		}
	}
}
func (c *Checker) checkStep1(rawBridgeEventTableName string, eventType int, clientAddress string) error {
	// 1. Query the latest unchecked message nonce
	latestUnCheckMessageNonce, err := c.rawBridgeEventOrm.GetMaxNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusUnChecked))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	checkStartMessageNonce := latestUnCheckMessageNonce
	checkEndMessageNonce := latestUnCheckMessageNonce + c.cfg.CheckerBatchSize

	// 1.1 Query the actual number of data entries
	actualCount, err := c.rawBridgeEventOrm.CountEventsByMessageNonceRange(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce)
	if err != nil {
		logrus.Errorf("Failed to count events: %v", err)
		return err
	}

	// 1.2 Calculate the expected number of data entries
	expectedCount := c.CalculateExpectedCount(checkStartMessageNonce, checkEndMessageNonce)

	// 1.3 Check if the actual count matches the expected count
	if actualCount != int64(expectedCount) {
		fmt.Printf("Actual count (%d) does not match expected count (%d)\n", actualCount, expectedCount)

		// 1.3.1 Query gaps
		gaps, err := c.rawBridgeEventOrm.FindMessageNonceGaps(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce)
		if err != nil {
			logrus.Errorf("Failed to find message nonce gaps: %v", err)
			return err
		}
		client, err := ethclient.Dial(clientAddress)
		if err != nil {
			logrus.Errorf("failed to connect to the Ethereum client: %v", err)
			return err
		}
		// 1.3.2 Process gaps
		for _, gap := range gaps {
			fmt.Printf("Gap from %d to %d\n", gap.StartGap, gap.EndGap)
			c.processGap(gap, client)
		}
	} else {
		// If the actual count matches the expected count, update the check_status to 1
		err := c.rawBridgeEventOrm.UpdateCheckStatusByNonceRange(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce, 1)
		if err != nil {
			log.Fatalf("Failed to update check status: %v", err)
		}
	}

	return nil
}

func (c *Checker) checkStep2(rawBridgeEventTableName string, crossMessageTableName string, eventType int) error {
	// 1. Query the latest unchecked message nonce
	latestUnCheckMessageNonce, err := c.rawBridgeEventOrm.GetMaxNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusCheckedStep1))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	checkStartMessageNonce := latestUnCheckMessageNonce
	checkEndMessageNonce := latestUnCheckMessageNonce + c.cfg.CheckerBatchSize

	// 1.1 Query the actual number of data entries
	rawBridgeEvents, err := c.rawBridgeEventOrm.GetEventsByMessageNonceRange(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce)
	if err != nil {
		logrus.Errorf("Failed to get events by message nonce range: %v", err)
		return err
	}

	// 1.2 Check if the corresponding crossMessage exists
	for _, event := range rawBridgeEvents {
		exists, err := c.crossMessageOrm.ExistsByMessageHash(crossMessageTableName, event.MessageHash)
		if err != nil {
			logrus.Errorf("Failed to check if cross message exists: %v", err)
			return err
		}
		if exists {
			// Update check_status to 2 if crossMessage exists
			err := c.rawBridgeEventOrm.UpdateCheckStatusByNonceRange(rawBridgeEventTableName, eventType, event.MessageNonce, event.MessageNonce, int(btypes.CheckStatusCheckedStep2))
			if err != nil {
				logrus.Errorf("Failed to update check status: %v", err)
				return err
			}
		} else {
			// Update check_fail_reason if crossMessage does not exist
			err := c.rawBridgeEventOrm.UpdateCheckFailReason(rawBridgeEventTableName, event.ID, int(btypes.CheckStatusCheckedStep1), "Cross message not found")
			if err != nil {
				logrus.Errorf("Failed to update check fail reason: %v", err)
				return err
			}
		}
	}

	return nil
}
func (c *Checker) processGap(gap orm.Gap, client *ethclient.Client) error {

	parentLayerContractAddress := common.HexToAddress(c.cfg.ParentLayerContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{parentLayerContractAddress},
		FromBlock: big.NewInt(int64(gap.StartBlockNumber)),
		ToBlock:   big.NewInt(int64(gap.EndBlockNumber)),
	}

	parentBridgeCoreFacetFilterer, err := contract.NewParentBridgeCoreFacetFilterer(parentLayerContractAddress, client)
	if err != nil {
		return nil
	}

	upwardMessageDispatcherFacetFilterer, err := contract.NewUpwardMessageDispatcherFacetFilterer(parentLayerContractAddress, client)
	if err != nil {
		return nil
	}
	var allBridgeEvents []*orm.RawBridgeEvent
	var queueEventCount, relayedEventCount int

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %v", err)
	}
	for _, vLog := range logs {
		switch vLog.Topics[0] {
		case backendabi.L1QueueTransactionEventSig:
			//fmt.Println("QueueTransaction event detected")
			event, err := parentBridgeCoreFacetFilterer.ParseQueueTransaction(vLog)
			if err != nil {
				return fmt.Errorf("failed to unpack event: %v", err)
			}
			bridgeEvents, err := c.l1EventParser.ParseDepositEventToRawBridgeEvents(context.Background(), event)
			if err != nil {
				return fmt.Errorf("failed to parse L1RelayedMessage: %v", err)
			}
			if len(bridgeEvents) > 0 {
				fmt.Println("bridgeEvents[0].MessageNonce", bridgeEvents[0].MessageNonce)
			}
			allBridgeEvents = append(allBridgeEvents, bridgeEvents...)
			queueEventCount++

		case backendabi.L1RelayedMessageEventSig:
			//fmt.Println("RelayedMessage event detected")
			event, err := upwardMessageDispatcherFacetFilterer.ParseRelayedMessage(vLog)
			if err != nil {
				return fmt.Errorf("failed to unpack event: %v", err)
			}
			bridgeEvents, err := c.l1EventParser.ParseL1RelayedMessageToRawBridgeEvents(context.Background(), event)
			if err != nil {
				return fmt.Errorf("failed to parse L1RelayedMessage: %v", err)
			}
			allBridgeEvents = append(allBridgeEvents, bridgeEvents...)
			relayedEventCount++
		}
	}
	err = c.rawBridgeEventOrm.InsertRawBridgeEvents(context.Background(), orm.TableRawBridgeEvents11155111, allBridgeEvents)
	if err != nil {
		return fmt.Errorf("failed to insert bridge events: %v", err)
	}
	return nil
}
