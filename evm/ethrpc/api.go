package ethrpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	yutypes "github.com/yu-org/yu/core/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/gasestimator"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
)

var (
	getLogsLbl    = "getLogs"
	gasPriceLbl   = "gasPrice"
	feeHistoryLbl = "feeHistory"

	statusSuccess = "success"
	statusErr     = "err"
	statusExceed  = "exceed"
)

// estimateGasErrorRatio is the amount of overestimation eth_estimateGas is
// allowed to produce in order to speed up calculations.
const estimateGasErrorRatio = 0.015

var errBlobTxNotSupported = errors.New("signing blob transactions not supported")

// EthereumAPI provides an API to access Ethereum related information.
type EthereumAPI struct {
	b Backend
}

// NewEthereumAPI creates a new Ethereum protocol API.
func NewEthereumAPI(b Backend) *EthereumAPI {
	return &EthereumAPI{b}
}

// GetLogs from *FilterAPI
func (s *EthereumAPI) GetLogs(ctx context.Context, crit FilterCriteria) (results []*types.Log, err error) {
	start := time.Now()
	defer func() {
		if err != nil {
			metrics.EthereumAPICounterHist.WithLabelValues(getLogsLbl, statusErr).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(getLogsLbl, statusErr).Inc()
		} else {
			metrics.EthereumAPICounterHist.WithLabelValues(getLogsLbl, statusSuccess).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(getLogsLbl, statusSuccess).Inc()
		}
	}()

	if len(crit.Topics) > maxTopics {
		return nil, errExceedMaxTopics
	}

	filter, err := newLogFilter(ctx, s.b, crit)
	if err != nil {
		return nil, err
	}

	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		return []*types.Log{}, nil
	}

	return logs, nil
}

// GasPrice returns a suggestion for a gas price for legacy transactions.
func (s *EthereumAPI) GasPrice(ctx context.Context) (result *hexutil.Big, err error) {
	start := time.Now()
	defer func() {
		if err != nil {
			metrics.EthereumAPICounterHist.WithLabelValues(gasPriceLbl, statusErr).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(gasPriceLbl, statusErr).Inc()
		} else {
			metrics.EthereumAPICounterHist.WithLabelValues(gasPriceLbl, statusSuccess).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(gasPriceLbl, statusSuccess).Inc()
		}
	}()
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	if head := s.b.CurrentHeader(); head.BaseFee != nil {
		tipcap.Add(tipcap, head.BaseFee)
	}
	return (*hexutil.Big)(tipcap), err
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (s *EthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	tipcap, err := s.b.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(tipcap), err
}

type feeHistoryResult struct {
	OldestBlock      *hexutil.Big     `json:"oldestBlock"`
	Reward           [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee          []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio     []float64        `json:"gasUsedRatio"`
	BlobBaseFee      []*hexutil.Big   `json:"baseFeePerBlobGas,omitempty"`
	BlobGasUsedRatio []float64        `json:"blobGasUsedRatio,omitempty"`
}

// FeeHistory returns the fee market history.
func (s *EthereumAPI) FeeHistory(ctx context.Context, blockCount math.HexOrDecimal64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (result *feeHistoryResult, err error) {
	start := time.Now()
	defer func() {
		if err != nil {
			metrics.EthereumAPICounterHist.WithLabelValues(feeHistoryLbl, statusErr).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(feeHistoryLbl, statusErr).Inc()
		} else {
			metrics.EthereumAPICounterHist.WithLabelValues(feeHistoryLbl, statusSuccess).Observe(float64(time.Since(start).Microseconds()))
			metrics.EthereumAPICounter.WithLabelValues(feeHistoryLbl, statusSuccess).Inc()
		}
	}()
	oldest, reward, baseFee, gasUsed, blobBaseFee, blobGasUsed, err := s.b.FeeHistory(ctx, uint64(blockCount), lastBlock, rewardPercentiles)
	if err != nil {
		return nil, err
	}
	results := &feeHistoryResult{
		OldestBlock:  (*hexutil.Big)(oldest),
		GasUsedRatio: gasUsed,
	}
	if reward != nil {
		results.Reward = make([][]*hexutil.Big, len(reward))
		for i, w := range reward {
			results.Reward[i] = make([]*hexutil.Big, len(w))
			for j, v := range w {
				results.Reward[i][j] = (*hexutil.Big)(v)
			}
		}
	}
	if baseFee != nil {
		results.BaseFee = make([]*hexutil.Big, len(baseFee))
		for i, v := range baseFee {
			results.BaseFee[i] = (*hexutil.Big)(v)
		}
	}
	if blobBaseFee != nil {
		results.BlobBaseFee = make([]*hexutil.Big, len(blobBaseFee))
		for i, v := range blobBaseFee {
			results.BlobBaseFee[i] = (*hexutil.Big)(v)
		}
	}
	if blobGasUsed != nil {
		results.BlobGasUsedRatio = blobGasUsed
	}
	return results, nil
}

// BlobBaseFee returns the base fee for blob gas at the current head.
func (s *EthereumAPI) BlobBaseFee(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(s.b.BlobBaseFee(ctx))
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up-to-date or has not
// yet received the latest block headers from its peers. In case it is synchronizing:
// - startingBlock: block number this node started to synchronize from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *EthereumAPI) Syncing() (interface{}, error) {
	progress := s.b.SyncProgress()

	// Return not syncing if the synchronization already completed
	if progress.Done() {
		return false, nil
	}
	// Otherwise gather the block sync stats
	return map[string]interface{}{
		"startingBlock":          hexutil.Uint64(progress.StartingBlock),
		"currentBlock":           hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":           hexutil.Uint64(progress.HighestBlock),
		"syncedAccounts":         hexutil.Uint64(progress.SyncedAccounts),
		"syncedAccountBytes":     hexutil.Uint64(progress.SyncedAccountBytes),
		"syncedBytecodes":        hexutil.Uint64(progress.SyncedBytecodes),
		"syncedBytecodeBytes":    hexutil.Uint64(progress.SyncedBytecodeBytes),
		"syncedStorage":          hexutil.Uint64(progress.SyncedStorage),
		"syncedStorageBytes":     hexutil.Uint64(progress.SyncedStorageBytes),
		"healedTrienodes":        hexutil.Uint64(progress.HealedTrienodes),
		"healedTrienodeBytes":    hexutil.Uint64(progress.HealedTrienodeBytes),
		"healedBytecodes":        hexutil.Uint64(progress.HealedBytecodes),
		"healedBytecodeBytes":    hexutil.Uint64(progress.HealedBytecodeBytes),
		"healingTrienodes":       hexutil.Uint64(progress.HealingTrienodes),
		"healingBytecode":        hexutil.Uint64(progress.HealingBytecode),
		"txIndexFinishedBlocks":  hexutil.Uint64(progress.TxIndexFinishedBlocks),
		"txIndexRemainingBlocks": hexutil.Uint64(progress.TxIndexRemainingBlocks),
	}, nil
}

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) (*keystore.KeyStore, error) {
	if ks := am.Backends(keystore.KeyStoreType); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore), nil
	}
	return nil, errors.New("local keystore not used")
}

