package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
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
		logrus.Fatalf("PackOrder BindJson txRequest(%s): %v", tx.TxnHash.String(), err)
		return false
	}
	oldNonce, ok := to.packingNonces[req.Origin]
	if !ok {
		to.packingNonces[req.Origin] = to.state.GetNonce(req.Origin)
	}

	logrus.Infof("PackOrder Address(%s), current Nonce(%d), request tx Nonce(%d)", req.Address.String(), oldNonce, req.Nonce)

	if req.Nonce == oldNonce {
		to.packingNonces[req.Origin]++
		return true
	}
	return false
}
