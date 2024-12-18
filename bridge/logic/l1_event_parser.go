package logic

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// L1EventParser the l1 event parser
type L1EventParser struct {
	cfg    *evm.GethConfig
	client *ethclient.Client
}

//	type CrossMessage struct {
//		Sender         string
//		Receiver       string
//		TokenType      int
//		L1TokenAddress string
//		MessageValue   int
//		TokenAmounts   string
//	}
type ETHLocked struct {
	ParentSender   common.Address
	ChildRecipient common.Address
	Amount         *big.Int
}

type ParentREDTokenLocked struct {
	TokenAddress   common.Address
	ParentSender   common.Address
	ChildRecipient common.Address
	Amount         *big.Int
}

type ParentERC20TokenLocked struct {
	TokenAddress   common.Address
	TokenName      string
	TokenSymbol    string
	Decimals       *big.Int
	ParentSender   common.Address
	ChildRecipient common.Address
	Amount         *big.Int
}

type ParentERC721TokenLocked struct {
	TokenAddress   common.Address
	TokenName      string
	TokenSymbol    string
	ParentSender   common.Address
	ChildRecipient common.Address
	TokenId        *big.Int
}

type ParentERC1155TokenLocked struct {
	TokenAddress   common.Address
	ParentSender   common.Address
	ChildRecipient common.Address
	TokenIds       []*big.Int
	Amounts        []*big.Int
}

// NewL1EventParser creates l1 event parser
func NewL1EventParser(cfg *evm.GethConfig, client *ethclient.Client) *L1EventParser {
	return &L1EventParser{
		cfg:    cfg,
		client: client,
	}
} // ParseL1CrossChainEventLogs parse l1 cross chain event logs
func (e *L1EventParser) ParseL1RelayMessagePayload(ctx context.Context, msg *contract.UpwardMessageDispatcherFacetRelayedMessage) ([]*orm.CrossMessage, error) {
	var l1RelayedMessages []*orm.CrossMessage

	l1RelayedMessages = append(l1RelayedMessages, &orm.CrossMessage{
		MessageHash:   common.BytesToHash(msg.MessageHash[:]).String(),
		L1BlockNumber: msg.Raw.BlockNumber,
		L1TxHash:      msg.Raw.TxHash.String(),
		TxStatus:      int(btypes.TxStatusTypeConsumed),
		MessageType:   int(btypes.MessageTypeL2SentMessage),
	})

	return l1RelayedMessages, nil
}

// ParseL1CrossChainEventLogs parse l1 cross chain event logs
func (e *L1EventParser) ParseL1CrossChainPayload(ctx context.Context, msg *contract.ParentBridgeCoreFacetDownwardMessage, tx *types.Transaction, receipt *types.Receipt) ([]*orm.CrossMessage, error) {
	l1CrossChainDepositMessages, err := e.ParseL1SingleCrossChainPayload(ctx, msg, tx, receipt)
	if err != nil {
		return nil, err
	}

	return l1CrossChainDepositMessages, nil
}

