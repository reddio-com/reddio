package relayer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/bridge/utils/database"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
	"gorm.io/gorm"
)

type L2ToL1Relayer struct {
	ctx             context.Context
	cfg             *evm.GethConfig
	l1Client        *ethclient.Client
	Solidity        *evm.Solidity `tripod:"solidity"`
	db              *gorm.DB
	crossMessageOrm *orm.CrossMessage
}

func NewL2ToL1Relayer(ctx context.Context, cfg *evm.GethConfig, l1Client *ethclient.Client) (*L2ToL1Relayer, error) {
	fmt.Println("cfg.BridgeDBConfig: ", cfg.BridgeDBConfig)
	db, err := database.InitDB(cfg.BridgeDBConfig)
	if err != nil {
		log.Fatal("failed to init db", "err", err)
	}
	//fmt.Println("db init success")
	// defer func() {
	// 	fmt.Println("defer close db")
	// 	if deferErr := database.CloseDB(db); deferErr != nil {
	// 		log.Fatal("failed to close db", "err", err)
	// 	}
	// }()
	return &L2ToL1Relayer{
		ctx:             ctx,
		cfg:             cfg,
		l1Client:        l1Client,
		db:              db,
		crossMessageOrm: orm.NewCrossMessage(db),
	}, nil
}
func LoadPrivateKey(envFilePath string) (string, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return "", err
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		return "", fmt.Errorf("PRIVATE_KEY not set in %s", envFilePath)
	}

	return privateKey, nil
}

// handleL2UpwardMessage
func (b *L2ToL1Relayer) HandleUpwardMessage(msgs []*orm.CrossMessage) error {
	// 1. parse upward message
	// 2. setup auth
	// 3. send upward message to parent layer contract by calling upwardMessageDispatcher.ReceiveUpwardMessages
	for _, msg := range msgs {
		metrics.UpwardMessageReceivedCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
	}
	// upwardMessageDispatcher, err := contract.NewUpwardMessageDispatcherFacet(common.HexToAddress(b.cfg.ParentLayerContractAddress), b.l1Client)
	// if err != nil {
	// 	return err
	// }
	privateKey, err := LoadPrivateKey("bridge/relayer/.sepolia.env")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}
	// testUserPK, err := crypto.HexToECDSA(privateKey)
	// if err != nil {
	// 	return err
	// }
	// l1ChainId, err := b.l1Client.ChainID(context.Background())
	// if err != nil {
	// 	return err
	// }

	// auth, err := bind.NewKeyedTransactorWithChainID(testUserPK, l1ChainId)
	// if err != nil {
	// 	log.Fatalf("Failed to create authorized transactor: %v", err)
	// }

	privateKeys := []string{
		privateKey,
	}
	for _, msg := range msgs {
		var upwardMessages []contract.UpwardMessage
		payloadBytes, err := hex.DecodeString(msg.MessagePayload)
		if err != nil {
			//fmt.Println("Failed to decode hex string:", err)
			return err
		}
		upwardMessages = append(upwardMessages, contract.UpwardMessage{
			PayloadType: uint32(msg.MessagePayloadType),
			Payload:     payloadBytes,
			Nonce:       utils.GenerateNonce(),
		})
		signaturesArray, err := generateUpwardMessageMultiSignatures(upwardMessages, privateKeys)
		if err != nil {
			log.Fatalf("Failed to generate multi-signatures: %v", err)
		}

		msg.MessageHash = utils.ComputeMessageHash(upwardMessages[0].PayloadType, upwardMessages[0].Payload, upwardMessages[0].Nonce).Hex()
		//fmt.Println("msg.MessageHash:", msg.MessageHash)
		msg.MessageNonce = upwardMessages[0].Nonce.String()
		var multiSignProofs []string
		for _, sig := range signaturesArray {
			multiSignProofs = append(multiSignProofs, "0x"+hex.EncodeToString(sig))
		}

		msg.MultiSignProof = strings.Join(multiSignProofs, ",")
		//fmt.Println("msg.MultiSignProof:", msg.MultiSignProof)
	}

	if msgs != nil {
		//fmt.Println("msgs:", msgs)
		err = b.crossMessageOrm.InsertOrUpdateL2Messages(context.Background(), msgs)
		if err != nil {
			fmt.Println("Failed to insert or update L2 messages:", err)
		}
	}
	// upwardMessagesJSON, err := json.MarshalIndent(upwardMessages, "", "  ")
	// if err != nil {
	// 	fmt.Printf("Error marshalling upwardMessages to JSON: %v\n", err)
	// 	return err
	// }

	// // Print JSON
	// fmt.Printf("UpwardMessages JSON:\n%s\n", string(upwardMessagesJSON))

	// signaturesArray, err := generateUpwardMessageMultiSignatures(upwardMessages, privateKeys)
	// if err != nil {
	// 	log.Fatalf("Failed to generate multi-signatures: %v", err)
	// }

	// for i, sig := range signaturesArray {
	// 	log.Printf("MultiSignature %d: %x\n", i+1, sig)
	// }

	// tx, err := upwardMessageDispatcher.ReceiveUpwardMessages(auth, upwardMessages, signaturesArray)
	// if err != nil {
	// 	log.Printf("Failed to send transaction: %v", err)
	// 	for _, msg := range msgs {
	// 		metrics.UpwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
	// 	}
	// 	return err
	// }

	// log.Printf("Transaction sent: %s", tx.Hash().Hex())
	// for _, msg := range msgs {
	// 	metrics.UpwardMessageSuccessCounter.WithLabelValues(fmt.Sprintf("%d", msg.PayloadType)).Inc()
	// }
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
func generateUpwardMessageMultiSignatures(upwardMessages []contract.UpwardMessage, privateKeys []string) ([][]byte, error) {

	dataHash, err := generateUpwardMessageToHash(upwardMessages)
	if err != nil {
		return nil, err
	}

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

	return signaturesArray, nil
}

func generateUpwardMessageToHash(upwardMessages []contract.UpwardMessage) (common.Hash, error) {
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
			{Type: abi.Type{T: abi.UintTy, Size: 32}}, // Use UintTy with size 32 for uint32
			{Type: abi.Type{T: abi.BytesTy}},
			{Type: abi.Type{T: abi.UintTy, Size: 256}}, // Use UintTy with size 256 for *big.Int
		}.Pack(msg.PayloadType, msg.Payload, msg.Nonce)
		if err != nil {
			//fmt.Printf("Failed to pack upwardMessages: %v\n", err)
			return common.Hash{}, err
		}
		data = append(data, packedData...)
	}

	//fmt.Printf("Encoded Data (Hex): %s\n", hex.EncodeToString(data))

	dataHash := crypto.Keccak256Hash(data)
	return dataHash, nil
}
