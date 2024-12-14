package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// TestLoadPrivateKey tests the loadPrivateKey function
func TestLoadPrivateKey(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	t.Logf("Current working directory: %s", wd)
	// Create a temporary .env file
	envFilePath := ".env.test"
	privateKey := "test_private_key"
	envContent := fmt.Sprintf("PRIVATE_KEY=%s\n", privateKey)
	err = os.WriteFile(envFilePath, []byte(envContent), 0644)
	assert.NoError(t, err)

	// Ensure the temporary .env file is removed after the test
	defer os.Remove(envFilePath)

	// Test case: Successfully load private key
	t.Run("Success", func(t *testing.T) {
		loadedKey, err := LoadPrivateKey(envFilePath)
		assert.NoError(t, err)
		assert.Equal(t, privateKey, loadedKey)
	})
	t.Run("Success2", func(t *testing.T) {
		_, err := LoadPrivateKey("../relayer/.sepolia.env")
		assert.NoError(t, err)

	})
	// Test case: .env file does not exist
	t.Run("FileNotExist", func(t *testing.T) {
		_, err := LoadPrivateKey("nonexistent.env")
		assert.Error(t, err)
	})

	// Test case: RELAYER_PRIVATE_KEY not set in .env file
}
func TestComputeMessageHash(t *testing.T) {
	payloadType := uint32(4)
	payloadHex := "000000000000000000000000b878927d79975bdb288ab53271f171534a49eb7d000000000000000000000000a90381616eebc94d89b11afde57b869705626968000000000000000000000000a90381616eebc94d89b11afde57b8697056269680000000000000000000000000000000000000000000000000de0b6b3a7640000"
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	nonce := big.NewInt(1733313681151894329)

	expectedHash := common.HexToHash("0x4c49d7969ff27718263e07f4a9c89d82c65a667191c1c93a2b0785df4bf7172a")

	hash, err := ComputeMessageHash(payloadType, payload, nonce)
	if err != nil {
		t.Fatalf("ComputeMessageHash failed: %v", err)
	}
	fmt.Println("hash:", hash)
	if hash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}