func (e *L1EventParser) ParseL1CrossChainPayloadToRefundMsg(ctx context.Context, msg *contract.ParentBridgeCoreFacetDownwardMessage, tx *types.Transaction, receipt *types.Receipt) ([]*orm.CrossMessage, error) {
	var refundMessages []*orm.CrossMessage

	switch utils.MessagePayloadType(msg.PayloadType) {
	case utils.ETH:
		payloadHex := hex.EncodeToString(msg.Payload)

		ethLocked, err := decodeETHLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ETHLocked", "err", err)
			return nil, err
		}

		refundMessages = append(refundMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.ETH),
			TxType:             int(btypes.TxTypeRefund),
			Sender:             ethLocked.ChildRecipient.String(),
			Receiver:           ethLocked.ParentSender.String(),
			MessagePayloadType: int(btypes.ETH),
			MessagePayload:     payloadHex,
			MessageFrom:        ethLocked.ChildRecipient.String(),
			MessageTo:          ethLocked.ParentSender.String(),
			MessageValue:       ethLocked.Amount.String(),
			TokenAmounts:       ethLocked.Amount.String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L1BlockNumber:      msg.Raw.BlockNumber,
			L1TxHash:           msg.Raw.TxHash.String(),
			L2TxHash:           receipt.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
		})
	case utils.ERC20:
		payloadHex := hex.EncodeToString(msg.Payload)

		erc20Locked, err := decodeERC20TokenLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ParentERC20TokenLocked", "err", err)
			return nil, err
		}
		refundMessages = append(refundMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.ERC20),
			TxType:             int(btypes.TxTypeRefund),
			L1TokenAddress:     erc20Locked.TokenAddress.String(),
			Sender:             erc20Locked.ChildRecipient.String(),
			Receiver:           erc20Locked.ParentSender.String(),
			MessagePayloadType: int(btypes.ERC20),
			MessagePayload:     payloadHex,
			MessageFrom:        erc20Locked.ChildRecipient.String(),
			MessageTo:          erc20Locked.ParentSender.String(),
			MessageValue:       erc20Locked.Amount.String(),
			TokenAmounts:       erc20Locked.Amount.String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L2TxHash:           receipt.TxHash.String(),
			L1BlockNumber:      msg.Raw.BlockNumber,
			L1TxHash:           msg.Raw.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
		})

	case utils.RED:
		payloadHex := hex.EncodeToString(msg.Payload)

		redLocked, err := decodeREDTokenLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ParentREDTokenLocked", "err", err)
			return nil, err
		}
		refundMessages = append(refundMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.RED),
			TxType:             int(btypes.TxTypeRefund),
			L1TokenAddress:     redLocked.TokenAddress.String(),
			Sender:             redLocked.ChildRecipient.String(),
			Receiver:           redLocked.ParentSender.String(),
			MessagePayloadType: int(btypes.RED),
			MessagePayload:     payloadHex,
			MessageFrom:        redLocked.ChildRecipient.String(),
			MessageTo:          redLocked.ParentSender.String(),
			MessageValue:       redLocked.Amount.String(),
			TokenAmounts:       redLocked.Amount.String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L2TxHash:           receipt.TxHash.String(),
			L1BlockNumber:      msg.Raw.BlockNumber,
			L1TxHash:           msg.Raw.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
		})

	}
	return refundMessages, nil

}

// ParseL1SingleCrossChainEventLogs parses L1 watched single cross chain events.
func (e *L1EventParser) ParseL1SingleCrossChainPayload(ctx context.Context, msg *contract.ParentBridgeCoreFacetDownwardMessage, tx *types.Transaction, receipt *types.Receipt) ([]*orm.CrossMessage, error) {
	var l1DepositMessages []*orm.CrossMessage

	switch utils.MessagePayloadType(msg.PayloadType) {
	case utils.ETH:
		payloadHex := hex.EncodeToString(msg.Payload)

		ethLocked, err := decodeETHLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ETHLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{

			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeConsumed),
			TokenType:          int(btypes.ETH),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             ethLocked.ParentSender.String(),
			Receiver:           ethLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.ETH),
			MessagePayload:     payloadHex,
			MessageFrom:        ethLocked.ParentSender.String(),
			MessageTo:          ethLocked.ChildRecipient.String(),
			MessageValue:       ethLocked.Amount.String(),
			TokenAmounts:       ethLocked.Amount.String(),
			L1TxHash:           msg.Raw.TxHash.String(),
			L1BlockNumber:      msg.Raw.BlockNumber,
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L2TxHash:           receipt.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
		})
	case utils.ERC20:
		payloadHex := hex.EncodeToString(msg.Payload)

		erc20Locked, err := decodeERC20TokenLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ParentERC20TokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeConsumed),
			TokenType:          int(btypes.ERC20),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             erc20Locked.ParentSender.String(),
			Receiver:           erc20Locked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.ERC20),
			MessagePayload:     payloadHex,
			L1TokenAddress:     erc20Locked.TokenAddress.String(),
			MessageFrom:        erc20Locked.ParentSender.String(),
			MessageTo:          erc20Locked.ChildRecipient.String(),
			MessageValue:       erc20Locked.Amount.String(),
			TokenAmounts:       erc20Locked.Amount.String(),
			L1BlockNumber:      msg.Raw.BlockNumber,
			L1TxHash:           msg.Raw.TxHash.String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L2TxHash:           receipt.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
		})
	case utils.RED:
		payloadHex := hex.EncodeToString(msg.Payload)

		redLocked, err := decodeREDTokenLocked(payloadHex)
		if err != nil {
			logrus.Error("Failed to decode ParentREDTokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeConsumed),
			TokenType:          int(btypes.RED),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             redLocked.ParentSender.String(),
			Receiver:           redLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.RED),
			MessagePayload:     payloadHex,
			L1TokenAddress:     redLocked.TokenAddress.String(),
			MessageFrom:        redLocked.ParentSender.String(),
			MessageTo:          redLocked.ChildRecipient.String(),
			MessageValue:       redLocked.Amount.String(),
			TokenAmounts:       redLocked.Amount.String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			BlockTimestamp:     uint64(tx.Time().Unix()),
			L2TxHash:           receipt.TxHash.String(),
			L2BlockNumber:      receipt.BlockNumber.Uint64(),
			L1BlockNumber:      msg.Raw.BlockNumber,
			L1TxHash:           msg.Raw.TxHash.String(),
		})

	}
	return l1DepositMessages, nil
}

