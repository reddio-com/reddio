package evm

import (
	// "github.com/yu-org/yu/common/yerror"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/yu-org/yu/common/yerror"
	"math/big"
	"net/http"
	"sync"

	yuConfig "github.com/reddio-com/reddio/evm/config"

	"github.com/reddio-com/reddio/evm/pending_state"

	"github.com/sirupsen/logrus"
	yu_common "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/context"
	"github.com/yu-org/yu/core/tripod"
	yu_types "github.com/yu-org/yu/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	"github.com/holiman/uint256"
)

type Solidity struct {
	sync.Mutex

	*tripod.Tripod
	ethState    *EthState
	cfg         *GethConfig
	stateConfig *yuConfig.Config
}

func (s *Solidity) StateDB() *state.StateDB {
	return s.ethState.StateDB()
}

func (s *Solidity) SetStateDB(d *state.StateDB) {
	s.ethState.SetStateDB(d)
}

func newEVM_copy(cfg *GethConfig, req *TxRequest) *vm.EVM {
	txContext := vm.TxContext{
		Origin:     req.Origin,
		GasPrice:   req.GasPrice,
		BlobHashes: cfg.BlobHashes,
		BlobFeeCap: cfg.BlobFeeCap,
	}
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     cfg.GetHashFn,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    req.GasLimit,
		BaseFee:     cfg.BaseFee,
		BlobBaseFee: cfg.BlobBaseFee,
		Random:      cfg.Random,
	}

	return vm.NewEVM(blockContext, txContext, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
}

func newEVM(cfg *GethConfig) *vm.EVM {
	txContext := vm.TxContext{
		Origin:     cfg.Origin,
		GasPrice:   cfg.GasPrice,
		BlobHashes: cfg.BlobHashes,
		BlobFeeCap: cfg.BlobFeeCap,
	}
	blockContext := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		GetHash:     cfg.GetHashFn,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		BaseFee:     cfg.BaseFee,
		BlobBaseFee: cfg.BlobBaseFee,
		Random:      cfg.Random,
	}

	return vm.NewEVM(blockContext, txContext, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
}

func (s *Solidity) InitChain(genesisBlock *yu_types.Block) {
	cfg := s.stateConfig
	genesis := DefaultGoerliGenesisBlock()

	logrus.Printf("Genesis GethConfig: %+v", genesis.Config)
	logrus.Println("Genesis Timestamp: ", genesis.Timestamp)
	logrus.Printf("Genesis ExtraData: %x", genesis.ExtraData)
	logrus.Println("Genesis GasLimit: ", genesis.GasLimit)
	logrus.Println("Genesis Difficulty: ", genesis.Difficulty.String())

	var lastStateRoot common.Hash
	block, err := s.GetCurrentBlock()
	if err != nil && err != yerror.ErrBlockNotFound {
		logrus.Fatal("get current block failed: ", err)
	}
	if block != nil {
		lastStateRoot = common.Hash(block.StateRoot)
	}

	ethState, err := NewEthState(cfg, lastStateRoot)
	if err != nil {
		logrus.Fatal("init NewEthState failed: ", err)
	}
	s.ethState = ethState
	s.cfg.State = ethState.stateDB

	chainConfig, _, err := SetupGenesisBlock(ethState, genesis)
	if err != nil {
		logrus.Fatal("SetupGenesisBlock failed: ", err)
	}

	// s.cfg.ChainConfig = chainConfig

	logrus.Println("Genesis SetupGenesisBlock chainConfig: ", chainConfig)
	logrus.Println("Genesis NewEthState cfg.DbPath: ", ethState.cfg.DbPath)
	logrus.Println("Genesis NewEthState ethState.cfg.NameSpace: ", ethState.cfg.NameSpace)
	logrus.Println("Genesis NewEthState ethState.StateDB.SnapshotCommits: ", ethState.stateDB)
	logrus.Println("Genesis NewEthState ethState.stateCache: ", ethState.stateCache)
	logrus.Println("Genesis NewEthState ethState.trieDB: ", ethState.trieDB)

	// commit genesis state
	genesisStateRoot, err := s.ethState.GenesisCommit()
	if err != nil {
		logrus.Fatal("genesis state commit failed: ", err)
	}

	genesisBlock.StateRoot = yu_common.Hash(genesisStateRoot)
}

