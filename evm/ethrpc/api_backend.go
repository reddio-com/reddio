package ethrpc

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

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
	"github.com/sirupsen/logrus"
	yucommon "github.com/yu-org/yu/common"
	yucore "github.com/yu-org/yu/core"
	"github.com/yu-org/yu/core/kernel"
	yutypes "github.com/yu-org/yu/core/types"

	"github.com/reddio-com/reddio/evm"
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

//func (e *EthAPIBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
//	//TODO implement me
//	panic("implement me")
//}

// Move to ethrpc/gasprice.go
//func (e *EthAPIBackend) FeeHistory(ctx context.Context, blockCount uint64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*big.Int, [][]*big.Int, []*big.Int, []float64, []*big.Int, []float64, error) {}

func (e *EthAPIBackend) BlobBaseFee(ctx context.Context) *big.Int {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) ExtRPCEnabled() bool {
	return true
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
		yuBlock, err = e.chain.Chain.GetEndCompactBlock()
	case rpc.LatestBlockNumber:
		yuBlock, err = e.chain.Chain.GetEndCompactBlock()
	case rpc.FinalizedBlockNumber, rpc.SafeBlockNumber:
		yuBlock, err = e.chain.Chain.LastFinalizedCompact()
	default:
		yuBlock, err = e.chain.Chain.GetCompactBlockByHeight(yucommon.BlockNum(number))
	}
	return yuHeader2EthHeader(yuBlock.Header), err
}

func (e *EthAPIBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	yuBlock, err := e.chain.Chain.GetCompactBlock(yucommon.Hash(hash))
	if err != nil {
		logrus.Error("ethrpc.api_backend.HeaderByHash() failed: ", err)
		return new(types.Header), err
	}

	return yuHeader2EthHeader(yuBlock.Header), err
}

func (e *EthAPIBackend) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return e.HeaderByNumber(ctx, blockNr)
	}

	if blockHash, ok := blockNrOrHash.Hash(); ok {
		return e.HeaderByHash(ctx, blockHash)
	}

	return nil, errors.New("invalid arguments; neither block number nor hash specified")
}

func (e *EthAPIBackend) CurrentHeader() *types.Header {
	yuBlock, err := e.chain.Chain.GetEndCompactBlock()

	if err != nil {
		logrus.Error("EthAPIBackend.CurrentBlock() failed: ", err)
		return new(types.Header)
	}

	return yuHeader2EthHeader(yuBlock.Header)
}

func (e *EthAPIBackend) CurrentBlock() *types.Header {
	yuBlock, err := e.chain.Chain.GetEndCompactBlock()

	if err != nil {
		logrus.Error("EthAPIBackend.CurrentBlock() failed: ", err)
		return new(types.Header)
	}

	return yuHeader2EthHeader(yuBlock.Header)
}

func (e *EthAPIBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	var (
		yuBlock *yutypes.Block
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
	if err != nil {
		return nil, err
	}
	return compactBlock2EthBlock(yuBlock), err
}

func (e *EthAPIBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	yuBlock, err := e.chain.Chain.GetBlock(yucommon.Hash(hash))

	return compactBlock2EthBlock(yuBlock), err
}

func (e *EthAPIBackend) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return e.BlockByNumber(ctx, blockNr)
	}

	if blockHash, ok := blockNrOrHash.Hash(); ok {
		return e.BlockByHash(ctx, blockHash)
	}

	return nil, errors.New("invalid arguments; neither block number nor hash specified")
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

func (e *EthAPIBackend) ChainDb() ethdb.Database {
	tri := e.chain.GetTripodInstance(SolidityTripod)
	solidityTri := tri.(*evm.Solidity)
	ethDB := solidityTri.GetEthDB()
	return ethDB
}

func (e *EthAPIBackend) AccountManager() *accounts.Manager {
	//TODO implement me
	return nil
}

func (e *EthAPIBackend) Pending() (*types.Block, types.Receipts, *state.StateDB) {
	//TODO implement me
	panic("implement me")
}

func (e *EthAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	//TODO implement me
	panic("implement me")
}

// Eth has changed to POS, Td(total difficulty) is for POW
func (e *EthAPIBackend) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	return nil
}

func (e *EthAPIBackend) GetEVM(ctx context.Context, msg *core.Message, state *state.StateDB, header *types.Header, vmConfig *vm.Config, blockCtx *vm.BlockContext) *vm.EVM {
	if vmConfig == nil {
		//vmConfig = e.chain.Chain.GetVMConfig()
		vmConfig = &vm.Config{
			EnablePreimageRecording: false, // TODO: replace with ctx.Bool()
		}
	}
	txContext := core.NewEVMTxContext(msg)
	var context vm.BlockContext
	if blockCtx != nil {
		context = *blockCtx
	} else {
		var b Backend
		context = core.NewEVMBlockContext(header, NewChainContext(ctx, b), nil)
	}
	return vm.NewEVM(context, txContext, state, e.ChainConfig(), *vmConfig)
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
		Hash:     signedTx.Hash(),
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

func yuTxn2EthTxn(yuSignedTxn *yutypes.SignedTxn) *types.Transaction {
	// Un-serialize wrCall.params to retrive datas:
	wrCallParams := yuSignedTxn.Raw.WrCall.Params
	var txReq = &evm.TxRequest{}
	json.Unmarshal([]byte(wrCallParams), txReq)

	// if nonce is assigned to signedTx.Raw.Nonce, then this is ok; otherwise it's nil:
	nonce := yuSignedTxn.Raw.Nonce

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: txReq.GasPrice,
		Gas:      txReq.GasLimit, // gasLimit: should be obtained from Block & Settings
		To:       &txReq.Address,
		Value:    txReq.Value,
		Data:     txReq.Input,
	})

	return tx
}
func (e *EthAPIBackend) GetTransaction(ctx context.Context, txHash common.Hash) (bool, *types.Transaction, common.Hash, uint64, uint64, error) {
	// Used to get txn from either txdb & txpool:
	stxn, err := e.chain.GetTxn(yucommon.Hash(txHash))
	if err != nil || stxn == nil {
		return false, nil, common.Hash{}, 0, 0, err
	}
	ethTxn := yuTxn2EthTxn(stxn)

	// Fixme: should return lookup.BlockHash, lookup.BlockIndex, lookup.Index
	blockHash := txHash
	blockIndex := uint64(0)
	index := uint64(0)

	return true, ethTxn, blockHash, blockIndex, index, nil
}

