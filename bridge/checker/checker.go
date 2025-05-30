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
	"gorm.io/gorm"
)

type Checker struct {
	cfg                 *evm.GethConfig
	l1EventParser       *logic.L1EventParser
	l2EventParser       *logic.L2EventParser
	rawBridgeEventOrm   *orm.RawBridgeEvent
	crossMessageOrm     *orm.CrossMessage
	ctx                 context.Context
	l1CheckingSemaphore chan struct{}
	l2CheckingSemaphore chan struct{}
}

// CalculateExpectedCount calculates the expected number of data entries between start and end (inclusive).
func (c *Checker) CalculateExpectedCount(start, end int) int {
	if start > end {
		return 0
	}
	return end - start + 1
}

// NewChecker creates a new Checker instance.
func NewChecker(ctx context.Context, cfg *evm.GethConfig, db *gorm.DB) *Checker {
	return &Checker{
		cfg:                 cfg,
		l1EventParser:       logic.NewL1EventParser(cfg),
		l2EventParser:       logic.NewL2EventParser(cfg),
		rawBridgeEventOrm:   orm.NewRawBridgeEvent(db, cfg),
		crossMessageOrm:     orm.NewCrossMessage(db),
		ctx:                 ctx,
		l1CheckingSemaphore: make(chan struct{}, 1),
		l2CheckingSemaphore: make(chan struct{}, 1),
	}
}
func (c *Checker) StartChecking() {
	// Ticker for Sepolia deposit
	tickerSepolia := time.NewTicker(time.Duration(c.cfg.BridgeCheckerConfig.SepoliaTickerInterval) * time.Second)
	defer tickerSepolia.Stop()

	// Ticker for L2 withdraw
	tickerReddio := time.NewTicker(time.Duration(c.cfg.BridgeCheckerConfig.ReddioTickerInterval) * time.Second)
	defer tickerReddio.Stop()

	for {
		select {
		// L1 checker
		case <-tickerSepolia.C:
			select {
			case c.l1CheckingSemaphore <- struct{}{}:
				go func() {
					defer func() { <-c.l1CheckingSemaphore }()
					if c.cfg.BridgeCheckerConfig.EnableL1CheckStep1 {
						if err := c.checkStep1(c.cfg.L1_RawBridgeEventsTableName, int(btypes.QueueTransaction), c.cfg.L1ClientAddress); err != nil {
							logrus.Errorf("checkStep1 for Sepolia deposit failed: %v", err)
						}
					}
					if c.cfg.BridgeCheckerConfig.EnableL1CheckStep2 {
						if err := c.checkStep2(c.cfg.L1_RawBridgeEventsTableName, int(btypes.QueueTransaction)); err != nil {
							logrus.Errorf("checkStep2 for Sepolia deposit failed: %v", err)
						}
					}
				}()
			default:
				// skip this round if semaphore is full
			}
		// L2 checker
		case <-tickerReddio.C:
			select {
			case c.l2CheckingSemaphore <- struct{}{}:
				go func() {
					defer func() { <-c.l2CheckingSemaphore }()
					if c.cfg.BridgeCheckerConfig.EnableL2CheckStep1 {
						if err := c.checkStep1(c.cfg.L2_RawBridgeEventsTableName, int(btypes.SentMessage), c.cfg.L2ClientAddress); err != nil {
							logrus.Errorf("checkStep1 for L2 withdraw failed: %v", err)
						}
					}
					if c.cfg.BridgeCheckerConfig.EnableL2CheckStep2 {
						if err := c.checkStep2(c.cfg.L2_RawBridgeEventsTableName, int(btypes.SentMessage)); err != nil {
							logrus.Errorf("checkStep2 for L2 withdraw failed: %v", err)
						}
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
	earliestUnCheckMessageNonce, err := c.rawBridgeEventOrm.GetMinNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusUnChecked))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	if earliestUnCheckMessageNonce == -1 {
		//fmt.Println("No unchecked message nonce found")
		return nil
	}
	maxMessageNonce, err := c.rawBridgeEventOrm.GetMaxNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusUnChecked))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	checkStartMessageNonce := earliestUnCheckMessageNonce
	checkEndMessageNonce := earliestUnCheckMessageNonce + c.cfg.BridgeCheckerConfig.CheckerBatchSize
	if checkEndMessageNonce > maxMessageNonce {
		checkEndMessageNonce = maxMessageNonce
	}
	//fmt.Printf("Checking message nonce range from %d to %d\n", checkStartMessageNonce, checkEndMessageNonce)
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
		//fmt.Printf("Actual count (%d) does not match expected count (%d)\n", actualCount, expectedCount)
		logrus.Infof("Actual count (%d) does not match expected count (%d)", actualCount, expectedCount)

		// 1.3.1 Query gaps
		gaps, err := c.rawBridgeEventOrm.FindMessageNonceGaps(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce)
		if err != nil {
			logrus.Errorf("Failed to find message nonce gaps: %v", err)
			return err
		}
		logrus.Infof("Found %d gaps", len(gaps))
		client, err := ethclient.Dial(clientAddress)
		if err != nil {
			logrus.Errorf("failed to connect to the Ethereum client: %v", err)
			return err
		}
		// 1.3.2 Process gaps
		for _, gap := range gaps {

			logrus.Infof("Gap from %d to %d:,starblock:%d,endblock:%d", gap.StartGap, gap.EndGap, gap.StartBlockNumber, gap.EndBlockNumber)

			//fmt.Printf("Gap from %d to %d\n", gap.StartGap, gap.EndGap)
			if rawBridgeEventTableName == c.cfg.L1_RawBridgeEventsTableName {
				logrus.Infof("Processing Sepolia deposit gap,start block number:%d,end block number:%d", gap.StartBlockNumber, gap.EndBlockNumber)
				//fmt.Println("Processing Sepolia deposit gap")
				err = c.processL1Gap(gap, client)
				if err != nil {
					logrus.Errorf("Failed to process L1 gap: %v", err)
					return err
				}
			} else if rawBridgeEventTableName == c.cfg.L2_RawBridgeEventsTableName {
				logrus.Infof("Processing L2 gap,start block number:%d,end block number:%d", gap.StartBlockNumber, gap.EndBlockNumber)
				//fmt.Println("Processing L2 withdraw gap")
				err = c.processL2Gap(gap, client)
				if err != nil {
					logrus.Errorf("Failed to process L2 gap: %v", err)
					return err
				}
			}
			//fmt.Println("Gap Processd,start gap", gap.StartGap, "end gap", gap.EndGap)
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

func (c *Checker) checkStep2(rawBridgeEventTableName string, eventType int) error {
	// 1. Query the latest unchecked message nonce
	earliestUnCheckMessageNonce, err := c.rawBridgeEventOrm.GetMinNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusCheckedStep1))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	if earliestUnCheckMessageNonce == -1 {
		//fmt.Println("checkStep2:No unchecked message nonce found")
		return nil
	}
	maxMessageNonce, err := c.rawBridgeEventOrm.GetMaxNonceByCheckStatus(rawBridgeEventTableName, eventType, int(btypes.CheckStatusCheckedStep1))
	if err != nil {
		logrus.Errorf("Failed to get max nonce by check status: %v", err)
		return err
	}
	checkStartMessageNonce := earliestUnCheckMessageNonce
	checkEndMessageNonce := earliestUnCheckMessageNonce + c.cfg.BridgeCheckerConfig.CheckerBatchSize
	if checkEndMessageNonce > maxMessageNonce {
		checkEndMessageNonce = maxMessageNonce
	}

	// 1.1 Query the actual number of data entries
	//fmt.Printf("checkStep2 Checking message nonce range from %d to %d\n", checkStartMessageNonce, checkEndMessageNonce)
	rawBridgeEvents, err := c.rawBridgeEventOrm.GetEventsByMessageNonceRange(rawBridgeEventTableName, eventType, checkStartMessageNonce, checkEndMessageNonce)
	if err != nil {
		logrus.Errorf("Failed to get events by message nonce range: %v", err)
		return err
	}

	// 1.2 Check if the corresponding crossMessage exists
	for _, event := range rawBridgeEvents {
		exists, err := c.crossMessageOrm.ExistsByMessageHash(event.MessageHash)
		if err != nil {
			logrus.Errorf("Failed to check if cross message exists: %v,tabelName:%s,MessageHash:%s", err, rawBridgeEventTableName, event.MessageHash)
			return err
		}
		if exists {
			// Update check_status to 2 if crossMessage exists
			err := c.rawBridgeEventOrm.UpdateCheckStatus(rawBridgeEventTableName, event.ID, int(btypes.CheckStatusCheckedStep2))
			if err != nil {
				logrus.Errorf("Failed to update check status: %v", err)
				return err
			}
		} else {
			// Update check_fail_reason if crossMessage does not exist
			err := c.rawBridgeEventOrm.UpdateCheckFailReason(rawBridgeEventTableName, event.ID, int(btypes.CheckStatusCheckedStep2), "checkStep2 failed:Cross message not found")
			if err != nil {
				logrus.Errorf("Failed to update check fail reason: %v", err)
				return err
			}
		}
	}

	return nil
}
func (c *Checker) processL1Gap(gap orm.Gap, client *ethclient.Client) error {
	//fmt.Println("Processing L1 gap，start block number", gap.StartBlockNumber, "end block number", gap.EndBlockNumber)
	parentLayerContractAddress := common.HexToAddress(c.cfg.BridgeCheckerConfig.CheckL1ContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{parentLayerContractAddress},
		FromBlock: big.NewInt(int64(gap.StartBlockNumber)),
		ToBlock:   big.NewInt(int64(gap.EndBlockNumber)),
	}

	parentBridgeCoreFacetFilterer, err := contract.NewParentBridgeCoreFacetFilterer(parentLayerContractAddress, client)
	if err != nil {
		//fmt.Println("failed to create parentBridgeCoreFacetFilterer")
		return nil
	}

	upwardMessageDispatcherFacetFilterer, err := contract.NewUpwardMessageDispatcherFacetFilterer(parentLayerContractAddress, client)
	if err != nil {
		//fmt.Println("failed to create upwardMessageDispatcherFacetFilterer")
		return nil
	}
	var allBridgeEvents []*orm.RawBridgeEvent
	var queueEventCount, relayedEventCount int

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %v", err)
	}
	//fmt.Println("logs length", len(logs))
	for _, vLog := range logs {
		switch vLog.Topics[0] {
		case backendabi.L1QueueTransactionEventSig:
			////fmt.Println("QueueTransaction event detected")
			event, err := parentBridgeCoreFacetFilterer.ParseQueueTransaction(vLog)
			if err != nil {
				return fmt.Errorf("failed to unpack event: %v", err)
			}
			bridgeEvents, err := c.l1EventParser.ParseDepositEventToRawBridgeEvents(context.Background(), event)
			if err != nil {
				return fmt.Errorf("failed to parse L1RelayedMessage: %v", err)
			}
			if len(bridgeEvents) > 0 {
				//fmt.Println("bridgeEvents[0].MessageNonce", bridgeEvents[0].MessageNonce)
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
	logrus.Infof("QueueTransaction event count: %d", queueEventCount)
	err = c.rawBridgeEventOrm.InsertRawBridgeEventsFromCheckStep1(context.Background(), c.cfg.L1_RawBridgeEventsTableName, allBridgeEvents)
	if err != nil {
		return fmt.Errorf("failed to insert bridge events: %v", err)
	}
	return nil
}
func (c *Checker) processL2Gap(gap orm.Gap, client *ethclient.Client) error {

	childLayerContractAddress := common.HexToAddress(c.cfg.BridgeCheckerConfig.CheckL2ContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{childLayerContractAddress},
		FromBlock: big.NewInt(int64(gap.StartBlockNumber)),
		ToBlock:   big.NewInt(int64(gap.EndBlockNumber)),
	}
	logrus.Infof("processL2Gap,start block number:%d,end block number:%d", gap.StartBlockNumber, gap.EndBlockNumber)
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %v", err)
	}
	l2WithdrawMessages, l2RelayedMessages, err := c.l2EventParser.ParseL2EventToRawBridgeEvents(context.Background(), logs)
	if err != nil {
		return fmt.Errorf("failed to parse L2 event: %v", err)
	}
	logrus.Infof("L2 withdraw messages count: %d", len(l2WithdrawMessages))
	err = c.rawBridgeEventOrm.InsertRawBridgeEventsFromCheckStep1(context.Background(), c.cfg.L2_RawBridgeEventsTableName, l2WithdrawMessages)
	if err != nil {
		return fmt.Errorf("failed to insert bridge l2WithdrawMessages: %v", err)
	}

	logrus.Infof("L2RelayedMessages messages count: %d", len(l2RelayedMessages))

	err = c.rawBridgeEventOrm.InsertRawBridgeEventsFromCheckStep1(context.Background(), c.cfg.L2_RawBridgeEventsTableName, l2RelayedMessages)
	if err != nil {
		return fmt.Errorf("failed to insert bridge l2RelayedMessages: %v", err)
	}
	return nil
}