func NewSolidity(gethConfig *GethConfig) *Solidity {
	ethStateConfig := setDefaultEthStateConfig()

	solidity := &Solidity{
		Tripod:      tripod.NewTripod(),
		cfg:         gethConfig,
		stateConfig: ethStateConfig,
		// network:       utils.Network(cfg.Network),
	}
	solidity.SetWritings(solidity.ExecuteTxn)
	solidity.SetReadings(
		solidity.Call, solidity.GetReceipt, solidity.GetReceipts,
		// solidity.GetClass, solidity.GetClassAt,
		// 	solidity.GetClassHashAt, solidity.GetNonce, solidity.GetStorage,
		// 	solidity.GetTransaction, solidity.GetTransactionStatus,
		// 	solidity.SimulateTransactions,
		// 	solidity.GetBlockWithTxs, solidity.GetBlockWithTxHashes,
	)

	return solidity
}

// region ---- Tripod Api ----

func (s *Solidity) StartBlock(block *yu_types.Block) {
	s.cfg.BlockNumber = big.NewInt(int64(block.Height))
	//s.cfg.GasLimit =
	s.cfg.Time = block.Timestamp
	s.cfg.Difficulty = big.NewInt(int64(block.Difficulty))
}

func (s *Solidity) EndBlock(block *yu_types.Block) {
	// nothing
}

func (s *Solidity) FinalizeBlock(block *yu_types.Block) {
	// nothing
}

func (s *Solidity) PreHandleTxn(txn *yu_types.SignedTxn) error {
	var txReq TxRequest
	param := txn.GetParams()
	err := json.Unmarshal([]byte(param), &txReq)
	if err != nil {
		return err
	}

	yuHash, err := ConvertHashToYuHash(txReq.Hash)
	if err != nil {
		return err
	}
	txn.TxnHash = yuHash

	return nil
}

// ExecuteTxn executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Execute sets up an in-memory, temporary, environment for the execution of
// the given code. It makes sure that it's restored to its original state afterwards.
func (s *Solidity) ExecuteTxn(ctx *context.WriteContext) (err error) {
	txReq := new(TxRequest)
	coinbase := common.BytesToAddress(s.cfg.Coinbase.Bytes())
	//s.Lock()
	err = ctx.BindJson(txReq)
	if err != nil {
		return err
	}
	cfg := s.cfg

	vmenv := newEVM_copy(cfg, txReq)
	pd := pending_state.NewPendingState(ctx.ExtraInterface.(*state.StateDB))
	vmenv.StateDB = pd
	s.ethState.setTxContext(common.Hash(ctx.GetTxnHash()), ctx.TxnIndex)

	vmenv.Context.BlockNumber = big.NewInt(int64(ctx.Block.Height))

	sender := vm.AccountRef(txReq.Origin)
	rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)
	//s.Unlock()

	if txReq.Address == nil {
		err = executeContractCreation(ctx, txReq, pd, txReq.Origin, coinbase, vmenv, sender, rules)
	} else {
		err = executeContractCall(ctx, txReq, pd, txReq.Origin, coinbase, vmenv, sender, rules)
	}

	if err != nil {
		return err
	}
	ctx.ExtraInterface = pd

	return nil
}