// BlockChainAPI provides an API to access Ethereum blockchain data.
type BlockChainAPI struct {
	b Backend
}

// NewBlockChainAPI creates a new Ethereum blockchain API.
func NewBlockChainAPI(b Backend) *BlockChainAPI {
	return &BlockChainAPI{b}
}

// ChainId is the EIP-155 replay-protection chain id for the current Ethereum chain config.
//
// Note, this method does not conform to EIP-695 because the configured chain ID is always
// returned, regardless of the current head block. We used to return an error when the chain
// wasn't synced up to a block where EIP-155 is enabled, but this behavior caused issues
// in CL clients.
func (api *BlockChainAPI) ChainId() *hexutil.Big {
	return (*hexutil.Big)(api.b.ChainConfig().ChainID)
}

// BlockNumber returns the block number of the chain head.
func (s *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return hexutil.Uint64(header.Number.Uint64())
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *BlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	b := state.GetBalance(address).ToBig()
	return (*hexutil.Big)(b), state.Error()
}

// AccountResult structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// proofList implements ethdb.KeyValueWriter and collects the proofs as
// hex-strings for delivery to rpc-caller.
type proofList []string

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, hexutil.Encode(value))
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys.
func (s *BlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	var (
		keys         = make([]common.Hash, len(storageKeys))
		keyLengths   = make([]int, len(storageKeys))
		storageProof = make([]StorageResult, len(storageKeys))
	)
	// Deserialize all keys. This prevents state access on invalid input.
	for i, hexKey := range storageKeys {
		var err error
		keys[i], keyLengths[i], err = decodeHash(hexKey)
		if err != nil {
			return nil, err
		}
	}
	statedb, header, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if statedb == nil || err != nil {
		return nil, err
	}
	codeHash := statedb.GetCodeHash(address)
	storageRoot := statedb.GetStorageRoot(address)

	if len(keys) > 0 {
		var storageTrie state.Trie
		if storageRoot != types.EmptyRootHash && storageRoot != (common.Hash{}) {
			id := trie.StorageTrieID(header.Root, crypto.Keccak256Hash(address.Bytes()), storageRoot)
			st, err := trie.NewStateTrie(id, statedb.Database().TrieDB())
			if err != nil {
				return nil, err
			}
			storageTrie = st
		}
		// Create the proofs for the storageKeys.
		for i, key := range keys {
			// Output key encoding is a bit special: if the input was a 32-byte hash, it is
			// returned as such. Otherwise, we apply the QUANTITY encoding mandated by the
			// JSON-RPC spec for getProof. This behavior exists to preserve backwards
			// compatibility with older client versions.
			var outputKey string
			if keyLengths[i] != 32 {
				outputKey = hexutil.EncodeBig(key.Big())
			} else {
				outputKey = hexutil.Encode(key[:])
			}
			if storageTrie == nil {
				storageProof[i] = StorageResult{outputKey, &hexutil.Big{}, []string{}}
				continue
			}
			var proof proofList
			if err := storageTrie.Prove(crypto.Keccak256(key.Bytes()), &proof); err != nil {
				return nil, err
			}
			value := (*hexutil.Big)(statedb.GetState(address, key).Big())
			storageProof[i] = StorageResult{outputKey, value, proof}
		}
	}
	// Create the accountProof.
	tr, err := trie.NewStateTrie(trie.StateTrieID(header.Root), statedb.Database().TrieDB())
	if err != nil {
		return nil, err
	}
	var accountProof proofList
	if err := tr.Prove(crypto.Keccak256(address.Bytes()), &accountProof); err != nil {
		return nil, err
	}
	balance := statedb.GetBalance(address).ToBig()
	return &AccountResult{
		Address:      address,
		AccountProof: accountProof,
		Balance:      (*hexutil.Big)(balance),
		CodeHash:     codeHash,
		Nonce:        hexutil.Uint64(statedb.GetNonce(address)),
		StorageHash:  storageRoot,
		StorageProof: storageProof,
	}, statedb.Error()
}

