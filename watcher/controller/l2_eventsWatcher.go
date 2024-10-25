package controller

import (
	"context"
	"fmt"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/logic"
	"github.com/reddio-com/reddio/watcher/relayer"
	yutypes "github.com/yu-org/yu/core/types"
)

type L2EventsWatcher struct {
	ctx context.Context
	cfg *evm.GethConfig
	//ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	l2toL1Relayer  relayer.L2ToL1RelayerInterface
	solidity       *evm.Solidity `tripod:"solidity"`
}

func NewL2EventsWatcher(ctx context.Context, cfg *evm.GethConfig, l2toL1Relayer relayer.L2ToL1RelayerInterface,
	solidity *evm.Solidity) (*L2EventsWatcher, error) {

	c := &L2EventsWatcher{
		ctx:            ctx,
		cfg:            cfg,
		l2WatcherLogic: logic.NewL2WatcherLogic(cfg, solidity),
		l2toL1Relayer:  l2toL1Relayer,
		solidity:       solidity,
	}
	return c, nil
}

func (w *L2EventsWatcher) WatchUpwardMessage(ctx context.Context, block *yutypes.Block, Solidity *evm.Solidity) error {

	upwardMessage, err := w.l2WatcherLogic.L2FetcherUpwardMessageFromLogs(ctx, block, w.cfg.L2BlockCollectionDepth)
	if err != nil {
		fmt.Println("Watcher L2FetcherUpwardMessageFromLogs error: ", err)
		return err
	}
	// print for test
	// jsonData, err := json.MarshalIndent(upwardMessage, "", "  ")
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal upwardMessage to JSON: %v", err)
	// }

	// fmt.Println("WatchUpwardMessage: ", string(jsonData))
	err = w.l2toL1Relayer.HandleUpwardMessage(upwardMessage)
	if err != nil {
		fmt.Println("Watcher HandleUpwardMessage error: ", err)
		return err
	}

	return nil

}
