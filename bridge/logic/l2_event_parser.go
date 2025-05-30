package logic

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"

	backendabi "github.com/reddio-com/reddio/bridge/abi"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/reddio-com/reddio/evm"
)

// L2ETHBurnt represent the structure of the event logs ParentEthBurnt from ChildTokenMessageTransmitterFacet.sol
type L2ETHBurnt struct {
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// L2REDBurnt represent the structure of the event logs ParentREDTokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2REDBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// L2ERC20TokenBurnt represent the structure of the event logs ParentERC20TokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2ERC20TokenBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	Amount          *big.Int       // uint256
}

// L2Erc721TokenBurnt represent the structure of the event logs ParentERC721TokenBurnt from ChildTokenMessageTransmitterFacet.sol
type L2Erc721TokenBurnt struct {
	TokenAddress    common.Address // address
	ChildSender     common.Address // address
	ParentRecipient common.Address // address
	TokenID         *big.Int       // uint256
}

// L2Erc1155BatchTokenBurnt represent the structure of the event logs ParentERC1155TokenBurnt from ChildTokenMessageTransmitterFacet.sol
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
func (e *L2EventParser) ParseL2EventLogs(ctx context.Context, logs []types.Log) ([]*orm.RawBridgeEvent, []*orm.RawBridgeEvent, error) {
	l2WithdrawMessages, l2RelayedMessages, err := e.ParseL2EventToRawBridgeEvents(ctx, logs)
	if err != nil {
		return nil, nil, err
	}

	return l2WithdrawMessages, l2RelayedMessages, nil
}
func (e *L2EventParser) ParseL2EventToRawBridgeEvents(ctx context.Context, logs []types.Log) ([]*orm.RawBridgeEvent, []*orm.RawBridgeEvent, error) {
	var l2WithdrawMessages []*orm.RawBridgeEvent
	var l2RelayedMessages []*orm.RawBridgeEvent

	for _, vlog := range logs {
		// vlogJson, err := json.MarshalIndent(vlog, "", "  ")
		// if err != nil {
		// 	logrus.Errorf("json.MarshalIndent(vlog) failed: %v", err)
		// 	return nil, nil, err
		// }
		//fmt.Printf("vlog: %s\n", vlogJson)
		if vlog.Topics[0] == backendabi.L2SentMessageEventSig {
			event := new(contract.ChildBridgeCoreFacetSentMessage)
			err := utils.UnpackLog(backendabi.IL2ChildBridgeCoreFacetABI, event, "SentMessage", vlog)
			if err != nil {
				logrus.Error("Failed to unpack UpwardMessage event", "err", err)
				return nil, nil, err
			}
			switch btypes.MessagePayloadType(event.PayloadType) {
			case btypes.PayloadTypeETH:
				payloadHex := hex.EncodeToString(event.Payload)
				l2ETHBurntMsg, err := decodeL2ETHBurnt(payloadHex)
				if err != nil {
					logrus.Error("Failed to decode ETHLocked", "err", err)
					return nil, nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.SentMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ChildLayerContractAddress,
					TokenType:       int(btypes.ETH),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             l2ETHBurntMsg.ChildSender.String(),
					Receiver:           l2ETHBurntMsg.ParentRecipient.String(),
					MessagePayloadType: int(btypes.ETH),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        l2ETHBurntMsg.ChildSender.String(),
					MessageTo:          l2ETHBurntMsg.ParentRecipient.String(),
					MessageValue:       l2ETHBurntMsg.Amount.String(),
					MessageHash:        common.BytesToHash(event.XDomainCalldataHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			case btypes.PayloadTypeERC20:
				payloadHex := hex.EncodeToString(event.Payload)

				l2ERC20BurntMsg, err := decodeERC20TokenBurnt(payloadHex)
				if err != nil {
					logrus.Error("Failed to decode ERC20TokenBurnt", "err", err)
					return nil, nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.SentMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ChildLayerContractAddress,
					TokenType:       int(btypes.ERC20),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             l2ERC20BurntMsg.ChildSender.String(),
					Receiver:           l2ERC20BurntMsg.ParentRecipient.String(),
					MessagePayloadType: int(btypes.ERC20),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        l2ERC20BurntMsg.ChildSender.String(),
					MessageTo:          l2ERC20BurntMsg.ParentRecipient.String(),
					MessageValue:       l2ERC20BurntMsg.Amount.String(),
					MessageHash:        common.BytesToHash(event.XDomainCalldataHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})

			case btypes.PayloadTypeRED:
				payloadHex := hex.EncodeToString(event.Payload)

				l2REDBurntMsg, err := decodeREDTokenBurnt(payloadHex)
				if err != nil {
					logrus.Error("Failed to decode REDTokenBurnt", "err", err)
					return nil, nil, err
				}
				l2WithdrawMessages = append(l2WithdrawMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.SentMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ChildLayerContractAddress,
					TokenType:       int(btypes.RED),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             l2REDBurntMsg.ChildSender.String(),
					Receiver:           l2REDBurntMsg.ParentRecipient.String(),
					MessagePayloadType: int(btypes.RED),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        l2REDBurntMsg.ChildSender.String(),
					MessageTo:          l2REDBurntMsg.ParentRecipient.String(),
					MessageValue:       l2REDBurntMsg.Amount.String(),
					MessageHash:        common.BytesToHash(event.XDomainCalldataHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			}
		} else if vlog.Topics[0] == backendabi.DownwardMessageDispatcherFacetABI.Events["RelayedMessage"].ID {
			//fmt.Println("find RelayedMessage!")
			event := new(contract.DownwardMessageDispatcherFacetRelayedMessage)
			err := utils.UnpackLog(backendabi.DownwardMessageDispatcherFacetABI, event, "RelayedMessage", vlog)
			if err != nil {
				logrus.Error("Failed to unpack event event", "err", err)
				return nil, nil, err
			}

			switch btypes.MessagePayloadType(event.PayloadType) {
			case btypes.PayloadTypeETH:
				payloadHex := hex.EncodeToString(event.Payload)

				ethLocked, err := decodeL2ETHBurnt(payloadHex)
				if err != nil {
					fmt.Errorf("Failed to decode ETHLocked: %v", err)
					return nil, nil, err
				}
				l2RelayedMessages = append(l2RelayedMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.L2RelayedMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ChildLayerContractAddress,
					TokenType:       int(btypes.ETH),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             ethLocked.ChildSender.String(),
					Receiver:           ethLocked.ParentRecipient.String(),
					MessagePayloadType: int(btypes.ETH),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        ethLocked.ChildSender.String(),
					MessageTo:          ethLocked.ParentRecipient.String(),
					MessageValue:       ethLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.MessageHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			case btypes.PayloadTypeRED:
				payloadHex := hex.EncodeToString(event.Payload)

				redLocked, err := decodeREDTokenBurnt(payloadHex)
				if err != nil {
					fmt.Errorf("Failed to decode redLocked: %v", err)
					return nil, nil, err
				}
				l2RelayedMessages = append(l2RelayedMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.L2RelayedMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ChildLayerContractAddress,
					TokenType:       int(btypes.RED),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             redLocked.ChildSender.String(),
					Receiver:           redLocked.ParentRecipient.String(),
					MessagePayloadType: int(btypes.RED),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        redLocked.ChildSender.String(),
					MessageTo:          redLocked.ParentRecipient.String(),
					MessageValue:       redLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.MessageHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			}
		}
	}
	return l2WithdrawMessages, l2RelayedMessages, nil
}

// ParseL2SingleCrossChainEventLogs parses L2 watched single cross chain events.
func (e *L2EventParser) ParseL2SingleRawBridgeEventToCrossChainMessage(ctx context.Context, bridgeEvent *orm.RawBridgeEvent) ([]*orm.CrossMessage, error) {
	var l2WithdrawMessages []*orm.CrossMessage

	switch btypes.MessagePayloadType(bridgeEvent.MessagePayloadType) {
	case btypes.PayloadTypeETH:
		ethLocked, err := decodeL2ETHBurnt(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ETHLocked", "err", err)
			return nil, err
		}
		l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{

			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeReadyForConsumption),
			TokenType:          int(btypes.ETH),
			TxType:             int(btypes.TxTypeWithdraw),
			Sender:             ethLocked.ChildSender.String(),
			Receiver:           ethLocked.ParentRecipient.String(),
			MessagePayloadType: int(btypes.ETH),
			MessagePayload:     bridgeEvent.MessagePayload,
			MessageFrom:        ethLocked.ChildSender.String(),
			MessageTo:          ethLocked.ParentRecipient.String(),
			MessageValue:       ethLocked.Amount.String(),
			TokenAmounts:       ethLocked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			L2TxHash:       bridgeEvent.TxHash,
			L2BlockNumber:  bridgeEvent.BlockNumber,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: bridgeEvent.Timestamp,
		})
	case btypes.PayloadTypeERC20:

		erc20Locked, err := decodeERC20TokenBurnt(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ParentERC20TokenLocked", "err", err)
			return nil, err
		}
		l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeReadyForConsumption),
			TokenType:          int(btypes.ERC20),
			TxType:             int(btypes.TxTypeWithdraw),
			Sender:             erc20Locked.ChildSender.String(),
			Receiver:           erc20Locked.ParentRecipient.String(),
			MessagePayloadType: int(btypes.ERC20),
			MessagePayload:     bridgeEvent.MessagePayload,
			L1TokenAddress:     erc20Locked.TokenAddress.String(),
			MessageFrom:        erc20Locked.ChildSender.String(),
			MessageTo:          erc20Locked.ParentRecipient.String(),
			MessageValue:       erc20Locked.Amount.String(),
			TokenAmounts:       erc20Locked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			L2BlockNumber:  bridgeEvent.BlockNumber,
			L2TxHash:       bridgeEvent.TxHash,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: bridgeEvent.Timestamp,
		})
	case btypes.PayloadTypeRED:

		redLocked, err := decodeREDTokenBurnt(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ParentREDTokenLocked", "err", err)
			return nil, err
		}
		l2WithdrawMessages = append(l2WithdrawMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL2SentMessage),
			TxStatus:           int(btypes.TxStatusTypeReadyForConsumption),
			TokenType:          int(btypes.RED),
			TxType:             int(btypes.TxTypeWithdraw),
			Sender:             redLocked.ChildSender.String(),
			Receiver:           redLocked.ParentRecipient.String(),
			MessagePayloadType: int(btypes.RED),
			MessagePayload:     bridgeEvent.MessagePayload,
			L1TokenAddress:     redLocked.TokenAddress.String(),
			MessageFrom:        redLocked.ChildSender.String(),
			MessageTo:          redLocked.ParentRecipient.String(),
			MessageValue:       redLocked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			TokenAmounts:   redLocked.Amount.String(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: bridgeEvent.Timestamp,
			L2BlockNumber:  bridgeEvent.BlockNumber,
			L2TxHash:       bridgeEvent.TxHash,
		})

	}
	return l2WithdrawMessages, nil
}
func (e *L2EventParser) ParseL2RawBridgeEventToCrossChainMessage(ctx context.Context, msg *orm.RawBridgeEvent) ([]*orm.CrossMessage, error) {
	l2CrossChainDepositMessages, err := e.ParseL2SingleRawBridgeEventToCrossChainMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	return l2CrossChainDepositMessages, nil
}