// decodeHash parses a hex-encoded 32-byte hash. The input may optionally
// be prefixed by 0x and can have a byte length up to 32.
func decodeHash(s string) (h common.Hash, inputLength int, err error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if (len(s) & 1) > 0 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return common.Hash{}, 0, errors.New("hex string invalid")
	}
	if len(b) > 32 {
		return common.Hash{}, len(b), errors.New("hex string too long, want at most 32 bytes")
	}
	return common.BytesToHash(b), len(b), nil
}

// GetHeaderByNumber returns the requested canonical block header.
//   - When blockNr is -1 the chain pending header is returned.
//   - When blockNr is -2 the chain latest header is returned.
//   - When blockNr is -3 the chain finalized header is returned.
//   - When blockNr is -4 the chain safe header is returned.
func (s *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, yuHeader, err := s.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		response := s.rpcMarshalHeader(ctx, header, yuHeader)
		if number == rpc.PendingBlockNumber {
			// Pending header need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetHeaderByHash returns the requested header by hash.
func (s *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	header, yuHeader, _ := s.b.HeaderByHash(ctx, hash)
	if header != nil {
		return s.rpcMarshalHeader(ctx, header, yuHeader)
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
//   - When blockNr is -1 the chain pending block is returned.
//   - When blockNr is -2 the chain latest block is returned.
//   - When blockNr is -3 the chain finalized block is returned.
//   - When blockNr is -4 the chain safe block is returned.
//   - When fullTx is true all transactions in the block are returned, otherwise
//     only the transaction hash is returned.
func (s *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, yuBlock, err := s.b.BlockByNumber(ctx, number)
	if block != nil && err == nil {
		response, err := s.rpcMarshalBlock(ctx, block, yuBlock, true, fullTx)
		if err == nil && number == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		if err != nil {
			logrus.Warnf("GetBlockByNumber Failed to marshal block %v: %v ", number, err)
		}
		return response, err
	}
	if err != nil {
		logrus.Errorf("Failed to get block by number %v: %v", number, err)
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, yuBlock, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		return s.rpcMarshalBlock(ctx, block, yuBlock, true, fullTx)
	}
	return nil, err
}

// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index.
func (s *BlockChainAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, yuBlock, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			logrus.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, yuBlock, false, false)
	}
	return nil, err
}

// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index.
func (s *BlockChainAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, yuBlock, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			logrus.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader(uncles[index])
		return s.rpcMarshalBlock(ctx, block, yuBlock, false, false)
	}
	return nil, err
}

// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
func (s *BlockChainAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
func (s *BlockChainAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (s *BlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (s *BlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	key, _, err := decodeHash(hexKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode storage key: %s", err)
	}
	res := state.GetState(address, key)
	return res[:], state.Error()
}

// GetBlockReceipts returns the block receipts for the given block hash or number or tag.
func (s *BlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	block, yuBlock, err := s.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if block == nil || err != nil {
		// When the block doesn't exist, the RPC method should return JSON null
		// as per specification.
		return nil, nil
	}
	receipts, err := s.b.GetReceipts(ctx, common.Hash(yuBlock.Hash))
	if err != nil {
		return nil, err
	}
	txs := block.Transactions()
	if len(txs) != len(receipts) {
		return nil, fmt.Errorf("receipts length mismatch: %d vs %d", len(txs), len(receipts))
	}

	// Derive the sender.
	signer := types.MakeSigner(s.b.ChainConfig(), block.Number(), block.Time())

	result := make([]map[string]interface{}, len(receipts))
	for i, receipt := range receipts {
		result[i] = marshalReceipt(receipt, common.Hash(yuBlock.Hash), block.NumberU64(), signer, txs[i], i)
	}

	return result, nil
}

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride map[common.Address]OverrideAccount

// Apply overrides the fields of specified accounts into the given state.
func (diff *StateOverride) Apply(statedb *state.StateDB) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			statedb.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			statedb.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			u256Balance, _ := uint256.FromBig((*big.Int)(*account.Balance))
			statedb.SetBalance(addr, u256Balance, tracing.BalanceChangeUnspecified)
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			statedb.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				statedb.SetState(addr, key, value)
			}
		}
	}
	// Now finalize the changes. Finalize is normally performed between transactions.
	// By using finalize, the overrides are semantically behaving as
	// if they were created in a transaction just before the tracing occur.
	statedb.Finalise(false)
	return nil
}

// BlockOverrides is a set of header fields to override.
type BlockOverrides struct {
	Number      *hexutil.Big
	Difficulty  *hexutil.Big
	Time        *hexutil.Uint64
	GasLimit    *hexutil.Uint64
	Coinbase    *common.Address
	Random      *common.Hash
	BaseFee     *hexutil.Big
	BlobBaseFee *hexutil.Big
}

// Apply overrides the given header fields into the given block context.
func (diff *BlockOverrides) Apply(blockCtx *vm.BlockContext) {
	if diff == nil {
		return
	}
	if diff.Number != nil {
		blockCtx.BlockNumber = diff.Number.ToInt()
	}
	if diff.Difficulty != nil {
		blockCtx.Difficulty = diff.Difficulty.ToInt()
	}
	if diff.Time != nil {
		blockCtx.Time = uint64(*diff.Time)
	}
	if diff.GasLimit != nil {
		blockCtx.GasLimit = uint64(*diff.GasLimit)
	}
	if diff.Coinbase != nil {
		blockCtx.Coinbase = *diff.Coinbase
	}
	if diff.Random != nil {
		blockCtx.Random = diff.Random
	}
	if diff.BaseFee != nil {
		blockCtx.BaseFee = diff.BaseFee.ToInt()
	}
	if diff.BlobBaseFee != nil {
		blockCtx.BlobBaseFee = diff.BlobBaseFee.ToInt()
	}
}

