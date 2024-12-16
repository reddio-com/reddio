package parallel

import (
	"time"

	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod/dev"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/pending_state"
	"github.com/reddio-com/reddio/metrics"
)

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
