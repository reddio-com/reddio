package logic

import (
	"context"
	"math/big"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	yucommon "github.com/yu-org/yu/common"
	yutypes "github.com/yu-org/yu/core/types"

	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/evm"
)

type L2WatcherLogic struct {
	cfg         *evm.GethConfig
	addressList []common.Address
	parser      *L2EventParser
	solidity    *evm.Solidity `tripod:"solidity"`
}

func NewL2WatcherLogic(cfg *evm.GethConfig, solidity *evm.Solidity) (*L2WatcherLogic, error) {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ChildLayerContractAddress),
	}
	f := &L2WatcherLogic{
		cfg:         cfg,
		addressList: contractAddressList,
		parser:      NewL2EventParser(cfg),
		solidity:    solidity,
	}

	return f, nil
}

// L2FetcherUpwardMessageFromLogs collects upward messages from the logs of the current block
// and the previous l2BlockCollectionDepth blocks.
func (f *L2WatcherLogic) L2FetcherBridgeEventsFromLogs(ctx context.Context, block *yutypes.Block, l2BlockCollectionDepth *big.Int) ([]*orm.RawBridgeEvent, []*orm.RawBridgeEvent, map[uint64]uint64, error) {
	var l2WithdrawMessagesAll []*orm.RawBridgeEvent
	var l2RelayedMessagesAll []*orm.RawBridgeEvent

	depth := int(l2BlockCollectionDepth.Int64())
	blockHeight := block.Height
	blockTimestampsMap := make(map[uint64]uint64)

	var err error
	startBlockHeight := int(blockHeight) - depth
	endBlockHeight := int(blockHeight) - 2*depth

	for height := startBlockHeight; height > endBlockHeight; height-- {
		//fmt.Println("Watcher GetCompactBlock startBlockHeight: ", startBlockHeight)
		//fmt.Println("Watcher GetCompactBlock endBlockHeight: ", endBlockHeight)
		block, err = f.GetBlockWithRetry(yucommon.BlockNum(height), 5, 1*time.Second)
		if err != nil {
			//fmt.Println("Watcher GetCompactBlock error: ", err)
			logrus.Error("Watcher GetCompactBlock ,Height:", height, "error:", err)
			return nil, nil, nil, err
		}
		blockTimestampsMap[uint64(height)] = block.Timestamp
		query := ethereum.FilterQuery{
			// FromBlock: new(big.Int).SetUint64(from), // inclusive
			// ToBlock:   new(big.Int).SetUint64(to),   // inclusive
			Addresses: f.addressList,
			Topics:    make([][]common.Hash, 1),
		}
		query.Topics[0] = make([]common.Hash, 2)
		query.Topics[0][0] = backendabi.L2SentMessageEventSig
		query.Topics[0][1] = backendabi.L2RelayedMessageEventSig

		eventLogs, err := f.FilterLogs(ctx, block, query)
		if err != nil {
			logrus.Error("FilterLogs err:", err)
			return nil, nil, nil, err
		}
		if len(eventLogs) == 0 {
			continue
		}
		l2WithdrawMessages, l2RelayedMessages, err := f.parser.ParseL2EventLogs(ctx, eventLogs)
		if err != nil {
			logrus.Error("Failed to parse L2 event logs 3", "err", err)
			return nil, nil, nil, err
		}
		l2WithdrawMessagesAll = append(l2WithdrawMessagesAll, l2WithdrawMessages...)
		l2RelayedMessagesAll = append(l2RelayedMessagesAll, l2RelayedMessages...)
		blockHeight--
	}
	return l2WithdrawMessagesAll, l2RelayedMessagesAll, blockTimestampsMap, nil
}
func (f *L2WatcherLogic) GetBlockWithRetry(height yucommon.BlockNum, retries int, delay time.Duration) (*yutypes.Block, error) {
	var block *yutypes.Block
	var err error
	for i := 0; i < retries; i++ {
		block, err = f.solidity.Chain.GetBlockByHeight(height)
		if err == nil {
			return block, nil
		}
		logrus.Warnf("Retrying to get block, attempt %d/%d, height: %d, error: %v", i+1, retries, height, err)
		time.Sleep(delay)
	}
	return nil, err
}
func (f *L2WatcherLogic) FilterLogs(ctx context.Context, block *yutypes.Block, criteria ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := f.getLogs(ctx, block)
	if err != nil {
		return nil, err
	}

	result := make([]types.Log, 0)
	var logIdx uint
	for i, txLogs := range logs {
		for _, vLog := range txLogs {
			vLog.BlockHash = common.Hash(block.Hash)
			vLog.BlockNumber = uint64(block.Height)
			vLog.TxIndex = uint(i)
			vLog.Index = logIdx
			logIdx++

			//TODO
			if f.checkMatches(ctx, vLog) {
				result = append(result, *vLog)
			}
		}
	}

	return result, nil
}

func (f *L2WatcherLogic) checkMatches(ctx context.Context, vLog *types.Log) bool {
	if len(f.addressList) > 0 {
		if !slices.Contains(f.addressList, vLog.Address) {
			return false
		}
	}

	// TODO: The logic for topic filtering is a bit complex; it will not be implemented for now.
	//if len(f.topics) > len(vLog.Topics) {
	//	return false
	//}
	//for i, sub := range f.topics {
	//	if len(sub) == 0 {
	//		continue // empty rule set == wildcard
	//	}
	//	if !slices.Contains(sub, vLog.Topics[i]) {
	//		return false
	//	}
	//}

	return true
}
func (f *L2WatcherLogic) getLogs(ctx context.Context, block *yutypes.Block) ([][]*types.Log, error) {

	receipts, err := f.getReceipts(ctx, block)
	if err != nil {
		return nil, err
	}

	result := [][]*types.Log{}
	for _, receipt := range receipts {
		logs := []*types.Log{}
		logs = append(logs, receipt.Logs...)
		result = append(result, logs)
	}

	return result, nil
}

// param hash just for test , its gonna be removed in the final version
func (f *L2WatcherLogic) getReceipts(ctx context.Context, block *yutypes.Block) (types.Receipts, error) {

	var receipts []*types.Receipt
	for _, tx := range block.Txns {
		receipt, err := f.GetEthReceiptWithRetry(common.Hash(tx.TxnHash), 5, 1*time.Second)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
func (f *L2WatcherLogic) GetEthReceiptWithRetry(txHash common.Hash, retries int, delay time.Duration) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	receipt, err = f.solidity.GetEthReceipt(txHash)
	if err == nil {
		return receipt, nil
	}
	return nil, err
}
