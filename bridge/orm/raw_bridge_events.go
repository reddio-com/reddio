package orm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	btypes "github.com/reddio-com/reddio/bridge/types"

	"gorm.io/gorm"
)

const (
	TableRawBridgeEvents11155111 = "raw_bridge_events_11155111"
	TableRawBridgeEvents50341    = "raw_bridge_events_50341"
)

// BridgeEvents represents a bridge event.
type RawBridgeEvent struct {
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
	ProcessStatus      int        `json:"process_status" gorm:"column:process_status"`                             // 1.processed 2.unprocessed
	ProcessFailReason  string     `json:"process_fail_reason" gorm:"column:process_fail_reason;type:varchar(256)"` // reason for process failure
	ProcessFailCount   int        `json:"process_fail_count" gorm:"column:process_fail_count"`                     // number of process failures
	CheckStatus        int        `json:"check_status" gorm:"column:check_status"`                                 // 1.checked1 2.checked2
}

// NewBridgeEvents creates a new instance of BridgeEvents.
func NewRawBridgeEvent(db *gorm.DB) *RawBridgeEvent {
	err := db.Table("raw_bridge_events_11155111").AutoMigrate(&RawBridgeEvent{})
	if err != nil {
		log.Fatal("failed to AutoMigrate db", "err", err)
	}
	err = db.Table("raw_bridge_events_50341").AutoMigrate(&RawBridgeEvent{})
	if err != nil {
		log.Fatal("failed to AutoMigrate db", "err", err)
	}
	return &RawBridgeEvent{db: db}
}

// InsertBridgeEvents inserts a new BridgeEvents record into the database.
func (b *RawBridgeEvent) InsertRawBridgeEvents(ctx context.Context, tableName string, bridgeEvents []*RawBridgeEvent) error {
	if len(bridgeEvents) == 0 {
		return nil
	}
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)

	for _, event := range bridgeEvents {
		if err := db.Create(event).Error; err != nil {
			if isDuplicateEntryError(err) {
				fmt.Errorf("Message with hash %s already exists, skipping insert.\n", event.MessageHash)
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

func (b *RawBridgeEvent) QueryUnProcessedBridgeEventsByEventType(ctx context.Context, tableName string, eventType int, limit int) ([]*RawBridgeEvent, error) {
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)

	var bridgeEvents []*RawBridgeEvent
	if err := db.Where("process_status = ? AND event_type = ?", btypes.UnProcessed, eventType).
		Order("block_number ASC, message_nonce ASC").
		Limit(limit).
		Find(&bridgeEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to query unprocessed bridge events: %w", err)
	}
	return bridgeEvents, nil
}
func (b *RawBridgeEvent) QueryUnProcessedBridgeEvents(ctx context.Context, tableName string, limit int) ([]*RawBridgeEvent, error) {
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)

	var bridgeEvents []*RawBridgeEvent
	if err := db.Where("process_status = ? ", btypes.UnProcessed).
		Order("block_number ASC, message_nonce ASC").
		Limit(limit).
		Find(&bridgeEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to query unprocessed bridge events: %w", err)
	}
	return bridgeEvents, nil
}

func (e *RawBridgeEvent) UpdateProcessStatus(tableName string, id uint64, newStatus int) error {
	db := e.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"process_status": newStatus,
		"updated_at":     time.Now().UTC(),
	}).Error
}

// UpdateProcessStatusBatch updates the process status of multiple RawBridgeEvents.
func (e *RawBridgeEvent) UpdateProcessStatusBatch(tableName string, ids []uint64, newStatus int) error {
	fmt.Println("ids:", ids)
	db := e.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"process_status": newStatus,
		"updated_at":     time.Now().UTC(),
	}).Error
}

// UpdateProcessFail updates the process fail reason and count of the RawBridgeEvent.
func (e *RawBridgeEvent) UpdateProcessFail(tableName string, id uint64, reason string) error {
	db := e.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"process_fail_reason": reason,
		"process_fail_count":  gorm.Expr("process_fail_count + ?", 1),
		"process_status":      int(btypes.ProcessFailed),
		"updated_at":          time.Now().UTC(),
	}).Error
}
