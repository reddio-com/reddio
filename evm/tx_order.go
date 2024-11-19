package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/yu-org/yu/core/types"
)

type TxOrdered struct {
	packingNonces map[common.Address]uint64
	state         *EthState
}

func NewTxOrdered(state *EthState) *TxOrdered {
	return &TxOrdered{
		packingNonces: make(map[common.Address]uint64),
		state:         state,
	}
}

func (to *TxOrdered) PackOrder(tx *types.SignedTxn) bool {
	req := new(TxRequest)
	err := tx.BindJson(req)
	if err != nil {
		return false
	}
	nonce, ok := to.packingNonces[req.Origin]
	if !ok {
		to.packingNonces[req.Origin] = to.state.GetNonce(req.Origin)
	}

	if req.Nonce == nonce+1 {
		to.packingNonces[req.Origin]++
		return true
	}
	return false
}
