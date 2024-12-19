package logic

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
)

// HistoryLogic services.
type HistoryLogic struct {
	crossMessageOrm *orm.CrossMessage

	singleFlight singleflight.Group
}

// NewHistoryLogic returns bridge history services.
func NewHistoryLogic(db *gorm.DB) *HistoryLogic {
	if err := db.AutoMigrate(&orm.CrossMessage{}); err != nil {
		logrus.Error("Failed to auto migrate: %v", err)
	}
	logic := &HistoryLogic{
		crossMessageOrm: orm.NewCrossMessage(db),
	}
	return logic
}

// GetL2UnclaimedWithdrawalsByAddress gets all unclaimed withdrawal txs under given address.
func (h *HistoryLogic) GetL2UnclaimedWithdrawalsByAddress(ctx context.Context, address string, page, pageSize uint64) ([]*types.TxHistoryInfo, uint64, error) {
	cacheKey := fmt.Sprintf("unclaimed_withdrawals_%s_%d_%d", address, page, pageSize)
	logrus.Info("cache miss", "cache key", cacheKey)
	var total uint64
	result, err, _ := h.singleFlight.Do(cacheKey, func() (interface{}, error) {
		var txHistoryInfos []*types.TxHistoryInfo
		crossMessages, totalCount, getErr := h.crossMessageOrm.GetL2UnclaimedWithdrawalsByAddress(ctx, address, page, pageSize)
		if getErr != nil {
			return nil, getErr
		}
		for _, message := range crossMessages {
			txHistoryInfos = append(txHistoryInfos, getTxHistoryInfoFromCrossMessage(message))
		}
		total = totalCount
		return txHistoryInfos, nil
	})
	if err != nil {
		logrus.Error("failed to get L2 claimable withdrawals by address", "address", address, "error", err)
		return nil, 0, err
	}

	txHistoryInfos, ok := result.([]*types.TxHistoryInfo)
	if !ok {
		logrus.Error("unexpected type", "expected", "[]*types.TxHistoryInfo", "got", reflect.TypeOf(result), "address", address)
		return nil, 0, errors.New("unexpected error")
	}

	return txHistoryInfos, uint64(total), nil
}

// GetL2UnclaimedWithdrawalsByAddress gets all unclaimed withdrawal txs under given address.
func (h *HistoryLogic) GetTxsByAddress(ctx context.Context, address string, page, pageSize uint64) ([]*types.TxHistoryInfo, uint64, error) {
	cacheKey := fmt.Sprintf("txs_by_address_%s_%d_%d", address, page, pageSize)
	logrus.Info("cache miss", "cache key", cacheKey)
	var total uint64
	result, err, _ := h.singleFlight.Do(cacheKey, func() (interface{}, error) {
		var txHistoryInfos []*types.TxHistoryInfo
		crossMessages, totalCount, getErr := h.crossMessageOrm.GetTxsByAddress(ctx, address, page, pageSize)
		if getErr != nil {
			return nil, getErr
		}
		for _, message := range crossMessages {
			txHistoryInfos = append(txHistoryInfos, getTxHistoryInfoFromCrossMessage(message))
		}
		total = totalCount
		return txHistoryInfos, nil
	})
	if err != nil {
		logrus.Error("failed to get L2 claimable withdrawals by address", "address", address, "error", err)
		return nil, 0, err
	}

	txHistoryInfos, ok := result.([]*types.TxHistoryInfo)
	if !ok {
		logrus.Error("unexpected type", "expected", "[]*types.TxHistoryInfo", "got", reflect.TypeOf(result), "address", address)
		return nil, 0, errors.New("unexpected error")
	}

	return txHistoryInfos, uint64(total), nil
}

func getTxHistoryInfoFromCrossMessage(message *orm.CrossMessage) *types.TxHistoryInfo {
	txHistory := &types.TxHistoryInfo{
		MessageHash:    message.MessageHash,
		TokenType:      types.TokenType(message.TokenType),
		TokenIDs:       utils.ConvertStringToStringArray(message.TokenIDs),
		TokenAmounts:   utils.ConvertStringToStringArray(message.TokenAmounts),
		L1TokenAddress: message.L1TokenAddress,
		L2TokenAddress: message.L2TokenAddress,
		MessageType:    types.MessageType(message.MessageType),
		TxStatus:       types.TxStatusType(message.TxStatus),
		TxType:         types.TxType(message.TxType),
		BlockTimestamp: message.BlockTimestamp,
	}

	txHistory.Hash = message.L2TxHash
	txHistory.BlockNumber = message.L2BlockNumber
	txHistory.ClaimInfo = &types.ClaimInfo{
		From:  message.MessageFrom,
		To:    message.MessageTo,
		Value: message.MessageValue,

		Message: types.Message{
			PayloadType: uint32(message.MessagePayloadType),
			Payload:     message.MessagePayload,
			Nonce:       message.MessageNonce,
		},
		Proof: types.L2MessageProof{
			MultiSignProof: message.MultiSignProof,
		},
	}

	return txHistory
}
