package relayer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/watcher/contract"
	yucommon "github.com/yu-org/yu/common"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/protocol"
)

type BridgeRelayer struct {
	ctx      context.Context
	cfg      *evm.GethConfig
	l1Client *ethclient.Client
	l2Client *ethclient.Client

	processedEvents map[string]bool
	l2chain         *kernel.Kernel
}

func NewBridgeRelayer(ctx context.Context, cfg *evm.GethConfig, l1Client *ethclient.Client, l2Client *ethclient.Client, l2chain *kernel.Kernel) (*BridgeRelayer, error) {
	return &BridgeRelayer{
		ctx:             ctx,
		cfg:             cfg,
		l1Client:        l1Client,
		l2Client:        l2Client,
		processedEvents: make(map[string]bool),
		l2chain:         l2chain,
	}, nil
}
func LoadPrivateKey(envFilePath string) (string, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return "", err
	}

	privateKey := os.Getenv("RELAYER_PRIVATE_KEY")
	if privateKey == "" {
		return "", fmt.Errorf("RELAYER_PRIVATE_KEY not set in %s", envFilePath)
	}

	return privateKey, nil
}

// handleDownwardMessage
func (b *BridgeRelayer) HandleDownwardMessage(msg *contract.ParentBridgeCoreFacetDownwardMessage) error {

	// 1. parse downward message
	// 2. setup auth
	// 3. send downward message to child layer contract by calling downwardMessageDispatcher.ReceiveDownwardMessages
	privateKey, err := LoadPrivateKey("relayer/.sepolia.env")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}
	downwardMessageDispatcher, err := contract.NewDownwardMessageDispatcherFacet(common.HexToAddress(b.cfg.ChildLayerContractAddress), b.l2Client)
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

