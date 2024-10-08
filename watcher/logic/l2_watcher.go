package logic

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/reddio-com/reddio/evm"
	backendabi "github.com/reddio-com/reddio/watcher/abi"
	"github.com/reddio-com/reddio/watcher/contract"
)

// for test

type L2WatcherLogic struct {
	cfg         *evm.GethConfig
	client      *ethclient.Client
	addressList []common.Address
	parser      *L2EventParser
}

func NewL2WatcherLogic(cfg *evm.GethConfig, client *ethclient.Client) *L2WatcherLogic {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ChildLayerContractAddress),
	}

	f := &L2WatcherLogic{
		cfg:         cfg,
		client:      client,
		addressList: contractAddressList,
		parser:      NewL2EventParser(cfg, client),
	}

	return f
}

func (f *L2WatcherLogic) L2FetcherUpwardMessageFromLogs(ctx context.Context, from, to uint64) (*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from), // inclusive
		ToBlock:   new(big.Int).SetUint64(to),   // inclusive
		Addresses: f.addressList,
		Topics:    make([][]common.Hash, 1),
	}
	query.Topics[0] = make([]common.Hash, 1)
	query.Topics[0][0] = backendabi.L2UpwardMessageEventSig
	//fmt.Println("query: ", query)

	eventLogs, err := f.client.FilterLogs(ctx, query)
	if err != nil {
		log.Error("Failed to filter L2 event logs 2", "from", from, "to", to, "err", err)
		return nil, err
	}
	//fmt.Println("eventLogs: ", eventLogs)
	if len(eventLogs) == 0 {
		log.Info("No event logs found", "from", from, "to", to)
		return nil, nil
	}
	upwardMessageEvent, err := f.parser.ParseL2EventLogs(ctx, eventLogs)
	if err != nil {
		log.Error("Failed to parse L2 event logs 3", "err", err)
		return nil, err
	}

	return upwardMessageEvent, nil
}
