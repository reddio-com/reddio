package parallel

import (
	"time"

	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
)

type SerialEvmExecutor struct {
	k          *ParallelEVM
	receipts   map[common.Hash]*types.Receipt
	txnCtxList []*txnCtx
}

func NewSerialEvmExecutor(evm *ParallelEVM) *SerialEvmExecutor {
	return &SerialEvmExecutor{
		k: evm,
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
	got := s.executeTxnCtxListInSerial(s.txnCtxList)
	for _, c := range got {
		s.receipts[c.txn.TxnHash] = c.receipt
	}
}

func (s *SerialEvmExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return s.receipts
}

func (s *SerialEvmExecutor) executeTxnCtxListInSerial(list []*txnCtx) []*txnCtx {
	defer func() {
		if config.GetGlobalConfig().AsyncCommit {
			s.k.updateTxnObjSub(list)
			s.k.Solidity.StateDB().PendingCommit(true, s.k.objectInc)
		}
	}()
	return s.k.executeTxnCtxListInOrder(s.k.Solidity.StateDB(), list, false)
}
