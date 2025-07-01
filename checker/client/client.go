package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
)

type NodeClient struct {
	host string
	http *http.Client
}

func NewNodeClient(host string) *NodeClient {
	return &NodeClient{
		host: host,
		http: &http.Client{},
	}
}

type rpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type rpcResponse struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (client *NodeClient) sendRequest(method string, params []interface{}) (json.RawMessage, error) {
	request := rpcRequest{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := client.http.Post(client.host, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var response rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("RPC error: %d - %s", response.Error.Code, response.Error.Message)
	}
	return response.Result, nil
}

type BlockTransaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    uint64 `json:"value"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
	Input    string `json:"input"`
	Nonce    uint64 `json:"nonce"`
}

func parseHexUint64(hexStr string) (uint64, error) {
	if hexStr == "" || hexStr == "0x" {
		return 0, nil
	}
	val := new(big.Int)
	if _, ok := val.SetString(hexStr[2:], 16); !ok {
		return 0, fmt.Errorf("invalid hex number: %s", hexStr)
	}
	if !val.IsUint64() {
		return 0, fmt.Errorf("value exceeds uint64 range: %s", hexStr)
	}
	return val.Uint64(), nil
}

func parseTransaction(data json.RawMessage) (*BlockTransaction, error) {
	var tx struct {
		Hash     string `json:"hash"`
		From     string `json:"from"`
		To       string `json:"to"`
		Value    string `json:"value"`
		Gas      string `json:"gas"`
		GasPrice string `json:"gasPrice"`
		Input    string `json:"input"`
		Nonce    string `json:"nonce"`
	}
	if err := json.Unmarshal(data, &tx); err != nil {
		return nil, err
	}

	value, err := parseHexUint64(tx.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse value: %v", err)
	}

	gas, err := parseHexUint64(tx.Gas)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gas: %v", err)
	}

	gasPrice, err := parseHexUint64(tx.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasPrice: %v", err)
	}

	nonce, err := parseHexUint64(tx.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nonce: %v", err)
	}

	return &BlockTransaction{
		Hash:     tx.Hash,
		From:     tx.From,
		To:       tx.To,
		Value:    value,
		Gas:      gas,
		GasPrice: gasPrice,
		Input:    tx.Input,
		Nonce:    nonce,
	}, nil
}

func (client *NodeClient) GetBlockByNumber(blockNumber uint64) ([]*BlockTransaction, error) {
	hexBlockNumber := "0x" + strconv.FormatUint(blockNumber, 16)
	result, err := client.sendRequest("eth_getBlockByNumber", []interface{}{hexBlockNumber, true})
	if err != nil {
		return nil, err
	}

	var block struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %v", err)
	}

	var txs []*BlockTransaction
	for _, txData := range block.Transactions {
		tx, err := parseTransaction(txData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse transaction: %v", err)
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func (client *NodeClient) GetLatestBlock() (uint64, error) {
	result, err := client.sendRequest("eth_blockNumber", nil)
	if err != nil {
		return 0, err
	}

	var hexBlockNumber string
	if err := json.Unmarshal(result, &hexBlockNumber); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %v", err)
	}

	blockNumber, err := strconv.ParseUint(hexBlockNumber[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %v", err)
	}

	return blockNumber, nil
}

func (client *NodeClient) GetBalanceByBlock(blockNumber uint64, address string) (uint64, error) {
	hexBlockNumber := "0x" + strconv.FormatUint(blockNumber, 16)
	result, err := client.sendRequest("eth_getBalance", []interface{}{address, hexBlockNumber})
	if err != nil {
		return 0, err
	}

	var hexBalance string
	if err := json.Unmarshal(result, &hexBalance); err != nil {
		return 0, fmt.Errorf("failed to unmarshal balance: %v", err)
	}

	balance := new(big.Int)
	balance.SetString(hexBalance[2:], 16)
	if !balance.IsUint64() {
		return 0, fmt.Errorf("balance is too large to fit in uint64")
	}

	return balance.Uint64(), nil
}
