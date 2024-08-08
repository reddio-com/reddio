package kernel

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/tripod/dev"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/pending_state"
)

type Kernel struct {
	kernel   *kernel.Kernel
	Solidity *evm.Solidity
}

func NewReddioKernel(k *kernel.Kernel, s *evm.Solidity) *Kernel {
	return &Kernel{
		kernel:   k,
		Solidity: s,
	}
}

func (k *Kernel) Execute(block *types.Block) error {
	stxns := block.Txns
	receipts := make(map[common.Hash]*types.Receipt)
	txnCtxList := make([]*txnCtx, 0)
	for _, stxn := range stxns {
		wrCall := stxn.Raw.WrCall
		ctx, err := context.NewWriteContext(stxn, block)
		if err != nil {
			receipt := k.kernel.HandleError(err, ctx, block, stxn)
			receipts[stxn.TxnHash] = receipt
			continue
		}
		req := &evm.TxRequest{}
		if err := ctx.BindJson(req); err != nil {
			receipt := k.kernel.HandleError(err, ctx, block, stxn)
			receipts[stxn.TxnHash] = receipt
			continue
		}
		writing, _ := k.kernel.Land.GetWriting(wrCall.TripodName, wrCall.FuncName)
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
	return k.kernel.PostExecute(block, receipts)
}

const (
	maxConcurrency = 4
)

func (k *Kernel) SplitTxnCtxList(list []*txnCtx) [][]*txnCtx {
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
		if len(curList) >= maxConcurrency {
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
		if compare.req.Address == curTxn.req.Address {
			return true
		}
		if compare.req.Address == curTxn.req.Origin {
			return true
		}
		if compare.req.Origin == curTxn.req.Address {
			return true
		}
		if compare.req.Origin == curTxn.req.Origin {
			return true
		}
	}
	return false
}

func (k *Kernel) executeTxnCtxList(list []*txnCtx) []*txnCtx {
	return k.executeTxnCtxListInConcurrency(k.Solidity.StateDB(), list)
}

func (k *Kernel) reExecuteTxnCtxListInOrder(originStateDB *state.StateDB, list []*txnCtx) []*txnCtx {
	for index, tctx := range list {
		if tctx.err != nil {
			list[index] = tctx
			continue
		}
		tctx.ctx.ExtraInterface = originStateDB
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.err = err
			tctx.r = k.kernel.HandleError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.r = k.kernel.HandleEvent(tctx.ctx, tctx.ctx.Block, tctx.txn)
			tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingState)
		}
		list[index] = tctx
	}
	return list
}

func (k *Kernel) executeTxnCtxListInConcurrency(originStateDB *state.StateDB, list []*txnCtx) []*txnCtx {
	copyedStateDBList := make([]*state.StateDB, 0)
	for i := 0; i < len(list); i++ {
		copyedStateDBList = append(copyedStateDBList, originStateDB.Copy())
	}
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
				tctx.r = k.kernel.HandleError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
			} else {
				tctx.r = k.kernel.HandleEvent(tctx.ctx, tctx.ctx.Block, tctx.txn)
				tctx.ps = tctx.ctx.ExtraInterface.(*pending_state.PendingState)
			}
			list[index] = tctx
		}(i, c, copyedStateDBList[i])
	}
	wg.Wait()
	curtCtx := pending_state.NewStateContext()
	conflict := false
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
		return k.reExecuteTxnCtxListInOrder(originStateDB, list)
	}
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
