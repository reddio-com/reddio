package parallel

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
)

type SerialEvmExecutor struct {
	db         *state.StateDB
	k          *ParallelEVM
	receipts   map[common.Hash]*types.Receipt
	txnCtxList []*txnCtx
	startNonce *big.Int
	endNonce   *big.Int
}

func NewSerialEvmExecutor(evm *ParallelEVM) *SerialEvmExecutor {
	return &SerialEvmExecutor{
		k:  evm,
		db: evm.db,
	}
}

func (s *SerialEvmExecutor) Prepare(block *types.Block) {
	s.k.prepareExecute()
	s.txnCtxList, s.receipts = s.k.prepareTxnList(block)
	s.k.updateTxnObjInc(s.txnCtxList)
}

func (s *SerialEvmExecutor) Execute(block *types.Block) {
	start := time.Now()
	defer func() {
		s.k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	s.startNonce = s.getCurrentNonce(s.db)
	got := s.executeTxnCtxListInSerial(s.txnCtxList)
	for _, c := range got {
		s.receipts[c.txn.TxnHash] = c.receipt
	}
	s.endNonce = s.getCurrentNonce(s.db)
	logrus.Infof("block: %d, startNonce: %v, endNonce: %v", block.Height, s.startNonce.String(), s.endNonce.String())
}

func (s *SerialEvmExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return s.receipts
}

func (s *SerialEvmExecutor) executeTxnCtxListInSerial(list []*txnCtx) []*txnCtx {
	defer func() {
		if config.GetGlobalConfig().AsyncCommit {
			s.k.updateTxnObjSub(list)
			//s.cpdb.PendingCommit(true, s.k.objectInc)
		}
	}()
	return s.k.executeTxnCtxListInOrder(s.db, list, false)
}

func (s *SerialEvmExecutor) getCurrentNonce(sdb *state.StateDB) *big.Int {
	messageNonceSlot := sdb.GetState(testBridgeContractAddress, testStorageSlotHash)
	currentMessageNonceSlot := new(big.Int).SetBytes(messageNonceSlot.Bytes())
	return currentMessageNonceSlot
}

func (s *SerialEvmExecutor) getNonceTxnCtxHash() []string {
	txnHash := make([]string, 0)
	for _, tctx := range s.txnCtxList {
		if tctx.req.Address != nil && *tctx.req.Address == testBridgeContractAddress {
			txnHash = append(txnHash, tctx.txn.TxnHash.String())
		}
	}
	return txnHash
}
