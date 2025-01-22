package relayer

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/reddio-com/reddio/bridge/contract"
	"github.com/stretchr/testify/assert"
)

/**
 * TestGenerateUpwardMessageMultiSignatures tests the generateUpwardMessageMultiSignatures function.
 *
 * Arrange:
 * - Create a mock list of upwardMessages.
 * - Create a mock list of privateKeys.
 *
 * Act:
 * - Call the generateUpwardMessageMultiSignatures function to generate signatures.
 *
 * Assert:
 * - Ensure no errors occurred.
 * - Ensure the number of generated signatures matches the expected count.
 * - Verify each signature is correct.
 */
func TestGenerateUpwardMessageMultiSignatures(t *testing.T) {
	//Arrange:
	// Create mock upward messages
	upwardMessages := []contract.UpwardMessage{
		{
			PayloadType: 0,
			Payload:     hexDecode("0x0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f5188645000000000000000000000000000000000000000000000000000000000000006f"),
			Nonce:       big.NewInt(0),
		},
		{
			PayloadType: 0,
			Payload:     hexDecode("0x0000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f51886450000000000000000000000007888b7b844b4b16c03f8dacacef7dda0f5188645000000000000000000000000000000000000000000000000000000000000006f"),
			Nonce:       big.NewInt(0),
		},
	}
	// Create mock private keys
	privateKeys := []string{
		"32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff",
		"78740b0ee70f3e8fda88f90da06d3852043c70235b6cd8b3a2337ddd37423dc5",
	}

	// Call the function to test
	signatures, err := generateUpwardMessageMultiSignatures(upwardMessages, privateKeys)
	if err != nil {
		t.Fatalf("Failed to generate multi-signatures: %v", err)
	}

	assert.NoError(t, err)

	// Assert the length of signatures
	assert.Equal(t, len(privateKeys), len(signatures))

	for i, sig := range signatures {
		privateKey, err := crypto.HexToECDSA(privateKeys[i])
		assert.NoError(t, err)
		t.Logf("signature[%d]: %x\n", i, sig)
		dataHash, err := generateUpwardMessageToHash(upwardMessages)
		if err != nil {
			t.Fatalf("Failed to generate data hash: %v", err)
		}
		t.Logf("dataHash: %s\n", dataHash)
		pubKey, err := crypto.Ecrecover(dataHash.Bytes(), sig)
		if err != nil {
			t.Fatalf("Failed to recover public key: %v", err)
		}

		if pubKey == nil {
			t.Fatalf("Recovered public key is nil")
		}
		// Convert public key to address
		publicKeyECDSA, err := crypto.UnmarshalPubkey(pubKey)
		if err != nil {
			t.Fatal(err)
		}
		recoveredAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
		expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

		assert.Equal(t, expectedAddr, recoveredAddr)
	}
}

// hexDecode decodes a hex string to a byte slice
func hexDecode(hexStr string) []byte {
	bytes, err := hex.DecodeString(hexStr[2:]) // Skip the "0x" prefix
	if err != nil {
		panic(err)
	}
	return bytes
}
