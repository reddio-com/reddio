package logic

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/evm"
)

// L1EventParser the l1 event parser
type L1EventParser struct {
	cfg    *evm.GethConfig
	client *ethclient.Client
}

// NewL1EventParser creates l1 event parser
func NewL1EventParser(cfg *evm.GethConfig, client *ethclient.Client) *L1EventParser {
	return &L1EventParser{
		cfg:    cfg,
		client: client,
	}
}
