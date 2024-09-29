package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/contract"
	"github.com/sirupsen/logrus"
	yucommon "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/protocol"
	yutypes "github.com/yu-org/yu/core/types"
)

// just for test
const privateKey = ""
const SolidityTripod = "solidity"

type EventsWatcher struct {
	l1Watcher *EthSubscriber
	l2Watcher *ReddioSubscriber
}

func NewEventsWatcher(chain *kernel.Kernel, cfg *evm.GethConfig) (*EventsWatcher, error) {

	l1Watcher, err := NewEthSubscriber(cfg.L1ClientAddress, common.HexToAddress(cfg.ParentLayerContractAddress))
	if err != nil {
		return nil, err
	}
	l2Watcher, err := NewReddioSubscriber(chain, cfg)
	if err != nil {
		return nil, err
	}
	return &EventsWatcher{
		l1Watcher: l1Watcher,
		l2Watcher: l2Watcher,
	}, nil
}

func StartupEventsWatcher(chain *kernel.Kernel, cfg *evm.GethConfig) {
	if cfg.EnableEventsWatcher {
		eventsWatcher, err := NewEventsWatcher(chain, cfg)
		if err != nil {
			logrus.Fatal("init L1 client failed: ", err)
		}
		err = eventsWatcher.Run(cfg, context.Background())
		if err != nil {
			logrus.Fatal("l1 client run failed: ", err)
		}

	}
}

func (w *EventsWatcher) Run(cfg *evm.GethConfig, ctx context.Context) error {
	downwardMsgChan := make(chan *contract.ParentBridgeCoreFacetDownwardMessage)
	upwardMsgChan := make(chan *contract.ChildBridgeCoreFacetUpwardMessage)

	// Monitor L1 chain
	if w.l1Watcher.ethClient.Client().SupportsSubscriptions() {
		sub, err := w.l1Watcher.WatchDownwardMessage(ctx, downwardMsgChan, nil)
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case msg := <-downwardMsgChan:
					fmt.Println("Listen for msgChan", msg)
					jsonData, err := json.Marshal(msg)
					if err != nil {
						logrus.Errorf("Error converting downwardMsgChan txn to JSON: %v", err)
						continue
					}
					fmt.Println("msg as JSON:", string(jsonData))
					fmt.Println("handleDownwardMessage")
					w.handleDownwardMessage(msg, cfg.ChildLayerContractAddress)
					fmt.Println("handleDownwardMessage end")
				case subErr := <-sub.Err():
					logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
					sub.Unsubscribe()

					sub, err = w.l1Watcher.WatchDownwardMessage(ctx, downwardMsgChan, nil)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
					}
				case <-ctx.Done():
					sub.Unsubscribe()
					return
				}
			}
		}()
	}
	// Monitor L2 chain
	if w.l2Watcher.ethClient.Client().SupportsSubscriptions() {
		sub, err := w.l2Watcher.WatchUpwardMessageWss(ctx, upwardMsgChan, nil)
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case msg := <-upwardMsgChan:
					fmt.Println("Listen for msgChan", msg)
					jsonData, err := json.Marshal(msg)
					if err != nil {
						logrus.Errorf("Error converting upwardMsgChan txn to JSON: %v", err)
						continue
					}
					fmt.Println("msg as JSON:", string(jsonData))
					fmt.Println("handleupwardMessage")
					fmt.Println("handleUpwardMessage end")
				case subErr := <-sub.Err():
					logrus.Errorf("L1 subscription failed: %v, Resubscribing...", subErr)
					sub.Unsubscribe()

					sub, err = w.l2Watcher.WatchUpwardMessageWss(ctx, upwardMsgChan, nil)
					if err != nil {
						logrus.Errorf("Resubscribe failed: %v", err)
					}
				case <-ctx.Done():
					sub.Unsubscribe()
					return
				}
			}
		}()
	} else {
		err := w.l2Watcher.WatchUpwardMessageHttp(ctx, upwardMsgChan, nil)
		if err != nil {
			return err
		}

		go func() {
			for {
				select {
				case msg := <-upwardMsgChan:
					fmt.Println("Listen for msgChan", msg)
					jsonData, err := json.Marshal(msg)
					if err != nil {
						logrus.Errorf("Error converting upwardMsgChan txn to JSON: %v", err)
						continue
					}
					fmt.Println("msg as JSON:", string(jsonData))
					fmt.Println("handleupwardMessage")
					fmt.Println("handleUpwardMessage end")
				case <-ctx.Done():
					fmt.Println("Context done, stopping event processing")
					return
				}
			}
		}()
	}

	return nil
}

