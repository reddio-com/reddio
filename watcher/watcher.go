package watcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/controller"
	"github.com/reddio-com/reddio/watcher/relayer"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/tripod"
	yutypes "github.com/yu-org/yu/core/types"
)

var (
	// ErrAlreadyKnown is returned if the transactions is already contained

	ErrNotFoundReceipt = errors.New("receipt not found")
)

type Watcher struct {
	cfg *evm.GethConfig
	*tripod.Tripod
	Solidity      *evm.Solidity `tripod:"solidity"` // Provides L2 on-chinteraction capabilities to the object
	l1Watcher     *controller.L1EventsWatcher
	l2Watcher     *controller.L2EventsWatcher
	l2toL1Relayer relayer.L2ToL1RelayerInterface
	evmBridgeDB   *pebble.DB
}
type ReceiptRequest struct {
	Hash common.Hash `json:"hash"`
}

type ReceiptResponse struct {
	Receipt *types.Receipt `json:"receipt"`
	Err     error          `json:"err"`
}
type ReceiptsRequest struct {
	Hashes []common.Hash `json:"hashes"`
}

type ReceiptsResponse struct {
	Receipts []*types.Receipt `json:"receipts"`
	Err      error            `json:"err"`
}

func NewWatcher(cfg *evm.GethConfig) *Watcher {

	tri := tripod.NewTripod()
	logrus.Info("NewWatcher tripod")
	w := &Watcher{
		cfg:    cfg,
		Tripod: tri,
	}
	//w.SetReadings(w.GetReceipt)
	return w
}

func (w *Watcher) InitChain(block *yutypes.Block) {
	if w.cfg.EnableL1Client {
		db, err := pebble.Open("evm_bridge_db", &pebble.Options{})
		if err != nil {
			logrus.Fatal("open db failed: ", err)
		}
		l1Client, err := ethclient.Dial(w.cfg.L1ClientAddress)
		if err != nil {
			log.Fatal("failed to connect to L1 geth", "endpoint", w.cfg.L1ClientAddress, "err", err)
		}
		l2toL1Relayer, err := relayer.NewL2ToL1Relayer(context.Background(), w.cfg, l1Client)
		if err != nil {
			logrus.Fatal("init bridge relayer failed: ", err)
		}

		l2Watcher, err := controller.NewL2EventsWatcher(context.Background(), w.cfg, l2toL1Relayer, w.Solidity)
		if err != nil {
			logrus.Fatal("init l2Watcher failed: ", err)
		}

		w.l2toL1Relayer = l2toL1Relayer
		w.l2Watcher = l2Watcher
		w.evmBridgeDB = db

	}
	logrus.Info("Watcher InitChain")
}

func (w *Watcher) StartBlock(block *yutypes.Block) {
	logrus.Info("Watcher StartBlock")
}

func (w *Watcher) EndBlock(block *yutypes.Block) {
	logrus.Info("Watcher EndBlock")
}

func (w *Watcher) FinalizeBlock(block *yutypes.Block) {
	if w.cfg.EnableL1Client {
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
			err := w.l2Watcher.WatchUpwardMessage(context.Background(), block, w.Solidity)
			if err != nil {
				fmt.Println("WatchUpwardMessage error", "err", err)
				return
			}
		}
		// upwardSequence++
		// w.SetSequence("upwardSequence", upwardSequence)
	}
}

func (w *Watcher) SetSequence(key string, value int64) error {
	valueBytes := []byte(fmt.Sprintf("%d", value))
	err := w.evmBridgeDB.Set([]byte(key), valueBytes, pebble.Sync)
	if err != nil {
		return fmt.Errorf("failed to set %s: %v", key, err)
	}
	return nil
}

func (w *Watcher) GetSequence(key string) (int64, error) {
	valueBytes, closer, err := w.evmBridgeDB.Get([]byte(key))
	if err != nil {
		if err == pebble.ErrNotFound {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get %s: %v", key, err)
	}
	defer closer.Close()

	value, err := strconv.ParseInt(string(valueBytes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s: %v", key, err)
	}
	return value, nil
}