func (e *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	// Similar to: e.chain.ChainEnv.Pool.GetTxn - ChainEnv can be ignored b/c txpool has index based on hxHash, therefore it's unique
	stxn, _ := e.chain.Pool.GetAllTxns() // will not return error here

	var ethTxns []*types.Transaction

	for _, yuSignedTxn := range stxn {
		ethTxn := yuTxn2EthTxn(yuSignedTxn)
		ethTxns = append(ethTxns, ethTxn)
	}

	return ethTxns, nil
}

// Similar to GetTransaction():
func (e *EthAPIBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	stxn, err := e.chain.Pool.GetTxn(yucommon.Hash(txHash)) // will not return error here
	if err != nil || stxn == nil {
		return nil
	}
	return yuTxn2EthTxn(stxn)
}

func (e *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	// Loop through all transactions to find matching Account Address, and return it's nonce (if have)
	allEthTxns, _ := e.GetPoolTransactions()

	for _, ethTxn := range allEthTxns {
		if *ethTxn.To() == addr {
			return ethTxn.Nonce(), nil
		}
	}

	return 0, nil
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
	return FakeEngine{}
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
		Difficulty:  new(big.Int).SetUint64(yuHeader.Difficulty),
		Number:      new(big.Int).SetUint64(uint64(yuHeader.Height)),
		GasLimit:    yuHeader.LeiLimit,
		GasUsed:     yuHeader.LeiUsed,
		Time:        yuHeader.Timestamp,
		Extra:       yuHeader.Extra,
		Nonce:       types.BlockNonce{},
		BaseFee:     nil,
	}
}

func compactBlock2EthBlock(yuBlock *yutypes.Block) *types.Block {
	//// Init default values for Eth.Block.Transactions.TxData:
	//var data []byte
	//var ethTxs []*types.Transaction
	//
	//nonce := uint64(0)
	//to := common.HexToAddress("")
	//gasLimit := yuBlock.Header.LeiLimit
	//gasPrice := big.NewInt(0)
	//
	//// Create Eth.Block.Transactions from yu.CompactBlock.Hashes:
	//for _, yuSignedTxn := range yuBlock.Txns {
	//	tx := types.NewTx(&types.LegacyTx{
	//		Nonce:    nonce,
	//		GasPrice: gasPrice,
	//		Gas:      gasLimit,
	//		To:       &to,
	//		Value:    big.NewInt(0),
	//		Data:     data,
	//	}, common.Hash(yuSignedTxn.TxnHash))
	//
	//	ethTxs = append(ethTxs, tx)
	//}
	//
	//// Create new Eth.Block using yu.Header & yu.Hashes:
	//return types.NewBlock(yuHeader2EthHeader(yuBlock.Header), ethTxs, nil, nil, nil)
	return types.NewBlock(yuHeader2EthHeader(yuBlock.Header), nil, nil, nil, nil)
}

// region ---- Fake Consensus Engine ----

type FakeEngine struct{}

// Author retrieves the Ethereum address of the account that minted the given block.
func (f FakeEngine) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (f FakeEngine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("Unimplemented fake engine method VerifyHeader")
}

// VerifyHeaders checks whether a batch of headers conforms to the consensus rules.
func (f FakeEngine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header) (chan<- struct{}, <-chan error) {
	panic("Unimplemented fake engine method VerifyHeaders")
}

// VerifyUncles verifies that the given block's uncles conform to the consensus rules.
func (f FakeEngine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	panic("Unimplemented fake engine method VerifyUncles")
}

// Prepare initializes the consensus fields of a block header.
func (f FakeEngine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	panic("Unimplemented fake engine method Prepare")
}

// Finalize runs any post-transaction state modifications.
func (f FakeEngine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, body *types.Body) {
	panic("Unimplemented fake engine method Finalize")
}

// FinalizeAndAssemble runs any post-transaction state modifications and assembles the final block.
func (f FakeEngine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, body *types.Body, receipts []*types.Receipt) (*types.Block, error) {
	panic("Unimplemented fake engine method FinalizeAndAssemble")
}

// Seal generates a new sealing request for the given input block.
func (f FakeEngine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	panic("Unimplemented fake engine method Seal")
}

// SealHash returns the hash of a block prior to it being sealed.
func (f FakeEngine) SealHash(header *types.Header) common.Hash {
	panic("Unimplemented fake engine method SealHash")
}

// CalcDifficulty is the difficulty adjustment algorithm.
func (f FakeEngine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	panic("Unimplemented fake engine method CalcDifficulty")
}

// APIs returns the RPC APIs this consensus engine provides.
func (f FakeEngine) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	panic("Unimplemented fake engine method APIs")
}

// Close terminates any background threads maintained by the consensus engine.
func (f FakeEngine) Close() error {
	panic("Unimplemented fake engine method Close")
}

// endregion  ---- Fake Consensus Engine ----
