package logic

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"
	"github.com/reddio-com/reddio/bridge/orm"
	"github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils"
)

const (
// cacheKeyPrefixBridgeHistory serves as a specific namespace for all Redis cache keys
// associated with the 'bridge-history' user. This prefix is used to enforce access controls
// in Redis, allowing permissions to be set such that only users with the appropriate
// access rights can read or write to keys starting with "bridge-history".
// cacheKeyPrefixBridgeHistory = "bridge-history-"

// cacheKeyPrefixL2ClaimableWithdrawalsByAddr = cacheKeyPrefixBridgeHistory + "l2ClaimableWithdrawalsByAddr:"
// cacheKeyPrefixL2WithdrawalsByAddr          = cacheKeyPrefixBridgeHistory + "l2WithdrawalsByAddr:"
// cacheKeyPrefixTxsByAddr                    = cacheKeyPrefixBridgeHistory + "txsByAddr:"
// cacheKeyPrefixQueryTxsByHashes             = cacheKeyPrefixBridgeHistory + "queryTxsByHashes:"
// cacheKeyExpiredTime                        = 1 * time.Minute
)

// HistoryLogic services.
type HistoryLogic struct {
	crossMessageOrm *orm.CrossMessage

	singleFlight singleflight.Group
}

// NewHistoryLogic returns bridge history services.
func NewHistoryLogic(db *gorm.DB) *HistoryLogic {
	if err := db.AutoMigrate(&orm.CrossMessage{}); err != nil {
		log.Error("Failed to auto migrate: %v", err)
	}
	logic := &HistoryLogic{
		crossMessageOrm: orm.NewCrossMessage(db),
	}
	return logic
}

// GetL2UnclaimedWithdrawalsByAddress gets all unclaimed withdrawal txs under given address.
func (h *HistoryLogic) GetL2UnclaimedWithdrawalsByAddress(ctx context.Context, address string, page, pageSize uint64) ([]*types.TxHistoryInfo, uint64, error) {
	cacheKey := fmt.Sprintf("unclaimed_withdrawals_%s_%d_%d", address, page, pageSize)
	log.Info("cache miss", "cache key", cacheKey)
	//fmt.Println("cache miss", "cache key", cacheKey)
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
		log.Error("failed to get L2 claimable withdrawals by address", "address", address, "error", err)
		return nil, 0, err
	}

	txHistoryInfos, ok := result.([]*types.TxHistoryInfo)
	if !ok {
		log.Error("unexpected type", "expected", "[]*types.TxHistoryInfo", "got", reflect.TypeOf(result), "address", address)
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
