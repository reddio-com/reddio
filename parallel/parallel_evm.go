package parallel

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	common2 "github.com/ethereum/go-ethereum/common"
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
	Solidity    *evm.Solidity `tripod:"solidity"`
	statManager *BlockTxnStatManager
	objectInc   map[common2.Address]int
}

func NewParallelEVM() *ParallelEVM {
	return &ParallelEVM{
		Tripod: tripod.NewTripod(),
	}
}

func (k *ParallelEVM) prepareExecute() {
	if config.GetGlobalConfig().AsyncCommit {
		k.Solidity.StateDB().ClearPendingCommitMark()
		k.clearObjInc()
	}
}

func (k *ParallelEVM) clearObjInc() {
	k.objectInc = make(map[common2.Address]int)
}

func (k *ParallelEVM) updateTxnObjSub(txns []*txnCtx) {
	if !config.GetGlobalConfig().AsyncCommit {
		return
	}
	sub := func(key common2.Address) {
		v, ok := k.objectInc[key]
		if ok {
			if v == 1 {
				delete(k.objectInc, key)
				return
			}
			k.objectInc[key] = v - 1
		}
	}
	for _, txn := range txns {
		addr1 := txn.req.Address
		if addr1 != nil {
			sub(*addr1)
		}
		sub(txn.req.Origin)
	}
}

func (k *ParallelEVM) updateTxnObjInc(txns []*txnCtx) {
	if !config.GetGlobalConfig().AsyncCommit {
		return
	}
	inc := func(key common2.Address) {
		v, ok := k.objectInc[key]
		if ok {
			k.objectInc[key] = v + 1
			return
		}
		k.objectInc[key] = 1
	}
	for _, txn := range txns {
		addr1 := txn.req.Address
		if addr1 != nil {
			inc(*addr1)
		}
		inc(txn.req.Origin)
	}
}

func (k *ParallelEVM) Execute(block *types.Block) error {
	k.statManager = &BlockTxnStatManager{TxnCount: len(block.Txns)}
	start := time.Now()
	defer func() {
		k.statManager.ExecuteDuration = time.Since(start)
		k.statManager.UpdateMetrics()
	}()
	k.prepareExecute()
	txnCtxList, receipts := k.prepareTxnList(block)
	k.updateTxnObjInc(txnCtxList)
	got := k.splitTxnCtxList(txnCtxList)
	got = k.executeAllTxn(got)
	for _, subList := range got {
		for _, c := range subList {
			receipts[c.txn.TxnHash] = c.receipt
		}
	}
	return k.Commit(block, receipts)
}

func (k *ParallelEVM) Commit(block *types.Block, receipts map[common.Hash]*types.Receipt) error {
	commitStart := time.Now()
	defer func() {
		k.statManager.CommitDuration = time.Since(commitStart)
	}()
	return k.PostExecute(block, receipts)
}

func (k *ParallelEVM) executeAllTxn(got [][]*txnCtx) [][]*txnCtx {
	start := time.Now()
	defer func() {
		k.statManager.ExecuteTxnDuration = time.Since(start)
	}()
	for index, subList := range got {
		k.executeTxnCtxList(subList)
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
	return got
}

func checkAddressConflict(curTxn *txnCtx, curList []*txnCtx) bool {
	for _, compare := range curList {

		if curTxn.req.Address != nil && compare.req.Address != nil {
			if *compare.req.Address == *curTxn.req.Address {
				return true
			}
		}

		if compare.req.Address != nil {
			if *compare.req.Address == curTxn.req.Origin {
				return true
			}
		}

		if curTxn.req.Address != nil {
			if compare.req.Origin == *curTxn.req.Address {
				return true
			}
		}

		if compare.req.Origin == curTxn.req.Origin {
			return true
		}

	}
	return false
}

func (k *ParallelEVM) executeTxnCtxList(list []*txnCtx) []*txnCtx {
	if config.GetGlobalConfig().IsParallel {
		defer func() {
			k.Solidity.FinaliseStateDB(true)
		}()
		metrics.BatchTxnSplitCounter.WithLabelValues(strconv.FormatInt(int64(len(list)), 10)).Inc()
		return k.executeTxnCtxListInConcurrency(k.Solidity.StateDB(), list)
	}
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
		tctx.ctx.ExtraInterface = currStateDb
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
			break
		}
	}

	if conflict {
		metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelRedo).Inc()
		return k.executeTxnCtxListInOrder(originStateDB, list, true)
	}
	metrics.BatchTxnCounter.WithLabelValues(batchTxnLabelSuccess).Inc()
	k.mergeStateDB(originStateDB, list)
	k.Solidity.SetStateDB(originStateDB)
	k.gcCopiedStateDB(copiedStateDBList, list)
	return list
}

