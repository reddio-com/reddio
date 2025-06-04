package orm

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/evm"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// Gap represents a gap in MessageNonce.
type Gap struct {
	StartGap         int `json:"start_gap"`
	EndGap           int `json:"end_gap"`
	StartBlockNumber int `json:"start_block_number"`
	EndBlockNumber   int `json:"end_block_number"`
}

// BridgeEvents represents a bridge event.
type RawBridgeEvent struct {
	db                 *gorm.DB   `gorm:"column:-"`
	ID                 uint64     `json:"id" gorm:"column:id;primary_key;autoIncrement"` // primary key in the database
	EventType          int        `json:"event_type" gorm:"column:event_type"`           // 1.QueueTransaction(L1DepositMsgSent) 2.L2RelayedMessage(L2DepositMsgConsumed) 3.SentMessage(L2withdrawMsgSent) 4.L1RelayedMessage(L2DepositMsgConsumed)
	ChainID            int        `json:"chain_id" gorm:"column:chain_id"`
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
	TokenAddress       string     `json:"token_address" gorm:"column:token_address;type:varchar(100);"`
	TokenName          string     `json:"token_name" gorm:"column:token_name;type:varchar(100);"`
	TokenSymbol        string     `json:"token_symbol" gorm:"column:token_symbol;type:varchar(100);"`
	Decimals           string     `json:"decimals" gorm:"column:decimals;type:varchar(10);"`                     // token decimals
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
	CheckStatus        int        `json:"check_status" gorm:"column:check_status"`                                 // 1.checkedStep1 2.checkedStep2
	CheckFailReason    string     `json:"check_fail_reason" gorm:"column:check_fail_reason;type:varchar(256)"`     // reason for check failure
}

// NewBridgeEvents creates a new instance of BridgeEvents.
func NewRawBridgeEvent(db *gorm.DB, cfg *evm.GethConfig) *RawBridgeEvent {
	err := db.Table(cfg.L1_RawBridgeEventsTableName).AutoMigrate(&RawBridgeEvent{})
	if err != nil {
		log.Fatal("failed to AutoMigrate db", "err", err)
	}
	err = db.Table(cfg.L2_RawBridgeEventsTableName).AutoMigrate(&RawBridgeEvent{})
	if err != nil {
		log.Fatal("failed to AutoMigrate db", "err", err)
	}
	return &RawBridgeEvent{db: db}
}

/***************
 *    Read     *
 ***************/
func (b *RawBridgeEvent) QueryUnProcessedBridgeEventsByEventType(ctx context.Context, tableName string, eventType int, limit int) ([]*RawBridgeEvent, error) {
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)

	var bridgeEvents []*RawBridgeEvent
	if err := db.Where("process_status = ? AND event_type = ? ", btypes.UnProcessed, eventType).
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

func (r *RawBridgeEvent) CountEventsByMessageNonceRange(tableName string, eventType, startNonce, endNonce int) (int64, error) {
	var count int64
	err := r.db.Table(tableName).Model(&RawBridgeEvent{}).
		Where("event_type = ? AND message_nonce BETWEEN ? AND ?", eventType, startNonce, endNonce).
		Count(&count).Error
	return count, err
}

