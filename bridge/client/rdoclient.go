package rdoclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	*ethclient.Client // ethclient.Client
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}

func (rc *Client) HeaderByNumberNoType(ctx context.Context, number *big.Int) (*map[string]interface{}, error) {
	var head map[string]interface{}
	err := rc.Client.Client().CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = ethereum.NotFound
	}
	return &head, err
}

// func (rc *Client) BlockByNumberNoType(ctx context.Context, number *big.Int) (map[string]interface{}, error) {
// 	var block map[string]interface{}
// 	err := rc.Client.Client().CallContext(ctx, &block, "eth_getBlockByNumber", toBlockNumArg(number), true)
// 	if err == nil && block == nil {
// 		err = ethereum.NotFound
// 	}
// 	return block, err
// }

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (ec *Client) RdoBlockByNumber(ctx context.Context, number *big.Int) (*RdoBlock, error) {
	return ec.getRdoBlock(ctx, "eth_getBlockByNumber", toBlockNumArg(number), true)
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}
type rpcBlock struct {
	Hash         common.Hash         `json:"hash"`
	Transactions []rpcTransaction    `json:"transactions"`
	UncleHashes  []common.Hash       `json:"uncles"`
	Withdrawals  []*types.Withdrawal `json:"withdrawals,omitempty"`
}

func (rc *Client) BlockByHashNoType(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	var block map[string]interface{}
	err := rc.Client.Client().CallContext(ctx, &block, "eth_getBlockByHash", hash, true)
	if err == nil && block == nil {
		err = ethereum.NotFound
	}
	return block, err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return "0x" + number.Text(16)
}

// func (rc *Client) HeaderByNumberNoType(ctx context.Context, number *big.Int) (*map[string]interface{}, error) {

//		// Call eth_getBlockByNumber to get the block details
//		var head map[string]interface{}
//		err := rc.ethclient.c.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
//		if err == nil && head == nil {
//			err = ethereum.NotFound
//		}
//		return header, nil
//	}
func (rc *Client) getRdoBlock(ctx context.Context, method string, args ...interface{}) (*RdoBlock, error) {
	var raw json.RawMessage
	err := rc.Client.Client().CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, err
	}

	// Decode header and transactions.
	var head *types.Header
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, ethereum.NotFound
	}

	var body rpcBlock
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == types.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != types.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == types.EmptyTxsHash && len(body.Transactions) > 0 {
		return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != types.EmptyTxsHash && len(body.Transactions) == 0 {
		return nil, errors.New("server returned empty transaction list but block header indicates transactions")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*types.Header
	if len(body.UncleHashes) > 0 {
		uncles = make([]*types.Header, len(body.UncleHashes))
		reqs := make([]rpc.BatchElem, len(body.UncleHashes))
		for i := range reqs {
			reqs[i] = rpc.BatchElem{
				Method: "eth_getUncleByBlockHashAndIndex",
				Args:   []interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
				Result: &uncles[i],
			}
		}
		if err := rc.Client.Client().BatchCallContext(ctx, reqs); err != nil {
			return nil, err
		}
		for i := range reqs {
			if reqs[i].Error != nil {
				return nil, reqs[i].Error
			}
			if uncles[i] == nil {
				return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
			}
		}
	}
	// Fill the sender cache of transactions in the block.
	txs := make([]*types.Transaction, len(body.Transactions))
	block := types.NewBlockWithHeader(head).WithBody(txs, uncles).WithWithdrawals(body.Withdrawals)
	rdoBlock := &RdoBlock{Block: block}
	rdoBlock.SetHash(body.Hash)

	return rdoBlock, nil
}

type RdoBlock struct {
	*types.Block
	hash common.Hash
}

func (rb *RdoBlock) SetHash(hash common.Hash) {
	rb.hash = hash
}

func (rb *RdoBlock) Hash() common.Hash {
	if rb.hash != (common.Hash{}) {
		return rb.hash
	}
	return rb.Block.Hash()
}
