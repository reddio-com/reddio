package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
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
		loadedKey, err := LoadPrivateKey(envFilePath, "PRIVATE_KEY")
		assert.NoError(t, err)
		assert.Equal(t, privateKey, loadedKey)
	})
	t.Run("Success2", func(t *testing.T) {
		_, err := LoadPrivateKey("../relayer/.sepolia.env", "PRIVATE_KEY")
		assert.NoError(t, err)

	})
	// Test case: .env file does not exist
	t.Run("FileNotExist", func(t *testing.T) {
		_, err := LoadPrivateKey("nonexistent.env", "PRIVATE_KEY")
		assert.Error(t, err)
	})

	// Test case: RELAYER_PRIVATE_KEY not set in .env file
}
func TestComputeMessageHash(t *testing.T) {
	payloadType := uint32(4)
	payloadHex := "0000000000000000000000002655fc00139e0274dd0d84270fd80150b5f25426000000000000000000000000994314a99177eee8554fb0d0a246f3a3ea4ef56c000000000000000000000000994314a99177eee8554fb0d0a246f3a3ea4ef56c0000000000000000000000000000000000000000000000000de0b6b3a7640000"
	payload, err := hex.DecodeString(payloadHex)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	nonce := big.NewInt(1579416)

	expectedHash := common.HexToHash("0xb68fe48d80c53ad8794b9f8a147d14ca9f9e8f181a8a2eecd5d8e239a74b34fa")

	hash, err := ComputeMessageHash(payloadType, payload, nonce)
	if err != nil {
		t.Fatalf("ComputeMessageHash failed: %v", err)
	}
	if hash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}
func TestStorage(t *testing.T) {
	key := "child.bridge.core.storage"

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(key))
	result := hash.Sum(nil)

	address := fmt.Sprintf("%x", result)
	t.Logf("Address: %s", address)
	expectedAddress := "some_incorrect_address"
	if address != expectedAddress {
		t.Errorf("Expected address %s, but got %s", expectedAddress, address)
	}
}