// FindMessageNonceGaps finds gaps in MessageNonce between the specified range.
func (r *RawBridgeEvent) FindMessageNonceGaps(tableName string, eventType, startNonce, endNonce int) ([]Gap, error) {
	var gaps []Gap
	db := r.db
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)
	// SQL query to find gaps within the specified range
	query := `
        SELECT t1.message_nonce + 1 AS start_gap, MIN(t2.message_nonce) - 1 AS end_gap, t1.block_number AS start_block_number, MIN(t2.block_number) AS end_block_number
        FROM ` + tableName + ` t1
        JOIN ` + tableName + ` t2 ON t1.message_nonce < t2.message_nonce
        WHERE t1.event_type = ? AND t2.event_type = ? AND t1.message_nonce BETWEEN ? AND ? AND t2.message_nonce BETWEEN ? AND ?
        GROUP BY t1.message_nonce, t1.block_number
        HAVING start_gap < MIN(t2.message_nonce);
    `
	err := r.db.Raw(query, eventType, eventType, startNonce, endNonce, startNonce, endNonce).Scan(&gaps).Error
	if err != nil {
		return nil, err
	}

	// Check for head gap
	var firstEvent RawBridgeEvent
	err = r.db.Table(tableName).Where("event_type = ? AND message_nonce >= ?", eventType, startNonce).
		Order("message_nonce ASC").First(&firstEvent).Error
	if err == nil && firstEvent.MessageNonce > startNonce {
		var prevEvent RawBridgeEvent
		err = r.db.Table(tableName).Where("event_type = ? AND message_nonce < ?", eventType, startNonce).
			Order("message_nonce DESC").First(&prevEvent).Error
		startBlockNumber := 0
		if err == nil {
			startBlockNumber = int(prevEvent.BlockNumber)
		} else {
			logrus.Errorf("failed to get prevEvent: %v", err)
			return nil, err
		}
		gaps = append([]Gap{{StartGap: startNonce, EndGap: firstEvent.MessageNonce - 1, StartBlockNumber: startBlockNumber, EndBlockNumber: int(firstEvent.BlockNumber)}}, gaps...)
	}

	// Check for tail gap
	var lastEvent RawBridgeEvent
	err = r.db.Table(tableName).Where("event_type = ? AND message_nonce <= ?", eventType, endNonce).
		Order("message_nonce DESC").First(&lastEvent).Error
	if err == nil && lastEvent.MessageNonce < endNonce {
		var nextEvent RawBridgeEvent
		err = r.db.Table(tableName).Where("event_type = ? AND message_nonce > ?", eventType, endNonce).
			Order("message_nonce ASC").First(&nextEvent).Error
		endBlockNumber := 0
		if err == nil {
			endBlockNumber = int(nextEvent.BlockNumber)
		} else {
			logrus.Errorf("failed to get nextEvent: %v", err)
			return nil, err
		}
		gaps = append(gaps, Gap{StartGap: lastEvent.MessageNonce + 1, EndGap: endNonce, StartBlockNumber: int(lastEvent.BlockNumber), EndBlockNumber: endBlockNumber})
	}

	return gaps, nil
}

// GetMinNonceByCheckStatus gets the minimum MessageNonce by check_status.
func (r *RawBridgeEvent) GetMinNonceByCheckStatus(tableName string, eventType, checkStatus int) (int, error) {
	var minNonce sql.NullInt64
	twoHoursAgo := time.Now().Add(-2 * time.Hour).Unix()

	err := r.db.Table(tableName).
		Where("event_type = ? AND check_status = ? AND timestamp < ?", eventType, checkStatus, twoHoursAgo).
		Select("MIN(message_nonce)").
		Scan(&minNonce).Error
	if err != nil {
		return -1, err
	}
	if minNonce.Valid {
		return int(minNonce.Int64), nil
	}
	return -1, nil
}

// GetMaxNonceByCheckStatus gets the maximum MessageNonce by check_status.
func (r *RawBridgeEvent) GetMaxNonceByCheckStatus(tableName string, eventType, checkStatus int) (int, error) {
	var maxNonce sql.NullInt64
	twoHoursAgo := time.Now().Add(-2 * time.Hour).Unix()

	err := r.db.Table(tableName).
		Where("event_type = ? AND check_status = ? AND timestamp < ?", eventType, checkStatus, twoHoursAgo).
		Select("MAX(message_nonce)").
		Scan(&maxNonce).Error
	if err != nil {
		return 0, err
	}
	if maxNonce.Valid {
		return int(maxNonce.Int64), nil
	}
	return 0, nil
}

// GetEventsByMessageNonceRange gets events by message nonce range.
func (r *RawBridgeEvent) GetEventsByMessageNonceRange(tableName string, eventType, startNonce, endNonce int) ([]RawBridgeEvent, error) {
	var events []RawBridgeEvent
	err := r.db.Table(tableName).Where(" event_type = ? AND message_nonce BETWEEN ? AND ?", eventType, startNonce, endNonce).Find(&events).Error
	return events, err
}

