package ethrpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/reddio-com/reddio/evm"
	"github.com/sirupsen/logrus"
	yucommon "github.com/yu-org/yu/common"
	yucore "github.com/yu-org/yu/core"
	"github.com/yu-org/yu/core/kernel"
	yutypes "github.com/yu-org/yu/core/types"
	"math/big"
	"time"
)

type EthAPIBackend struct {
	allowUnprotectedTxs bool
	ethChainCfg         *params.ChainConfig
	chain               *kernel.Kernel
}

func (e *EthAPIBackend) SyncProgress() ethereum.SyncProgress {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, []*big.Int, []float64, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) BlobBaseFee(ctx context.Context) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) ChainDb() ethdb.Database {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) AccountManager() *accounts.Manager {
	//TODO implement me
	return nil
	//panic("implement me")
}

func (e *EthAPIBackend) ExtRPCEnabled() bool {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) RPCGasCap() uint64 {
	return 50000000
}

func (e *EthAPIBackend) RPCEVMTimeout() time.Duration {
	return 5 * time.Second
}

func (e *EthAPIBackend) RPCTxFeeCap() float64 {
	return 1
}

func (e *EthAPIBackend) UnprotectedAllowed() bool {
	return e.allowUnprotectedTxs
}

func (e *EthAPIBackend) SetHead(number uint64) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	var (
		yuBlock *yutypes.CompactBlock
		err     error
	)
	switch number {
	case rpc.PendingBlockNumber:
		// FIXME
		yuBlock, err = e.chain.Chain.GetEndBlock()
	case rpc.LatestBlockNumber:
		yuBlock, err = e.chain.Chain.GetEndBlock()
	case rpc.FinalizedBlockNumber, rpc.SafeBlockNumber:
		yuBlock, err = e.chain.Chain.LastFinalized()
	default:
		yuBlock, err = e.chain.Chain.GetBlockByHeight(yucommon.BlockNum(number))
	}
	return yuHeader2EthHeader(yuBlock.Header), err
}

func (e *EthAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) CurrentHeader() *types.Header {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) CurrentBlock() *types.Header {
	yuBlock, err := e.chain.Chain.GetEndBlock()
	if err != nil {
		logrus.Error("EthAPIBackend.CurrentBlock() failed: ", err)
		return new(types.Header)
	}
	return &types.Header{
		ParentHash:  common.Hash(yuBlock.PrevHash),
		Coinbase:    common.Address{}, // FIXME
		Root:        common.Hash(yuBlock.StateRoot),
		TxHash:      common.Hash(yuBlock.TxnRoot),
		ReceiptHash: common.Hash(yuBlock.ReceiptRoot),
		Number:      new(big.Int).SetUint64(uint64(yuBlock.Height)),
		GasLimit:    yuBlock.LeiLimit,
		GasUsed:     yuBlock.LeiUsed,
		Time:        yuBlock.Timestamp,
		Extra:       yuBlock.Extra,
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}
}

func (e *EthAPIBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	header, err := e.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	tri := e.chain.GetTripodInstance(SolidityTripod)
	solidityTri := tri.(*evm.Solidity)
	stateDB, err := solidityTri.StateAt(header.Root)
	if err != nil {
		return nil, nil, err
	}
	return stateDB, header, nil
}

func (e *EthAPIBackend) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*state.StateDB, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return e.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		yuBlock, err := e.chain.Chain.GetBlock(yucommon.Hash(hash))
		if err != nil {
			return nil, nil, err
		}
		tri := e.chain.GetTripodInstance(SolidityTripod)
		solidityTri := tri.(*evm.Solidity)
		stateDB, err := solidityTri.StateAt(common.Hash(yuBlock.StateRoot))
		if err != nil {
			return nil, nil, err
		}
		return stateDB, yuHeader2EthHeader(yuBlock.Header), nil
	}
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (e *EthAPIBackend) Pending() (*types.Block, types.Receipts, *state.StateDB) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) *vm.EVM {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) Call(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	err := args.setDefaults(ctx, e, true)
	if err != nil {
		return nil, err
	}

	// byt, _ := json.Marshal(args)
	callRequest := evm.CallRequest{
		Origin:   *args.From,
		Address:  *args.To,
		Input:    *args.Data,
		Value:    args.Value.ToInt(),
		GasLimit: uint64(*args.Gas),
		GasPrice: args.GasPrice.ToInt(),
	}

	requestByt, _ := json.Marshal(callRequest)
	rdCall := new(yucommon.RdCall)
	rdCall.TripodName = SolidityTripod
	rdCall.FuncName = "Call"
	rdCall.Params = string(requestByt)

	response, err := e.chain.HandleRead(rdCall)
	if err != nil {
		return nil, err
	}

	resp := response.DataInterface.(*evm.CallResponse)
	return resp.Ret, nil
}

func (e *EthAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	signer := types.NewEIP155Signer(e.ethChainCfg.ChainID)
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return err
	}
	txReq := &evm.TxRequest{
		Input:    signedTx.Data(),
		Origin:   sender,
		GasLimit: signedTx.Gas(),
		GasPrice: signedTx.GasPrice(),
		Value:    signedTx.Value(),
	}
	if signedTx.To() != nil {
		txReq.Address = *signedTx.To()
	}
	byt, err := json.Marshal(txReq)
	logrus.Printf("SendTx, Request=%+v\n", string(byt))
	if err != nil {
		return err
	}
	signedWrCall := &yucore.SignedWrCall{
		Call: &yucommon.WrCall{
			TripodName: SolidityTripod,
			FuncName:   "ExecuteTxn",
			Params:     string(byt),
		},
	}
	return e.chain.HandleTxn(signedWrCall)
}

func (e *EthAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (bool, *types.Transaction, common.Hash, uint64, uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) Stats() (pending int, queued int) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) TxPoolContent() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) TxPoolContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeNewTxsEvent(events chan<- core.NewTxsEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) ChainConfig() *params.ChainConfig {
	return e.ethChainCfg
}

func (e *EthAPIBackend) Engine() consensus.Engine {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetBody(ctx context.Context, hash common.Hash, number rpc.BlockNumber) (*types.Body, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetLogs(ctx context.Context, blockHash common.Hash, number uint64) ([][]*types.Log, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) BloomStatus() (uint64, uint64) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	//TODO implement me
	panic("implement me")
}

func yuHeader2EthHeader(yuHeader *yutypes.Header) *types.Header {
	return &types.Header{
		ParentHash:  common.Hash(yuHeader.PrevHash),
		Coinbase:    common.Address{}, // FIXME
		Root:        common.Hash(yuHeader.StateRoot),
		TxHash:      common.Hash(yuHeader.TxnRoot),
		ReceiptHash: common.Hash(yuHeader.ReceiptRoot),
		Number:      new(big.Int).SetUint64(uint64(yuHeader.Height)),
		GasLimit:    yuHeader.LeiLimit,
		GasUsed:     yuHeader.LeiUsed,
		Time:        yuHeader.Timestamp,
		Extra:       yuHeader.Extra,
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}
}
