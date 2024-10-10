package logic

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/evm"
)

type L1WatcherLogic struct {
	cfg         *evm.GethConfig
	client      *ethclient.Client
	addressList []common.Address
	parser      *L1EventParser
}

func NewL1WatcherLogic(cfg *evm.GethConfig, client *ethclient.Client) *L1WatcherLogic {
	contractAddressList := []common.Address{
		common.HexToAddress(cfg.ParentLayerContractAddress),
	}

	f := &L1WatcherLogic{
		cfg:         cfg,
		client:      client,
		addressList: contractAddressList,
		parser:      NewL1EventParser(cfg, client),
	}
	return f
}
