package parallel

import (
	"math/big"
	"time"

	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm/pending_state"
)

type SerialEvmExecutor struct {
	k          *TxnEVMProcessor
	receipts   map[common.Hash]*types.Receipt
	txnCtxList []*txnCtx
	startNonce *big.Int
	endNonce   *big.Int
}

func NewSerialEvmExecutor(evm *TxnEVMProcessor) *SerialEvmExecutor {
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
			//s.cpdb.PendingCommit(true, s.k.objectInc)
		}
	}()
	return s.executeTxnCtxListInOrder(list, false)
}

func (s *SerialEvmExecutor) executeTxnCtxListInOrder(list []*txnCtx, isRedo bool) []*txnCtx {
	for index, tctx := range list {
		if tctx.err != nil {
			tctx.receipt = s.k.handleTxnError(tctx.err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			continue
		}
		tctx.ctx.ExtraInterface = pending_state.NewPendingStateWrapper(pending_state.NewStateDBWrapper(s.k.Solidity.StateDB()), pending_state.NewStateContext(false), int64(index))
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.receipt = s.k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.receipt = s.k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, isRedo)
		}
		tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)
		list[index] = tctx
	}
	s.k.gcCopiedStateDB(nil, list)
	return list
}
