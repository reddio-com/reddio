package logic

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/reddio-com/reddio/evm"
	backendabi "github.com/reddio-com/reddio/watcher/abi"
	"github.com/reddio-com/reddio/watcher/contract"
	"github.com/reddio-com/reddio/watcher/utils"
)

// L2EventParser the l1 event parser
type L2EventParser struct {
	cfg    *evm.GethConfig
	client *ethclient.Client
}

// NewL2EventParser creates l1 event parser
func NewL2EventParser(cfg *evm.GethConfig, client *ethclient.Client) *L2EventParser {
	return &L2EventParser{
		cfg:    cfg,
		client: client,
	}
}

// ParseL2EventLogs parses L2 watchedevents
func (e *L2EventParser) ParseL2EventLogs(ctx context.Context, logs []types.Log) (*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	upwardMessageEvent, err := e.ParseL2UpwardMessageEventEventLogs(ctx, logs)
	if err != nil {
		return nil, err
	}
	return upwardMessageEvent, nil
}

// ParseL2UpwardMessageEventEventLogs parses L2 watched events
func (e *L2EventParser) ParseL2UpwardMessageEventEventLogs(ctx context.Context, logs []types.Log) (*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	event := contract.ChildBridgeCoreFacetUpwardMessage{}
	for _, vlog := range logs {
		switch vlog.Topics[0] {
		case backendabi.L2UpwardMessageEventSig:
			err := utils.UnpackLog(backendabi.IL1ParentBridgeCoreFacetABI, &event, "UpwardMessage", vlog)
			if err != nil {
				log.Error("Failed to unpack WithdrawETH event", "err", err)
				return nil, err
			}
			event.Raw = vlog

		}
	}
	return &event, nil
}
