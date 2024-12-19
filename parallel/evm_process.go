package parallel

import (
	"github.com/sirupsen/logrus"
	"time"

	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/yu-org/yu/core/tripod"

	"github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
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
	processor   EvmProcessor
}

func NewParallelEVM() *ParallelEVM {
	evm := &ParallelEVM{
		Tripod: tripod.NewTripod(),
	}
	return evm
}

func (k *ParallelEVM) setupProcessor() {
	if config.GetGlobalConfig().IsParallel {
		k.processor = NewParallelEvmExecutor(k)
	} else {
		k.processor = NewSerialEvmExecutor(k)
	}
}

func (k *ParallelEVM) Execute(block *types.Block) error {
	k.statManager = &BlockTxnStatManager{TxnCount: len(block.Txns)}
	k.setupProcessor()
	start := time.Now()
	defer func() {
		k.statManager.ExecuteDuration = time.Since(start)
		k.statManager.UpdateMetrics()
	}()
	k.processor.Prepare(block)
	logrus.Info("ParallelEVM.processor.Execute block --- ")
	k.processor.Execute(block)
	logrus.Info("ParallelEVM.processor.Receipts block --- ")
	receipts := k.processor.Receipts(block)
	logrus.Info("ParallelEVM.processor.Commit block --- ")
	return k.Commit(block, receipts)
}

func (k *ParallelEVM) Commit(block *types.Block, receipts map[common.Hash]*types.Receipt) error {
	commitStart := time.Now()
	defer func() {
		k.statManager.CommitDuration = time.Since(commitStart)
	}()
	return k.PostExecute(block, receipts)
}

type EvmProcessor interface {
	Prepare(block *types.Block)
	Execute(block *types.Block)
	Receipts(block *types.Block) map[common.Hash]*types.Receipt
}