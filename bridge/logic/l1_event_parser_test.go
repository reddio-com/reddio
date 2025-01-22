package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/reddio-com/reddio/bridge/orm"
	btypes "github.com/reddio-com/reddio/bridge/types"
	"github.com/stretchr/testify/assert"
)

func TestETHParseL1SingleCrossChainEventLogs(t *testing.T) {
	parser := &L1EventParser{}

	payloadBase64 := "AAAAAAAAAAAAAAAAeIi3uES0sWwD+NrKzvfdoPUYhkUAAAAAAAAAAAAAAAB4iLe4RLSxbAP42srO992g9RiGRQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABk"
	t.Log("Payload Base64:", payloadBase64)
	payload, err := base64.StdEncoding.DecodeString(payloadBase64)
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}
	t.Log("Payload:", payload)
	payloadHex := fmt.Sprintf("0x%s", hex.EncodeToString(payload))
	t.Log("Payload Hex:", payloadHex)

	tx := types.NewTransaction(
		1,                          // nonce
		common.HexToAddress("0x0"), // to address
		big.NewInt(0),              // value
		21000,                      // gas limit
		big.NewInt(1),              // gas price
		nil,                        // data
	)
	msg := &orm.RawBridgeEvent{
		MessagePayloadType: int(btypes.ETH),
		MessagePayload:     string(payload),
	}

	l1DepositMessages, err := parser.ParseL1SingleRawBridgeEventToCrossChainMessage(context.Background(), msg, tx)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(btypes.ETH)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}
func TestREDParseL1SingleCrossChainEventLogs(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x000000000000000000000000b878927d79975bdb288ab53271f171534a49eb7d0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000064"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	tx := types.NewTransaction(
		1,                          // nonce
		common.HexToAddress("0x0"), // to address
		big.NewInt(0),              // value
		21000,                      // gas limit
		big.NewInt(1),              // gas price
		nil,                        // data
	)
	msg := &orm.RawBridgeEvent{
		MessagePayloadType: int(btypes.ETH),
		MessagePayload:     string(payload),
	}

	l1DepositMessages, err := parser.ParseL1SingleRawBridgeEventToCrossChainMessage(context.Background(), msg, tx)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0xB878927d79975BDb288ab53271f171534A49eb7D", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(btypes.RED)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}

func TestERC20ParseL1SingleCrossChainEventLogs(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000009627e313c18be25fc03100bbd3bf48743b4dee7000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000080000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000b577261707065642042544300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045742544300000000000000000000000000000000000000000000000000000000"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}
	tx := types.NewTransaction(
		1,                          // nonce
		common.HexToAddress("0x0"), // to address
		big.NewInt(0),              // value
		21000,                      // gas limit
		big.NewInt(1),              // gas price
		nil,                        // data
	)
	msg := &orm.RawBridgeEvent{
		MessagePayloadType: int(btypes.ETH),
		MessagePayload:     string(payload),
	}

	l1DepositMessages, err := parser.ParseL1SingleRawBridgeEventToCrossChainMessage(context.Background(), msg, tx)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0xF1E77FF9A4d4fc09CD955EfC44cB843617C73F23", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x0CC0cD4A9024A2d15BbEdd348Fbf7Cd69B5489bA", lastMessage.Sender)
	assert.Equal(t, "0x0CC0cD4A9024A2d15BbEdd348Fbf7Cd69B5489bA", lastMessage.Receiver)
	assert.Equal(t, int(uint32(btypes.ERC20)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}

func TestERC20ParseL1SingleCrossChainEventLogs2(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000009627e313c18be25fc03100bbd3bf48743b4dee7000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000080000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f518864500000000000000000000000000000000000000000000000000000000000003e8000000000000000000000000000000000000000000000000000000000000000b577261707065642042544300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045742544300000000000000000000000000000000000000000000000000000000"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	tx := types.NewTransaction(
		1,                          // nonce
		common.HexToAddress("0x0"), // to address
		big.NewInt(0),              // value
		21000,                      // gas limit
		big.NewInt(1),              // gas price
		nil,                        // data
	)
	msg := &orm.RawBridgeEvent{
		MessagePayloadType: int(btypes.ETH),
		MessagePayload:     string(payload),
	}

	l1DepositMessages, err := parser.ParseL1SingleRawBridgeEventToCrossChainMessage(context.Background(), msg, tx)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0x9627E313C18be25fC03100bbD3bf48743B4dee70", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(btypes.ERC20)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}

func TestERC20ParseL1SingleCrossChainEventLogs3(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000009627e313c18be25fc03100bbd3bf48743b4dee7000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000080000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000b577261707065642042544300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045742544300000000000000000000000000000000000000000000000000000000"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	tx := types.NewTransaction(
		1,                          // nonce
		common.HexToAddress("0x0"), // to address
		big.NewInt(0),              // value
		21000,                      // gas limit
		big.NewInt(1),              // gas price
		nil,                        // data
	)
	msg := &orm.RawBridgeEvent{
		MessagePayloadType: int(btypes.ETH),
		MessagePayload:     string(payload),
	}

	l1DepositMessages, err := parser.ParseL1SingleRawBridgeEventToCrossChainMessage(context.Background(), msg, tx)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0x9627E313C18be25fC03100bbD3bf48743B4dee71", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(btypes.ERC20)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}
