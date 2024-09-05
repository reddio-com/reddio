package parallel

import (
	"fmt"
	"sync"
	"time"

	"github.com/yu-org/yu/core/tripod"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod/dev"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/pending_state"
	"github.com/reddio-com/reddio/metrics"
)

const (
	txnLabelRedoExecute    = "redo"
	txnLabelExecuteSuccess = "success"
	txnLabelErrExecute     = "err"

	batchTxnLabelSuccess = "success"
	batchTxnLabelRedo    = "redo"
)

type ParallelEVM struct {
	*tripod.Tripod
	Solidity *evm.Solidity `tripod:"solidity"`
}

func NewParallelEVM() *ParallelEVM {
	return &ParallelEVM{
		Tripod: tripod.NewTripod(),
	}
}

func (k *ParallelEVM) Execute(block *types.Block) error {
	stxns := block.Txns
	receipts := make(map[common.Hash]*types.Receipt)
	txnCtxList := make([]*txnCtx, 0)

	start := time.Now()
	defer func() {
		metrics.TxsExecutePerBlockDuration.WithLabelValues().Observe(time.Since(start).Seconds())
	}()
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
		txnCtxList = append(txnCtxList, stxnCtx)
	}
	got := k.SplitTxnCtxList(txnCtxList)
	for index, subList := range got {
		k.executeTxnCtxList(subList)
		got[index] = subList
	}
	for _, subList := range got {
		for _, c := range subList {
			receipts[c.txn.TxnHash] = c.r
		}
	}
	return k.PostExecute(block, receipts)
}

func (k *ParallelEVM) SplitTxnCtxList(list []*txnCtx) [][]*txnCtx {
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
	return got
}

func checkAddressConflict(curTxn *txnCtx, curList []*txnCtx) bool {
	for _, compare := range curList {
		if *compare.req.Address == *curTxn.req.Address {
			return true
		}
		if *compare.req.Address == curTxn.req.Origin {
			return true
		}
		if compare.req.Origin == *curTxn.req.Address {
			return true
		}
		if compare.req.Origin == curTxn.req.Origin {
			return true
		}
	}
	return false
}

func (k *ParallelEVM) executeTxnCtxList(list []*txnCtx) []*txnCtx {
	if config.GetGlobalConfig().IsParallel {
		return k.executeTxnCtxListInConcurrency(k.Solidity.StateDB(), list)
	}
	return k.executeTxnCtxListInOrder(k.Solidity.StateDB(), list, false)
}

func (k *ParallelEVM) executeTxnCtxListInOrder(originStateDB *state.StateDB, list []*txnCtx, isRedo bool) []*txnCtx {
	currStateDb := originStateDB
	for index, tctx := range list {
		if tctx.err != nil {
			list[index] = tctx
			continue
		}
		tctx.ctx.ExtraInterface = currStateDb
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.r = k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.r = k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, isRedo)
			tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingState)
			currStateDb = tctx.ps.GetStateDB()
		}
		list[index] = tctx
	}
	k.Solidity.SetStateDB(currStateDb)
	return list
}

func (k *ParallelEVM) executeTxnCtxListInConcurrency(originStateDB *state.StateDB, list []*txnCtx) []*txnCtx {
	copiedStateDBList := make([]*state.StateDB, 0)
	conflict := false
	start := time.Now()
	defer func() {
		end := time.Now()
		metrics.BatchTxnDuration.WithLabelValues(fmt.Sprintf("%v", conflict)).Observe(end.Sub(start).Seconds())
	}()
	for i := 0; i < len(list); i++ {
		copiedStateDBList = append(copiedStateDBList, originStateDB.Copy())
	}
	metrics.StatedbCopyDuration.WithLabelValues().Observe(time.Since(start).Seconds())
	wg := sync.WaitGroup{}
	for i, c := range list {
		wg.Add(1)
		go func(index int, tctx *txnCtx, cpDb *state.StateDB) {
			defer func() {
				wg.Done()
			}()
			tctx.ctx.ExtraInterface = cpDb
			err := tctx.writing(tctx.ctx)
			if err != nil {
				tctx.err = err
				tctx.r = k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			} else {
				tctx.r = k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, false)
				tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingState)
			}
			list[index] = tctx
		}(i, c, copiedStateDBList[i])
	}
	wg.Wait()
	curtCtx := pending_state.NewStateContext()
	for _, tctx := range list {
		if tctx.err != nil {
			continue
		}
		if curtCtx.IsConflict(tctx.ps.GetCtx()) {
			conflict = true
			break
		}
	}
	if conflict {
		metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelRedo).Inc()
		return k.executeTxnCtxListInOrder(originStateDB, list, true)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	for _, tctx := range list {
		if tctx.err != nil {
			continue
		}
		tctx.ps.MergeInto(originStateDB)
	}
	k.Solidity.SetStateDB(originStateDB)
	return list
}

type txnCtx struct {
	ctx     *context.WriteContext
	txn     *types.SignedTxn
	r       *types.Receipt
	writing dev.Writing
	req     *evm.TxRequest
	err     error
	ps      *pending_state.PendingState
}

func (k *ParallelEVM) handleTxnError(err error, ctx *context.WriteContext, block *types.Block, stxn *types.SignedTxn) *types.Receipt {
	metrics.TxnCounter.WithLabelValues(txnLabelErrExecute).Inc()
	return k.HandleError(err, ctx, block, stxn)
}

func (k *ParallelEVM) handleTxnEvent(ctx *context.WriteContext, block *types.Block, stxn *types.SignedTxn, isRedo bool) *types.Receipt {
	metrics.TxnCounter.WithLabelValues(txnLabelExecuteSuccess).Inc()
	if isRedo {
		metrics.TxnCounter.WithLabelValues(txnLabelRedoExecute).Inc()
	}
	return k.HandleEvent(ctx, block, stxn)
}
