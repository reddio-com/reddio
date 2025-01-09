package logic

import (
	"context"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	yucommon "github.com/yu-org/yu/common"
	yucontext "github.com/yu-org/yu/core/context"
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
func (f *L2WatcherLogic) L2FetcherUpwardMessageFromLogs(ctx context.Context, block *yutypes.Block, l2BlockCollectionDepth *big.Int) ([]*orm.CrossMessage, map[uint64]uint64, error) {
	var allL2CrossMessages []*orm.CrossMessage

	depth := int(l2BlockCollectionDepth.Int64())
	blockHeight := block.Height
	blockTimestampsMap := make(map[uint64]uint64)

	var err error
	startBlockHeight := int(blockHeight) - depth
	endBlockHeight := int(blockHeight) - 2*depth

	for height := startBlockHeight; height > endBlockHeight; height-- {

		block, err = f.solidity.Chain.GetBlockByHeight(yucommon.BlockNum(height))
		if err != nil {
			//fmt.Println("Watcher GetCompactBlock error: ", err)
			return nil, nil, err
		}
		blockTimestampsMap[uint64(height)] = block.Timestamp
		query := ethereum.FilterQuery{
			// FromBlock: new(big.Int).SetUint64(from), // inclusive
			// ToBlock:   new(big.Int).SetUint64(to),   // inclusive
			Addresses: f.addressList,
			Topics:    make([][]common.Hash, 1),
		}
		query.Topics[0] = make([]common.Hash, 1)
		query.Topics[0][0] = backendabi.L2SentMessageEventSig

		eventLogs, err := f.FilterLogs(ctx, block, query)
		if err != nil {
			logrus.Error("FilterLogs err:", err)
			return nil, nil, err
		}
		if len(eventLogs) == 0 {
			continue
		}
		upwardMessages, err := f.parser.ParseL2EventLogs(ctx, eventLogs)
		if err != nil {
			logrus.Error("Failed to parse L2 event logs 3", "err", err)
			return nil, nil, err
		}
		allL2CrossMessages = append(allL2CrossMessages, upwardMessages...)
		blockHeight--
	}
	return allL2CrossMessages, blockTimestampsMap, nil
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
		receipt, err := f.solidity.GetEthReceipt(common.Hash(tx.TxnHash))
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (f *L2WatcherLogic) HandleRead(rdCall *yucommon.RdCall) (*yucontext.ResponseData, error) {
	ctx, err := yucontext.NewReadContext(rdCall)
	if err != nil {
		return nil, err
	}

	rd, err := f.solidity.Land.GetReading(rdCall.TripodName, rdCall.FuncName)
	if err != nil {
		return nil, err
	}
	rd(ctx)
	return ctx.Response(), nil
}