// ParseL2RelayMessagePayload parses a single L2 relay message payload
func (e *L2EventParser) ParseL2RelayMessagePayload(ctx context.Context, msg *orm.RawBridgeEvent) (*orm.CrossMessage, error) {
	l2RelayedMessage := &orm.CrossMessage{
		MessageHash:   msg.MessageHash,
		L2BlockNumber: msg.BlockNumber,
		L2TxHash:      msg.TxHash,
		TxStatus:      int(btypes.TxStatusTypeConsumed),
	}

	return l2RelayedMessage, nil
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

func decodeERC721TokenBurnt(payloadHex string) (*L2Erc721TokenBurnt, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) != 32+32+32+32 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	erc721Burnt := &L2Erc721TokenBurnt{
		TokenAddress:    common.BytesToAddress(payload[0:32]),
		ChildSender:     common.BytesToAddress(payload[32:64]),
		ParentRecipient: common.BytesToAddress(payload[64:96]),
		TokenID:         new(big.Int).SetBytes(payload[96:128]),
	}

	return erc721Burnt, nil
}
func decodeERC1155BatchTokenBurnt(payloadHex string) (*L2Erc1155BatchTokenBurnt, error) {
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	if len(payload) < 96 {
		return nil, fmt.Errorf("invalid payload length: %d", len(payload))
	}

	erc1155Burnt := &L2Erc1155BatchTokenBurnt{
		TokenAddress:    common.BytesToAddress(payload[0:32]),
		ChildSender:     common.BytesToAddress(payload[32:64]),
		ParentRecipient: common.BytesToAddress(payload[64:96]),
	}

	offset := 96
	for offset < len(payload) {
		if offset+32 > len(payload) {
			return nil, fmt.Errorf("invalid payload length for TokenIDs: %d", len(payload))
		}
		tokenID := new(big.Int).SetBytes(payload[offset : offset+32])
		erc1155Burnt.TokenIDs = append(erc1155Burnt.TokenIDs, tokenID)
		offset += 32
	}

	for offset < len(payload) {
		if offset+32 > len(payload) {
			return nil, fmt.Errorf("invalid payload length for Amounts: %d", len(payload))
		}
		amount := new(big.Int).SetBytes(payload[offset : offset+32])
		erc1155Burnt.Amounts = append(erc1155Burnt.Amounts, amount)
		offset += 32
	}

	return erc1155Burnt, nil
}
