package utils

import (
	"fmt"
	"os"
	"testing"

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
