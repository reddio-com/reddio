package logic

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// L2EventParser the l1 event parser
type L2EventParser struct {
	cfg *evm.GethConfig
}

// NewL2EventParser creates l1 event parser
func NewL2EventParser(cfg *evm.GethConfig) *L2EventParser {
	return &L2EventParser{
		cfg: cfg,
	}
}

// ParseL2EventLogs parses L2 watchedevents
func (e *L2EventParser) ParseL2EventLogs(ctx context.Context, logs []types.Log) ([]*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	upwardMessageEvent, err := e.ParseL2UpwardMessageEventEventLogs(ctx, logs)
	if err != nil {
		return nil, err
	}
	return upwardMessageEvent, nil
}

// ParseL2UpwardMessageEventEventLogs parses L2 watched events
func (e *L2EventParser) ParseL2UpwardMessageEventEventLogs(ctx context.Context, logs []types.Log) ([]*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	events := []*contract.ChildBridgeCoreFacetUpwardMessage{}
	for _, vlog := range logs {
		switch vlog.Topics[0] {
		case backendabi.L2UpwardMessageEventSig:
			event := new(contract.ChildBridgeCoreFacetUpwardMessage)
			err := utils.UnpackLog(backendabi.IL2ChildBridgeCoreFacetABI, event, "UpwardMessage", vlog)
			if err != nil {
				fmt.Println("Failed to unpack UpwardMessage event", "err", err)
				log.Error("Failed to unpack UpwardMessage event", "err", err)
				return nil, err
			}
			event.Raw = vlog
			events = append(events, event)

		}
	}
	return events, nil
}
