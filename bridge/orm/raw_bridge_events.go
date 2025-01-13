package orm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// BridgeEvents represents a bridge event.
type RawBridgeEvents struct {
	db                 *gorm.DB   `gorm:"column:-"`
	ID                 uint64     `json:"id" gorm:"column:id;primary_key;autoIncrement"` // primary key in the database
	EventType          int        `json:"event_type" gorm:"column:event_type"`           // 1.QueueTransaction(L1DepositMsgSent) 2.L2RelayedMessage(L2DepositMsgConsumed) 3.SentMessage(L2withdrawMsgSent) 4.L1RelayedMessage(L2DepositMsgConsumed)
	ChainID            int        `json:"chain_id" gorm:"column:chain_id"`               // L1:1 11155111  L2:2 50341
	ContractAddress    string     `json:"contract_address" gorm:"column:contract_address"`
	TokenType          int        `json:"token_type" gorm:"column:token_type"`
	TxHash             string     `json:"tx_hash" gorm:"column:tx_hash"`
	GasPriced          string     `json:"gas_priced" gorm:"column:gas_priced"`
	BlockNumber        uint64     `json:"block_number" gorm:"column:block_number"`
	GasUsed            uint64     `json:"gas_used" gorm:"column:gas_used"`
	MsgValue           string     `json:"msg_value" gorm:"column:msg_value"`
	Timestamp          uint64     `json:"timestamp" gorm:"column:timestamp"`
	Sender             string     `json:"sender" gorm:"column:sender"` // sender address
	Receiver           string     `json:"receiver" gorm:"column:receiver"`
	MessageHash        string     `json:"message_hash" gorm:"column:message_hash;type:varchar(256);uniqueIndex"` // unique message hash
	MessagePayloadType int        `json:"message_payloadtype" gorm:"column:message_payloadtype"`
	MessagePayload     string     `json:"message_payload" gorm:"column:message_payload"`
	MessageNonce       int        `json:"message_nonce" gorm:"column:message_nonce"`
	MessageFrom        string     `json:"message_from" gorm:"column:message_from;index"`
	MessageTo          string     `json:"message_to" gorm:"column:message_to"`
	MessageValue       string     `json:"message_value" gorm:"column:message_value"`
	CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	Remark             string     `json:"remark" gorm:"column:remark"`
	Status             int        `json:"status" gorm:"column:status"`             // 1.success 2.failed
	CheckStatus        int        `json:"check_status" gorm:"column:check_status"` // 1.checked1 2.checked2
}

// TableName returns the table name for the BridgeEvents model.
func (b *RawBridgeEvents) TableName() string {
	return "raw_bridge_events"
}

// NewBridgeEvents creates a new instance of BridgeEvents.
func NewBridgeEvents(db *gorm.DB) *RawBridgeEvents {
	return &RawBridgeEvents{db: db}
}

// InsertBridgeEvents inserts a new BridgeEvents record into the database.
func (b *RawBridgeEvents) InsertBridgeEvents(ctx context.Context, bridgeEvents []*RawBridgeEvents) error {
	if len(bridgeEvents) == 0 {
		return nil
	}
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvents{})
	// 'tx_status' column is not explicitly assigned during the update to prevent a later status from being overwritten back to "sent".

	for _, event := range bridgeEvents {
		if err := db.Create(event).Error; err != nil {
			if isDuplicateEntryError(err) {
				fmt.Printf("Message with hash %s already exists, skipping insert.\n", event.MessageHash)
				continue
			}
			return fmt.Errorf("failed to insert message, error: %w", err)
		}
	}
	return nil
}
func isDuplicateEntryError(err error) bool {
	return strings.Contains(err.Error(), "Error 1062")
}