// handleDownwardMessage
func (w *EventsWatcher) handleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage, childLayerContractAddress string) error {

	//fmt.Println("Handle downward message", msg)
	client := w.l2Watcher.ethClient
	chainId, err := client.ChainID(context.Background())
	fmt.Println("ChainID", chainId)

	fmt.Println("childLayerContractAddress", childLayerContractAddress)
	downwardMessageDispatcher, err := contract.NewDownwardMessageDispatcherFacet(common.HexToAddress(childLayerContractAddress), client)
	if err != nil {
		return err
	}
	fmt.Println("downwardMessageDispatcher", downwardMessageDispatcher)

	testUserPK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(testUserPK, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	downwardMessages := []contract.DownwardMessageDispatcherFacetDownwardMessage{
		{
			Sequence:    msg.Sequence,
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
		},
	}
	log.Printf("Sending downward messages: %v", downwardMessages)
	nonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(3e6)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	contractABI, err := abi.JSON(strings.NewReader(contract.DownwardMessageDispatcherFacetABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	data, err := contractABI.Pack("receiveDownwardMessages", downwardMessages)

	tx := types.NewTransaction(nonce, common.HexToAddress(childLayerContractAddress), value, gasLimit, gasPrice, data)
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return err
	}
	fmt.Println("signedTx: ", signedTx)
	txJSON, err := json.MarshalIndent(struct {
		Nonce    uint64   `json:"nonce"`
		To       string   `json:"to"`
		Value    *big.Int `json:"value"`
		GasLimit uint64   `json:"gasLimit"`
		GasPrice *big.Int `json:"gasPrice"`
		Data     string   `json:"data"`
	}{
		Nonce:    tx.Nonce(),
		To:       tx.To().Hex(),
		Value:    tx.Value(),
		GasLimit: tx.Gas(),
		GasPrice: tx.GasPrice(),
		Data:     common.Bytes2Hex(tx.Data()),
	}, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal transaction to JSON: %v", err)
	}
	fmt.Println(string(txJSON))
	w.systemCall(context.Background(), signedTx)
	// tx, err := downwardMessageDispatcher.ReceiveDownwardMessages(auth, downwardMessages)
	// if err != nil {
	// 	log.Printf("Failed to send transaction: %v", err)
	// 	return err
	// }

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	return nil
}

// handleDownwardMessage
func (w *EventsWatcher) handleDownwardMessage(msg *contract.ParentBridgeCoreFacetDownwardMessage, childLayerContractAddress string) error {

	//fmt.Println("Handle downward message", msg)
	client := w.l2Watcher.ethClient
	chainId, err := client.ChainID(context.Background())
	fmt.Println("ChainID", chainId)

	fmt.Println("childLayerContractAddress", childLayerContractAddress)
	downwardMessageDispatcher, err := contract.NewDownwardMessageDispatcherFacet(common.HexToAddress(childLayerContractAddress), client)
	if err != nil {
		return err
	}
	testUserPK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(testUserPK, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	downwardMessages := []contract.DownwardMessageDispatcherFacetDownwardMessage{
		{
			Sequence:    msg.Sequence,
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
		},
	}
	log.Printf("Sending downward messages: %v", downwardMessages)

	tx, err := downwardMessageDispatcher.ReceiveDownwardMessages(auth, downwardMessages)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
		return err
	}

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	return nil
}

func yuHeader2EthHeader(yuHeader *yutypes.Header) *types.Header {
	return &types.Header{
		ParentHash:  common.Hash(yuHeader.PrevHash),
		Coinbase:    common.Address{}, // FIXME
		Root:        common.Hash(yuHeader.StateRoot),
		TxHash:      common.Hash(yuHeader.TxnRoot),
		ReceiptHash: common.Hash(yuHeader.ReceiptRoot),
		Difficulty:  new(big.Int).SetUint64(yuHeader.Difficulty),
		Number:      new(big.Int).SetUint64(uint64(yuHeader.Height)),
		GasLimit:    yuHeader.LeiLimit,
		GasUsed:     yuHeader.LeiUsed,
		Time:        yuHeader.Timestamp,
		Extra:       yuHeader.Extra,
		Nonce:       types.BlockNonce{},
		BaseFee:     big.NewInt(params.InitialBaseFee),
	}
}

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

func NewTxArgsFromTx(tx *types.Transaction) *TransactionArgs {
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
func (w *EventsWatcher) systemCall(ctx context.Context, signedTx *types.Transaction) error {
	// Check if this tx has been created
	// Create Tx
	yuBlock, err := w.l2Watcher.chain.Chain.GetEndCompactBlock()
	if err != nil {
		log.Fatalf("EthAPIBackend.CurrentBlock() failed: ", err)
		return nil
	}
	head := yuHeader2EthHeader(yuBlock.Header)
	signer := types.MakeSigner(w.l2Watcher.chainConfig, head.Number, head.Time)

	// Get v, r, s values
	v, r, s := signedTx.RawSignatureValues()
	fmt.Printf("v: %s\n", v.String())
	fmt.Printf("r: %s\n", r.String())
	fmt.Printf("s: %s\n", s.String())
	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		log.Fatalf("types.Sender() failed: ", err)
		return err
	}
	v, r, s = signedTx.RawSignatureValues()
	txArg := NewTxArgsFromTx(signedTx)
	txArgByte, _ := json.Marshal(txArg)
	txReq := &evm.TxRequest{
		Input:    signedTx.Data(),
		Origin:   sender,
		Address:  signedTx.To(),
		GasLimit: signedTx.Gas(),
		GasPrice: signedTx.GasPrice(),
		Value:    signedTx.Value(),
		Hash:     signedTx.Hash(),
		Nonce:    signedTx.Nonce(),
		V:        v,
		R:        r,
		S:        s,

		OriginArgs: txArgByte,
	}
	jsonData, err := json.MarshalIndent(txReq, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal txReq to JSON: %v", err)
	}

	fmt.Println("jsonData:", string(jsonData))
	byt, err := json.Marshal(txReq)
	if err != nil {
		log.Fatalf("json.Marshal(txReq) failed: ", err)
		return err
	}
	signedWrCall := &protocol.SignedWrCall{
		Call: &yucommon.WrCall{
			TripodName: SolidityTripod,
			FuncName:   "ExecuteTxn",
			Params:     string(byt),
		},
	}
	fmt.Println("signedWrCall", signedWrCall)
	fmt.Println("signedWrCall", signedWrCall.Call.Params)

	err = w.l2Watcher.chain.HandleTxn(signedWrCall)
	if err != nil {
		log.Fatalf("HandleTxn() failed: ", err)
		return err
	}
	return nil

}
