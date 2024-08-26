package evm

import (
	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	yuConfig "github.com/reddio-com/reddio/evm/config"
	"github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type GethConfig struct {
	ChainConfig *params.ChainConfig

	// BlockContext provides the EVM with auxiliary information. Once provided
	// it shouldn't be modified.
	GetHashFn   func(n uint64) common.Hash
	Coinbase    common.Address
	GasLimit    uint64
	BlockNumber *big.Int
	Time        uint64
	Difficulty  *big.Int
	BaseFee     *big.Int
	BlobBaseFee *big.Int
	Random      *common.Hash

	// TxContext provides the EVM with information about a transaction.
	// All fields can change between transactions.
	Origin     common.Address
	GasPrice   *big.Int
	BlobHashes []common.Hash
	BlobFeeCap *big.Int

	// StateDB gives access to the underlying state
	State *state.StateDB

	// Unknown
	Value     *big.Int
	Debug     bool
	EVMConfig vm.Config

	// Global config
	EnableEthRPC bool   `toml:"enable_eth_rpc"`
	EthHost      string `toml:"eth_host"`
	EthPort      string `toml:"eth_port"`
}

func (gc *GethConfig) Copy() *GethConfig {
	return &GethConfig{
		ChainConfig:  gc.ChainConfig,
		GetHashFn:    gc.GetHashFn,
		Coinbase:     gc.Coinbase,
		GasLimit:     gc.GasLimit,
		BlockNumber:  gc.BlockNumber,
		Time:         gc.Time,
		Difficulty:   gc.Difficulty,
		BaseFee:      gc.BaseFee,
		BlobBaseFee:  gc.BlobBaseFee,
		Random:       gc.Random,
		Origin:       gc.Origin,
		GasPrice:     gc.GasPrice,
		BlobHashes:   gc.BlobHashes,
		BlobFeeCap:   gc.BlobFeeCap,
		State:        gc.State,
		Value:        gc.Value,
		Debug:        gc.Debug,
		EVMConfig:    gc.EVMConfig,
		EnableEthRPC: gc.EnableEthRPC,
		EthHost:      gc.EthHost,
		EthPort:      gc.EthPort,
	}
}

// sets defaults on the config
func SetDefaultGethConfig() *GethConfig {
	cfg := &GethConfig{
		ChainConfig: params.AllEthashProtocolChanges,
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
		BaseFee:     big.NewInt(params.InitialBaseFee), // 1 gwei
		BlobBaseFee: big.NewInt(params.BlobTxMinBlobGasprice),
		BlobHashes:  []common.Hash{},
		BlobFeeCap:  big.NewInt(0),
		Random:      &common.Hash{},
		State:       nil,
		GetHashFn: func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		},
	}

	cfg.ChainConfig.ChainID = big.NewInt(50341)

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

func setDefaultEthStateConfig() *yuConfig.Config {
	return &yuConfig.Config{
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
