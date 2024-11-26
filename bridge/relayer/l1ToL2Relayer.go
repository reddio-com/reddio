package relayer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
	yucommon "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/protocol"
)

type L1ToL2Relayer struct {
	ctx           context.Context
	cfg           *evm.GethConfig
	l1Client      *ethclient.Client
	l2Client      *ethclient.Client
	chain         *kernel.Kernel
	l1EventParser *logic.L1EventParser
}

const (
	maxRetries              = 10
	waitForConfirmationTime = 5 * time.Second
)

// TransactionArgs represents the arguments to construct a new transaction
// or a message call.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	// Issue detail: https://github.com/ethereum/go-ethereum/issues/15628
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`

	// For BlobTxType
	BlobFeeCap *hexutil.Big  `json:"maxFeePerBlobGas"`
	BlobHashes []common.Hash `json:"blobVersionedHashes,omitempty"`

	// For BlobTxType transactions with blob sidecar
	Blobs       []kzg4844.Blob       `json:"blobs"`
	Commitments []kzg4844.Commitment `json:"commitments"`
	Proofs      []kzg4844.Proof      `json:"proofs"`

	// This configures whether blobs are allowed to be passed.
	blobSidecarAllowed bool
}

func NewL1ToL2Relayer(ctx context.Context, cfg *evm.GethConfig, l1Client *ethclient.Client, l2Client *ethclient.Client, chain *kernel.Kernel) (*L1ToL2Relayer, error) {
	l1EventParser := logic.NewL1EventParser(cfg, l2Client)
	return &L1ToL2Relayer{
		ctx:           ctx,
		cfg:           cfg,
		l1Client:      l1Client,
		l2Client:      l2Client,
		chain:         chain,
		l1EventParser: l1EventParser,
	}, nil
}

// handleDownwardMessage
func (b *L1ToL2Relayer) HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error {
	// 1. parse downward message
	// 2. setup auth
	// 3. send downward message to child layer contract by calling downwardMessageDispatcher.ReceiveDownwardMessages

	downwardMessages := []contract.DownwardMessage{
		{
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
			Nonce:       utils.GenerateNonce(),
		},
	}
	metrics.DownwardMessageReceivedCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
	log.Printf("Sending downward messages: %v", downwardMessages)
	jsonData, err := json.MarshalIndent(downwardMessages, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal downward messages: %v", err)
	}

	fmt.Printf("Downward messages in JSON format:\n%s\n", string(jsonData))
	txNonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(6e6)
	gasPrice, err := b.l2Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(contract.DownwardMessageDispatcherFacetABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	data, err := contractABI.Pack("receiveDownwardMessages", downwardMessages)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}
	// fmt.Printf("Packed data: %s\n", hex.EncodeToString(data))

	tx := types.NewTransaction(txNonce, common.HexToAddress(b.cfg.ChildLayerContractAddress), value, gasLimit, gasPrice, data)

	crossMessages, err := b.l1EventParser.ParseL1CrossChainPayload(b.ctx, msg)
	if err != nil {
		log.Printf("Failed to parse L1 cross chain payload: %v", err)
	}
	fmt.Println("crossMessages L1TokenAddress", crossMessages[0].L1TokenAddress)
	fmt.Println("crossMessages sender", crossMessages[0].Sender)
	fmt.Println("crossMessages receiver", crossMessages[0].Receiver)
	fmt.Println("crossMessages tokenAmounts", crossMessages[0].TokenAmounts)
	fmt.Println("crossMessages tokenType", crossMessages[0].TokenType)

	err = b.systemCall(context.Background(), tx)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
		return err
	}
	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	success, err := waitForConfirmation(b.l2Client, tx.Hash())
	if err != nil {
		log.Printf("Failed to wait for confirmation: %v", err)
		b.refund(crossMessages)
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
	} else if !success {
		log.Printf("Transaction failed: %s", tx.Hash().Hex())
		b.refund(crossMessages)
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
	} else if success {
		metrics.DownwardMessageSuccessCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()

	}

	return nil
}

func newTxArgsFromTx(tx *types.Transaction) *TransactionArgs {
	args := TransactionArgs{}

	nonce := hexutil.Uint64(tx.Nonce())
	args.Nonce = &nonce

	gas := hexutil.Uint64(tx.Gas())
	args.Gas = &gas
	if tx.Value() != nil {
		value := hexutil.Big(*tx.Value())
		args.Value = &value
	}
	if tx.To() != nil {
		to := *tx.To()
		args.To = &to
	}
	if tx.Data() != nil {
		data := hexutil.Bytes(tx.Data())
		args.Data = &data
	}
	if tx.ChainId() != nil {
		chainId := hexutil.Big(*tx.ChainId())
		args.ChainID = &chainId
	}
	if tx.GasPrice() != nil {
		gasPrice := hexutil.Big(*tx.GasPrice())
		args.GasPrice = &gasPrice
	}

	if tx.Type() > 0 {
		if tx.AccessList() != nil {
			accessList := tx.AccessList()
			args.AccessList = &accessList
		}

		if tx.Type() > 1 {
			if tx.GasFeeCap() != nil {
				maxFeePerGas := hexutil.Big(*tx.GasFeeCap())
				args.MaxFeePerGas = &maxFeePerGas
			}
			if tx.GasTipCap() != nil {
				maxPriorityFeePerGas := hexutil.Big(*tx.GasTipCap())
				args.MaxPriorityFeePerGas = &maxPriorityFeePerGas
			}

			if tx.Type() == 3 {
				if tx.BlobHashes() != nil {
					args.BlobHashes = tx.BlobHashes()
				}
				if tx.BlobGasFeeCap() != nil {
					blobFeeCap := hexutil.Big(*tx.BlobGasFeeCap())
					args.BlobFeeCap = &blobFeeCap
				}
				if tx.BlobTxSidecar() != nil {
					args.Blobs = tx.BlobTxSidecar().Blobs
					args.Commitments = tx.BlobTxSidecar().Commitments
					args.Proofs = tx.BlobTxSidecar().Proofs
				}
			}
		}
	}

	return &args
}
func (b *L1ToL2Relayer) systemCall(ctx context.Context, signedTx *types.Transaction) error {
	// Check if this tx has been created
	// Create Tx
	//v, r, s := signedTx.RawSignatureValues()
	// fmt.Printf("v: %s\n", v.String())
	// fmt.Printf("r: %s\n", r.String())
	// fmt.Printf("s: %s\n", s.String())

	v, r, s := signedTx.RawSignatureValues()
	txArg := newTxArgsFromTx(signedTx)
	txArgByte, _ := json.Marshal(txArg)
	txReq := &evm.TxRequest{
		Input:          signedTx.Data(),
		Origin:         common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Address:        signedTx.To(),
		GasLimit:       signedTx.Gas(),
		GasPrice:       signedTx.GasPrice(),
		Value:          signedTx.Value(),
		Hash:           signedTx.Hash(),
		V:              v,
		R:              r,
		S:              s,
		IsInternalCall: true,

		OriginArgs: txArgByte,
	}

	txNonce, err := b.l2Client.PendingNonceAt(context.Background(), txReq.Origin)
	if err != nil {
		log.Printf("Failed to get nonce: %v", err)
	}
	txReq.Nonce = txNonce
	jsonData, err := json.MarshalIndent(txReq, "", "    ")
	if err != nil {
		log.Printf("Failed to marshal txReq to JSON: %v", err)
	}

	fmt.Println("systemCall jsonData:", string(jsonData))
	byt, err := json.Marshal(txReq)
	if err != nil {
		log.Printf("json.Marshal(txReq) failed: %v", err)
		return err
	}
	signedWrCall := &protocol.SignedWrCall{
		Call: &yucommon.WrCall{
			TripodName: "solidity",
			FuncName:   "ExecuteTxn",
			Params:     string(byt),
		},
	}
	// fmt.Println("signedWrCall", signedWrCall)
	// fmt.Println("signedWrCall", signedWrCall.Call.Params)

	err = b.chain.HandleTxn(signedWrCall)
	if err != nil {
		log.Printf("json.Marshal(txReq) failed: %v", err)
		return err
	}
	return nil

}

func waitForConfirmation(client *ethclient.Client, txHash common.Hash) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return true, nil
			}
			return false, fmt.Errorf("transaction failed with status: %v", receipt.Status)
		}
		time.Sleep(waitForConfirmationTime)
	}
	return false, fmt.Errorf("transaction was not confirmed after %d retries", maxRetries)
}

func (b *L1ToL2Relayer) refund(crossMessages []*logic.CrossMessage) error {
	fmt.Println("Start refund")

	for _, msg := range crossMessages {
		switch utils.MessagePayloadType(msg.TokenType) {
		case utils.ETH:
			err := b.refundETH(msg)
			if err != nil {
				return fmt.Errorf("failed to refund ETH: %v", err)
			}
		case utils.ERC20:
			err := b.refundERC20(msg)
			if err != nil {
				return fmt.Errorf("failed to refund ERC20: %v", err)
			}
		case utils.RED:
			err := b.refundRED(msg)
			if err != nil {
				return fmt.Errorf("failed to refund RED: %v", err)
			}
		default:
			return fmt.Errorf("unsupported token type: %d", msg.TokenType)
		}
	}

	return nil
}

func (b *L1ToL2Relayer) refundETH(msg *logic.CrossMessage) error {
	// TODO: implement
	return nil
}

func (b *L1ToL2Relayer) refundERC20(msg *logic.CrossMessage) error {
	// TODO: implement
	return nil
}

func (b *L1ToL2Relayer) refundRED(msg *logic.CrossMessage) error {
	// TODO: implement
	return nil
}
