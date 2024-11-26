package logic

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/reddio-com/reddio/bridge/utils"
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

	msg := &contract.ParentBridgeCoreFacetDownwardMessage{
		PayloadType: uint32(utils.ETH),
		Payload:     payload,
	}
	l1DepositMessages, err := parser.ParseL1SingleCrossChainPayload(context.Background(), msg)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(utils.ETH)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}
func TestREDParseL1SingleCrossChainEventLogs(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x000000000000000000000000b878927d79975bdb288ab53271f171534a49eb7d0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000000000000000000000000000000000000000000064"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	msg := &contract.ParentBridgeCoreFacetDownwardMessage{
		PayloadType: uint32(utils.RED),
		Payload:     payload,
	}
	l1DepositMessages, err := parser.ParseL1SingleCrossChainPayload(context.Background(), msg)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0xB878927d79975BDb288ab53271f171534A49eb7D", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Sender)
	assert.Equal(t, "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645", lastMessage.Receiver)
	assert.Equal(t, int(uint32(utils.RED)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}

func TestERC20ParseL1SingleCrossChainEventLogs(t *testing.T) {
	parser := &L1EventParser{}

	payloadHex := "0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000f1e77ff9a4d4fc09cd955efc44cb843617c73f2300000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000cc0cd4a9024a2d15bbedd348fbf7cd69b5489ba0000000000000000000000000cc0cd4a9024a2d15bbedd348fbf7cd69b5489ba0000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000954455354546f6b656e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000025454000000000000000000000000000000000000000000000000000000000000"
	payload, err := hex.DecodeString(payloadHex[2:])
	if err != nil {
		t.Fatalf("Failed to decode payload: %v", err)
	}

	msg := &contract.ParentBridgeCoreFacetDownwardMessage{
		PayloadType: uint32(utils.ERC20),
		Payload:     payload,
	}
	l1DepositMessages, err := parser.ParseL1SingleCrossChainPayload(context.Background(), msg)

	assert.NoError(t, err)
	assert.NotNil(t, l1DepositMessages)
	assert.Len(t, l1DepositMessages, 1)

	lastMessage := l1DepositMessages[0]
	assert.Equal(t, "0xF1E77FF9A4d4fc09CD955EfC44cB843617C73F23", lastMessage.L1TokenAddress)
	assert.Equal(t, "0x0CC0cD4A9024A2d15BbEdd348Fbf7Cd69B5489bA", lastMessage.Sender)
	assert.Equal(t, "0x0CC0cD4A9024A2d15BbEdd348Fbf7Cd69B5489bA", lastMessage.Receiver)
	assert.Equal(t, int(uint32(utils.ERC20)), lastMessage.TokenType)
	assert.Equal(t, "100", lastMessage.TokenAmounts)
}
