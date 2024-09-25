package parallel

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	yucommon "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/types"
	"math/big"
	"time"
)

// SimulateBlock simulates a block for test.
func SimulateBlock(txs []*ethtypes.Transaction) (*types.Block, error) {
	var txns []*types.SignedTxn
	for _, tx := range txs {
		txn, err := ConvertTx(tx)
		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}
	return &types.Block{
		Header: &types.Header{
			Height: 10,
			Hash:   yucommon.Hash(ethtypes.EmptyRootHash),
		},
		Txns: txns,
	}, nil
}

// ConvertTx converts eth tx to SignedTx.
func ConvertTx(tx *ethtypes.Transaction) (*types.SignedTxn, error) {
	signer := ethtypes.MakeSigner(params.AllEthashProtocolChanges, new(big.Int).SetUint64(0), uint64(time.Now().Unix()))
	sender, err := ethtypes.Sender(signer, tx)
	if err != nil {
		return nil, err
	}
	v, r, s := tx.RawSignatureValues()
	txArg := ethrpc.NewTxArgsFromTx(tx)
	txArgByte, _ := json.Marshal(txArg)
	txReq := &evm.TxRequest{
		Input:    tx.Data(),
		Origin:   sender,
		Address:  tx.To(),
		GasLimit: tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Hash:     tx.Hash(),
		Nonce:    tx.Nonce(),
		V:        v,
		R:        r,
		S:        s,

		OriginArgs: txArgByte,
	}
	byt, err := json.Marshal(txReq)
	if err != nil {
		return nil, err
	}
	wrCall := &yucommon.WrCall{
		TripodName: ethrpc.SolidityTripod,
		FuncName:   "ExecuteTxn",
		Params:     string(byt),
	}

	return types.NewSignedTxn(wrCall, nil, nil, nil)
}

func MakeEthTX(privateKeyHex string, toAddress string, amount uint64, data []byte, nonce uint64) (*ethtypes.Transaction, error) {
	to := common.HexToAddress(toAddress)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)

	tx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    big.NewInt(int64(amount)),
		Data:     data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	chainID := params.AllEthashProtocolChanges.ChainID
	return ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
}
