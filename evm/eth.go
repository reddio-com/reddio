package evm

import (
	// "github.com/yu-org/yu/common/yerror"

	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/NethermindEth/juno/encoder"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/yu-org/yu/common/yerror"

	"github.com/reddio-com/reddio/evm/config"

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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"

	"github.com/holiman/uint256"
)

type Solidity struct {
	sync.Mutex

	*tripod.Tripod
	ethState    *EthState
	cfg         *GethConfig
	stateConfig *config.Config
}

func (s *Solidity) StateDB() *state.StateDB {
	return s.ethState.StateDB()
}

func (s *Solidity) SetStateDB(d *state.StateDB) {
	s.ethState.SetStateDB(d)
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

type GethConfig struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        uint64
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config
	BaseFee     *big.Int
	BlobBaseFee *big.Int
	BlobHashes  []common.Hash
	BlobFeeCap  *big.Int
	Random      *common.Hash

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash

	EnableEthRPC bool   `toml:"enable_eth_rpc"`
	EthHost      string `toml:"eth_host"`
	EthPort      string `toml:"eth_port"`
}

// sets defaults on the config
func SetDefaultGethConfig() *GethConfig {
	cfg := defaultGethConfig()
	cfg.ChainConfig.ChainID = big.NewInt(1)
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			PetersburgBlock:     new(big.Int),
			IstanbulBlock:       new(big.Int),
			MuirGlacierBlock:    new(big.Int),
			BerlinBlock:         new(big.Int),
			LondonBlock:         new(big.Int),
		}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
	if cfg.BaseFee == nil {
		cfg.BaseFee = big.NewInt(params.InitialBaseFee)
	}
	if cfg.BlobBaseFee == nil {
		cfg.BlobBaseFee = big.NewInt(params.BlobTxMinBlobGasprice)
	}

	return cfg
}

func LoadEvmConfig(fpath string) *GethConfig {
	cfg := SetDefaultGethConfig()
	_, err := toml.DecodeFile(fpath, cfg)
	if err != nil {
		logrus.Fatalf("load config file failed: %v", err)
	}
	return cfg
}

func defaultGethConfig() *GethConfig {
	return &GethConfig{
		ChainConfig: params.GoerliChainConfig,
		Difficulty:  big.NewInt(1),
		Origin:      common.HexToAddress("0x0"),
		Coinbase:    common.HexToAddress("0x0"),
		BlockNumber: big.NewInt(0),
		Time:        0,
		GasLimit:    8000000,
		GasPrice:    big.NewInt(1),
		Value:       big.NewInt(0),
		Debug:       false,
		EVMConfig:   vm.Config{},
		BaseFee:     big.NewInt(1000000000), // 1 gwei
		BlobBaseFee: big.NewInt(0),
		BlobHashes:  []common.Hash{},
		BlobFeeCap:  big.NewInt(0),
		Random:      &common.Hash{},

		State:     nil,
		GetHashFn: nil,
	}
}

func setDefaultEthStateConfig() *config.Config {
	return &config.Config{
		VMTrace:                 "",
		VMTraceConfig:           "",
		EnablePreimageRecording: false,
		Recovery:                false,
		NoBuild:                 false,
		SnapshotWait:            false,
		SnapshotCache:           128,              // Default cache size
		TrieCleanCache:          256,              // Default Trie cleanup cache size
		TrieDirtyCache:          256,              // Default Trie dirty cache size
		TrieTimeout:             60 * time.Second, // Default Trie timeout
		Preimages:               false,
		NoPruning:               false,
		NoPrefetch:              false,
		StateHistory:            0,                   // By default, there is no state history
		StateScheme:             "hash",              // Default state scheme
		DbPath:                  "reddio_db",         // Default database path
		DbType:                  "pebble",            // Default database type
		NameSpace:               "eth/db/chaindata/", // Default namespace
		Ancient:                 "ancient",           // Default ancient data path
		Cache:                   512,                 // Default cache size
		Handles:                 64,                  // Default number of handles
	}
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
		solidity.Call, solidity.GetReceipt,
		// solidity.GetClass, solidity.GetClassAt,
		// 	solidity.GetClassHashAt, solidity.GetNonce, solidity.GetStorage,
		// 	solidity.GetTransaction, solidity.GetTransactionStatus,
		// 	solidity.SimulateTransactions,
		// 	solidity.GetBlockWithTxs, solidity.GetBlockWithTxHashes,
	)

	return solidity
}

// region ---- Tripod Api ----