func decodeETHLocked(payloadHex string) (*ETHLocked, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	ethLocked := &ETHLocked{
		ParentSender:   common.BytesToAddress(payload[0:32]),
		ChildRecipient: common.BytesToAddress(payload[32:64]),
		Amount:         new(big.Int).SetBytes(payload[64:96]),
	}

	return ethLocked, nil
}

func decodeREDTokenLocked(payloadHex string) (*ParentREDTokenLocked, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	redLocked := &ParentREDTokenLocked{
		TokenAddress:   common.BytesToAddress(payload[0:32]),
		ParentSender:   common.BytesToAddress(payload[32:64]),
		ChildRecipient: common.BytesToAddress(payload[64:96]),
		Amount:         new(big.Int).SetBytes(payload[96:128]),
	}

	return redLocked, nil
}

func decodeERC20TokenLocked(payloadHex string) (*ParentERC20TokenLocked, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}
	firstElementOffset := new(big.Int).SetBytes(payload[0:32]).Int64()
	// TokenAddress
	tokenAddress := common.BytesToAddress(payload[firstElementOffset : firstElementOffset+32])

	nameOffset := firstElementOffset + new(big.Int).SetBytes(payload[firstElementOffset+32:firstElementOffset+64]).Int64()
	symbolOffset := firstElementOffset + new(big.Int).SetBytes(payload[firstElementOffset+64:firstElementOffset+96]).Int64()
	// Decimals
	decimals := new(big.Int).SetBytes(payload[firstElementOffset+96 : firstElementOffset+128])

	// ParentSender
	parentSender := common.BytesToAddress(payload[firstElementOffset+128 : firstElementOffset+160])

	// ChildRecipient
	childRecipient := common.BytesToAddress(payload[firstElementOffset+160 : firstElementOffset+192])

	// Amount
	amount := new(big.Int).SetBytes(payload[firstElementOffset+192 : firstElementOffset+224])
	// TokenName
	nameLength := new(big.Int).SetBytes(payload[nameOffset : nameOffset+32]).Int64()
	tokenName := string(payload[nameOffset+32 : nameOffset+32+nameLength])
	// TokenSymbol
	symbolLength := new(big.Int).SetBytes(payload[symbolOffset : symbolOffset+32]).Int64()
	tokenSymbol := string(payload[symbolOffset+32 : symbolOffset+32+symbolLength])

	erc20Locked := &ParentERC20TokenLocked{
		TokenAddress:   tokenAddress,
		TokenName:      tokenName,
		TokenSymbol:    tokenSymbol,
		Decimals:       decimals,
		ParentSender:   parentSender,
		ChildRecipient: childRecipient,
		Amount:         amount,
	}

	return erc20Locked, nil
}
