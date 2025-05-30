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

// L1EventParser the l1 event parser
type L1EventParser struct {
	cfg *evm.GethConfig
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
func NewL1EventParser(cfg *evm.GethConfig) *L1EventParser {
	return &L1EventParser{
		cfg: cfg,
	}
}

/*****************************
 *    [CrossMessages]       *
 *****************************/
// ParseL1CrossChainEventLogs parse l1 cross chain event logs
func (e *L1EventParser) ParseL1RelayMessagePayload(ctx context.Context, msg *orm.RawBridgeEvent) (*orm.CrossMessage, error) {
	l1RelayedMessage := &orm.CrossMessage{
		MessageHash:   msg.MessageHash,
		L1BlockNumber: msg.BlockNumber,
		L1TxHash:      msg.TxHash,
		TxStatus:      int(btypes.TxStatusTypeConsumed),
	}
	return l1RelayedMessage, nil
}

// ParseL1CrossChainEventLogs parse l1 cross chain event logs
func (e *L1EventParser) ParseL1RawBridgeEventToCrossChainMessage(ctx context.Context, msg *orm.RawBridgeEvent, tx *types.Transaction) ([]*orm.CrossMessage, error) {
	l1CrossChainDepositMessages, err := e.ParseL1SingleRawBridgeEventToCrossChainMessage(ctx, msg, tx)
	if err != nil {
		return nil, err
	}

	return l1CrossChainDepositMessages, nil
}

func (e *L1EventParser) ParseL1CrossChainPayloadToRefundMsg(ctx context.Context, msg *orm.CrossMessage, receipt *types.Receipt) ([]*orm.CrossMessage, error) {
	var refundMessages []*orm.CrossMessage

	switch btypes.MessagePayloadType(msg.MessagePayloadType) {
	case btypes.PayloadTypeETH:
		payloadHex := msg.MessagePayload

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
			BlockTimestamp:     uint64(time.Now().Unix()),
			RefundTxHash:       receipt.TxHash.String(),
			//L2TxHash:           receipt.TxHash.String(),
		})
	case btypes.PayloadTypeERC20:
		payloadHex := msg.MessagePayload

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
			BlockTimestamp:     uint64(time.Now().Unix()),
			RefundTxHash:       receipt.TxHash.String(),
			//L2TxHash:           receipt.TxHash.String(),
			//L1TxHash:           msg.Raw.TxHash.String(),
		})

	case btypes.PayloadTypeRED:
		payloadHex := msg.MessagePayload

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
			BlockTimestamp:     uint64(time.Now().Unix()),
			RefundTxHash:       receipt.TxHash.String(),
			//L1TxHash:           msg.Raw.TxHash.String(),
			//L2TxHash:           receipt.TxHash.String(),
		})

	}
	return refundMessages, nil

}