func (s *Solidity) CheckTxn(txn *yu_types.SignedTxn) error {
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
	origin := common.BytesToAddress(txReq.Origin.Bytes())
	s.Lock()
	err = ctx.BindJson(txReq)
	if err != nil {
		return err
	}
	//logrus.Printf("ExecuteTxn: %+v\n", txReq)
	zeroAddress := common.Address{}
	gasLimit := txReq.GasLimit
	gasPrice := txReq.GasPrice
	value := txReq.Value

	cfg := s.cfg
	//ethstate := s.ethState

	cfg.Origin = origin
	cfg.GasLimit = gasLimit
	cfg.GasPrice = gasPrice
	cfg.Value = value

	vmenv := newEVM(cfg)
	pd := pending_state.NewPendingState(ctx.ExtraInterface.(*state.StateDB))
	vmenv.StateDB = pd

	//logrus.Println("ExecuteTxn vmenv: ", vmenv)

	sender := vm.AccountRef(txReq.Origin)
	rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber, vmenv.Context.Random != nil, vmenv.Context.Time)
	s.Unlock()

	if txReq.Address == zeroAddress {
		err = executeContractCreation(txReq, pd, cfg, origin, coinbase, vmenv, sender, rules)
	} else {
		err = executeContractCall(txReq, pd, cfg, origin, coinbase, vmenv, sender, rules)
	}
	if err != nil {
		return err
	}
	ctx.ExtraInterface = pd
	return nil
}

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

	vmenv.StateDB = pending_state.NewPendingState(s.ethState.stateDB)

	if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
		cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: &address, Data: input, Value: value, Gas: gasLimit}), origin)
	}
	// Execute the preparatory steps for state transition which includes:
	// - prepare accessList(post-berlin)
	// - reset transient storage(eip 1153)
	ethState.Prepare(rules, origin, cfg.Coinbase, &address, vm.ActivePrecompiles(rules), nil)

	println("Call Request sender:", sender.Address().Hex())
	println("Call Request address:", address.Hex())
	println("Call Request input:", hex.EncodeToString(input))
	println("Call Request gasLimit:", gasLimit)
	println("Call Request value :", value.String())

	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		gasLimit,
		uint256.MustFromBig(value),
	)
	println("Call Return ret value:", ret)
	println("Call Return ret value:", hex.EncodeToString(ret))
	retBigInt := new(big.Int).SetBytes(ret)
	println("Call Return ret value:", retBigInt.String())
	println("Call Return leftOverGas value:", leftOverGas)

	if err != nil {
		ctx.Json(http.StatusInternalServerError, &CallResponse{Err: err})
		return
	}

	ctx.JsonOk(&CallResponse{Ret: ret, LeftOverGas: leftOverGas})
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

func executeContractCreation(txReq *TxRequest, stateDB *pending_state.PendingState, cfg *GethConfig, origin, coinBase common.Address, vmenv *vm.EVM, sender vm.AccountRef, rules params.Rules) error {
	if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
		cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
	}

	stateDB.Prepare(rules, origin, coinBase, nil, vm.ActivePrecompiles(rules), nil)

	code, address, leftOverGas, err := vmenv.Create(sender, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
	if err != nil {
		return err
	}

	println("Return code value:", code)
	println("Return code value:", hex.EncodeToString(code))
	println("Return address value:", address.Hex())
	println("Return leftOverGas value:", leftOverGas)
	println("Contract deployment successful!")

	return nil
}

func executeContractCall(txReq *TxRequest, ethState *pending_state.PendingState, cfg *GethConfig, origin, coinBase common.Address, vmenv *vm.EVM, sender vm.AccountRef, rules params.Rules) error {
	if cfg.EVMConfig.Tracer != nil && cfg.EVMConfig.Tracer.OnTxStart != nil {
		cfg.EVMConfig.Tracer.OnTxStart(vmenv.GetVMContext(), types.NewTx(&types.LegacyTx{To: &txReq.Address, Data: txReq.Input, Value: txReq.Value, Gas: txReq.GasLimit}), txReq.Origin)
	}

	ethState.Prepare(rules, origin, coinBase, &txReq.Address, vm.ActivePrecompiles(rules), nil)
	ethState.SetNonce(txReq.Origin, ethState.GetNonce(sender.Address())+1)

	logrus.Printf("before transfer: account %s balance %d \n", sender.Address(), ethState.GetBalance(sender.Address()))

	_, _, err := vmenv.Call(sender, txReq.Address, txReq.Input, txReq.GasLimit, uint256.MustFromBig(txReq.Value))
	if err != nil {
		return err
	}

	logrus.Printf("after transfer: account %s balance %d \n", sender.Address(), ethState.GetBalance(sender.Address()))

	//println("Return ret value:", ret)
	//println("Return leftOverGas value:", leftOverGas)

	return nil
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
		return nil, errors.New("no receipt found")
	}
	receipt := new(types.Receipt)
	err = encoder.Unmarshal(yuReceipt.Extra, receipt)
	return receipt, err
}

// endregion ---- Tripod Api ----