// handleDownwardMessage
func (b *BridgeRelayer) HandleDownwardMessageWithSystemCall(msg *contract.ParentBridgeCoreFacetDownwardMessage) error {
	// 1. parse downward message
	// 2. setup auth
	// 3. send downward message to child layer contract by calling downwardMessageDispatcher.ReceiveDownwardMessages

	// downwardMessageDispatcher, err := contract.NewDownwardMessageDispatcherFacet(common.HexToAddress(b.cfg.ChildLayerContractAddress), b.l2Client)
	// if err != nil {
	// 	return err
	// }
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

	tx := types.NewTransaction(nonce, common.HexToAddress(b.cfg.ChildLayerContractAddress), value, gasLimit, gasPrice, data)

	b.systemCall(context.Background(), tx)

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	return nil
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
func (b *BridgeRelayer) systemCall(ctx context.Context, signedTx *types.Transaction) error {
	// Check if this tx has been created
	// Create Tx
	v, r, s := signedTx.RawSignatureValues()
	// fmt.Printf("v: %s\n", v.String())
	// fmt.Printf("r: %s\n", r.String())
	// fmt.Printf("s: %s\n", s.String())
	v, r, s = signedTx.RawSignatureValues()
	txArg := NewTxArgsFromTx(signedTx)
	txArgByte, _ := json.Marshal(txArg)
	txReq := &evm.TxRequest{
		Input:          signedTx.Data(),
		Origin:         common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Address:        signedTx.To(),
		GasLimit:       signedTx.Gas(),
		GasPrice:       signedTx.GasPrice(),
		Value:          signedTx.Value(),
		Hash:           signedTx.Hash(),
		Nonce:          signedTx.Nonce(),
		V:              v,
		R:              r,
		S:              s,
		IsInternalCall: true,

		OriginArgs: txArgByte,
	}
	// jsonData, err := json.MarshalIndent(txReq, "", "    ")
	// if err != nil {
	// 	log.Fatalf("Failed to marshal txReq to JSON: %v", err)
	// }

	// fmt.Println("jsonData:", string(jsonData))
	byt, err := json.Marshal(txReq)
	if err != nil {
		log.Fatalf("json.Marshal(txReq) failed: ", err)
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

	err = b.l2chain.HandleTxn(signedWrCall)
	if err != nil {
		log.Fatalf("HandleTxn() failed: ", err)
		return err
	}
	return nil

}

// handleL2UpwardMessage
func (b *BridgeRelayer) HandleUpwardMessage(msg *contract.ChildBridgeCoreFacetUpwardMessage) error {
	// 1. parse upward message
	// 2. setup auth
	// 3. send upward message to parent layer contract by calling upwardMessageDispatcher.ReceiveUpwardMessages

	upwardMessageDispatcher, err := contract.NewUpwardMessageDispatcherFacet(common.HexToAddress(b.cfg.ParentLayerContractAddress), b.l1Client)
	if err != nil {
		return err
	}
	privateKey, err := LoadPrivateKey("relayer/.sepolia.env")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}
	testUserPK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(testUserPK, big.NewInt(11155111))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	//for test ,set  msg.Sequence manually///
	//msg.Sequence = big.NewInt(1)
	////////////////////////////////////////

	upwardMessages := []contract.UpwardMessage{
		{
			Sequence:    msg.Sequence,
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
		},
	}

	privateKeys := []string{
		privateKey,
	}

	signaturesArray, err := GenerateUpwardMessageMultiSignatures(upwardMessages, privateKeys)
	if err != nil {
		log.Fatalf("Failed to generate multi-signatures: %v", err)
	}

	for i, sig := range signaturesArray {
		log.Printf("MutiSignature %d: %x\n", i+1, sig)
	}

	tx, err := upwardMessageDispatcher.ReceiveUpwardMessages(auth, upwardMessages, signaturesArray)
	if err != nil {
		log.Printf("Failed to send transaction: %v", err)
		return err
	}

	log.Printf("Transaction sent: %s", tx.Hash().Hex())
	return nil
}

/**
 * GenerateUpwardMessageMultiSignatures generates multi-signatures for upward messages.
 * The signature hash generation process includes the message header to ensure the integrity and authenticity of the message.
 * The message header typically contains the following metadata:
 * - Initial offset: Points to the first element (array) offset, usually fixed at 32 bytes.
 * - Array length: The number of upward messages in the array.
 * - Tuple offset: Points to the offset of the tuple.
 *
 * Parameters:
 * - upwardMessages: A slice of UpwardMessage structs containing the messages to be signed.
 * - privateKeys: A slice of strings containing the private keys used for signing.
 *
 * Returns:
 * - A slice of byte slices containing the generated signatures.
 * - An error if the signature generation fails.
 */
func GenerateUpwardMessageMultiSignatures(upwardMessages []contract.UpwardMessage, privateKeys []string) ([][]byte, error) {

	arrayLength := big.NewInt(int64(len(upwardMessages)))
	initialOffset := big.NewInt(32)
	headerData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(initialOffset)
	if err != nil {
		log.Fatalf("Failed to pack initial offset: %v", err)
	}

	lengthData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(arrayLength)
	if err != nil {
		log.Fatalf("Failed to pack array length: %v", err)
	}

	tupleOffset := big.NewInt(32)
	tupleOffsetData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(tupleOffset)
	if err != nil {
		log.Fatalf("Failed to pack tuple offset: %v", err)
	}

	var data []byte
	data = append(data, headerData...)
	data = append(data, lengthData...)
	data = append(data, tupleOffsetData...)

	for _, msg := range upwardMessages {
		packedData, err := abi.Arguments{
			{Type: abi.Type{T: abi.UintTy, Size: 256}}, // Use UintTy with size 256 for *big.Int
			{Type: abi.Type{T: abi.UintTy, Size: 32}},  // Use UintTy with size 32 for uint32
			{Type: abi.Type{T: abi.BytesTy}},
		}.Pack(msg.Sequence, msg.PayloadType, msg.Payload)
		if err != nil {
			//fmt.Printf("Failed to pack upwardMessages: %v\n", err)
			return nil, err
		}
		data = append(data, packedData...)
	}

	//fmt.Printf("Encoded Data (Hex): %s\n", hex.EncodeToString(data))

	dataHash := crypto.Keccak256Hash(data)
	//fmt.Println("newdataHash:", dataHash)
	// Generate multiple signatures
	var signaturesArray [][]byte
	for _, pk := range privateKeys {
		privateKey, err := crypto.HexToECDSA(pk)
		if err != nil {
			return nil, err
		}

		signature, err := crypto.Sign(dataHash.Bytes(), privateKey)
		if err != nil {
			return nil, err
		}

		signaturesArray = append(signaturesArray, signature)
	}

	// for print
	// Recover the public key
	// sigPublicKey, err := crypto.Ecrecover(dataHash.Bytes(), signaturesArray[0])
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Convert public key to address
	// publicKeyECDSA, err := crypto.UnmarshalPubkey(sigPublicKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// fmt.Printf("Signed hash: %x\n", signaturesArray[0])
	// fmt.Printf("Signer address: %s\n", address.Hex())
	// Ensure v value is 27 or 28
	if signaturesArray[0][64] < 27 {
		//fmt.Println("v value is less than 27")
		signaturesArray[0][64] += 27
	}
	return signaturesArray, nil
}