func (r *RawBridgeEvent) GetMaxBlockNumber(ctx context.Context, tableName string) (uint64, error) {
	var maxBlockNumber uint64
	db := r.db.WithContext(ctx)
	err := db.Table(tableName).Select("COALESCE(MAX(block_number), 0)").Scan(&maxBlockNumber).Error
	if err != nil {
		return 0, err
	}
	return maxBlockNumber, nil
}

/****************
 *    Write     *
 ****************/
// InsertBridgeEvents inserts a new BridgeEvents record into the database.
func (b *RawBridgeEvent) InsertRawBridgeEvents(ctx context.Context, tableName string, bridgeEvents []*RawBridgeEvent) error {
	if len(bridgeEvents) == 0 {
		return nil
	}
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)
	tx := db.Begin()
	if tx.Error != nil {
		logrus.Errorf("Failed to begin transaction: %v", tx.Error)
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	for _, event := range bridgeEvents {
		result := tx.Create(event)
		if result.Error != nil {
			if isDuplicateEntryError(result.Error) {
				logrus.Errorf("Message with hash %s already exists, skipping insert.\n", event.MessageHash)
				continue
			}
			logrus.Errorf("Failed to insert message: %v", result.Error)
			tx.Rollback()
			return fmt.Errorf("failed to insert message, error: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			logrus.Warnf("No rows affected for message with hash %s, skipping insert.\n", event.MessageHash)
			continue
		}
	}

	if err := tx.Commit().Error; err != nil {
		logrus.Errorf("Failed to commit transaction: %v", err)
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logrus.Infof("Transaction committed successfully for table %s", tableName)
	return nil
}

func (b *RawBridgeEvent) InsertRawBridgeEventsFromCheckStep1(ctx context.Context, tableName string, bridgeEvents []*RawBridgeEvent) error {
	if len(bridgeEvents) == 0 {
		return nil
	}
	db := b.db
	db = db.WithContext(ctx)
	db = db.Model(&RawBridgeEvent{})
	db = db.Table(tableName)
	//fmt.Println("InsertRawBridgeEvents: tableName: , bridgeEvents:", tableName, bridgeEvents)
	return db.Transaction(func(tx *gorm.DB) error {
		for _, event := range bridgeEvents {
			event.Remark = "checkStep1 inserted"
			result := tx.Create(event)
			if result.Error != nil {
				if isDuplicateEntryError(result.Error) {
					logrus.Errorf("Message with hash %s already exists, skipping insert.\n", event.MessageHash)
					continue
				}
				return fmt.Errorf("failed to insert message, error: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				logrus.Warnf("No rows affected for message with hash %s, skipping insert.\n", event.MessageHash)
				continue
			}
		}
		return nil
	})
}

func isDuplicateEntryError(err error) bool {
	return strings.Contains(err.Error(), "Error 1062")
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

// UpdateCheckStatus updates the CheckStatus of the RawBridgeEvent.
func (r *RawBridgeEvent) UpdateCheckStatus(tableName string, id uint64, newStatus int) error {
	db := r.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"check_status": newStatus,
		"updated_at":   time.Now().UTC(),
	}).Error
}

// UpdateCheckStatusByNonceRange updates the CheckStatus of RawBridgeEvents within a range of MessageNonce and eventType.
func (r *RawBridgeEvent) UpdateCheckStatusByNonceRange(tableName string, eventType, startNonce, endNonce, newStatus int) error {
	db := r.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("event_type = ? AND message_nonce BETWEEN ? AND ?", eventType, startNonce, endNonce).Updates(map[string]interface{}{
		"check_status": newStatus,
		"updated_at":   time.Now().UTC(),
	}).Error
}
func (r *RawBridgeEvent) UpdateCheckFailReason(tableName string, id uint64, newStatus int, reason string) error {
	db := r.db.Table(tableName)
	return db.Model(&RawBridgeEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"check_status":      newStatus,
		"check_fail_reason": reason,
		"updated_at":        time.Now().UTC(),
	}).Error
}
