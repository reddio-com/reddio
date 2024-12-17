package parallel

import (
	"time"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm/pending_state"
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
	txnCtxList, receipts := s.k.prepareTxnList(block)
	s.receipts = receipts
	s.k.updateTxnObjInc(txnCtxList)
}

func (s *SerialEvmExecutor) Execute(block *types.Block) {
	start := time.Now()
	defer func() {
		s.k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	got := s.k.executeTxnCtxListInSerial(s.txnCtxList)
	for _, c := range got {
		s.receipts[c.txn.TxnHash] = c.receipt
	}
}

func (s *SerialEvmExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return s.receipts
}

func (k *ParallelEVM) executeTxnCtxListInSerial(list []*txnCtx) []*txnCtx {
	defer func() {
		if config.GetGlobalConfig().AsyncCommit {
			k.updateTxnObjSub(list)
			k.Solidity.StateDB().PendingCommit(true, k.objectInc)
		}
	}()
	return k.executeTxnCtxListInOrder(k.Solidity.StateDB(), list, false)
}

func (k *ParallelEVM) executeTxnCtxListInOrder(originStateDB *state.StateDB, list []*txnCtx, isRedo bool) []*txnCtx {
	currStateDb := originStateDB
	for index, tctx := range list {
		if tctx.err != nil {
			list[index] = tctx
			continue
		}
		tctx.ctx.ExtraInterface = pending_state.NewPendingStateWrapper(pending_state.NewPendingState(currStateDb), 0)
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.receipt = k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.receipt = k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, isRedo)
		}
		tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)
		currStateDb = tctx.ps.GetStateDB()
		list[index] = tctx
	}
	k.Solidity.SetStateDB(currStateDb)
	k.gcCopiedStateDB(nil, list)
	return list
}
