package parallel

import (
	"math/big"
	"time"

	common2 "github.com/ethereum/go-ethereum/common"
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

var (
	//0xeC054c6ee2DbbeBC9EbCA50CdBF94A94B02B2E40
	testBridgeContractAddress = common2.HexToAddress("0xA3ED8915aE346bF85E56B6BB6b723091716f58b4")
	testStorageSlotHash       = common2.HexToHash("0x65d63ba7ddf496eb3f7a10c642f6d20487aee1873bcad4890f640b167eab1069")
	lastMessageNonceSlot      = big.NewInt(0)
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
		//k.cpdb.ClearPendingCommitMark()
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

func (k *ParallelEVM) prepareTxnList(block *types.Block) ([]*txnCtx, map[common.Hash]*types.Receipt) {
	start := time.Now()
	defer func() {
		k.statManager.PrepareDuration = time.Since(start)
	}()
	stxns := block.Txns
	receipts := make(map[common.Hash]*types.Receipt)
	txnCtxList := make([]*txnCtx, 0)
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
	return txnCtxList, receipts
}

func (k *ParallelEVM) executeTxnCtxListInOrder(sdb *state.StateDB, list []*txnCtx, isRedo bool) []*txnCtx {
	for index, tctx := range list {
		if tctx.err != nil {
			tctx.receipt = k.handleTxnError(tctx.err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			continue
		}
		tctx.ctx.ExtraInterface = pending_state.NewPendingStateWrapper(pending_state.NewStateDBWrapper(sdb), pending_state.NewStateContext(false), int64(index))
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.receipt = k.handleTxnError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.receipt = k.handleTxnEvent(tctx.ctx, tctx.ctx.Block, tctx.txn, isRedo)
		}
		tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingStateWrapper)
		list[index] = tctx
	}
	k.gcCopiedStateDB(nil, list)
	return list
}

func (k *ParallelEVM) gcCopiedStateDB(copiedStateDBList []*pending_state.PendingStateWrapper, list []*txnCtx) {
	copiedStateDBList = nil
	for _, ctx := range list {
		ctx.ctx.ExtraInterface = nil
		ctx.ps = nil
	}
}