// ChainContextBackend provides methods required to implement ChainContext.
type ChainContextBackend interface {
	Engine() consensus.Engine
	HeaderByNumber(context.Context, rpc.BlockNumber) (*types.Header, *yutypes.Header, error)
}

// ChainContext is an implementation of core.ChainContext. It's main use-case
// is instantiating a vm.BlockContext without having access to the BlockChain object.
type ChainContext struct {
	b   ChainContextBackend
	ctx context.Context
}

// NewChainContext creates a new ChainContext object.
func NewChainContext(ctx context.Context, backend ChainContextBackend) *ChainContext {
	return &ChainContext{ctx: ctx, b: backend}
}

func (context *ChainContext) Engine() consensus.Engine {
	return context.b.Engine()
}

func (context *ChainContext) GetHeader(hash common.Hash, number uint64) *types.Header {
	// This method is called to get the hash for a block number when executing the BLOCKHASH
	// opcode. Hence no need to search for non-canonical blocks.
	header, _, err := context.b.HeaderByNumber(context.ctx, rpc.BlockNumber(number))
	if err != nil || header.Hash() != hash {
		return nil
	}
	return header
}

// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
func (s *BlockChainAPI) Call(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	return s.b.Call(ctx, args, blockNrOrHash, overrides, blockOverrides)
}

// DoEstimateGas returns the lowest possible gas limit that allows the transaction to run
// successfully at block `blockNrOrHash`. It returns error if the transaction would revert, or if
// there are unexpected failures. The gas limit is capped by both `args.Gas` (if non-nil &
// non-zero) and `gasCap` (if non-zero).
func DoEstimateGas(ctx context.Context, b Backend, args TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, gasCap uint64) (hexutil.Uint64, error) {
	// Retrieve the base stateDb and mutate it with any overrides
	stateDb, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if stateDb == nil || err != nil {
		return 0, err
	}
	if err = overrides.Apply(stateDb); err != nil {
		return 0, err
	}
	// Construct the gas estimator option from the user input
	opts := &gasestimator.Options{
		Config:     b.ChainConfig(),
		Chain:      NewChainContext(ctx, b),
		Header:     header,
		State:      stateDb,
		ErrorRatio: estimateGasErrorRatio,
	}
	originGasLimit := args.Gas
	args.Gas = nil

	if err := args.CallDefaults(gasCap, header.BaseFee, b.ChainConfig().ChainID); err != nil {
		return 0, err
	}
	call := args.ToMessage(header.BaseFee)
	// Run the gas estimation andwrap any revertals into a custom return
	estimate, revert, err := gasestimator.Estimate(ctx, call, opts, gasCap)
	if err != nil {
		if len(revert) > 0 {
			return 0, evm.NewRevertError(revert)
		}
		return 0, err
	}

	if originGasLimit != nil && uint64(*originGasLimit) < estimate {
		return 0, fmt.Errorf("gas required exceeds allowance (%d)", uint64(*originGasLimit))
	}
	return hexutil.Uint64(estimate), nil
}