// ParseL1SingleCrossChainEventLogs parses L1 watched single cross chain events.
func (e *L1EventParser) ParseL1SingleRawBridgeEventToCrossChainMessage(ctx context.Context, bridgeEvent *orm.RawBridgeEvent, tx *types.Transaction) ([]*orm.CrossMessage, error) {
	var l1DepositMessages []*orm.CrossMessage

	switch btypes.MessagePayloadType(bridgeEvent.MessagePayloadType) {
	case btypes.PayloadTypeETH:
		ethLocked, err := decodeETHLocked(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ETHLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{

			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.ETH),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             ethLocked.ParentSender.String(),
			Receiver:           ethLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.ETH),
			MessagePayload:     bridgeEvent.MessagePayload,
			MessageFrom:        ethLocked.ParentSender.String(),
			MessageTo:          ethLocked.ChildRecipient.String(),
			MessageValue:       ethLocked.Amount.String(),
			TokenAmounts:       ethLocked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			L1TxHash:       bridgeEvent.TxHash,
			L2TxHash:       tx.Hash().String(),
			L1BlockNumber:  bridgeEvent.BlockNumber,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: uint64(tx.Time().Unix()),
		})
	case btypes.PayloadTypeERC20:

		erc20Locked, err := decodeERC20TokenLocked(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ParentERC20TokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.ERC20),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             erc20Locked.ParentSender.String(),
			Receiver:           erc20Locked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.ERC20),
			MessagePayload:     bridgeEvent.MessagePayload,
			L1TokenAddress:     erc20Locked.TokenAddress.String(),
			MessageFrom:        erc20Locked.ParentSender.String(),
			MessageTo:          erc20Locked.ChildRecipient.String(),
			MessageValue:       erc20Locked.Amount.String(),
			TokenAmounts:       erc20Locked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			L1BlockNumber:  bridgeEvent.BlockNumber,
			L1TxHash:       bridgeEvent.TxHash,
			L2TxHash:       tx.Hash().String(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: uint64(tx.Time().Unix()),
		})
	case btypes.PayloadTypeRED:

		redLocked, err := decodeREDTokenLocked(bridgeEvent.MessagePayload)
		if err != nil {
			logrus.Error("Failed to decode ParentREDTokenLocked", "err", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.CrossMessage{
			MessageType:        int(btypes.MessageTypeL1SentMessage),
			TxStatus:           int(btypes.TxStatusTypeSent),
			TokenType:          int(btypes.RED),
			TxType:             int(btypes.TxTypeDeposit),
			Sender:             redLocked.ParentSender.String(),
			Receiver:           redLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.RED),
			MessagePayload:     bridgeEvent.MessagePayload,
			L1TokenAddress:     redLocked.TokenAddress.String(),
			MessageFrom:        redLocked.ParentSender.String(),
			MessageTo:          redLocked.ChildRecipient.String(),
			MessageValue:       redLocked.Amount.String(),
			//toDo: change to message nonce to uint64
			MessageNonce:   fmt.Sprintf("%d", bridgeEvent.MessageNonce),
			MessageHash:    bridgeEvent.MessageHash,
			TokenAmounts:   redLocked.Amount.String(),
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			BlockTimestamp: uint64(tx.Time().Unix()),
			L1BlockNumber:  bridgeEvent.BlockNumber,
			L1TxHash:       bridgeEvent.TxHash,
			L2TxHash:       tx.Hash().String(),
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

/*****************************
 *    [RawBridgeEvent]       *
 *****************************/
func (e *L1EventParser) ParseL1EventToRawBridgeEvents(ctx context.Context, logs []types.Log) ([]*orm.RawBridgeEvent, []*orm.RawBridgeEvent, error) {
	var l1DepositMessages []*orm.RawBridgeEvent
	var l1RelayedMessages []*orm.RawBridgeEvent

	for _, vlog := range logs {
		// vlogJson, err := json.MarshalIndent(vlog, "", "  ")
		// if err != nil {
		// 	logrus.Errorf("json.MarshalIndent(vlog) failed: %v", err)
		// 	return nil, nil, err
		// }
		//fmt.Printf("vlog: %s\n", vlogJson)
		if vlog.Topics[0] == backendabi.L1QueueTransactionEventSig {
			event := new(contract.ParentBridgeCoreFacetQueueTransaction)
			err := utils.UnpackLog(backendabi.IL1ParentBridgeCoreFacetABI, event, "QueueTransaction", vlog)
			if err != nil {
				logrus.Error("Failed to unpack UpwardMessage event", "err", err)
				return nil, nil, err
			}
			switch btypes.MessagePayloadType(event.PayloadType) {
			case btypes.PayloadTypeETH:
				payloadHex := hex.EncodeToString(event.Payload)
				ethLocked, err := decodeETHLocked(payloadHex)
				if err != nil {
					logrus.Error("Failed to decode ETHLocked", "err", err)
					return nil, nil, err
				}
				l1DepositMessages = append(l1DepositMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.QueueTransaction),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ParentLayerContractAddress,
					TokenType:       int(btypes.ETH),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             ethLocked.ParentSender.String(),
					Receiver:           ethLocked.ChildRecipient.String(),
					MessagePayloadType: int(btypes.ETH),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.QueueIndex),
					MessageFrom:        ethLocked.ParentSender.String(),
					MessageTo:          ethLocked.ChildRecipient.String(),
					MessageValue:       ethLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.Hash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})

			case btypes.PayloadTypeRED:
				payloadHex := hex.EncodeToString(event.Payload)

				redLocked, err := decodeREDTokenLocked(payloadHex)
				if err != nil {
					logrus.Error("Failed to decode REDTokenBurnt", "err", err)
					return nil, nil, err
				}
				l1DepositMessages = append(l1DepositMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.QueueTransaction),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ParentLayerContractAddress,
					TokenType:       int(btypes.RED),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             redLocked.ParentSender.String(),
					Receiver:           redLocked.ChildRecipient.String(),
					MessagePayloadType: int(btypes.RED),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.QueueIndex),
					MessageFrom:        redLocked.ParentSender.String(),
					MessageTo:          redLocked.ChildRecipient.String(),
					MessageValue:       redLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.Hash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			}
		} else if vlog.Topics[0] == backendabi.L1RelayedMessageEventSig {
			//fmt.Println("find RelayedMessage!")
			event := new(contract.UpwardMessageDispatcherFacetRelayedMessage)
			err := utils.UnpackLog(backendabi.UpwardMessageDispatcherFacetABI, event, "RelayedMessage", vlog)
			if err != nil {
				logrus.Error("Failed to unpack event event", "err", err)
				return nil, nil, err
			}

			switch btypes.MessagePayloadType(event.PayloadType) {
			case btypes.PayloadTypeETH:
				payloadHex := hex.EncodeToString(event.Payload)

				ethLocked, err := decodeETHLocked(payloadHex)
				if err != nil {
					fmt.Errorf("Failed to decode ETHLocked: %v", err)
					return nil, nil, err
				}
				l1RelayedMessages = append(l1RelayedMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.L1RelayedMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ParentLayerContractAddress,
					TokenType:       int(btypes.PayloadTypeETH),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             ethLocked.ParentSender.String(),
					Receiver:           ethLocked.ChildRecipient.String(),
					MessagePayloadType: int(btypes.PayloadTypeETH),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        ethLocked.ParentSender.String(),
					MessageTo:          ethLocked.ChildRecipient.String(),
					MessageValue:       ethLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.MessageHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			case btypes.PayloadTypeRED:
				payloadHex := hex.EncodeToString(event.Payload)

				redLocked, err := decodeREDTokenLocked(payloadHex)
				if err != nil {
					fmt.Errorf("Failed to decode ETHLocked: %v", err)
					return nil, nil, err
				}
				l1RelayedMessages = append(l1RelayedMessages, &orm.RawBridgeEvent{
					EventType:       int(btypes.L1RelayedMessage),
					ChainID:         int(e.cfg.ChainID),
					ContractAddress: e.cfg.ParentLayerContractAddress,
					TokenType:       int(btypes.PayloadTypeRED),
					TxHash:          vlog.TxHash.String(),
					//GasPriced:
					//GasUsed:
					//MsgValue:           msg.Raw
					Timestamp:          uint64(time.Now().Unix()),
					BlockNumber:        vlog.BlockNumber,
					Sender:             redLocked.ParentSender.String(),
					Receiver:           redLocked.ChildRecipient.String(),
					MessagePayloadType: int(btypes.PayloadTypeRED),
					MessagePayload:     payloadHex,
					MessageNonce:       int(event.Nonce.Int64()),
					MessageFrom:        redLocked.ParentSender.String(),
					MessageTo:          redLocked.ChildRecipient.String(),
					MessageValue:       redLocked.Amount.String(),
					MessageHash:        common.BytesToHash(event.MessageHash[:]).String(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          time.Now().UTC(),
					ProcessStatus:      int(btypes.UnProcessed),
				})
			}
		}
	}
	return l1DepositMessages, l1RelayedMessages, nil
}
func (e *L1EventParser) ParseDepositEventToRawBridgeEvents(ctx context.Context, msg *contract.ParentBridgeCoreFacetQueueTransaction) ([]*orm.RawBridgeEvent, error) {
	var l1DepositMessages []*orm.RawBridgeEvent

	switch btypes.MessagePayloadType(msg.PayloadType) {
	case btypes.PayloadTypeETH:
		payloadHex := hex.EncodeToString(msg.Payload)

		ethLocked, err := decodeETHLocked(payloadHex)
		if err != nil {
			fmt.Errorf("Failed to decode ETHLocked: %v", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.RawBridgeEvent{
			EventType:       int(btypes.QueueTransaction),
			ChainID:         int(e.cfg.ChainID),
			ContractAddress: e.cfg.ParentLayerContractAddress,
			TokenType:       int(btypes.ETH),
			TxHash:          msg.Raw.TxHash.String(),
			//GasPriced:
			//GasUsed:
			//MsgValue:           msg.Raw
			Timestamp:          uint64(time.Now().Unix()),
			BlockNumber:        msg.Raw.BlockNumber,
			Sender:             ethLocked.ParentSender.String(),
			Receiver:           ethLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.ETH),
			MessagePayload:     payloadHex,
			MessageNonce:       int(msg.QueueIndex),
			MessageFrom:        ethLocked.ParentSender.String(),
			MessageTo:          ethLocked.ChildRecipient.String(),
			MessageValue:       ethLocked.Amount.String(),
			MessageHash:        common.BytesToHash(msg.Hash[:]).String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			ProcessStatus:      int(btypes.UnProcessed),
		})
	case btypes.PayloadTypeRED:
		payloadHex := hex.EncodeToString(msg.Payload)

		redLocked, err := decodeREDTokenLocked(payloadHex)
		if err != nil {
			fmt.Errorf("Failed to decode redLocked: %v", err)
			return nil, err
		}
		l1DepositMessages = append(l1DepositMessages, &orm.RawBridgeEvent{
			EventType:       int(btypes.QueueTransaction),
			ChainID:         int(e.cfg.ChainID),
			ContractAddress: e.cfg.ParentLayerContractAddress,
			TokenType:       int(btypes.RED),
			TxHash:          msg.Raw.TxHash.String(),
			//GasPriced:
			//GasUsed:
			//MsgValue:           msg.Raw
			Timestamp:          uint64(time.Now().Unix()),
			BlockNumber:        msg.Raw.BlockNumber,
			Sender:             redLocked.ParentSender.String(),
			Receiver:           redLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.RED),
			MessagePayload:     payloadHex,
			MessageNonce:       int(msg.QueueIndex),
			MessageFrom:        redLocked.ParentSender.String(),
			MessageTo:          redLocked.ChildRecipient.String(),
			MessageValue:       redLocked.Amount.String(),
			MessageHash:        common.BytesToHash(msg.Hash[:]).String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			ProcessStatus:      int(btypes.UnProcessed),
		})

	}
	return l1DepositMessages, nil
}

