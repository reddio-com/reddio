package logic

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// L1EventParser the l1 event parser
type L1EventParser struct {
	cfg    *evm.GethConfig
	client *ethclient.Client
}

type CrossMessage struct {
	Sender         string
	Receiver       string
	TokenType      int
	L1TokenAddress string
	MessageValue   int
	TokenAmounts   string
}
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
}

// ParseL1CrossChainEventLogs parse l1 cross chain event logs
func (e *L1EventParser) ParseL1CrossChainPayload(ctx context.Context, msg *contract.ParentBridgeCoreFacetDownwardMessage) ([]*CrossMessage, error) {
	l1CrossChainDepositMessages, err := e.ParseL1SingleCrossChainPayload(ctx, msg)
	if err != nil {
		return nil, err
	}

	return l1CrossChainDepositMessages, nil
}

// ParseL1SingleCrossChainEventLogs parses L1 watched single cross chain events.
func (e *L1EventParser) ParseL1SingleCrossChainPayload(ctx context.Context, msg *contract.ParentBridgeCoreFacetDownwardMessage) ([]*CrossMessage, error) {
	var l1DepositMessages []*CrossMessage

	switch utils.MessagePayloadType(msg.PayloadType) {
	case utils.ETH:
		payloadHex := hex.EncodeToString(msg.Payload)

		ethLocked, err := decodeETHLocked(payloadHex)
		if err != nil {
			log.Error("Failed to decode ETHLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &CrossMessage{
			Sender:       ethLocked.ParentSender.String(),
			Receiver:     ethLocked.ChildRecipient.String(),
			TokenType:    int(msg.PayloadType),
			TokenAmounts: ethLocked.Amount.String(),
		})
	case utils.ERC20:
		payloadHex := hex.EncodeToString(msg.Payload)

		erc20Locked, err := decodeERC20TokenLocked(payloadHex)
		if err != nil {
			log.Error("Failed to decode ParentERC20TokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &CrossMessage{
			Sender:         erc20Locked.ParentSender.String(),
			Receiver:       erc20Locked.ChildRecipient.String(),
			TokenType:      int(msg.PayloadType),
			TokenAmounts:   erc20Locked.Amount.String(),
			L1TokenAddress: erc20Locked.TokenAddress.String(),
		})

	case utils.RED:
		payloadHex := hex.EncodeToString(msg.Payload)

		redLocked, err := decodeREDTokenLocked(payloadHex)
		if err != nil {
			log.Error("Failed to decode ParentREDTokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &CrossMessage{
			Sender:         redLocked.ParentSender.String(),
			Receiver:       redLocked.ChildRecipient.String(),
			TokenType:      int(msg.PayloadType),
			TokenAmounts:   redLocked.Amount.String(),
			L1TokenAddress: redLocked.TokenAddress.String(),
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

	nameOffset := new(big.Int).SetBytes(payload[firstElementOffset+32 : firstElementOffset+64]).Int64()
	symbolOffset := new(big.Int).SetBytes(payload[firstElementOffset+64 : firstElementOffset+96]).Int64()

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
