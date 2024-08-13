package ethrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	yutypes "github.com/yu-org/yu/core/types"
	"math/big"
	"slices"
)

var (
	errInvalidTopic           = errors.New("invalid topic(s)")
	errFilterNotFound         = errors.New("filter not found")
	errInvalidBlockRange      = errors.New("invalid block range params")
	errPendingLogsUnsupported = errors.New("pending logs are not supported")
	errExceedMaxTopics        = errors.New("exceed max topics")
)

const maxTopics = 100
const maxSubTopics = 1000

// FilterCriteria represents a request to create a new filter.
// Same as ethereum.FilterQuery but with UnmarshalJSON() method.
type FilterCriteria ethereum.FilterQuery

// UnmarshalJSON sets *args fields with given data.
func (args *FilterCriteria) UnmarshalJSON(data []byte) error {
	type input struct {
		BlockHash *common.Hash     `json:"blockHash"`
		FromBlock *rpc.BlockNumber `json:"fromBlock"`
		ToBlock   *rpc.BlockNumber `json:"toBlock"`
		Addresses interface{}      `json:"address"`
		Topics    []interface{}    `json:"topics"`
	}

	var raw input
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.BlockHash != nil {
		if raw.FromBlock != nil || raw.ToBlock != nil {
			// BlockHash is mutually exclusive with FromBlock/ToBlock criteria
			return errors.New("cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")
		}
		args.BlockHash = raw.BlockHash
	} else {
		if raw.FromBlock != nil {
			args.FromBlock = big.NewInt(raw.FromBlock.Int64())
		}

		if raw.ToBlock != nil {
			args.ToBlock = big.NewInt(raw.ToBlock.Int64())
		}
	}

	args.Addresses = []common.Address{}

	if raw.Addresses != nil {
		// raw.Address can contain a single address or an array of addresses
		switch rawAddr := raw.Addresses.(type) {
		case []interface{}:
			for i, addr := range rawAddr {
				if strAddr, ok := addr.(string); ok {
					addr, err := decodeAddress(strAddr)
					if err != nil {
						return fmt.Errorf("invalid address at index %d: %v", i, err)
					}
					args.Addresses = append(args.Addresses, addr)
				} else {
					return fmt.Errorf("non-string address at index %d", i)
				}
			}
		case string:
			addr, err := decodeAddress(rawAddr)
			if err != nil {
				return fmt.Errorf("invalid address: %v", err)
			}
			args.Addresses = []common.Address{addr}
		default:
			return errors.New("invalid addresses in query")
		}
	}
	//if len(raw.Topics) > maxTopics {
	//	return errExceedMaxTopics
	//}

	// topics is an array consisting of strings and/or arrays of strings.
	// JSON null values are converted to common.Hash{} and ignored by the filter manager.
	if len(raw.Topics) > 0 {
		args.Topics = make([][]common.Hash, len(raw.Topics))
		for i, t := range raw.Topics {
			switch topic := t.(type) {
			case nil:
				// ignore topic when matching logs

			case string:
				// match specific topic
				top, err := decodeTopic(topic)
				if err != nil {
					return err
				}
				args.Topics[i] = []common.Hash{top}

			case []interface{}:
				// or case e.g. [null, "topic0", "topic1"]
				//if len(topic) > maxSubTopics {
				//	return errExceedMaxTopics
				//}
				for _, rawTopic := range topic {
					if rawTopic == nil {
						// null component, match all
						args.Topics[i] = nil
						break
					}
					if topic, ok := rawTopic.(string); ok {
						parsed, err := decodeTopic(topic)
						if err != nil {
							return err
						}
						args.Topics[i] = append(args.Topics[i], parsed)
					} else {
						return errInvalidTopic
					}
				}
			default:
				return errInvalidTopic
			}
		}
	}

	return nil
}

func decodeAddress(s string) (common.Address, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != common.AddressLength {
		err = fmt.Errorf("hex has invalid length %d after decoding; expected %d for address", len(b), common.AddressLength)
	}
	return common.BytesToAddress(b), err
}

func decodeTopic(s string) (common.Hash, error) {
	b, err := hexutil.Decode(s)
	if err == nil && len(b) != common.HashLength {
		err = fmt.Errorf("hex has invalid length %d after decoding; expected %d for topic", len(b), common.HashLength)
	}
	return common.BytesToHash(b), err
}

type LogFilter struct {
	b Backend

	addresses []common.Address
	topics    [][]common.Hash

	block      *common.Hash // Block hash if filtering a single block
	begin, end int64        // Range interval if filtering multiple blocks
}

func newLogFilter(ctx context.Context, b Backend, crit FilterCriteria) (*LogFilter, error) {
	var filter *LogFilter
	if crit.BlockHash != nil {
		filter = &LogFilter{
			b:         b,
			block:     crit.BlockHash,
			addresses: crit.Addresses,
			topics:    crit.Topics,
		}
	} else {
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}

		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}

		if begin == rpc.PendingBlockNumber.Int64() || end == rpc.PendingBlockNumber.Int64() {
			return nil, errPendingLogsUnsupported
		}

		_, hdr, _ := b.HeaderByNumber(ctx, rpc.LatestBlockNumber)
		if begin == rpc.LatestBlockNumber.Int64() {
			begin = int64(hdr.Height)
		}
		if end == rpc.LatestBlockNumber.Int64() {
			end = int64(hdr.Height)
		}

		if begin > 0 && end > 0 && begin > end {
			return nil, errInvalidBlockRange
		}

		filter = &LogFilter{
			b:         b,
			begin:     begin,
			end:       end,
			addresses: crit.Addresses,
			topics:    crit.Topics,
		}
	}

	return filter, nil
}

func (f *LogFilter) Logs(ctx context.Context) ([]*types.Log, error) {
	if f.block != nil {
		_, yuHeader, err := f.b.HeaderByHash(ctx, *f.block)
		if err != nil {
			return nil, err
		}

		return f.FilterLogs(ctx, yuHeader)
	} else {
		var result []*types.Log
		for ; f.begin < f.end; f.begin++ {
			_, yuHeader, err := f.b.HeaderByNumber(ctx, rpc.BlockNumber(f.begin))
			if err != nil {
				logrus.Errorf("[GetLog] Failed to getHeaderByNumber %v", f.begin)
				return nil, err
			}
			logs, err := f.FilterLogs(ctx, yuHeader)
			if err != nil {
				return nil, err
			}
			result = append(result, logs...)
		}
		return result, nil
	}
}

func (f *LogFilter) FilterLogs(ctx context.Context, yuHeader *yutypes.Header) ([]*types.Log, error) {
	logs, err := f.b.GetLogs(ctx, common.Hash(yuHeader.Hash), uint64(yuHeader.Height))
	if err != nil {
		return nil, err
	}

	result := make([]*types.Log, 0)
	var logIdx uint
	for i, txLogs := range logs {
		for _, vLog := range txLogs {
			vLog.BlockHash = common.Hash(yuHeader.Hash)
			vLog.BlockNumber = uint64(yuHeader.Height)
			vLog.TxIndex = uint(i)
			vLog.Index = logIdx
			logIdx++

			if f.checkMatches(ctx, vLog) {
				result = append(result, vLog)
			}
		}
	}

	return result, nil
}

func (f *LogFilter) checkMatches(ctx context.Context, vLog *types.Log) bool {
	if len(f.addresses) > 0 {
		if !slices.Contains(f.addresses, vLog.Address) {
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
