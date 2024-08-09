package evm

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	yu_common "github.com/yu-org/yu/common"
)

func ConvertHashToYuHash(hash common.Hash) (yu_common.Hash, error) {
	var yuHash [yu_common.HashLen]byte
	if len(hash.Bytes()) == yu_common.HashLen {
		copy(yuHash[:], hash.Bytes())
		return yuHash, nil
	} else {
		return yu_common.Hash{}, errors.New(fmt.Sprintf("Expected hash to be 32 bytes long, but got %d bytes", len(hash.Bytes())))
	}
}
