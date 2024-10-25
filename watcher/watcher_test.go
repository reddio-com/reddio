package watcher

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/reddio-com/reddio/evm"
	"github.com/sirupsen/logrus"
	yutypes "github.com/yu-org/yu/core/types"
)

// Mock configuration for testing
var mockConfig = &evm.GethConfig{
	L1ClientAddress: "http://localhost:8545",
	EnableL1Client:  false,
}

// Mock Solidity instance
var mockSolidity = &evm.Solidity{}

func TestInitChain(t *testing.T) {
	// Create a new Watcher instance
	watcher := &Watcher{
		cfg:      mockConfig,
		Solidity: mockSolidity,
	}

	// Create a mock block
	block := &yutypes.Block{}

	// Call InitChain method
	watcher.InitChain(block)

	// Check if evmBridgeDB is set
	if watcher.evmBridgeDB == nil {
		t.Fatal("evmBridgeDB is not set")
	}

	// Check if evmBridgeDB is a valid PebbleDB instance
	_, closer, err := watcher.evmBridgeDB.Get([]byte("test"))
	if err != nil && err != pebble.ErrNotFound {
		t.Fatalf("evmBridgeDB is not a valid PebbleDB instance: %v", err)
	}
	if closer != nil {
		closer.Close()
	}

	// Test data to store and retrieve
	key := []byte("test-key")
	value := []byte("test-value")

	// Store data
	err = watcher.evmBridgeDB.Set(key, value, pebble.Sync)
	if err != nil {
		t.Fatalf("failed to set data in evmBridgeDB: %v", err)
	}

	// Retrieve data
	retrievedValue, closer, err := watcher.evmBridgeDB.Get(key)
	if err != nil {
		t.Fatalf("failed to get data from evmBridgeDB: %v", err)
	}
	t.Logf("retrieved value: %s", retrievedValue)
	if closer != nil {
		closer.Close()
	}

	// Check if the retrieved value matches the stored value
	if string(retrievedValue) != string(value) {
		t.Fatalf("retrieved value does not match stored value: got %s, want %s", retrievedValue, value)
	}

	// Close the database
	err = watcher.evmBridgeDB.Close()
	if err != nil {
		t.Fatalf("failed to close evmBridgeDB: %v", err)
	}

	logrus.Info("TestInitChain passed")
}
