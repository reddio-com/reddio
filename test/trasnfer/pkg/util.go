package pkg

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func generatePrivateKey() (string, string) {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", ""
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return hexutil.Encode(privateKeyBytes)[2:], address
}

func sendRequest(hostAddress string, dataString string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s", hostAddress), bytes.NewBuffer([]byte(dataString)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code err: %v", resp.StatusCode)
	}
	return data, nil
}
