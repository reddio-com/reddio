package evm

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	yu_common "github.com/yu-org/yu/common"
)

func TestConvertHashToYuHash(t *testing.T) {
	hash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000006835c46")

	yuHash, err := ConvertHashToYuHash(hash)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if len(yuHash) != yu_common.HashLen {
		t.Fatalf("Expected hash length %d, but got %d", yu_common.HashLen, len(yuHash))
	}

	invalidHash := common.HexToHash("0x123")
	yuHash, err = ConvertHashToYuHash(invalidHash)
	if err == nil {
		t.Log("yuHash", yuHash.String())
		t.Fatalf("Expected error, but got nil")
	}

	leadingZerosHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000123")
	yuHash, err = ConvertHashToYuHash(leadingZerosHash)
	if err == nil {
		t.Log("yuHash", yuHash.String())
		t.Fatalf("Expected error for leading zeros, but got nil")
	}
}

func TestValidateTxHash(t *testing.T) {
	validHash := "0x3d169ccf2d92ede8d507540c40a158cf7bf0aff09239b1966d5dda6a01ad965a"
	if !ValidateTxHash(validHash) {
		t.Fatalf("Expected valid hash, but got invalid: %s", validHash)
	}

	invalidHashShort := "0x123"
	if ValidateTxHash(invalidHashShort) {
		t.Fatalf("Expected invalid hash due to short length, but got valid: %s", invalidHashShort)
	}

	invalidHashNoPrefix := "0000000000000000000000000000000000000000000000000000000006835c46"
	if ValidateTxHash(invalidHashNoPrefix) {
		t.Fatalf("Expected invalid hash due to missing prefix, but got valid: %s", invalidHashNoPrefix)
	}

	allZeroHash := "0x0000000000000000000000000000000000000000000000000000000000000000"
	if ValidateTxHash(allZeroHash) {
		t.Fatalf("Expected invalid hash due to all zeros, but got valid: %s", allZeroHash)
	}

	leadingZerosHash := "0x0000000000000000000000000000000000000000000000000000000000000123"
	if ValidateTxHash(leadingZerosHash) {
		t.Fatalf("Expected invalid hash due to leading zeros, but got valid: %s", leadingZerosHash)
	}
}
