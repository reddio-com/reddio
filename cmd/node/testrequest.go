package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
	
	"github.com/reddio-com/reddio/evm"
)

// How To Use
// Add this line in main.go (before `app.StartupChain()`)
// go testSendTransaction(gethCfg, true)

func testEthCall(exit bool) {
	time.Sleep(5 * time.Second)
	requestBody := `{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_call",
		"params": [{
			"from": "0x123456789abcdef123456789abcdef123456789a",
			"to": "0x9d7bA953587B87c474a10beb65809Ea489F026bD",
			"data": "0x70a082310000000000000000000000006E0d01A76C3Cf4288372a29124A26D4353EE51BE"
		}, "latest"]
	}`
	sendRequest(requestBody)

	if exit {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func TestSendTransaction(gethCfg *evm.GethConfig, exit bool) {
	// A random private key. address = 0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02
	privateKeyStr := "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
	nonce := uint64(0)
	to := common.HexToAddress("0x2Efe24c33f049Ffec693ec1D809A45Fff14e9527")
	amount := big.NewInt(1)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	data := []byte{}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    amount,
		Data:     data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		logrus.Fatal(err)
	}

	// signer := types.MakeSigner(gethCfg, new(big.Int).SetUint64(uint64(block.Height)), block.Timestamp)

	signer := types.LatestSigner(gethCfg.ChainConfig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("SignedTx = %+v", signedTx)

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		logrus.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	requestBody := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0x%x"] 
	}`, rawTxBytes)

	sendRequest(requestBody)

	if exit {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func TestCreateContract(gethCfg *evm.GethConfig, exit bool) {
	// A random private key. address = 0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02
	privateKeyStr := "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
	nonce := uint64(0)
	amount := big.NewInt(0)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	data, err := hex.DecodeString("608060405234801561001057600080fd5b506040518060400160405280600981526020016805465737445726332360bc1b8152506040518060400160405280600381526020016215115560ea1b815250816003908161005e9190610114565b50600461006b8282610114565b5050506101d3565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061009d57607f821691505b6020821081036100bd57634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561010f576000816000526020600020601f850160051c810160208610156100ec5750805b601f850160051c820191505b8181101561010b578281556001016100f8565b5050505b505050565b81516001600160401b0381111561012d5761012d610073565b6101418161013b8454610089565b846100c3565b602080601f831160018114610176576000841561015e5750858301515b600019600386901b1c1916600185901b17855561010b565b600085815260208120601f198616915b828110156101a557888601518255948401946001909101908401610186565b50858210156101c35787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b610785806101e26000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806340c10f191161006657806340c10f191461011857806370a082311461012d57806395d89b4114610156578063a9059cbb1461015e578063dd62ed3e1461017157600080fd5b806306fdde03146100a3578063095ea7b3146100c157806318160ddd146100e457806323b872dd146100f6578063313ce56714610109575b600080fd5b6100ab6101aa565b6040516100b891906105ce565b60405180910390f35b6100d46100cf366004610639565b61023c565b60405190151581526020016100b8565b6002545b6040519081526020016100b8565b6100d4610104366004610663565b610256565b604051601281526020016100b8565b61012b610126366004610639565b61027a565b005b6100e861013b36600461069f565b6001600160a01b031660009081526020819052604090205490565b6100ab610288565b6100d461016c366004610639565b610297565b6100e861017f3660046106c1565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6060600380546101b9906106f4565b80601f01602080910402602001604051908101604052809291908181526020018280546101e5906106f4565b80156102325780601f1061020757610100808354040283529160200191610232565b820191906000526020600020905b81548152906001019060200180831161021557829003601f168201915b5050505050905090565b60003361024a8185856102a5565b60019150505b92915050565b6000336102648582856102b7565b61026f85858561033a565b506001949350505050565b6102848282610399565b5050565b6060600480546101b9906106f4565b60003361024a81858561033a565b6102b283838360016103cf565b505050565b6001600160a01b038381166000908152600160209081526040808320938616835292905220546000198114610334578181101561032557604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064015b60405180910390fd5b610334848484840360006103cf565b50505050565b6001600160a01b03831661036457604051634b637e8f60e11b81526000600482015260240161031c565b6001600160a01b03821661038e5760405163ec442f0560e01b81526000600482015260240161031c565b6102b28383836104a4565b6001600160a01b0382166103c35760405163ec442f0560e01b81526000600482015260240161031c565b610284600083836104a4565b6001600160a01b0384166103f95760405163e602df0560e01b81526000600482015260240161031c565b6001600160a01b03831661042357604051634a1406b160e11b81526000600482015260240161031c565b6001600160a01b038085166000908152600160209081526040808320938716835292905220829055801561033457826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161049691815260200190565b60405180910390a350505050565b6001600160a01b0383166104cf5780600260008282546104c4919061072e565b909155506105419050565b6001600160a01b038316600090815260208190526040902054818110156105225760405163391434e360e21b81526001600160a01b0385166004820152602481018290526044810183905260640161031c565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b03821661055d5760028054829003905561057c565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516105c191815260200190565b60405180910390a3505050565b60006020808352835180602085015260005b818110156105fc578581018301518582016040015282016105e0565b506000604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b038116811461063457600080fd5b919050565b6000806040838503121561064c57600080fd5b6106558361061d565b946020939093013593505050565b60008060006060848603121561067857600080fd5b6106818461061d565b925061068f6020850161061d565b9150604084013590509250925092565b6000602082840312156106b157600080fd5b6106ba8261061d565b9392505050565b600080604083850312156106d457600080fd5b6106dd8361061d565b91506106eb6020840161061d565b90509250929050565b600181811c9082168061070857607f821691505b60208210810361072857634e487b7160e01b600052602260045260246000fd5b50919050565b8082018082111561025057634e487b7160e01b600052601160045260246000fdfea26469706673582212200ca760b3f238aafac761f0cfb50021e2b4cc6064fb1cfc66ad8b866114ebe25b64736f6c63430008180033")
	if err != nil {
		logrus.Fatal(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       nil,
		Value:    amount,
		Data:     data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		logrus.Fatal(err)
	}

	signer := types.LatestSigner(gethCfg.ChainConfig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("SignedTx = %+v", signedTx)

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		logrus.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	requestBody := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0x%x"] 
	}`, rawTxBytes)

	sendRequest(requestBody)

	if exit {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func testSendRawTransaction(exit bool) {
	time.Sleep(5 * time.Second)

	requestBody := `{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0xf86c808506fc23ac00825208947bd36074b61cfe75a53e1b9df7678c96e6463b02880de0b6b3a76400008026a0b5050757a8005286d85c8ae9408a933ca1126400a6749ca64e415f77db41b439a0294f3d15727ed8231d061db0ec1014ef1bf767f1665731d62e0327464cf8ad3e"] 
	}`

	sendRequest(requestBody)

	if exit {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
}

func CreateRandomWallet(gethCfg *evm.GethConfig, count int64) []string {
	time.Sleep(5 * time.Second)
	result := make([]string, 0)
	addressList := make([]string, 0)
	for i := int64(0); i < count; i++ {
		privateKey, address := generatePrivateKey()
		result = append(result, privateKey)
		addressList = append(addressList, address)

		// A random private key. address = 0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02
		requestBody := GenerateTransferEthRequest(gethCfg, "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff", address, 100)
		sendRequest(requestBody)
		fmt.Printf("privateKey: %s, address: %s\n", privateKey, address)
		time.Sleep(3 * time.Second)
	}

	fmt.Printf("---- privatekey and address list ----\n")
	for i := int64(0); i < count; i++ {
		fmt.Printf("privateKey: %s, address: %s\n", result[i], addressList[i])
	}

	return result
}

func GenerateTransferEthRequest(gethCfg *evm.GethConfig, privateKeyHex string, toAddress string, amount int64) string {
	nonce := uint64(0)
	to := common.HexToAddress(toAddress)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)
	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    big.NewInt(amount),
		Data:     data,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		logrus.Fatal(err)
	}

	signer := types.LatestSigner(gethCfg.ChainConfig)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("SignedTx = %+v\n", signedTx)

	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		logrus.Fatal(err)
	}

	requestBody := fmt.Sprintf(
		`	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0x%x"] 
	}`, rawTxBytes)
	return requestBody
}

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

func sendRequest(dataString string) {
	req, err := http.NewRequest("POST", "http://localhost:9092", bytes.NewBuffer([]byte(dataString)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("curl --location 'localhost:9092' --header 'Content-Type: application/json' --data '%s'\n", dataString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			logrus.Errorf("could not close response body, err:%v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	logrus.Infof("Response [%v] : %v", resp.Status, string(body))
	return
}
