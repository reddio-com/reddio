package logic

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// represent the structure of the event logs ParentEthBurnt from ChildTokenMessageTransmitterFacet.sol
type L2ETHBurnt struct {
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// represent the structure of the event logs ParentREDTokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2REDBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// represent the structure of the event logs ParentERC20TokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2ERC20TokenBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// represent the structure of the event logs ParentERC721TokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2Erc721TokenBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	TokenID         *big.Int       // uint256
}

// represent the structure of the event logs ParentERC1155TokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2Erc1155BatchTokenBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	TokenIDs        []*big.Int     // uint256[]
	Amounts         []*big.Int     // uint256[]
}

// L2EventParser the l1 event parser
type L2EventParser struct {
	cfg *evm.GethConfig
}

// NewL2EventParser creates l1 event parser
func NewL2EventParser(cfg *evm.GethConfig) *L2EventParser {
	return &L2EventParser{
		cfg: cfg,
	}
}

// ParseL2EventLogs parses L2 watchedevents
func (e *L2EventParser) ParseL2EventLogs(ctx context.Context, logs []types.Log) ([]*orm.CrossMessage, error) {
	l2CrossMessage, err := e.ParseL2SingleCrossChainEventLogs(ctx, logs)
	if err != nil {
		return nil, err
	}
	return l2CrossMessage, nil
}

// ParseL2UpwardMessageEventEventLogs parses L2 watched events
func (e *L2EventParser) ParseL2UpwardMessageEventEventLogs(ctx context.Context, logs []types.Log) ([]*contract.ChildBridgeCoreFacetUpwardMessage, error) {
	events := []*contract.ChildBridgeCoreFacetUpwardMessage{}
	for _, vlog := range logs {
		switch vlog.Topics[0] {
		case backendabi.L2UpwardMessageEventSig:
			//fmt.Println("catch L2UpwardMessageEventSig")
			event := new(contract.ChildBridgeCoreFacetUpwardMessage)
			err := utils.UnpackLog(backendabi.IL2ChildBridgeCoreFacetABI, event, "UpwardMessage", vlog)
			if err != nil {
				log.Error("Failed to unpack UpwardMessage event", "err", err)
				return nil, err
			}
			event.Raw = vlog
			events = append(events, event)

		}
	}
	return events, nil
}