func (k *ParallelEVM) gcCopiedStateDB(copiedStateDBList []*pending_state.PendingStateWrapper, list []*txnCtx) {
	copiedStateDBList = nil
	for _, ctx := range list {
		ctx.ctx.ExtraInterface = nil
		ctx.ps = nil
	}
}

func (k *ParallelEVM) mergeStateDB(originStateDB *state.StateDB, list []*txnCtx) {
	k.Solidity.Lock()
	for _, tctx := range list {
		tctx.ps.MergeInto(originStateDB, tctx.req.Origin)
	}
	k.Solidity.Unlock()
}

func (k *ParallelEVM) CopyStateDb(originStateDB *state.StateDB, list []*txnCtx) []*pending_state.PendingStateWrapper {
	copiedStateDBList := make([]*pending_state.PendingStateWrapper, 0)
	k.Solidity.Lock()
	start := time.Now()
	defer func() {
		k.statManager.CopyDuration += time.Since(start)
		k.Solidity.Unlock()
	}()
	for i := 0; i < len(list); i++ {
		needCopy := make(map[common2.Address]struct{})
		if list[i].req.Address != nil {
			needCopy[*list[i].req.Address] = struct{}{}
		}
		needCopy[list[i].req.Origin] = struct{}{}
		copiedStateDBList = append(copiedStateDBList, pending_state.NewPendingStateWrapper(pending_state.NewPendingState(originStateDB.SimpleCopy(needCopy)), int64(i)))
	}
	return copiedStateDBList
}

type txnCtx struct {
	ctx     *context.WriteContext
	txn     *types.SignedTxn
	writing dev.Writing
	req     *evm.TxRequest
	err     error
	ps      *pending_state.PendingStateWrapper
	receipt *types.Receipt
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

type BlockTxnStatManager struct {
	TxnCount           int
	ExecuteDuration    time.Duration
	ExecuteTxnDuration time.Duration
	PrepareDuration    time.Duration
	CommitDuration     time.Duration
	CopyDuration       time.Duration
}

func (stat *BlockTxnStatManager) UpdateMetrics() {
	metrics.BlockExecuteTxnCountGauge.WithLabelValues().Set(float64(stat.TxnCount))
	metrics.BlockExecuteTxnDurationGauge.WithLabelValues().Set(float64(stat.ExecuteDuration.Seconds()))
	metrics.BlockTxnAllExecuteDurationGauge.WithLabelValues().Set(float64(stat.ExecuteTxnDuration.Seconds()))
	metrics.BlockTxnPrepareDurationGauge.WithLabelValues().Set(float64(stat.PrepareDuration.Seconds()))
	metrics.BlockTxnCommitDurationGauge.WithLabelValues().Set(float64(stat.CommitDuration.Seconds()))
	if config.GlobalConfig.IsBenchmarkMode {
		log.Printf("execute %v txn, total:%v, execute cost:%v, prepare:%v, copy:%v, commit:%v",
			stat.TxnCount, stat.ExecuteDuration.String(), stat.ExecuteTxnDuration.String(),
			stat.PrepareDuration.String(), stat.CopyDuration.String(), stat.CommitDuration.String())
	}
}
