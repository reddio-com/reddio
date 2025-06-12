package parallel

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm/pending_state"
	"github.com/reddio-com/reddio/metrics"
)

type ParallelEvmExecutor struct {
	k          *ParallelEVM
	cpdb       *state.StateDB
	receipts   map[common.Hash]*types.Receipt
	subTxnList [][]*txnCtx
}

func NewParallelEvmExecutor(evm *ParallelEVM) *ParallelEvmExecutor {
	return &ParallelEvmExecutor{
		k: evm,
	}
}

func (e *ParallelEvmExecutor) Prepare(block *types.Block) {
	e.k.prepareExecute()
	txnCtxList, receipts := e.k.prepareTxnList(block)
	e.receipts = receipts
	e.k.updateTxnObjInc(txnCtxList)
	e.subTxnList = e.splitTxnCtxList(txnCtxList)
	e.cpdb = e.k.Solidity.StateDBCopy()
}

func (e *ParallelEvmExecutor) Execute(block *types.Block) {
	got := e.executeAllTxn(e.subTxnList)
	for _, subList := range got {
		for _, c := range subList {
			e.receipts[c.txn.TxnHash] = c.receipt
		}
	}
}

func (e *ParallelEvmExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return e.receipts
}

func (e *ParallelEvmExecutor) executeAllTxn(got [][]*txnCtx) [][]*txnCtx {
	start := time.Now()
	defer func() {
		e.k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	for index, subList := range got {
		e.executeTxnCtxListInParallel(subList)
		got[index] = subList
	}
	e.k.Solidity.SetStateDB(e.cpdb)
	return got
}

func (e *ParallelEvmExecutor) splitTxnCtxList(list []*txnCtx) [][]*txnCtx {
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

func (e *ParallelEvmExecutor) executeTxnCtxListInParallel(list []*txnCtx) []*txnCtx {
	metrics.BatchTxnSplitCounter.WithLabelValues(strconv.FormatInt(int64(len(list)), 10)).Inc()
	return e.executeTxnCtxListInConcurrency(list)
}

func (e *ParallelEvmExecutor) executeTxnCtxListInConcurrency(list []*txnCtx) []*txnCtx {
	conflict := false
	start := time.Now()
	defer func() {
		end := time.Now()
		metrics.BatchTxnDuration.WithLabelValues(fmt.Sprintf("%v", conflict)).Observe(end.Sub(start).Seconds())
	}()
	copiedStateDBList := e.CopyStateDb(list)
	wg := sync.WaitGroup{}
	for i, c := range list {
		wg.Add(1)
		go func(index int, tctx *txnCtx, cpDb *pending_state.PendingStateWrapper) {
			defer func() {
				wg.Done()
			}()
			tctx.ctx.ExtraInterface = cpDb
			err := tctx.writing(tctx.ctx)
			if err != nil {
				tctx.err = err
				tctx.receipt = e.k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			} else {
				tctx.receipt = e.k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, false)
			}
			tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)
			list[index] = tctx
		}(i, c, copiedStateDBList[i])
	}
	wg.Wait()
	curtCtx := pending_state.NewStateContext(false)
	for _, tctx := range list {
		if curtCtx.IsConflict(tctx.ps.GetCtx()) {
			conflict = true
			e.k.statManager.ConflictCount++
			break
		}
	}
	if conflict && !config.GetGlobalConfig().IgnoreConflict {
		e.k.statManager.TxnBatchRedoCount++
		metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelRedo).Inc()
		return e.k.executeTxnCtxListInOrder(e.cpdb, list, true)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	e.mergeStateDB(list)
	e.k.gcCopiedStateDB(copiedStateDBList, list)
	return list
}

func (e *ParallelEvmExecutor) mergeStateDB(list []*txnCtx) {
	for _, tctx := range list {
		tctx.ps.MergeInto(e.cpdb, tctx.req.Origin)
	}
}

func (e *ParallelEvmExecutor) CopyStateDb(list []*txnCtx) []*pending_state.PendingStateWrapper {
	copiedStateDBList := make([]*pending_state.PendingStateWrapper, 0)
	start := time.Now()
	defer func() {
		e.k.statManager.CopyDuration += time.Since(start)
	}()
	for i := 0; i < len(list); i++ {
		needCopy := make(map[common2.Address]struct{})
		if list[i].req.Address != nil {
			needCopy[*list[i].req.Address] = struct{}{}
		}
		needCopy[list[i].req.Origin] = struct{}{}
		//copiedStateDBList = append(copiedStateDBList, pending_state.NewPendingStateWrapper(pending_state.NewStateDBWrapper(e.cpdb.SimpleCopy(needCopy)), pending_state.NewStateContext(false), int64(i)))
		copiedStateDBList = append(copiedStateDBList, pending_state.NewPendingStateWrapper(pending_state.NewStateDBWrapper(e.cpdb.Copy()), pending_state.NewStateContext(false), int64(i)))

	}
	return copiedStateDBList
}
