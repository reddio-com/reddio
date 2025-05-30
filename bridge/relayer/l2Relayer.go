package relayer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
)

type L2Relayer struct {
	ctx               context.Context
	cfg               *evm.GethConfig
	crossMessageOrm   *orm.CrossMessage
	rawBridgeEventOrm *orm.RawBridgeEvent
	l2EventParser     *logic.L2EventParser
	pollingSemaphore  chan struct{}
	sigPrivateKeys    []string
}

func NewL2Relayer(ctx context.Context, cfg *evm.GethConfig, db *gorm.DB) (*L2Relayer, error) {

	privateKeys, err := LoadPrivateKeyArray(cfg.MultisigEnvFile, cfg.MultisigEnvVar)
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	return &L2Relayer{
		ctx:               ctx,
		cfg:               cfg,
		crossMessageOrm:   orm.NewCrossMessage(db),
		rawBridgeEventOrm: orm.NewRawBridgeEvent(db, cfg),
		l2EventParser:     logic.NewL2EventParser(cfg),
		sigPrivateKeys:    privateKeys,
		pollingSemaphore:  make(chan struct{}, 1), // 1 means only one polling goroutine can run at a time
	}, nil
}
func LoadPrivateKey(envFilePath string, envVarName string) (string, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return "", err
	}

	privateKey := os.Getenv(envVarName)
	if privateKey == "" {
		return "", fmt.Errorf("PRIVATE_KEY not set in %s", envFilePath)
	}

	return privateKey, nil
}

func LoadPrivateKeyArray(envFilePath string, envVarName string) ([]string, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return nil, err
	}

	privateKeysStr := os.Getenv(envVarName)
	if privateKeysStr == "" {
		return nil, fmt.Errorf("%s not set in %s", envVarName, envFilePath)
	}

	privateKeys := strings.Split(privateKeysStr, ",")
	for i := range privateKeys {
		privateKeys[i] = strings.TrimSpace(privateKeys[i])
	}
	return privateKeys, nil
}

// HandleUpwardMessage handle L2 Upward Message
func (b *L2Relayer) HandleUpwardMessage(ctx context.Context, bridgeEvent *orm.RawBridgeEvent) error {
	// 1. parse upward message
	// 2. setup auth
	// 3. send upward message to parent layer contract by calling upwardMessageDispatcher.ReceiveUpwardMessages
	msgs, err := b.l2EventParser.ParseL2RawBridgeEventToCrossChainMessage(ctx, bridgeEvent)
	if err != nil {
		logrus.Errorf("Failed to parse L2 raw bridge event to cross chain message: %v", err)
	}

	for _, msg := range msgs {
		metrics.UpwardMessageReceivedCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()

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

	}

	if msgs != nil {
		err := b.crossMessageOrm.InsertOrUpdateCrossMessages(context.Background(), msgs)
		if err != nil {
			logrus.Errorf("Failed to insert or update L2 messages: %v", err)
		}
		err = b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L2_RawBridgeEventsTableName, bridgeEvent.ID, int(btypes.Processed))
		if err != nil {
			logrus.Errorf("Failed to update process status of raw bridge events: %v", err)
		}

	}

	return nil
}

func (b *L2Relayer) StartPolling() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case b.pollingSemaphore <- struct{}{}:
				go func() {
					defer func() { <-b.pollingSemaphore }()
					b.pollUnProcessedMessages()
				}()
			default:
				// skip this round if semaphore is full
			}
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *L2Relayer) pollUnProcessedMessages() {
	ctx := context.Background()
	//messages, err := r.crossMessageOrm.QueryL1UnConsumedMessages(ctx, btypes.TxTypeDeposit)
	bridgeEvents, err := b.rawBridgeEventOrm.QueryUnProcessedBridgeEvents(ctx, b.cfg.L2_RawBridgeEventsTableName, b.cfg.RelayerBatchSize)
	if err != nil {
		log.Printf("Failed to query unconsumed messages: %v", err)
		return
	}
	//1.proceeding the L1 unprocessed  messages
	// QueueTransaction
	//1.1 generate cross message
	//1.2 update the status of the raw bridge events to processed
	//2.check L2 message if it is consumed
	//2.1 if it is consumed, update the status of the L1 message to consumed
	for _, bridgeEvent := range bridgeEvents {
		if bridgeEvent.EventType == int(btypes.SentMessage) {
			//1.1 generate cross message

			b.HandleUpwardMessage(ctx, bridgeEvent)
		} else if bridgeEvent.EventType == int(btypes.L2RelayedMessage) {

			b.HandleL2RelayerMessage(ctx, bridgeEvent)

		}
	}
}
func (b *L2Relayer) HandleL2RelayerMessage(ctx context.Context, bridgeEvent *orm.RawBridgeEvent) error {
	//fmt.Println("HandleL2RelayerMessage")
	relayedMessage, err := b.l2EventParser.ParseL2RelayMessagePayload(ctx, bridgeEvent)
	//fmt.Println("relayedMessages:", relayedMessage.MessageHash)
	if err != nil {
		logrus.Infof("Failed to parse L1 cross chain payload: %v", err)
		b.rawBridgeEventOrm.UpdateProcessFail(b.cfg.L2_RawBridgeEventsTableName, bridgeEvent.ID, err.Error())
		return err
	}
	rowsAffected, err := b.crossMessageOrm.UpdateL1MessageConsumedStatus(b.ctx, relayedMessage)
	//fmt.Println("UpdateL1MessageConsumedStatus")
	if err != nil {
		logrus.Infof("Failed to update L2 message consumed status: %v", err)
		b.rawBridgeEventOrm.UpdateProcessFail(b.cfg.L2_RawBridgeEventsTableName, bridgeEvent.ID, err.Error())
		return err
	}
	if rowsAffected == 0 {
		logrus.Warn("L1 message cant be found: ", relayedMessage.MessageHash)
		return nil
	}

	b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L2_RawBridgeEventsTableName, bridgeEvent.ID, int(btypes.Processed))
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
