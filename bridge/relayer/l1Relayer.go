package relayer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/HyperService-Consortium/go-hexutil"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/kernel"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/logic"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/metrics"
)

type L1Relayer struct {
	ctx               context.Context
	cfg               *evm.GethConfig
	l2Client          *ethclient.Client
	chain             *kernel.Kernel
	l1EventParser     *logic.L1EventParser
	crossMessageOrm   *orm.CrossMessage
	rawBridgeEventOrm *orm.RawBridgeEvent
	pollingSemaphore  chan struct{}
}

const (
	MAX_RETRIES                = 10
	WAIT_FOR_CONFIRMATION_TIME = 5 * time.Second
	ZERO_ADDRESS               = "0x0000000000000000000000000000000000000000"
	ZERO_HASH                  = "0x0000000000000000000000000000000000000000000000000000000000000000"
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
	// "input" is the newer name and should be preferred by clients .
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

func NewL1Relayer(ctx context.Context, cfg *evm.GethConfig, l2Client *ethclient.Client, chain *kernel.Kernel, db *gorm.DB) (*L1Relayer, error) {
	l1EventParser := logic.NewL1EventParser(cfg)

	relayer := &L1Relayer{
		ctx:               ctx,
		cfg:               cfg,
		l2Client:          l2Client,
		chain:             chain,
		l1EventParser:     l1EventParser,
		crossMessageOrm:   orm.NewCrossMessage(db),
		rawBridgeEventOrm: orm.NewRawBridgeEvent(db, cfg),
		pollingSemaphore:  make(chan struct{}, 1), // 1 means only one polling goroutine can run at a time

	}

	return relayer, nil
}
func (b *L1Relayer) StartPolling() {
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
func (b *L1Relayer) pollUnProcessedMessages() {
	ctx := context.Background()
	//messages, err := r.crossMessageOrm.QueryL1UnConsumedMessages(ctx, btypes.TxTypeDeposit)
	bridgeEvents, err := b.rawBridgeEventOrm.QueryUnProcessedBridgeEvents(ctx, b.cfg.L1_RawBridgeEventsTableName, b.cfg.RelayerBatchSize)
	if err != nil {
		logrus.Error("Failed to query unconsumed messages: %v", err)
		return
	}
	//1.proceeding the L1 unprocessed  messages
	// QueueTransaction
	//1.1 generate cross message
	//1.2 update the status of the raw bridge events to processed
	//2.check L2 message if it is consumed
	//2.1 if it is consumed, update the status of the L1 message to consumed
	for _, bridgeEvent := range bridgeEvents {
		if bridgeEvent.EventType == int(btypes.QueueTransaction) {
			//1.1 generate cross message
			//b.HandleDownwardMessageWithSystemCall(bridgeEvent)
			b.HandleDownwardMessage(bridgeEvent)
		} else if bridgeEvent.EventType == int(btypes.L1RelayedMessage) {
			b.HandleL1RelayerMessage(bridgeEvent)

		}
	}
}

func (b *L1Relayer) HandleL1RelayerMessage(msg *orm.RawBridgeEvent) error {
	relayedMessage, err := b.l1EventParser.ParseL1RelayMessagePayload(b.ctx, msg)
	if err != nil {
		logrus.Infof("Failed to parse L1 cross chain payload: %v", err)
		b.rawBridgeEventOrm.UpdateProcessFail(b.cfg.L1_RawBridgeEventsTableName, msg.ID, err.Error())
		return err
	}
	rowsAffected, err := b.crossMessageOrm.UpdateL2MessageConsumedStatus(b.ctx, relayedMessage)
	if err != nil {
		logrus.Infof("Failed to update L2 message consumed status: %v", err)
		b.rawBridgeEventOrm.UpdateProcessFail(b.cfg.L1_RawBridgeEventsTableName, msg.ID, err.Error())
		return err
	}
	if rowsAffected == 0 {
		logrus.Warn("L2 message cant be found: ", relayedMessage.MessageHash)
		return nil
	}
	b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L1_RawBridgeEventsTableName, msg.ID, int(btypes.Processed))
	return nil
}

// HandleDownwardMessageWithSystemCall handles the downward message
func (b *L1Relayer) HandleDownwardMessageWithSystemCall(msg *orm.RawBridgeEvent) error {
	// 1. parse downward message
	// 2. setup auth
	// 3. send downward message to child layer contract by calling downwardMessageDispatcher.ReceiveDownwardMessages
	payloadBytes, err := hex.DecodeString(msg.MessagePayload)
	if err != nil {
		logrus.Errorf("Failed to decode hex string: %v", err)
		return err
	}
	downwardMessages := []contract.DownwardMessage{
		{
			PayloadType: uint32(msg.MessagePayloadType),
			Payload:     payloadBytes,
			Nonce:       big.NewInt(int64(msg.MessageNonce)),
		},
	}
	metrics.DownwardMessageReceivedCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()

	txNonce := uint64(0)
	value := big.NewInt(0)
	gasLimit := uint64(6e6)
	//gasPrice := big.NewInt(0)
	gasPrice, err := b.l2Client.SuggestGasPrice(context.Background())
	if err != nil {
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		logrus.Errorf("Failed to suggest gas price: %v", err)
		return err
	}

	contractABI, err := abi.JSON(strings.NewReader(contract.DownwardMessageDispatcherFacetABI))
	if err != nil {
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		logrus.Errorf("Failed to parse contract ABI: %v", err)
		return err
	}

	data, err := contractABI.Pack("receiveDownwardMessages", downwardMessages)
	if err != nil {
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		logrus.Errorf("Failed to pack data: %v", err)
		return err
	}

	tx := types.NewTransaction(txNonce, common.HexToAddress(b.cfg.ChildLayerContractAddress), value, gasLimit, gasPrice, data)
	// fmt.Printf("tx: %v\n", tx.Hash().Hex())
	// fmt.Printf("tx Time: %v\n", tx.Time().Unix())
	crossMessages, err := b.l1EventParser.ParseL1RawBridgeEventToCrossChainMessage(b.ctx, msg, tx)
	if err != nil {
		logrus.Errorf("Failed to parse L1 cross chain payload, err: %v, tx: %v", err, tx.Hash())
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		return err
	}

	err = b.insertDepositMessage(crossMessages)
	if err != nil {
		logrus.Errorf("Failed to insert deposit: %v, tx: %v", err, tx.Hash())
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		return err
	}

	b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L1_RawBridgeEventsTableName, msg.ID, int(btypes.Processed))

	return nil
}

