package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	yu_common "github.com/yu-org/yu/common"
)

// Loop Run the f func periodically.
func Loop(ctx context.Context, period time.Duration, f func()) {
	tick := time.NewTicker(period)
	defer tick.Stop()
	for ; ; <-tick.C {
		select {
		case <-ctx.Done():
			return
		default:
			f()
		}
	}
}

// UnpackLog unpacks a retrieved log into the provided output structure.
// @todo: add unit test.
func UnpackLog(c *abi.ABI, out interface{}, event string, log types.Log) error {
	if log.Topics[0] != c.Events[event].ID {
		return errors.New("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := c.UnpackIntoInterface(out, event, log.Data); err != nil {
			fmt.Println("log.Data ", log.Data)
			fmt.Println("event ", event)
			fmt.Println("Failed to UnpackIntoInterface", "err", err)
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range c.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

func ConvertHashToYuHash(hash common.Hash) (yu_common.Hash, error) {
	var yuHash [yu_common.HashLen]byte
	if len(hash.Bytes()) == yu_common.HashLen {
		copy(yuHash[:], hash.Bytes())
		return yuHash, nil
	} else {
		return yu_common.Hash{}, errors.New(fmt.Sprintf("Expected hash to be 32 bytes long, but got %d bytes", len(hash.Bytes())))
	}
}

func ConvertBigIntToUint256(b *big.Int) *uint256.Int {
	if b == nil {
		return nil
	}
	u, _ := uint256.FromBig(b)
	return u
}

func ObjToJson(obj interface{}) string {
	byt, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error marshalling obj to json: %v\n", err)
		return ""
	}
	return string(byt)
}
