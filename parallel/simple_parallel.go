package parallel

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/reddio-com/reddio/evm"
	yucommon "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/kernel"
	yutypes "github.com/yu-org/yu/core/types"
)

type ParallelExecutor struct {
	*kernel.Kernel
}

func (p *ParallelExecutor) SimpleParallelExecute(block *yutypes.Block) error {
	receipts := make(map[yucommon.Hash]*yutypes.Receipt)

	// key: sender address
	parallelTxns := make(map[common.Address][]*yutypes.SignedTxn)

	for _, stxn := range block.Txns {
		txReq := new(evm.TxRequest)
		err := stxn.BindJson(txReq)
		if err != nil {
			return err
		}
		parallelTxns[txReq.Origin] = append(parallelTxns[txReq.Origin], stxn)
	}

	for _, txns := range parallelTxns {
		go func(txns []*yutypes.SignedTxn) {
			// TODO: execute txn
		}(txns)
	}

	// TODO: solve the conflicts

	return p.PostExecute(block, receipts)
}