func (b *L1Relayer) HandleDownwardMessage(msg *orm.RawBridgeEvent) error {
	// 1. parse downward message
	// 2. setup auth
	// 3. send downward message to child layer contract by calling downwardMessageDispatcher.ReceiveDownwardMessages
	payloadBytes, err := hex.DecodeString(msg.MessagePayload)
	if err != nil {
		logrus.Errorf("Failed to decode hex string: %v", err)
		return err
	}
	chainId, err := b.l2Client.ChainID(b.ctx)
	if err != nil {
		logrus.Errorf("HandleDownwardMessage Failed to get chain ID: %v", err)
		return err
	}
	downwardMessages := []contract.DownwardMessage{
		{
			PayloadType: uint32(msg.MessagePayloadType),
			Payload:     payloadBytes,
			Nonce:       big.NewInt(int64(msg.MessageNonce)),
		},
	}
	metrics.DownwardMessageReceivedCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()

	relayerPkStr, err := LoadPrivateKey(b.cfg.RelayerEnvFile, b.cfg.RelayerEnvVar)
	if err != nil {
		logrus.Fatalf("Error loading private key: %v", err)
	}
	relayerPk, err := crypto.HexToECDSA(relayerPkStr)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(relayerPk, chainId)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	if err != nil {
		log.Fatal("Failed to estimate gas:", err)
	}

	contractAddress := common.HexToAddress(b.cfg.ChildLayerContractAddress)
	downwardMessageDispatcher, err := contract.NewDownwardMessageDispatcherFacet(contractAddress, b.l2Client)

	session := &contract.DownwardMessageDispatcherFacetSession{
		Contract:     downwardMessageDispatcher,
		TransactOpts: *auth,
	}

	tx, err := session.ReceiveDownwardMessages(downwardMessages)
	if err != nil {
		if strings.Contains(err.Error(), "Message was already successfully executed") {
			b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L1_RawBridgeEventsTableName, msg.ID, int(btypes.Processed))
			return nil
		}
		logrus.Errorf("Failed to send downward messages: %v", err)
		b.rawBridgeEventOrm.UpdateProcessFail(b.cfg.L1_RawBridgeEventsTableName, msg.ID, err.Error())
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		return err
	}

	crossMessages, err := b.l1EventParser.ParseL1RawBridgeEventToCrossChainMessage(b.ctx, msg, tx)
	if err != nil {
		logrus.Errorf("Failed to parse L1 cross chain payload, err: %v, tx: %v", err, tx.Hash())
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		return err
	}

	err = b.insertDepositMessage(crossMessages)
	if err != nil {
		logrus.Errorf("Failed to insert deposit: %v, tx: %v", err, tx.Hash())
		metrics.DownwardMessageFailureCounter.WithLabelValues(fmt.Sprintf("%d", msg.MessagePayloadType)).Inc()
		return err
	}

	b.rawBridgeEventOrm.UpdateProcessStatus(b.cfg.L1_RawBridgeEventsTableName, msg.ID, int(btypes.Processed))

	return nil
}