// EstimateGas returns the lowest possible gas limit that allows the transaction to run
// successfully at block `blockNrOrHash`, or the latest block if `blockNrOrHash` is unspecified. It
// returns error if the transaction would revert or if there are unexpected failures. The returned
// value is capped by both `args.Gas` (if non-nil & non-zero) and the backend's RPCGasCap
// configuration (if non-zero).
// Note: Required blob gas is not computed in this method.
func (s *BlockChainAPI) EstimateGas(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, overrides, s.b.RPCGasCap())
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader(head *types.Header) map[string]interface{} {
	result := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
	if head.BaseFee != nil {
		result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
	}
	if head.WithdrawalsHash != nil {
		result["withdrawalsRoot"] = head.WithdrawalsHash
	}
	if head.BlobGasUsed != nil {
		result["blobGasUsed"] = hexutil.Uint64(*head.BlobGasUsed)
	}
	if head.ExcessBlobGas != nil {
		result["excessBlobGas"] = hexutil.Uint64(*head.ExcessBlobGas)
	}
	if head.ParentBeaconRoot != nil {
		result["parentBeaconBlockRoot"] = head.ParentBeaconRoot
	}
	return result
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func RPCMarshalBlock(block *types.Block, yuBlock *yutypes.Block, inclTx bool, fullTx bool, config *params.ChainConfig) map[string]interface{} {
	fields := RPCMarshalHeader(block.Header())
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(idx int, tx *types.Transaction) interface{} {
			return tx.Hash()
		}
		if fullTx {
			formatTx = func(idx int, tx *types.Transaction) interface{} {
				return newRPCTransactionFromBlockIndex(block, yuBlock, uint64(idx), config)
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		for i, tx := range txs {
			transactions[i] = formatTx(i, tx)
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes
	if block.Header().WithdrawalsHash != nil {
		fields["withdrawals"] = block.Withdrawals()
	}
	return fields
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `BlockchainAPI`.
func (s *BlockChainAPI) rpcMarshalHeader(ctx context.Context, header *types.Header, yuHeader *yutypes.Header) map[string]interface{} {
	fields := RPCMarshalHeader(header)
	fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, header.Hash()))
	fields["hash"] = yuHeader.Hash
	return fields
}

// rpcMarshalBlock uses the generalized output filler, then adds the total difficulty field, which requires
// a `BlockchainAPI`.
func (s *BlockChainAPI) rpcMarshalBlock(ctx context.Context, b *types.Block, yuBlock *yutypes.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields := RPCMarshalBlock(b, yuBlock, inclTx, fullTx, s.b.ChainConfig())
	if inclTx {
		fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, b.Hash()))
	}
	fields["hash"] = yuBlock.Hash
	return fields, nil
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash           *common.Hash      `json:"blockHash"`
	BlockNumber         *hexutil.Big      `json:"blockNumber"`
	From                common.Address    `json:"from"`
	Gas                 hexutil.Uint64    `json:"gas"`
	GasPrice            *hexutil.Big      `json:"gasPrice"`
	GasFeeCap           *hexutil.Big      `json:"maxFeePerGas,omitempty"`
	GasTipCap           *hexutil.Big      `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas    *hexutil.Big      `json:"maxFeePerBlobGas,omitempty"`
	Hash                common.Hash       `json:"hash"`
	Input               hexutil.Bytes     `json:"input"`
	Nonce               hexutil.Uint64    `json:"nonce"`
	To                  *common.Address   `json:"to"`
	TransactionIndex    *hexutil.Uint64   `json:"transactionIndex"`
	Value               *hexutil.Big      `json:"value"`
	Type                hexutil.Uint64    `json:"type"`
	Accesses            *types.AccessList `json:"accessList,omitempty"`
	ChainID             *hexutil.Big      `json:"chainId,omitempty"`
	BlobVersionedHashes []common.Hash     `json:"blobVersionedHashes,omitempty"`
	V                   *hexutil.Big      `json:"v"`
	R                   *hexutil.Big      `json:"r"`
	S                   *hexutil.Big      `json:"s"`
	YParity             *hexutil.Uint64   `json:"yParity,omitempty"`
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(tx *types.Transaction, blockHash common.Hash, blockNumber uint64, blockTime uint64, index uint64, baseFee *big.Int, config *params.ChainConfig) *RPCTransaction {
	signer := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	switch tx.Type() {
	case types.LegacyTxType:
		// if a legacy transaction has an EIP-155 chain id, include it explicitly
		if id := tx.ChainId(); id.Sign() != 0 {
			result.ChainID = (*hexutil.Big)(id)
		}

	case types.AccessListTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign())
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity

	case types.DynamicFeeTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign())
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			// price = min(gasTipCap + baseFee, gasFeeCap)
			result.GasPrice = (*hexutil.Big)(effectiveGasPrice(tx, baseFee))
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}

	case types.BlobTxType:
		al := tx.AccessList()
		yparity := hexutil.Uint64(v.Sign())
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.YParity = &yparity
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (common.Hash{}) {
			result.GasPrice = (*hexutil.Big)(effectiveGasPrice(tx, baseFee))
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
		result.MaxFeePerBlobGas = (*hexutil.Big)(tx.BlobGasFeeCap())
		result.BlobVersionedHashes = tx.BlobHashes()
	}
	return result
}

// effectiveGasPrice computes the transaction gas fee, based on the given basefee value.
//
//	price = min(gasTipCap + baseFee, gasFeeCap)
func effectiveGasPrice(tx *types.Transaction, baseFee *big.Int) *big.Int {
	fee := tx.GasTipCap()
	fee = fee.Add(fee, baseFee)
	if tx.GasFeeCapIntCmp(fee) < 0 {
		return tx.GasFeeCap()
	}
	return fee
}

// NewRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func NewRPCPendingTransaction(tx *types.Transaction, current *types.Header, config *params.ChainConfig) *RPCTransaction {
	var (
		baseFee     *big.Int
		blockNumber = uint64(0)
		blockTime   = uint64(0)
	)
	if current != nil {
		baseFee = eip1559.CalcBaseFee(config, current)
		blockNumber = current.Number.Uint64()
		blockTime = current.Time
	}
	return newRPCTransaction(tx, common.Hash{}, blockNumber, blockTime, 0, baseFee, config)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex(b *types.Block, yuBlock *yutypes.Block, index uint64, config *params.ChainConfig) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	//return newRPCTransaction(txs[index], b.Hash(), b.NumberU64(), b.Time(), index, b.BaseFee(), config)
	return newRPCTransaction(txs[index], common.Hash(yuBlock.Hash), b.NumberU64(), b.Time(), index, b.BaseFee(), config)
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	// TODO: If this function need yuBlock?
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := txs[index].MarshalBinary()
	return blob
}

// accessListResult returns an optional accesslist
// It's the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates an EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	acl, gasUsed, vmerr, err := AccessList(ctx, s.b, bNrOrHash, args)
	if err != nil {
		return nil, err
	}
	result := &accessListResult{Accesslist: &acl, GasUsed: hexutil.Uint64(gasUsed)}
	if vmerr != nil {
		result.Error = vmerr.Error()
	}
	return result, nil
}

// AccessList creates an access list for the given transaction.
// If the accesslist creation fails an error is returned.
// If the transaction itself fails, an vmErr is returned.
func AccessList(ctx context.Context, b Backend, blockNrOrHash rpc.BlockNumberOrHash, args TransactionArgs) (acl types.AccessList, gasUsed uint64, vmErr error, err error) {
	// Retrieve the execution context
	db, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if db == nil || err != nil {
		return nil, 0, nil, err
	}

	// Ensure any missing fields are filled, extract the recipient and input data
	if err := args.setDefaults(ctx, b, true); err != nil {
		return nil, 0, nil, err
	}
	var to common.Address
	if args.To != nil {
		to = *args.To
	} else {
		to = crypto.CreateAddress(args.from(), uint64(*args.Nonce))
	}
	isPostMerge := header.Difficulty.Sign() == 0
	// Retrieve the precompiles since they don't need to be added to the access list
	precompiles := vm.ActivePrecompiles(b.ChainConfig().Rules(header.Number, isPostMerge, header.Time))

	// Create an initial tracer
	prevTracer := logger.NewAccessListTracer(nil, args.from(), to, precompiles)
	if args.AccessList != nil {
		prevTracer = logger.NewAccessListTracer(*args.AccessList, args.from(), to, precompiles)
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, 0, nil, err
		}
		// Retrieve the current access list to expand
		accessList := prevTracer.AccessList()
		logrus.Trace("Creating access list", "input", accessList)

		// Copy the original db so we don't modify it
		statedb := db.Copy()
		// Set the accesslist to the last al
		args.AccessList = &accessList
		msg := args.ToMessage(header.BaseFee)

		// Apply the transaction with the access list tracer
		tracer := logger.NewAccessListTracer(accessList, args.from(), to, precompiles)
		config := vm.Config{Tracer: tracer.Hooks(), NoBaseFee: true}
		vmenv := b.GetEVM(ctx, msg, statedb, header, &config, nil)
		res, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.GasLimit))
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to apply transaction: %v err: %v", args.ToTransaction(nil, nil, nil).Hash(), err)
		}
		if tracer.Equal(prevTracer) {
			return accessList, res.UsedGas, res.Err, nil
		}
		prevTracer = tracer
	}
}

// TransactionAPI exposes methods for reading and creating transaction data.
type TransactionAPI struct {
	b         Backend
	nonceLock *AddrLocker
	signer    types.Signer
}

// NewTransactionAPI creates a new RPC service with methods for interacting with transactions.
func NewTransactionAPI(b Backend, nonceLock *AddrLocker) *TransactionAPI {
	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	signer := types.LatestSigner(b.ChainConfig())
	return &TransactionAPI{b, nonceLock, signer}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	TransactionAPICounter.WithLabelValues("GetBlockTransactionCountByNumber").Inc()
	if block, _, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	TransactionAPICounter.WithLabelValues("GetBlockTransactionCountByHash").Inc()
	if block, _, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *TransactionAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	TransactionAPICounter.WithLabelValues("GetTransactionByBlockNumberAndIndex").Inc()
	if block, yuBlock, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, yuBlock, uint64(index), s.b.ChainConfig())
	}
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *TransactionAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	TransactionAPICounter.WithLabelValues("GetTransactionByBlockHashAndIndex").Inc()
	if block, yuBlock, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, yuBlock, uint64(index), s.b.ChainConfig())
	}
	return nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *TransactionAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	TransactionAPICounter.WithLabelValues("GetRawTransactionByBlockNumberAndIndex").Inc()
	if block, _, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *TransactionAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	TransactionAPICounter.WithLabelValues("GetRawTransactionByBlockHashAndIndex").Inc()
	if block, _, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *TransactionAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	TransactionAPICounter.WithLabelValues("GetTransactionCount").Inc()
	// Ask transaction pool for the nonce which includes pending transactions
	nonce := uint64(0)
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		poolNonce, err := s.b.GetPoolNonce(ctx, address)
		if err != nil {
			return nil, err
		}
		nonce += poolNonce
	}
	// Resolve block number and use its state to ask for the nonce
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	nonce += state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}

// GetTransactionByHash returns the transaction for the given hash
func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	TransactionAPICounter.WithLabelValues("GetTransactionByHash").Inc()
	// Try to return an already finalized transaction
	found, tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if !found {
		// In this case, it indicates that GetReceipt returned a specific error, not just that it wasn't found.
		if err != nil && !errors.Is(err, evm.ErrNotFoundReceipt) {
			return nil, err
		}

		// No finalized transaction, try to retrieve it from the pool
		tx, err = s.b.GetPoolTransaction(hash)
		if err != nil {
			return nil, err
		}
		if tx != nil {
			return NewRPCPendingTransaction(tx, s.b.CurrentHeader(), s.b.ChainConfig()), nil
		}

		if err != nil {
			logrus.Warnf("[GetTransactionByHash] Failed to get tx %s, err: %v", hash.String(), err)
		}
		return nil, nil
		// Only when err==nil should it return nil, and when err!=nil, it should return NewTxIndexingError.
		// However, this would cause errors with packages such as Hardhat. In order to handle compatibility,
		// it currently returns nil directly.
		// return nil, NewTxIndexingError()
	}
	header, yuHeader, err := s.b.HeaderByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	return newRPCTransaction(tx, common.Hash(yuHeader.Hash), blockNumber, header.Time, index, header.BaseFee, s.b.ChainConfig()), nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *TransactionAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	TransactionAPICounter.WithLabelValues("GetRawTransactionByHash").Inc()
	// Retrieve a finalized transaction, or a pooled otherwise
	found, tx, _, _, _, err := s.b.GetTransaction(ctx, hash)
	if !found {
		// In this case, it indicates that GetReceipt returned a specific error, not just that it wasn't found.
		if err != nil && !errors.Is(err, evm.ErrNotFoundReceipt) {
			return nil, err
		}

		tx, err = s.b.GetPoolTransaction(hash)
		if err != nil {
			return nil, err
		}
		if tx != nil {
			return tx.MarshalBinary()
		}
		if err != nil {
			logrus.Warnf("[GetRawTransactionByHash] Failed to get tx %s, err: %v", hash.String(), err)
		}
		// Only when err==nil should it return nil, and when err!=nil, it should return NewTxIndexingError.
		// However, this would cause errors with packages such as Hardhat. In order to handle compatibility,
		// it currently returns nil directly.
		return nil, nil
		// return nil, NewTxIndexingError()
	}
	return tx.MarshalBinary()
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	TransactionAPICounter.WithLabelValues("GetTransactionReceipt").Inc()
	found, tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		// Only when err==nil should it return nil, and when err!=nil, it should return NewTxIndexingError.
		// However, this would cause errors with packages such as Hardhat. In order to handle compatibility,
		// it currently returns nil directly.
		logrus.Warnf("[GetTransactionReceipt] Failed to get tx %s, err: %v", hash.String(), err)
		// return nil, NewTxIndexingError() // transaction is not fully indexed
	}
	if !found {
		return nil, nil // transaction is not existent or reachable
	}
	header, yuHeader, err := s.b.HeaderByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	receipt, err := s.b.GetReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}
	// Derive the sender.
	signer := types.MakeSigner(s.b.ChainConfig(), header.Number, header.Time)
	return marshalReceipt(receipt, common.Hash(yuHeader.Hash), blockNumber, signer, tx, int(index)), nil
}

// marshalReceipt marshals a transaction receipt into a JSON object.
func marshalReceipt(receipt *types.Receipt, blockHash common.Hash, blockNumber uint64, signer types.Signer, tx *types.Transaction, txIndex int) map[string]interface{} {
	from, _ := types.Sender(signer, tx)

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   tx.Hash(),
		"transactionIndex":  hexutil.Uint64(txIndex),
		"from":              from,
		"to":                tx.To(),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(tx.Type()),
		"effectiveGasPrice": (*hexutil.Big)(receipt.EffectiveGasPrice),
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = []*types.Log{}
	}

	if tx.Type() == types.BlobTxType {
		fields["blobGasUsed"] = hexutil.Uint64(receipt.BlobGasUsed)
		fields["blobGasPrice"] = (*hexutil.Big)(receipt.BlobGasPrice)
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields
}

// sign is a helper function that signs a transaction with the private key of the given address.
func (s *TransactionAPI) sign(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTx(account, tx, s.b.ChainConfig().ChainID)
}

// SubmitTransaction is a helper function that submits tx to txPool and logs a message.
func SubmitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
	// If the transaction fee cap is already specified, ensure the
	// fee of the given transaction is _reasonable_.
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	if !b.UnprotectedAllowed() && !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}
	// Print a log with full tx details for manual investigations and interventions
	head := b.CurrentBlock()
	signer := types.MakeSigner(b.ChainConfig(), head.Number, head.Time)
	from, err := types.Sender(signer, tx)
	if err != nil {
		return common.Hash{}, err
	}

	if tx.To() == nil {
		addr := crypto.CreateAddress(from, tx.Nonce())
		logrus.Debugf("Submitted contract creation: Hash=%s, From=%s, Nonce=%v, ContractAdd=%v, Value=%v", tx.Hash().Hex(), from, tx.Nonce(), addr.Hex(), tx.Value())
	} else {
		logrus.Debugf("Submitted transaction: Hash=%s, From=%s, Nonce=%v, To=%v, Value=%v", tx.Hash().Hex(), from, tx.Nonce(), tx.To(), tx.Value())
	}

	return tx.Hash(), nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *TransactionAPI) SendTransaction(ctx context.Context, args TransactionArgs) (common.Hash, error) {
	panic("implement me")
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (s *TransactionAPI) FillTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	panic("implement me")
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (txHash common.Hash, err error) {
	TransactionAPICounter.WithLabelValues("SendRawTransaction").Inc()
	defer func() {
		if err != nil {
			if !config.GetGlobalConfig().IsBenchmarkMode {
				logrus.Errorf("SendRawTransaction failed: %v, txHash(%s)", err, txHash.String())
			}
		}
	}()
	tx := new(types.Transaction)
	if err = tx.UnmarshalBinary(input); err != nil {
		return
	}
	txHash, err = SubmitTransaction(ctx, s.b, tx)
	return
}

type BatchTx struct {
	TxsBytes []hexutil.Bytes `json:"txs_bytes"`
}

func (s *TransactionAPI) SendBatchRawTransactions(ctx context.Context, inputs hexutil.Bytes) (txHashes []common.Hash, err error) {
	TransactionAPICounter.WithLabelValues("SendBatchRawTransactions").Inc()
	batchTx := new(BatchTx)
	err = json.Unmarshal(inputs, batchTx)
	if err != nil {
		return
	}
	for _, input := range batchTx.TxsBytes {
		var txHash common.Hash
		txHash, err = s.SendRawTransaction(ctx, input)
		if err != nil {
			return
		}
		txHashes = append(txHashes, txHash)
	}
	return
}

// Sign calculates an ECDSA signature for:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *TransactionAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Sign the requested hash with the wallet
	signature, err := wallet.SignText(account, data)
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, err
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *TransactionAPI) SignTransaction(ctx context.Context, args TransactionArgs) (*SignTransactionResult, error) {
	panic("implement me")
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (s *TransactionAPI) PendingTransactions() ([]*RPCTransaction, error) {
	TransactionAPICounter.WithLabelValues("PendingTransactions").Inc()
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := make(map[common.Address]struct{})
	for _, wallet := range s.b.AccountManager().Wallets() {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = struct{}{}
		}
	}
	curHeader := s.b.CurrentHeader()
	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from, _ := types.Sender(s.signer, tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, NewRPCPendingTransaction(tx, curHeader, s.b.ChainConfig()))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *TransactionAPI) Resend(ctx context.Context, sendArgs TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	panic("implement me")
}

// DebugAPI is the collection of Ethereum APIs exposed over the debugging
// namespace.
type DebugAPI struct {
	b Backend
}

// NewDebugAPI creates a new instance of DebugAPI.
func NewDebugAPI(b Backend) *DebugAPI {
	return &DebugAPI{b: b}
}

// GetRawHeader retrieves the RLP encoding for a single header.
func (api *DebugAPI) GetRawHeader(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	var hash common.Hash
	if h, ok := blockNrOrHash.Hash(); ok {
		hash = h
	} else {
		_, yuBlock, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return nil, err
		}
		hash = common.Hash(yuBlock.Hash)
	}
	header, _, _ := api.b.HeaderByHash(ctx, hash)
	if header == nil {
		return nil, fmt.Errorf("header #%d not found", hash)
	}
	return rlp.EncodeToBytes(header)
}

// GetRawBlock retrieves the RLP encoded for a single block.
func (api *DebugAPI) GetRawBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	var hash common.Hash
	if h, ok := blockNrOrHash.Hash(); ok {
		hash = h
	} else {
		_, yuBlock, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return nil, err
		}
		hash = common.Hash(yuBlock.Hash)
	}
	block, _, _ := api.b.BlockByHash(ctx, hash)
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", hash)
	}
	return rlp.EncodeToBytes(block)
}

// GetRawReceipts retrieves the binary-encoded receipts of a single block.
func (api *DebugAPI) GetRawReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]hexutil.Bytes, error) {
	var hash common.Hash
	if h, ok := blockNrOrHash.Hash(); ok {
		hash = h
	} else {
		_, yuBlock, err := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return nil, err
		}
		hash = common.Hash(yuBlock.Hash)
	}
	receipts, err := api.b.GetReceipts(ctx, hash)
	if err != nil {
		return nil, err
	}
	result := make([]hexutil.Bytes, len(receipts))
	for i, receipt := range receipts {
		b, err := receipt.MarshalBinary()
		if err != nil {
			return nil, err
		}
		result[i] = b
	}
	return result, nil
}

// GetRawTransaction returns the bytes of the transaction for the given hash.
func (api *DebugAPI) GetRawTransaction(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	// Retrieve a finalized transaction, or a pooled otherwise
	found, tx, _, _, _, err := api.b.GetTransaction(ctx, hash)
	if !found {
		// In this case, it indicates that GetReceipt returned a specific error, not just that it wasn't found.
		if err != nil && !errors.Is(err, evm.ErrNotFoundReceipt) {
			return nil, err
		}
		tx, err = api.b.GetPoolTransaction(hash)
		if err != nil {
			return nil, err
		}
		if tx != nil {
			return tx.MarshalBinary()
		}
		if err != nil {
			logrus.Warnf("[GetRawTransaction] Failed to get tx %s, err: %v", hash.String(), err)
		}
		return nil, nil
		// Only when err==nil should it return nil, and when err!=nil, it should return NewTxIndexingError.
		// However, this would cause errors with packages such as Hardhat. In order to handle compatibility,
		// it currently returns nil directly.
		// return nil, NewTxIndexingError()
	}
	return tx.MarshalBinary()
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *DebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return spew.Sdump(block), nil
}

// ChaindbProperty returns leveldb properties of the key-value database.
func (api *DebugAPI) ChaindbProperty(property string) (string, error) {
	return api.b.ChainDb().Stat(property)
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *DebugAPI) ChaindbCompact() error {
	cstart := time.Now()
	for b := 0; b <= 255; b++ {
		var (
			start = []byte{byte(b)}
			end   = []byte{byte(b + 1)}
		)
		if b == 255 {
			end = nil
		}
		logrus.Info("Compacting database", "range", fmt.Sprintf("%#X-%#X", start, end), "elapsed", common.PrettyDuration(time.Since(cstart)))
		if err := api.b.ChainDb().Compact(start, end); err != nil {
			logrus.Error("Database compaction failed", "err", err)
			return err
		}
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *DebugAPI) SetHead(number hexutil.Uint64) {
	api.b.SetHead(uint64(number))
}

// NetAPI offers network related RPC methods
type NetAPI struct {
	net            *p2p.Server
	networkVersion uint64
}

// NewNetAPI creates a new net API instance.
func NewNetAPI(net *p2p.Server, networkVersion uint64) *NetAPI {
	return &NetAPI{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *NetAPI) Listening() bool {
	return false // always listening
}

// PeerCount returns the number of connected peers
func (s *NetAPI) PeerCount() hexutil.Uint {
	return hexutil.Uint(0)
}

// Version returns the current ethereum protocol version.
func (s *NetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.Ether)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f ether) exceeds the configured cap (%.2f ether)", feeFloat, cap)
	}
	return nil
}

func NewWeb3API(b Backend) *Web3API {
	return &Web3API{b: b}
}

type Web3API struct {
	b Backend
}

func (s *Web3API) ClientVersion() string {
	clientVersion := "Reddio/devnet-v2.2.2-45f4681/linux-amd64/go1.23.3"
	return clientVersion
}