// L2->L1 ParseL2SingleCrossChainEventLogs parses L2 watched events
func (e *L2EventParser) ParseL2SingleCrossChainEventLogs(ctx context.Context, logs []types.Log) ([]*orm.CrossMessage, error) {
	var l2WithdrawMessages []*orm.CrossMessage

	for _, vlog := range logs {
		if vlog.Topics[0] == backendabi.L2UpwardMessageEventSig {
			event := new(contract.ChildBridgeCoreFacetUpwardMessage)
			err := utils.UnpackLog(backendabi.IL2ChildBridgeCoreFacetABI, event, "UpwardMessage", vlog)
			if err != nil {
				log.Error("Failed to unpack UpwardMessage event", "err", err)
				return nil, err
			}
			switch utils.MessagePayloadType(event.PayloadType) {
			case utils.ETH:
				payloadHex := hex.EncodeToString(event.Payload)
				fmt.Println("payloadHex: ", payloadHex)
				l2ETHBurntMsg, err := decodeL2ETHBurnt(payloadHex)
				if err != nil {
					log.Error("Failed to decode ETHLocked", "err", err)
					return nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{
					MessageType:        int(btypes.MessageTypeL2SentMessage),
					TxStatus:           int(btypes.TxStatusTypeSent),
					TokenType:          int(btypes.ETH),
					Sender:             l2ETHBurntMsg.ChildSender.String(),
					Receiver:           l2ETHBurntMsg.ParentRecipient.String(),
					L2TxHash:           vlog.TxHash.String(),
					L2BlockNumber:      vlog.BlockNumber,
					MessagePayloadType: int(btypes.ETH),
					MessagePayload:     payloadHex,
					MessageFrom:        l2ETHBurntMsg.ChildSender.String(),
					MessageTo:          l2ETHBurntMsg.ParentRecipient.String(),
					MessageValue:       l2ETHBurntMsg.Amount.String(),
					//MessageNonce: "",
					//MultiSignProof: "",
					TokenAmounts: l2ETHBurntMsg.Amount.String(),
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				})
			case utils.ERC20:
				payloadHex := hex.EncodeToString(event.Payload)

				l2ERC20BurntMsg, err := decodeERC20TokenBurnt(payloadHex)
				if err != nil {
					log.Error("Failed to decode ERC20TokenBurnt", "err", err)
					return nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{
					MessageType:        int(btypes.MessageTypeL2SentMessage),
					TxStatus:           int(btypes.TxStatusTypeSent),
					TokenType:          int(btypes.ERC20),
					L1TokenAddress:     l2ERC20BurntMsg.TokenAddress.String(),
					Sender:             l2ERC20BurntMsg.ChildSender.String(),
					Receiver:           l2ERC20BurntMsg.ParentRecipient.String(),
					L2TxHash:           vlog.TxHash.String(),
					L2BlockNumber:      vlog.BlockNumber,
					MessagePayloadType: int(btypes.ERC20),
					MessagePayload:     payloadHex,
					MessageFrom:        l2ERC20BurntMsg.ChildSender.String(),
					MessageTo:          l2ERC20BurntMsg.ParentRecipient.String(),
					MessageValue:       l2ERC20BurntMsg.Amount.String(),
					//MessageNonce: "",
					//MultiSignProof: "",
					TokenAmounts: l2ERC20BurntMsg.Amount.String(),
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				})

			case utils.RED:
				payloadHex := hex.EncodeToString(event.Payload)

				l2REDBurntMsg, err := decodeREDTokenBurnt(payloadHex)
				if err != nil {
					log.Error("Failed to decode REDTokenBurnt", "err", err)
					return nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{
					MessageType:        int(btypes.MessageTypeL2SentMessage),
					TxStatus:           int(btypes.TxStatusTypeSent),
					TokenType:          int(btypes.RED),
					L1TokenAddress:     l2REDBurntMsg.TokenAddress.String(),
					Sender:             l2REDBurntMsg.ChildSender.String(),
					Receiver:           l2REDBurntMsg.ParentRecipient.String(),
					L2TxHash:           vlog.TxHash.String(),
					L2BlockNumber:      vlog.BlockNumber,
					MessagePayloadType: int(btypes.RED),
					MessagePayload:     payloadHex,
					MessageFrom:        l2REDBurntMsg.ChildSender.String(),
					MessageTo:          l2REDBurntMsg.ParentRecipient.String(),
					MessageValue:       l2REDBurntMsg.Amount.String(),
					//MessageNonce: "",
					//MultiSignProof: "",
					TokenAmounts: l2REDBurntMsg.Amount.String(),
					CreatedAt:    time.Now().UTC(),
					UpdatedAt:    time.Now().UTC(),
				})
			}
		}
	}
	return l2WithdrawMessages, nil
}
func decodeL2ETHBurnt(payloadHex string) (*L2ETHBurnt, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	ethLocked := &L2ETHBurnt{
		ChildSender:     common.BytesToAddress(payload[0:32]),
		ParentRecipient: common.BytesToAddress(payload[32:64]),
		Amount:          new(big.Int).SetBytes(payload[64:96]),
	}

	return ethLocked, nil
}
func decodeREDTokenBurnt(payloadHex string) (*L2REDBurnt, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	redBurnt := &L2REDBurnt{
		TokenAddress:    common.BytesToAddress(payload[0:32]),
		ChildSender:     common.BytesToAddress(payload[32:64]),
		ParentRecipient: common.BytesToAddress(payload[64:96]),
		Amount:          new(big.Int).SetBytes(payload[96:128]),
	}

	return redBurnt, nil
}

func decodeERC20TokenBurnt(payloadHex string) (*L2ERC20TokenBurnt, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	erc20Burnt := &L2ERC20TokenBurnt{
		TokenAddress:    common.BytesToAddress(payload[0:32]),
		ChildSender:     common.BytesToAddress(payload[32:64]),
		ParentRecipient: common.BytesToAddress(payload[64:96]),
		Amount:          new(big.Int).SetBytes(payload[96:128]),
	}

	return erc20Burnt, nil

}