func GetCurrentBaseFee(client *ethclient.Client) (*big.Int, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	if header.BaseFee == nil {
		return nil, errors.New("chain does not support EIP-1559")
	}

	return header.BaseFee, nil
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

func (b *L1Relayer) createRefundMessage(msgs []*orm.CrossMessage) error {
	privateKeys, err := LoadPrivateKeyArray(b.cfg.MultisigEnvFile, b.cfg.MultisigEnvVar)
	if err != nil {
		logrus.Fatalf("Error loading private key: %v", err)
	}

	for _, msg := range msgs {
		var upwardMessages []contract.UpwardMessage
		payloadBytes, err := hex.DecodeString(msg.MessagePayload)
		if err != nil {
			return err
		}
		upwardMessages = append(upwardMessages, contract.UpwardMessage{
			PayloadType: uint32(msg.MessagePayloadType),
			Payload:     payloadBytes,
			Nonce:       utils.GenerateNonce(),
		})
		signaturesArray, err := generateUpwardMessageMultiSignatures(upwardMessages, privateKeys)
		if err != nil {
			logrus.Fatalf("Failed to generate multi-signatures: %v", err)
			return err
		}

		messageHash, err := utils.ComputeMessageHash(upwardMessages[0].PayloadType, upwardMessages[0].Payload, upwardMessages[0].Nonce)
		if err != nil {
			logrus.Fatalf("Failed to compute message hash: %v", err)
			return err
		}
		msg.MessageHash = messageHash.Hex()
		msg.MessageNonce = upwardMessages[0].Nonce.String()
		var multiSignProofs []string
		for _, sig := range signaturesArray {
			multiSignProofs = append(multiSignProofs, "0x"+hex.EncodeToString(sig))
		}

		msg.MultiSignProof = strings.Join(multiSignProofs, ",")
	}

	if msgs != nil {
		err = b.crossMessageOrm.InsertOrUpdateCrossMessages(context.Background(), msgs)
		if err != nil {
			logrus.Info("Failed to insert or update L2 messages:", err)
			return err
		}
	}
	return nil
}

func (b *L1Relayer) insertDepositMessage(msgs []*orm.CrossMessage) error {

	if msgs != nil {
		err := b.crossMessageOrm.InsertOrUpdateCrossMessages(context.Background(), msgs)
		if err != nil {
			logrus.Info("Failed to insert or update L2 messages:", err)
		}
	}
	return nil
}
