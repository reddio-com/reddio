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
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
)

type L2ToL1Relayer struct {
	ctx             context.Context
	cfg             *evm.GethConfig
	l1Client        *ethclient.Client
	Solidity        *evm.Solidity `tripod:"solidity"`
	crossMessageOrm *orm.CrossMessage
	sigPrivateKeys  []string
}

func NewL2ToL1Relayer(ctx context.Context, cfg *evm.GethConfig, l1Client *ethclient.Client, db *gorm.DB) (*L2ToL1Relayer, error) {

	privateKey, err := LoadPrivateKey("bridge/relayer/.sepolia.env")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}
	privateKeys := []string{
		privateKey,
	}

	return &L2ToL1Relayer{
		ctx:             ctx,
		cfg:             cfg,
		l1Client:        l1Client,
		crossMessageOrm: orm.NewCrossMessage(db),
		sigPrivateKeys:  privateKeys,
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

// HandleUpwardMessage handle L2 Upward Message
func (b *L2ToL1Relayer) HandleUpwardMessage(msgs []*orm.CrossMessage, blockTimestampsMap map[uint64]uint64) error {
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
	for _, msg := range msgs {
		var upwardMessages []contract.UpwardMessage
		payloadBytes, err := hex.DecodeString(msg.MessagePayload)
		if err != nil {
			logrus.Errorf("Error decoding payload: %v", err)
			return err
		}
		nonce := new(big.Int)
		nonce, ok := nonce.SetString(msg.MessageNonce, 10)
		if !ok {
			log.Fatalf("Failed to convert MessageNonce to *big.Int: %s", msg.MessageNonce)
		}
		upwardMessages = append(upwardMessages, contract.UpwardMessage{
			PayloadType: uint32(msg.MessagePayloadType),
			Payload:     payloadBytes,
			Nonce:       nonce,
		})

		signaturesArray, err := generateUpwardMessageMultiSignatures(upwardMessages, b.sigPrivateKeys)
		if err != nil {
			logrus.Fatalf("Failed to generate multi-signatures: %v", err)
		}

		var multiSignProofs []string
		for _, sig := range signaturesArray {
			multiSignProofs = append(multiSignProofs, "0x"+hex.EncodeToString(sig))
		}

		msg.MultiSignProof = strings.Join(multiSignProofs, ",")
		msg.BlockTimestamp = blockTimestampsMap[msg.L2BlockNumber]
	}

	if msgs != nil {
		err := b.crossMessageOrm.InsertOrUpdateL2Messages(context.Background(), msgs)
		if err != nil {
			logrus.Errorf("Failed to insert or update L2 messages: %v", err)
		}
	}
	
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

	return signaturesArray, nil
}

func generateUpwardMessageToHash(upwardMessages []contract.UpwardMessage) (common.Hash, error) {
	arrayLength := big.NewInt(int64(len(upwardMessages)))
	initialOffset := big.NewInt(32)
	headerData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(initialOffset)
	if err != nil {
		logrus.Fatalf("Failed to pack initial offset: %v", err)
	}

	lengthData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(arrayLength)
	if err != nil {
		logrus.Fatalf("Failed to pack array length: %v", err)
	}

	tupleOffset := big.NewInt(32)
	tupleOffsetData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 256}},
	}.Pack(tupleOffset)
	if err != nil {
		logrus.Fatalf("Failed to pack tuple offset: %v", err)
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
			return common.Hash{}, err
		}
		data = append(data, packedData...)
	}

	dataHash := crypto.Keccak256Hash(data)
	return dataHash, nil
}