//func emitReceipt(ctx *context.WriteContext, vmEvm *vm.EVM, txReq *TxRequest, contractAddr common.Address, leftOverGas uint64, err error) error {
//	evmReceipt := makeEvmReceipt(vmEvm, ctx.Txn, ctx.Block, contractAddr, leftOverGas, err)
//	receiptByt, err := json.Marshal(evmReceipt)
//	if err != nil {
//		return err
//	}
//	ctx.ExtraInterface = pd
//	return nil
//}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
func (s *Solidity) Call(ctx *context.ReadContext) {
	callReq := new(CallRequest)
	err := ctx.BindJson(callReq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &CallResponse{Err: err})
		return
	}

	cfg := s.cfg
	address := callReq.Address
	input := callReq.Input
	origin := callReq.Origin
	gasLimit := callReq.GasLimit
	gasPrice := callReq.GasPrice
	value := callReq.Value

	cfg.Origin = origin
	cfg.GasLimit = gasLimit
	cfg.GasPrice = gasPrice
	cfg.Value = value

	var (
		vmenv    = newEVM(cfg)
		sender   = vm.AccountRef(origin)
		ethState = s.ethState
		rules    = cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)
	)

	logrus.Printf("[StateDB] %v", s.ethState.stateDB == s.cfg.State)
	vmenv.StateDB = s.ethState.stateDB

	if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
		cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: &address, Data: input, Value: value, Gas: gasLimit}), origin)
	}
	// Execute the preparatory steps for state transition which includes:
	// - prepare accessList(post-berlin)
	// - reset transient storage(eip 1153)
	ethState.Prepare(rules, origin, cfg.Coinbase, &address, vm.ActivePrecompiles(rules), nil)

	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		gasLimit,
		uint256.MustFromBig(value),
	)

	logrus.Printf("[Call] Request from = %v, to = %v, gasLimit = %v, value = %v, input = %v", sender.Address().Hex(), address.Hex(), gasLimit, value.Uint64(), hex.EncodeToString(input))
	logrus.Printf("[Call] Response: Origin Code = %v, Hex Code = %v, String Code = %v, LeftOverGas = %v", ret, hex.EncodeToString(ret), new(big.Int).SetBytes(ret).String(), leftOverGas)

	if err != nil {
		ctx.Json(http.StatusInternalServerError, &CallResponse{Err: err})
		return
	}
	result := CallResponse{Ret: ret, LeftOverGas: leftOverGas}
	json, _ := json.Marshal(result)
	fmt.Printf("[ETH_CALL] eth return result is %v\n", string(json))
	ctx.JsonOk(&result)
}

func (s *Solidity) Commit(block *yu_types.Block) {
	blockNumber := uint64(block.Height)
	stateRoot, err := s.ethState.Commit(blockNumber)
	if err != nil {
		logrus.Errorf("Solidity commit failed on Block(%d), error: %v", blockNumber, err)
		return
	}
	block.StateRoot = AdaptHash(stateRoot)
}

func AdaptHash(ethHash common.Hash) yu_common.Hash {
	var yuHash yu_common.Hash
	copy(yuHash[:], ethHash[:])
	return yuHash
}

