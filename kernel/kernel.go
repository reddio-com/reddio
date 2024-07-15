package kernel

import (
	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/tripod/dev"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/evm"
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
	//got := k.executeTxnCtxList(txnCtxList)
	//for _, c := range got {
	//	receipts[c.txn.TxnHash] = c.r
	//}
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
	for i, c := range list {
		index := i
		tctx := c
		k.Solidity.SetStateDB(k.Solidity.StateDB().Copy())
		err := tctx.writing(tctx.ctx)
		if err != nil {
			tctx.r = k.kernel.HandleError(err, tctx.ctx, tctx.ctx.Block, tctx.txn)
		} else {
			tctx.r = k.kernel.HandleEvent(tctx.ctx, tctx.ctx.Block, tctx.txn)
		}
		list[index] = tctx
	}
	return list
}

type txnCtx struct {
	ctx     *context.WriteContext
	txn     *types.SignedTxn
	r       *types.Receipt
	writing dev.Writing
	req     *evm.TxRequest
}
