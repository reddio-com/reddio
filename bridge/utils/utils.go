package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/joho/godotenv"
	yu_common "github.com/yu-org/yu/common"
)

type MessagePayloadType int

const (
	ETH MessagePayloadType = iota
	ERC20
	ERC721
	ERC1155
	RED
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
			// fmt.Println("log.Data ", log.Data)
			// fmt.Println("event ", event)
			// fmt.Println("Failed to UnpackIntoInterface", "err", err)
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
func LoadPrivateKey(envFilePath string) (string, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return "", err
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		return "", fmt.Errorf("PRIVATE_KEY not set in %s", envFilePath)
	}

	return privateKey, nil
}

func GenerateNonce() *big.Int {
	return big.NewInt(time.Now().UnixNano())
}

func NowUTC() time.Time {
	utc, _ := time.LoadLocation("")
	return time.Now().In(utc)
}

// ConvertStringToStringArray takes a string with values separated by commas and returns a slice of strings
func ConvertStringToStringArray(s string) []string {
	if s == "" {
		return []string{}
	}
	stringParts := strings.Split(s, ",")
	for i, part := range stringParts {
		stringParts[i] = strings.TrimSpace(part)
	}
	return stringParts
}

func ComputeMessageHash(payloadType uint32, payload []byte, nonce *big.Int) (common.Hash, error) {
	packedData, err := abi.Arguments{
		{Type: abi.Type{T: abi.UintTy, Size: 32}}, // Use UintTy with size 32 for uint32
		{Type: abi.Type{T: abi.BytesTy}},
		{Type: abi.Type{T: abi.UintTy, Size: 256}}, // Use UintTy with size 256 for *big.Int
	}.Pack(payloadType, payload, nonce)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	dataHash := crypto.Keccak256Hash(packedData)
	return dataHash, nil
}