func executeContractCreation(ctx *context.WriteContext, txReq *TxRequest, stateDB *pending_state.PendingState, origin, coinBase common.Address, vmenv *vm.EVM, sender vm.AccountRef, rules params.Rules) error {
	//if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
	//	cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
	//}

	stateDB.Prepare(rules, origin, coinBase, nil, vm.ActivePrecompiles(rules), nil)

	code, address, leftOverGas, err := vmenv.Create(sender, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
	if err != nil {
		// byt, _ := json.Marshal(txReq)
		//logrus.Printf("[Execute Txn] Create contract Failed. err = %v. Request = %v", err, string(byt))
		_ = emitReceipt(ctx, vmenv, txReq, code, address, leftOverGas, err)
		return err
	}

	//logrus.Printf("[Execute Txn] Create contract success. Oringin code = %v, Hex Code = %v, Address = %v, Left Gas = %v", code, hex.EncodeToString(code), address.Hex(), leftOverGas)
	return emitReceipt(ctx, vmenv, txReq, code, address, leftOverGas, err)
}

func makeEvmReceipt(ctx *context.WriteContext, vmEvm *vm.EVM, code []byte, signedTx *yu_types.SignedTxn, block *yu_types.Block, address common.Address, leftOverGas uint64, err error) *types.Receipt {
	wrCallParams := signedTx.Raw.WrCall.Params
	var txReq = &TxRequest{}
	_ = json.Unmarshal([]byte(wrCallParams), txReq)

	txArgs := &TempTransactionArgs{}
	_ = json.Unmarshal(txReq.OriginArgs, txArgs)
	originTx := txArgs.ToTransaction(txReq.V, txReq.R, txReq.S)

	stateDb := vmEvm.StateDB.(*pending_state.PendingState).GetStateDB()
	usedGas := originTx.Gas() - leftOverGas

	blockNumber := big.NewInt(int64(block.Height))
	txHash := common.Hash(signedTx.TxnHash)
	effectiveGasPrice := big.NewInt(1000000000) // 1 GWei

	status := types.ReceiptStatusFailed
	if err == nil {
		status = types.ReceiptStatusSuccessful
	}
	var root []byte
	//stateDB := vmEvm.StateDB.(*pending_state.PendingState)
	//if vmEvm.ChainConfig().IsByzantium(blockNumber) {
	//	stateDB.Finalise(true)
	//} else {
	//	root = stateDB.IntermediateRoot(vmEvm.ChainConfig().IsEIP158(blockNumber)).Bytes()
	//}

	// TODO: 1. root is nil; 2. CumulativeGasUsed not; 3. logBloom is empty

	receipt := &types.Receipt{
		Type:              originTx.Type(),
		Status:            status,
		PostState:         root,
		CumulativeGasUsed: leftOverGas,
		TxHash:            txHash,
		ContractAddress:   address,
		GasUsed:           usedGas,
		EffectiveGasPrice: effectiveGasPrice,
	}

	if originTx.Type() == types.BlobTxType {
		receipt.BlobGasUsed = uint64(len(originTx.BlobHashes()) * params.BlobTxBlobGasPerBlob)
		receipt.BlobGasPrice = vmEvm.Context.BlobBaseFee
	}

	receipt.Logs = stateDb.GetLogs(txHash, blockNumber.Uint64(), common.Hash(block.Hash))
	receipt.Bloom = types.CreateBloom(types.Receipts{})
	receipt.BlockHash = common.Hash(block.Hash)
	receipt.BlockNumber = blockNumber
	receipt.TransactionIndex = uint(ctx.TxnIndex)

	logrus.Printf("[Receipt] log = %v", receipt.Logs)
	//spew.Dump("[Receipt] log = %v", stateDB.Logs())
	//logrus.Printf("[Receipt] log is nil = %v", receipt.Logs == nil)
	if receipt.Logs == nil {
		receipt.Logs = []*types.Log{}
	}

	for idx, txn := range block.Txns {
		if common.Hash(txn.TxnHash) == txHash {
			receipt.TransactionIndex = uint(idx)
		}
	}
	logrus.Printf("[Receipt] statedb txIndex = %v, actual txIndex = %v", ctx.TxnIndex, receipt.TransactionIndex)

	return receipt
}

func executeContractCall(ctx *context.WriteContext, txReq *TxRequest, ethState *pending_state.PendingState, origin, coinBase common.Address, vmenv *vm.EVM, sender vm.AccountRef, rules params.Rules) error {
	ethState.Prepare(rules, origin, coinBase, txReq.Address, vm.ActivePrecompiles(rules), nil)
	ethState.SetNonce(txReq.Origin, ethState.GetNonce(sender.Address())+1)

	logrus.Printf("before transfer: account %s balance %d \n", sender.Address(), ethState.GetBalance(sender.Address()))

	code, leftOverGas, err := vmenv.Call(sender, *txReq.Address, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
	//logrus.Printf("after transfer: account %s balance %d \n", sender.Address(), ethState.GetBalance(sender.Address()))
	if err != nil {
		//byt, _ := json.Marshal(txReq)
		//logrus.Printf("[Execute Txn] SendTx Failed. err = %v. Request = %v", err, string(byt))
		_ = emitReceipt(ctx, vmenv, txReq, code, common.Address{}, leftOverGas, err)
		return err
	}

	//logrus.Printf("[Execute Txn] SendTx success. Oringin code = %v, Hex Code = %v, Left Gas = %v", code, hex.EncodeToString(code), leftOverGas)
	return emitReceipt(ctx, vmenv, txReq, code, common.Address{}, leftOverGas, err)
}

func (s *Solidity) StateAt(root common.Hash) (*state.StateDB, error) {
	return s.ethState.StateAt(root)
}

func (s *Solidity) GetEthDB() ethdb.Database {
	return s.ethState.ethDB
}

type ReceiptRequest struct {
	Hash common.Hash `json:"hash"`
}

type ReceiptResponse struct {
	Receipt *types.Receipt `json:"receipt"`
	Err     error          `json:"err"`
}

type ReceiptsRequest struct {
	Hashes []common.Hash `json:"hashes"`
}

type ReceiptsResponse struct {
	Receipts []*types.Receipt `json:"receipts"`
	Err      error            `json:"err"`
}

func (s *Solidity) GetReceipt(ctx *context.ReadContext) {
	var rq ReceiptRequest
	err := ctx.BindJson(&rq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &ReceiptResponse{Err: err})
		return
	}

	receipt, err := s.getReceipt(rq.Hash)
	if err != nil {
		ctx.Json(http.StatusInternalServerError, &ReceiptResponse{Err: err})
		return
	}

	ctx.JsonOk(&ReceiptResponse{Receipt: receipt})
}

func (s *Solidity) getReceipt(hash common.Hash) (*types.Receipt, error) {
	yuHash, err := ConvertHashToYuHash(hash)
	if err != nil {
		return nil, err
	}
	yuReceipt, err := s.TxDB.GetReceipt(yuHash)
	if err != nil {
		return nil, err
	}
	if yuReceipt == nil {
		return nil, ErrNotFoundReceipt
	}
	receipt := new(types.Receipt)
	err = json.Unmarshal(yuReceipt.Extra, receipt)
	return receipt, err
}

func (s *Solidity) GetReceipts(ctx *context.ReadContext) {
	var rq ReceiptsRequest
	err := ctx.BindJson(&rq)
	if err != nil {
		ctx.Json(http.StatusBadRequest, &ReceiptsResponse{Err: err})
		return
	}

	receipts := make([]*types.Receipt, 0, len(rq.Hashes))
	for _, hash := range rq.Hashes {
		receipt, err := s.getReceipt(hash)
		if err != nil {
			ctx.Json(http.StatusInternalServerError, &ReceiptsResponse{Err: err})
			return
		}

		receipts = append(receipts, receipt)
	}

	ctx.JsonOk(&ReceiptsResponse{Receipts: receipts})
}

func emitReceipt(ctx *context.WriteContext, vmEmv *vm.EVM, txReq *TxRequest, code []byte, contractAddr common.Address, leftOverGas uint64, err error) error {
	evmReceipt := makeEvmReceipt(ctx, vmEmv, code, ctx.Txn, ctx.Block, contractAddr, leftOverGas, err)
	receiptByt, err := json.Marshal(evmReceipt)
	if err != nil {
		txReqByt, _ := json.Marshal(txReq)
		logrus.Errorf("Receipt marshal err: %v. Tx: %v", err, string(txReqByt))
		return err
	}
	ctx.EmitExtra(receiptByt)
	return nil
}

// endregion ---- Tripod Api ----
