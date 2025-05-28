package parallel

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm/pending_state"
	"github.com/reddio-com/reddio/metrics"
)

type ParallelEvmSingleStateDBExecutor struct {
	k          *TxnEVMProcessor
	receipts   map[common.Hash]*types.Receipt
	subTxnList [][]*txnCtx
}

// one evm
func NewParallelEvmSingleStateDBExecutor(evm *TxnEVMProcessor) *ParallelEvmSingleStateDBExecutor {
	return &ParallelEvmSingleStateDBExecutor{
		k: evm,
	}
}

func (e *ParallelEvmSingleStateDBExecutor) Prepare(block *types.Block) {
	e.k.prepareExecute()
	txnCtxList, receipts := e.k.prepareTxnList(block)
	e.receipts = receipts
	e.k.updateTxnObjInc(txnCtxList)
	e.subTxnList = e.splitTxnCtxList(txnCtxList)
}

func (e *ParallelEvmSingleStateDBExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return e.receipts
}

func (e *ParallelEvmSingleStateDBExecutor) Execute(block *types.Block) {
	got := e.executeAllTxn(e.subTxnList)
	for _, subList := range got {
		for _, c := range subList {
			e.receipts[c.txn.TxnHash] = c.receipt
		}
	}
}

func (e *ParallelEvmSingleStateDBExecutor) splitTxnCtxList(list []*txnCtx) [][]*txnCtx {
	cur := 0
	curList := make([]*txnCtx, 0)
	got := make([][]*txnCtx, 0)
	for cur < len(list) {
		curTxnCtx := list[cur]
		if checkAddressConflict(curTxnCtx, curList) {
			got = append(got, curList)
			curList = make([]*txnCtx, 0)
			continue
		}
		curList = append(curList, curTxnCtx)
		if len(curList) >= config.GetGlobalConfig().MaxConcurrency {
			got = append(got, curList)
			curList = make([]*txnCtx, 0)
		}
		cur++
	}
	if len(curList) > 0 {
		got = append(got, curList)
	}
	e.k.statManager.TxnBatchCount = len(got)
	return got
}

func (e *ParallelEvmSingleStateDBExecutor) executeAllTxn(got [][]*txnCtx) [][]*txnCtx {
	start := time.Now()
	defer func() {
		e.k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	for index, subList := range got {
		e.executeTxnCtxListInConcurrency(subList)
		got[index] = subList
	}
	return got
}

func (e *ParallelEvmSingleStateDBExecutor) executeTxnCtxListInConcurrency(list []*txnCtx) []*txnCtx {
	conflict := false
	start := time.Now()
	defer func() {
		end := time.Now()
		metrics.BatchTxnDuration.WithLabelValues(fmt.Sprintf("%v", conflict)).Observe(end.Sub(start).Seconds())
	}()
	version := e.k.Solidity.Snapshot()
	wrapperList := e.prepareStateDbWrapper(list)
	wg := sync.WaitGroup{}
	for i, c := range list {
		wg.Add(1)
		go func(index int, tctx *txnCtx, wrapper *pending_state.PendingStateWrapper) {
			defer func() {
				wg.Done()
			}()
			tctx.ctx.ExtraInterface = wrapper
			err := tctx.writing(tctx.ctx)
			if err != nil {
				if strings.Contains(err.Error(), "conflict") {
					conflict = true
				}
				tctx.err = err
				tctx.receipt = e.k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			} else {
				tctx.receipt = e.k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, false)
			}

			tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)

			list[index] = tctx
		}(i, c, wrapperList[i])
	}
	wg.Wait()
	if conflict {
		e.k.statManager.TxnBatchRedoCount++
		metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelRedo).Inc()
		e.k.Solidity.RevertToSnapshot(version)
		return e.executeTxnCtxListInOrder(list)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	e.k.gcCopiedStateDB(wrapperList, list)
	return list
}

func (e *ParallelEvmSingleStateDBExecutor) prepareStateDbWrapper(list []*txnCtx) []*pending_state.PendingStateWrapper {
	copiedStateDBList := make([]*pending_state.PendingStateWrapper, 0)
	start := time.Now()
	defer func() {
		e.k.statManager.CopyDuration += time.Since(start)
	}()
	stateDBWrapper := pending_state.NewStateDBWrapper(e.k.Solidity.StateDB())
	sctx := pending_state.NewStateContext(true)
	for i := 0; i < len(list); i++ {
		copiedStateDBList = append(copiedStateDBList, pending_state.NewPendingStateWrapper(stateDBWrapper, sctx, int64(i)))
	}
	return copiedStateDBList
}

func (e *ParallelEvmSingleStateDBExecutor) executeTxnCtxListInOrder(list []*txnCtx) []*txnCtx {
	for index, tctx := range list {
		if tctx.err != nil {
			tctx.receipt = e.k.handleTxnError(tctx.err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			continue
		}
		tctx.ctx.ExtraInterface = pending_state.NewPendingStateWrapper(pending_state.NewStateDBWrapper(e.k.Solidity.StateDB()), pending_state.NewStateContext(true), int64(index))
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.receipt = e.k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.receipt = e.k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, true)
		}
		tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)
		list[index] = tctx
	}
	e.k.gcCopiedStateDB(nil, list)
	return list
}