func (e *L1EventParser) ParseL1RelayedMessageToRawBridgeEvents(ctx context.Context, msg *contract.UpwardMessageDispatcherFacetRelayedMessage) ([]*orm.RawBridgeEvent, error) {
	var l1RelayedMessages []*orm.RawBridgeEvent

	switch btypes.MessagePayloadType(msg.PayloadType) {
	case btypes.PayloadTypeETH:
		payloadHex := hex.EncodeToString(msg.Payload)

		ethLocked, err := decodeETHLocked(payloadHex)
		if err != nil {
			fmt.Errorf("Failed to decode ETHLocked: %v", err)
			return nil, err
		}
		l1RelayedMessages = append(l1RelayedMessages, &orm.RawBridgeEvent{
			EventType:       int(btypes.L1RelayedMessage),
			ChainID:         int(e.cfg.ChainID),
			ContractAddress: e.cfg.ParentLayerContractAddress,
			TokenType:       int(btypes.PayloadTypeETH),
			TxHash:          msg.Raw.TxHash.String(),
			//GasPriced:
			//GasUsed:
			//MsgValue:           msg.Raw
			Timestamp:          uint64(time.Now().Unix()),
			BlockNumber:        msg.Raw.BlockNumber,
			Sender:             ethLocked.ParentSender.String(),
			Receiver:           ethLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.PayloadTypeETH),
			MessagePayload:     payloadHex,
			MessageNonce:       int(msg.Nonce.Int64()),
			MessageFrom:        ethLocked.ParentSender.String(),
			MessageTo:          ethLocked.ChildRecipient.String(),
			MessageValue:       ethLocked.Amount.String(),
			MessageHash:        common.BytesToHash(msg.MessageHash[:]).String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			ProcessStatus:      int(btypes.UnProcessed),
		})
	case btypes.PayloadTypeRED:
		payloadHex := hex.EncodeToString(msg.Payload)

		redLocked, err := decodeREDTokenLocked(payloadHex)
		if err != nil {
			fmt.Errorf("Failed to decode redLocked: %v", err)
			return nil, err
		}
		l1RelayedMessages = append(l1RelayedMessages, &orm.RawBridgeEvent{
			EventType:       int(btypes.L1RelayedMessage),
			ChainID:         int(e.cfg.ChainID),
			ContractAddress: e.cfg.ParentLayerContractAddress,
			TokenType:       int(btypes.PayloadTypeRED),
			TxHash:          msg.Raw.TxHash.String(),
			//GasPriced:
			//GasUsed:
			//MsgValue:           msg.Raw
			Timestamp:          uint64(time.Now().Unix()),
			BlockNumber:        msg.Raw.BlockNumber,
			Sender:             redLocked.ParentSender.String(),
			Receiver:           redLocked.ChildRecipient.String(),
			MessagePayloadType: int(btypes.PayloadTypeRED),
			MessagePayload:     payloadHex,
			MessageNonce:       int(msg.Nonce.Int64()),
			MessageFrom:        redLocked.ParentSender.String(),
			MessageTo:          redLocked.ChildRecipient.String(),
			MessageValue:       redLocked.Amount.String(),
			MessageHash:        common.BytesToHash(msg.MessageHash[:]).String(),
			CreatedAt:          time.Now().UTC(),
			UpdatedAt:          time.Now().UTC(),
			ProcessStatus:      int(btypes.UnProcessed),
		})

	}
	return l1RelayedMessages, nil
}
