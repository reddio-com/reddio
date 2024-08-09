package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type CallRequest struct {
	Input    []byte         `json:"input"`
	Address  common.Address `json:"address"`
	Origin   common.Address `json:"origin"`
	GasLimit uint64         `json:"gasLimit"`
	GasPrice *big.Int       `json:"gasPrice"`
	Value    *big.Int       `json:"value"`
}

type CallResponse struct {
	Ret         []byte `json:"ret"`
	LeftOverGas uint64 `json:"leftOverGas"`
	Err         error  `json:"err"`
}

type TxRequest struct {
	Input    []byte         `json:"input"`
	Address  common.Address `json:"address"`
	Origin   common.Address `json:"origin"`
	GasLimit uint64         `json:"gasLimit"`
	GasPrice *big.Int       `json:"gasPrice"`
	Value    *big.Int       `json:"value"`
	Hash     common.Hash    `json:"hash"`
	V        *big.Int       `json:"v"`
	R        *big.Int       `json:"r"`
	S        *big.Int       `json:"s"`
}

type CreateRequest struct {
	Input    []byte         `json:"input"`
	Origin   common.Address `json:"origin"`
	GasLimit uint64         `json:"gasLimit"`
	GasPrice *big.Int       `json:"gasPrice"`
	Value    *big.Int       `json:"value"`
}
