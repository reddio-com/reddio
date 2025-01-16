package parallel

import (
	"fmt"
	"strconv"
	"strings"
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

type ParallelEvmSingleStateDBExecutor struct {
	k          *ParallelEVM
	cpdb       *state.StateDB
	receipts   map[common.Hash]*types.Receipt
	subTxnList [][]*txnCtx
}

func NewParallelEvmSingleStateDBExecutor(evm *ParallelEVM) *ParallelEvmSingleStateDBExecutor {
	return &ParallelEvmSingleStateDBExecutor{
		k:    evm,
		cpdb: evm.cpdb,
	}
}

func (e *ParallelEvmSingleStateDBExecutor) Prepare(block *types.Block) {
	e.k.prepareExecute()
	txnCtxList, receipts := e.k.prepareTxnList(block)
	e.receipts = receipts
	e.k.updateTxnObjInc(txnCtxList)
	e.subTxnList = e.k.splitTxnCtxList(txnCtxList)
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

func (e *ParallelEvmSingleStateDBExecutor) executeAllTxn(got [][]*txnCtx) [][]*txnCtx {
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

func (e *ParallelEvmSingleStateDBExecutor) executeTxnCtxListInParallel(list []*txnCtx) []*txnCtx {
	defer func() {
		e.cpdb.Finalise(true)
		if config.GetGlobalConfig().AsyncCommit {
			e.k.updateTxnObjSub(list)
			e.cpdb.PendingCommit(true, e.k.objectInc)
		}
	}()
	metrics.BatchTxnSplitCounter.WithLabelValues(strconv.FormatInt(int64(len(list)), 10)).Inc()
	return e.executeTxnCtxListInConcurrency(list)
}

func (e *ParallelEvmSingleStateDBExecutor) executeTxnCtxListInConcurrency(list []*txnCtx) []*txnCtx {
	conflict := false
	start := time.Now()
	defer func() {
		end := time.Now()
		metrics.BatchTxnDuration.WithLabelValues(fmt.Sprintf("%v", conflict)).Observe(end.Sub(start).Seconds())
	}()
	wrapperList := e.CopyStateDbWrapper(list)
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
		return e.k.executeTxnCtxListInOrder(e.cpdb, list, true)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	e.k.gcCopiedStateDB(wrapperList, list)
	return list
}

func (e *ParallelEvmSingleStateDBExecutor) CopyStateDbWrapper(list []*txnCtx) []*pending_state.PendingStateWrapper {
	copiedStateDBList := make([]*pending_state.PendingStateWrapper, 0)
	start := time.Now()
	defer func() {
		e.k.statManager.CopyDuration += time.Since(start)
	}()
	dbWrapper := pending_state.NewStateDBWrapper(e.cpdb)
	sctx := pending_state.NewStateContext(true)
	for i := 0; i < len(list); i++ {
		needCopy := make(map[common2.Address]struct{})
		if list[i].req.Address != nil {
			needCopy[*list[i].req.Address] = struct{}{}
		}
		needCopy[list[i].req.Origin] = struct{}{}
		copiedStateDBList = append(copiedStateDBList, pending_state.NewPendingStateWrapper(dbWrapper, sctx, int64(i)))
	}
	return copiedStateDBList
}
