package parallel

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/pending_state"
	"github.com/reddio-com/reddio/metrics"
)

type ParallelEvmExecutor struct {
	k          *ParallelEVM
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
	e.subTxnList = e.k.splitTxnCtxList(txnCtxList)
}

func (e *ParallelEvmExecutor) Execute(block *types.Block) {
	got := e.k.executeAllTxn(e.subTxnList)
	for _, subList := range got {
		for _, c := range subList {
			e.receipts[c.txn.TxnHash] = c.receipt
		}
	}
}

func (e *ParallelEvmExecutor) Receipts(block *types.Block) map[common.Hash]*types.Receipt {
	return e.receipts
}

func (k *ParallelEVM) executeAllTxn(got [][]*txnCtx) [][]*txnCtx {
	start := time.Now()
	defer func() {
		k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	for index, subList := range got {
		k.executeTxnCtxListInParallel(subList)
		got[index] = subList
	}
	return got
}

func (k *ParallelEVM) prepareTxnList(block *types.Block) ([]*txnCtx, map[common.Hash]*types.Receipt) {
	start := time.Now()
	defer func() {
		k.statManager.PrepareDuration = time.Since(start)
	}()
	stxns := block.Txns
	receipts := make(map[common.Hash]*types.Receipt)
	txnCtxList := make([]*txnCtx, len(stxns), len(stxns))
	for index, stxn := range stxns {
		wrCall := stxn.Raw.WrCall
		ctx, err := context.NewWriteContext(stxn, block, index)
		if err != nil {
			receipt := k.handleTxnError(err, ctx, block, stxn)
			receipts[stxn.TxnHash] = receipt
			continue
		}
		req := &evm.TxRequest{}
		if err := ctx.BindJson(req); err != nil {
			receipt := k.handleTxnError(err, ctx, block, stxn)
			receipts[stxn.TxnHash] = receipt
			continue
		}
		writing, _ := k.Land.GetWriting(wrCall.TripodName, wrCall.FuncName)
		stxnCtx := &txnCtx{
			ctx:     ctx,
			txn:     stxn,
			writing: writing,
			req:     req,
		}
		txnCtxList[index] = stxnCtx
	}
	return txnCtxList, receipts
}

func (k *ParallelEVM) splitTxnCtxList(list []*txnCtx) [][]*txnCtx {
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
	k.statManager.TxnBatchCount = len(got)
	return got
}

func (k *ParallelEVM) executeTxnCtxListInParallel(list []*txnCtx) []*txnCtx {
	defer func() {
		k.Solidity.FinaliseStateDB(true)
		if config.GetGlobalConfig().AsyncCommit {
			k.updateTxnObjSub(list)
			k.Solidity.StateDB().PendingCommit(true, k.objectInc)
		}
	}()
	metrics.BatchTxnSplitCounter.WithLabelValues(strconv.FormatInt(int64(len(list)), 10)).Inc()
	return k.executeTxnCtxListInConcurrency(k.Solidity.StateDB(), list)
}

func (k *ParallelEVM) executeTxnCtxListInConcurrency(originStateDB *state.StateDB, list []*txnCtx) []*txnCtx {
	conflict := false
	start := time.Now()
	defer func() {
		end := time.Now()
		metrics.BatchTxnDuration.WithLabelValues(fmt.Sprintf("%v", conflict)).Observe(end.Sub(start).Seconds())
	}()
	copiedStateDBList := k.CopyStateDb(originStateDB, list)
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
				tctx.receipt = k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			} else {
				tctx.receipt = k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, false)
			}
			tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)

			list[index] = tctx
		}(i, c, copiedStateDBList[i])
	}
	wg.Wait()
	curtCtx := pending_state.NewStateContext()
	for _, tctx := range list {
		if curtCtx.IsConflict(tctx.ps.GetCtx()) {
			conflict = true
			k.statManager.ConflictCount++
			break
		}
	}

	if conflict && !config.GetGlobalConfig().IgnoreConflict {
		k.statManager.TxnBatchRedoCount++
		metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelRedo).Inc()
		return k.executeTxnCtxListInOrder(originStateDB, list, true)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	k.mergeStateDB(originStateDB, list)
	k.Solidity.SetStateDB(originStateDB)
	k.gcCopiedStateDB(copiedStateDBList, list)
	return list
}
