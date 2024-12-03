package orm

import (
	"context"
	"fmt"
	"time"

	btypes "github.com/reddio-com/reddio/bridge/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CrossMessage represents a cross message.
type CrossMessage struct {
	db *gorm.DB `gorm:"column:-"`

	ID                 uint64     `json:"id" gorm:"column:id;primary_key;autoIncrement"` // primary key in the database
	MessageType        int        `json:"message_type" gorm:"column:message_type"`       //0:MessageTypeUnknown, 1: MessageTypeL1SentMessage, 2: MessageTypeL2SentMessage
	TxStatus           int        `json:"tx_status" gorm:"column:tx_status"`
	TokenType          int        `json:"token_type" gorm:"column:token_type"` // 0: ETH, 1: ERC20, 2: ERC721, 3: ERC1155, 4: RED
	Sender             string     `json:"sender" gorm:"column:sender"`         // sender address
	Receiver           string     `json:"receiver" gorm:"column:receiver"`
	L1TxHash           string     `json:"l1_tx_hash" gorm:"column:l1_tx_hash"` // initial tx hash, if MessageType is MessageTypeL1SentMessage.
	L2TxHash           string     `json:"l2_tx_hash" gorm:"column:l2_tx_hash"` // initial tx hash, if MessageType is MessageTypeL2SentMessage.
	L1BlockNumber      uint64     `json:"l1_block_number" gorm:"column:l1_block_number"`
	L2BlockNumber      uint64     `json:"l2_block_number" gorm:"column:l2_block_number"`
	L1TokenAddress     string     `json:"l1_token_address" gorm:"column:l1_token_address"`
	L2TokenAddress     string     `json:"l2_token_address" gorm:"column:l2_token_address"`
	TokenIDs           string     `json:"token_ids" gorm:"column:token_ids"`
	TokenAmounts       string     `json:"token_amounts" gorm:"column:token_amounts"`
	BlockTimestamp     uint64     `json:"block_timestamp" gorm:"column:block_timestamp"`
	MessageHash        string     `json:"message_hash" gorm:"column:message_hash;uniqueIndex"` // unique message hash
	MessagePayloadType int        `json:"message_payloadtype" gorm:"column:message_payloadtype"`
	MessagePayload     string     `json:"message_payload" gorm:"column:message_payload"`
	MessageFrom        string     `json:"message_from" gorm:"column:message_from;index"`
	MessageTo          string     `json:"message_to" gorm:"column:message_to"`
	MessageValue       string     `json:"message_value" gorm:"column:message_value"`
	MessageNonce       string     `json:"message_nonce" gorm:"column:message_nonce"`
	MultiSignProof     string     `json:"multiSign_proof" gorm:"column:multiSign_proof"`
	CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
}

// TableName returns the table name for the CrossMessage model.
func (*CrossMessage) TableName() string {
	return "cross_messages"
}

// NewCrossMessage returns a new instance of CrossMessage.
func NewCrossMessage(db *gorm.DB) *CrossMessage {
	return &CrossMessage{db: db}
}

// InsertOrUpdateL2Messages inserts or updates a list of L2 cross messages into the database.
func (c *CrossMessage) InsertOrUpdateL2Messages(ctx context.Context, messages []*CrossMessage) error {
	if len(messages) == 0 {
		return nil
	}
	db := c.db
	db = db.WithContext(ctx)
	db = db.Model(&CrossMessage{})
	// 'tx_status' column is not explicitly assigned during the update to prevent a later status from being overwritten back to "sent".
	db = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "message_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{"sender", "receiver", "token_type", "l2_block_number", "l2_tx_hash", "l1_token_address", "l2_token_address", "token_ids", "token_amounts", "message_type", "block_timestamp", "message_from", "message_to", "message_value", "message_payload", "message_payloadtype", "message_nonce"}),
	})
	if err := db.Create(messages).Error; err != nil {
		return fmt.Errorf("failed to insert message, error: %w", err)
	}
	return nil
}

// GetL2UnclaimedWithdrawalsByAddress retrieves all L2 unclaimed withdrawal messages for a given sender address.
func (c *CrossMessage) GetL2UnclaimedWithdrawalsByAddress(ctx context.Context, sender string) ([]*CrossMessage, error) {
	var messages []*CrossMessage
	db := c.db.WithContext(ctx)
	db = db.Model(&CrossMessage{})
	db = db.Where("message_type = ?", btypes.MessageTypeL2SentMessage)
	db = db.Where("tx_status = ?", btypes.TxStatusTypeSent)
	db = db.Where("message_from = ?", sender)
	db = db.Order("block_timestamp desc")
	db = db.Limit(500)
	if err := db.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get L2 claimable withdrawal messages by sender address, message_from: %v, error: %w", sender, err)
	}
	return messages, nil
}
