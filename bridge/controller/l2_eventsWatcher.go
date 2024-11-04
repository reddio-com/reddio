package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/relayer"
	"github.com/reddio-com/reddio/evm"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/tripod"
	yutypes "github.com/yu-org/yu/core/types"
)

type L2EventsWatcher struct {
	//ctx context.Context
	cfg *evm.GethConfig
	//ethClient      *ethclient.Client
	l2WatcherLogic *logic.L2WatcherLogic
	l2toL1Relayer  relayer.L2ToL1RelayerInterface
	*tripod.Tripod
	solidity *evm.Solidity `tripod:"solidity"`
}

func NewL2EventsWatcher(cfg *evm.GethConfig,
) *L2EventsWatcher {
	tri := tripod.NewTripod()
	logrus.Info("NewWatcher tripod")
	c := &L2EventsWatcher{
		//ctx:            ctx,
		cfg:    cfg,
		Tripod: tri,
	}
	return c
}

func (w *L2EventsWatcher) WatchUpwardMessage(ctx context.Context, block *yutypes.Block, Solidity *evm.Solidity) error {

	upwardMessage, err := w.l2WatcherLogic.L2FetcherUpwardMessageFromLogs(ctx, block, w.cfg.L2BlockCollectionDepth)
	if err != nil {
		fmt.Println("Watcher L2FetcherUpwardMessageFromLogs error: ", err)
		return err
	}

	if len(upwardMessage) == 0 {
		fmt.Println("No upward messages found")
		return nil
	}
	// print for test
	jsonData, err := json.MarshalIndent(upwardMessage, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal upwardMessage to JSON: %v", err)
	}

	fmt.Println("WatchUpwardMessage: ", string(jsonData))
	err = w.l2toL1Relayer.HandleUpwardMessage(upwardMessage)
	if err != nil {
		fmt.Println("Watcher HandleUpwardMessage error: ", err)
		return err
	}
	return nil

}

func (w *L2EventsWatcher) InitChain(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		// db, err := pebble.Open("evm_bridge_db", &pebble.Options{})
		// if err != nil {
		// 	logrus.Fatal("open db failed: ", err)
		// }
		l1Client, err := ethclient.Dial(w.cfg.L1ClientAddress)
		if err != nil {
			log.Fatal("failed to connect to L1 geth", "endpoint", w.cfg.L1ClientAddress, "err", err)
		}
		l2toL1Relayer, err := relayer.NewL2ToL1Relayer(context.Background(), w.cfg, l1Client)
		if err != nil {
			logrus.Fatal("init bridge relayer failed: ", err)
		}
		l2WatcherLogic, err := logic.NewL2WatcherLogic(w.cfg, w.solidity)
		if err != nil {
			logrus.Fatal("init l2WatcherLogic failed: ", err)
		}

		// l2Watcher, err := controller.NewL2EventsWatcher(context.Background(), w.cfg, l2toL1Relayer, w.Solidity)
		// if err != nil {
		// 	logrus.Fatal("init l2Watcher failed: ", err)
		// }

		w.l2toL1Relayer = l2toL1Relayer
		w.l2WatcherLogic = l2WatcherLogic
		// w.l2Watcher = l2Watcher
		//w.evmBridgeDB = db

	}
	logrus.Info("Watcher InitChain")
}

func (w *L2EventsWatcher) StartBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcher) EndBlock(block *yutypes.Block) {
}

func (w *L2EventsWatcher) FinalizeBlock(block *yutypes.Block) {
	if w.cfg.EnableBridge {
		// upwardSequence, err := w.GetSequence("upwardSequence")
		// if err != nil {
		// 	fmt.Println("GetSequence error", "err", err)
		// 	return
		// }
		// if upwardSequence == 0 {
		// 	fmt.Println("upwardSequence is 0")
		// }
		// fmt.Println("upwardSequence", upwardSequence)
		//watch upward message
		blockHeightBigInt := big.NewInt(int64(block.Header.Height))
		if big.NewInt(0).Mod(blockHeightBigInt, w.cfg.L2BlockCollectionDepth).Cmp(big.NewInt(0)) == 0 {
			err := w.WatchUpwardMessage(context.Background(), block, w.solidity)
			if err != nil {
				fmt.Println("WatchUpwardMessage error", "err", err)
				return
			}
		}
		// upwardSequence++
		// w.SetSequence("upwardSequence", upwardSequence)
	}
}
