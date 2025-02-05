package orm

import (
	"context"
	"database/sql"
	"testing"
	"time"

	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/reddio-com/reddio/bridge/utils/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var MockConfig = &database.Config{
	DSN:        "testuser:123456@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
	DriverName: "mysql",
	MaxOpenNum: 10,
	MaxIdleNum: 5,
}

func MockPing(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	return sqlDB, sqlDB.Ping()
}

func TestCreateAndGetCrossMessage(t *testing.T) {
	db, err := database.InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()

	if err := db.AutoMigrate(&CrossMessage{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	crossMessage := CrossMessage{
		MessageType:        1,
		TxStatus:           1,
		TokenType:          1,
		Sender:             "sender",
		Receiver:           "receiver",
		TxType:             1,
		MessageHash:        "message_hash",
		L1TxHash:           "l1_tx_hash",
		L2TxHash:           "l2_tx_hash",
		L1BlockNumber:      100,
		L2BlockNumber:      200,
		L1TokenAddress:     "l1_token_address",
		L2TokenAddress:     "l2_token_address",
		TokenIDs:           "token_ids",
		TokenAmounts:       "token_amounts",
		BlockTimestamp:     1234567890,
		MessagePayloadType: 1,
		MessagePayload:     "payload",
		MessageFrom:        "message_from",
		MessageTo:          "message_to",
		MessageValue:       "message_value",
		MessageNonce:       "1",
		MultiSignProof:     "multisign_proof",
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := db.Create(&crossMessage).Error; err != nil {
		t.Fatalf("Failed to create cross message: %v", err)
	}

	var retrievedCrossMessage CrossMessage
	if err := db.First(&retrievedCrossMessage, crossMessage.ID).Error; err != nil {
		t.Fatalf("Failed to get cross message: %v", err)
	}

	if retrievedCrossMessage.MessageHash != crossMessage.MessageHash {
		t.Errorf("Expected message hash to be %s, got %s", crossMessage.MessageHash, retrievedCrossMessage.MessageHash)
	}
	if retrievedCrossMessage.Sender != crossMessage.Sender {
		t.Errorf("Expected sender to be %s, got %s", crossMessage.Sender, retrievedCrossMessage.Sender)
	}
	if retrievedCrossMessage.Receiver != crossMessage.Receiver {
		t.Errorf("Expected receiver to be %s, got %s", crossMessage.Receiver, retrievedCrossMessage.Receiver)
	}
}

func TestUpsertCrossMessage(t *testing.T) {
	db, err := database.InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()

	if err := db.AutoMigrate(&CrossMessage{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	crossMessage := CrossMessage{
		ID:                 1,
		MessageType:        1,
		TxStatus:           1,
		TokenType:          1,
		Sender:             "sender",
		Receiver:           "receiver",
		MessageHash:        "message_hash",
		L1TxHash:           "l1_tx_hash",
		L2TxHash:           "l2_tx_hash",
		L1BlockNumber:      100,
		L2BlockNumber:      200,
		L1TokenAddress:     "l1_token_address",
		L2TokenAddress:     "l2_token_address",
		TokenIDs:           "token_ids",
		TokenAmounts:       "token_amounts",
		BlockTimestamp:     1234567890,
		MessagePayloadType: 1,
		MessagePayload:     "payload",
		MessageFrom:        "message_from",
		MessageTo:          "message_to",
		MessageValue:       "message_value",
		MessageNonce:       "1",
		MultiSignProof:     "multisign_proof",
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	if err := db.Clauses(clause.OnConflict{
		UpdateAll: true,
		DoUpdates: clause.AssignmentColumns([]string{"sender", "receiver", "token_type", "l2_block_number", "l2_tx_hash", "l1_token_address", "l2_token_address", "token_ids", "token_amounts", "message_type", "block_timestamp", "message_payloadtype", "message_payload", "message_from", "message_to", "message_value", "message_data", "message_nonce", "multisign_proof"}),
	}).Create(&crossMessage).Error; err != nil {
		t.Fatalf("Failed to upsert cross message: %v", err)
	}

	var retrievedCrossMessage CrossMessage
	if err := db.First(&retrievedCrossMessage, crossMessage.ID).Error; err != nil {
		t.Fatalf("Failed to get cross message: %v", err)
	}

	if retrievedCrossMessage.MessageHash != crossMessage.MessageHash {
		t.Errorf("Expected message hash to be %s, got %s", crossMessage.MessageHash, retrievedCrossMessage.MessageHash)
	}
	if retrievedCrossMessage.Sender != crossMessage.Sender {
		t.Errorf("Expected sender to be %s, got %s", crossMessage.Sender, retrievedCrossMessage.Sender)
	}
	if retrievedCrossMessage.Receiver != crossMessage.Receiver {
		t.Errorf("Expected receiver to be %s, got %s", crossMessage.Receiver, retrievedCrossMessage.Receiver)
	}
}
func TestInsertOrUpdateL2Messages(t *testing.T) {
	db, err := database.InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			t.Fatalf("Failed to close database: %v", err)
		}
	}()

	if err := db.AutoMigrate(&CrossMessage{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	crossMessage := &CrossMessage{
		MessageType:        1,
		TxStatus:           1,
		TokenType:          1,
		Sender:             "sender",
		Receiver:           "receiver",
		MessageHash:        "message_hash",
		L1TxHash:           "l1_tx_hash",
		L2TxHash:           "l2_tx_hash",
		L1BlockNumber:      100,
		L2BlockNumber:      200,
		L1TokenAddress:     "l1_token_address",
		L2TokenAddress:     "l2_token_address",
		TokenIDs:           "token_ids",
		TokenAmounts:       "token_amounts",
		BlockTimestamp:     1234567890,
		MessagePayloadType: 1,
		MessagePayload:     "payload",
		MessageFrom:        "message_from",
		MessageTo:          "message_to",
		MessageValue:       "message_value",
		MessageNonce:       "1",
		MultiSignProof:     "multisign_proof",
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	c := &CrossMessage{db: db}
	err = c.InsertOrUpdateCrossMessages(context.Background(), []*CrossMessage{crossMessage})
	if err != nil {
		t.Fatalf("Failed to insert or update cross message: %v", err)
	}

	var retrievedCrossMessage CrossMessage
	if err := db.Where("message_hash = ?", crossMessage.MessageHash).First(&retrievedCrossMessage).Error; err != nil {
		t.Fatalf("Failed to retrieve cross message: %v", err)
	}

	if retrievedCrossMessage.MessageHash != crossMessage.MessageHash {
		t.Errorf("Expected message hash to be %s, got %s", crossMessage.MessageHash, retrievedCrossMessage.MessageHash)
	}
	if retrievedCrossMessage.Sender != crossMessage.Sender {
		t.Errorf("Expected sender to be %s, got %s", crossMessage.Sender, retrievedCrossMessage.Sender)
	}
	if retrievedCrossMessage.Receiver != crossMessage.Receiver {
		t.Errorf("Expected receiver to be %s, got %s", crossMessage.Receiver, retrievedCrossMessage.Receiver)
	}
}

func TestGetL2UnclaimedWithdrawalsByAddress(t *testing.T) {
	db, err := database.InitDB(MockConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.AutoMigrate(&CrossMessage{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	sender := "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645"
	// crossMessages := []*CrossMessage{
	// 	{
	// 		MessageType:        1,
	// 		TxStatus:           1,
	// 		TokenType:          1,
	// 		Sender:             "sender",
	// 		Receiver:           "receiver",
	// 		MessageHash:        "message_hash1",
	// 		L1TxHash:           "l1_tx_hash",
	// 		L2TxHash:           "l2_tx_hash",
	// 		L1BlockNumber:      100,
	// 		L2BlockNumber:      200,
	// 		L1TokenAddress:     "l1_token_address",
	// 		L2TokenAddress:     "l2_token_address",
	// 		TokenIDs:           "token_ids",
	// 		TokenAmounts:       "token_amounts",
	// 		BlockTimestamp:     1234567890,
	// 		MessagePayloadType: 1,
	// 		MessagePayload:     "payload",
	// 		MessageFrom:        "sender",
	// 		MessageTo:          "receiver",
	// 		MessageValue:       "message_value",
	// 		MessageNonce:       "1",
	// 		MultiSignProof:     "multisign_proof",
	// 		CreatedAt:          time.Now().UTC(),
	// 		UpdatedAt:          time.Now().UTC(),
	// 	},
	// 	{
	// 		MessageType:        2,
	// 		TxStatus:           0,
	// 		TokenType:          0,
	// 		Sender:             sender,
	// 		Receiver:           sender,
	// 		MessageHash:        "0x79df0b41ed1d6d0f2b2748da13849fad7d140e41e8b87434511286706ac64fb7",
	// 		L1TxHash:           "",
	// 		L2TxHash:           "0x95cf843c68af0db5ccfb19187a9c661b6bc46ee1b27c788fcf8922d962dbd2a3",
	// 		L1BlockNumber:      0,
	// 		L2BlockNumber:      44,
	// 		L1TokenAddress:     "",
	// 		L2TokenAddress:     "",
	// 		TokenIDs:           "",
	// 		TokenAmounts:       "50",
	// 		BlockTimestamp:     0,
	// 		MessagePayloadType: 0,
	// 		MessagePayload:     "0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000032",
	// 		MessageFrom:        sender,
	// 		MessageTo:          sender,
	// 		MessageValue:       "50",
	// 		MessageNonce:       "1733120884468899841",
	// 		MultiSignProof:     "0x5d1376022cd357dc9c830ffcc944bf9b8458fc3d1acc119f77b0bdcea3c4a2e65f589282df235243aa492c070471f7b5c58ced8dd0d3e51819c2d6f216140f1801",
	// 		CreatedAt:          time.Now().UTC(),
	// 		UpdatedAt:          time.Now().UTC(),
	// 	},
	// 	{
	// 		MessageType:        2,
	// 		TxStatus:           0,
	// 		TokenType:          0,
	// 		Sender:             sender,
	// 		Receiver:           sender,
	// 		MessageHash:        "0xcf17b5dc50789e18aff92dad8ccb4279271b1c90ad277e8b0c5aa87aec1483c4",
	// 		L1TxHash:           "",
	// 		L2TxHash:           "0x63ca588b0d2d7965315d065323ca5640e55cedd21bebe9bb19507ef4e264eda2",
	// 		L1BlockNumber:      0,
	// 		L2BlockNumber:      90,
	// 		L1TokenAddress:     "",
	// 		L2TokenAddress:     "",
	// 		TokenIDs:           "",
	// 		TokenAmounts:       "50",
	// 		BlockTimestamp:     0,
	// 		MessagePayloadType: 0,
	// 		MessagePayload:     "0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000032",
	// 		MessageFrom:        sender,
	// 		MessageTo:          sender,
	// 		MessageValue:       "50",
	// 		MessageNonce:       "1733121029872386955",
	// 		MultiSignProof:     "0x4ed471902c17c533f4a5dedb531bc4fb2a8b5e52c615fabca1916ebc2103476539a6f1eb86e00497c0666b0c5e6a4dccfd48c4825a3ca9d7d86a47011f677cc201",
	// 		CreatedAt:          time.Now().UTC(),
	// 		UpdatedAt:          time.Now().UTC(),
	// 	},
	// }

	// if err := db.Create(&crossMessages).Error; err != nil {
	// 	t.Fatalf("Failed to create cross messages: %v", err)
	// }

	c := &CrossMessage{db: db}
	messages, total, err := c.GetL2UnclaimedWithdrawalsByAddress(context.Background(), sender, 1, 2)
	if err != nil {
		t.Fatalf("Failed to get L2 unclaimed withdrawal messages: %v", err)
	}
	t.Log("total", total)
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	for _, msg := range messages {
		if msg.MessageFrom != sender {
			t.Errorf("Expected message_from to be %s, got %s", sender, msg.MessageFrom)
		}
		if msg.TxStatus != int(btypes.TxStatusTypeSent) {
			t.Errorf("Expected tx_status to be %d, got %d", int(btypes.TxStatusTypeSent), msg.TxStatus)
		}
		if msg.MessageType != int(btypes.MessageTypeL2SentMessage) {
			t.Errorf("Expected message_type to be %d, got %d", int(btypes.MessageTypeL2SentMessage), msg.MessageType)
		}
	}
	t.Log(messages)
	//t.Error("messages：", messages[0].L2TxHash)
}
