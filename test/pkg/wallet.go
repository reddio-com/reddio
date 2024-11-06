package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
)

const (
	GenesisPrivateKey = "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff"
)

type EthWallet struct {
	PK      string `json:"pk"`
	Address string `json:"address"`
}

func (e *EthWallet) Copy() *EthWallet {
	return &EthWallet{
		PK:      e.PK,
		Address: e.Address,
	}
}

type WalletManager struct {
	cfg         *evm.GethConfig
	hostAddress string
}

func NewWalletManager(cfg *evm.GethConfig, hostAddress string) *WalletManager {
	return &WalletManager{
		cfg:         cfg,
		hostAddress: hostAddress,
	}
}

func (m *WalletManager) GenerateRandomWallets(count int, initialEthCount uint64) ([]*EthWallet, error) {
	wallets := make([]*EthWallet, 0)
	for i := 0; i < count; i++ {
		wallet, err := m.createEthWallet(initialEthCount)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	// wait block ready
	time.Sleep(4 * time.Second)
	return wallets, nil
}

func (m *WalletManager) BatchGenerateRandomWallets(count int, initialEthCount uint64) ([]*EthWallet, error) {
	var wallets []*EthWallet
	var txs []*RawTxReq
	for i := 0; i < count; i++ {
		privateKey, address := generatePrivateKey()
		wallets = append(wallets, &EthWallet{
			PK:      privateKey,
			Address: address,
		})
		txs = append(txs, &RawTxReq{
			privateKeyHex: GenesisPrivateKey,
			toAddress:     address,
			amount:        initialEthCount,
		})
	}
	err := m.batchTransferEth(txs)
	return wallets, err
}

func (m *WalletManager) createEthWallet(initialEthCount uint64) (*EthWallet, error) {
	privateKey, address := generatePrivateKey()
	return m.CreateEthWalletByAddress(initialEthCount, privateKey, address)
}

func (m *WalletManager) CreateEthWalletByAddress(initialEthCount uint64, privateKey, address string) (*EthWallet, error) {
	if err := m.transferEth(GenesisPrivateKey, address, initialEthCount, 0); err != nil {
		return nil, err
	}
	// log.Println(fmt.Sprintf("create wallet %v", address))
	return &EthWallet{PK: privateKey, Address: address}, nil
}

func (m *WalletManager) TransferEth(from, to *EthWallet, amount, nonce uint64) error {
	// log.Println(fmt.Sprintf("transfer %v eth from %v to %v", amount, from.Address, to.Address))
	if err := m.transferEth(from.PK, to.Address, amount, nonce); err != nil {
		return err
	}
	return nil
}

func (m *WalletManager) BatchTransferETH(steps []*Step) error {
	var batch []*RawTxReq
	for _, step := range steps {
		batch = append(batch, &RawTxReq{
			privateKeyHex: step.From.PK,
			toAddress:     step.To.Address,
			amount:        step.Count,
		})
	}
	return m.batchTransferEth(batch)
}

func (m *WalletManager) QueryEth(wallet *EthWallet) (uint64, error) {
	requestBody := fmt.Sprintf(
		`	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_getBalance",
		"params": ["%s","latest"] 
	}`, wallet.Address)
	d, err := sendRequest(m.hostAddress, requestBody)
	if err != nil {
		return 0, err
	}
	resp := &queryResponse{}
	if err := json.Unmarshal(d, resp); err != nil {
		return 0, nil
	}
	return parse(resp.Result)
}

func parse(v string) (uint64, error) {
	if !strings.HasPrefix(v, "0x") {
		return 0, fmt.Errorf("%v should start with 0v", v)
	}
	value, err := strconv.ParseUint(v[2:], 16, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

type queryResponse struct {
	Result string `json:"result"`
}

func (m *WalletManager) transferEth(privateKeyHex string, toAddress string, amount, nonce uint64) error {
	return m.sendRawTx(privateKeyHex, toAddress, amount, nonce)
}

func (m *WalletManager) batchTransferEth(rawTxs []*RawTxReq) error {
	return m.sendBatchRawTxs(rawTxs)
}

var counter = uint64(0)

// sendRawTx is used by transferring and contract creation/invocation.
func (m *WalletManager) sendRawTx(privateKeyHex string, toAddress string, amount uint64, nonce uint64) error {
	to := common.HexToAddress(toAddress)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(0)

	counter++
	nonce = nonce + counter

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &to,
		Value:    big.NewInt(int64(amount)),
		Data:     nil,
	})

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	chainID := m.cfg.ChainConfig.ChainID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	rawTxBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}

	requestBody := fmt.Sprintf(
		`	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendRawTransaction",
		"params": ["0x%x"] 
	}`, rawTxBytes)
	_, err = sendRequest(m.hostAddress, requestBody)
	return err
}

type RawTxReq struct {
	privateKeyHex string
	toAddress     string
	amount        uint64
	data          []byte
	nonce         uint64
}

func (m *WalletManager) sendBatchRawTxs(rawTxs []*RawTxReq) error {
	batchTx := new(ethrpc.BatchTx)
	nonceMap := make(map[string]uint64)
	for _, rawTx := range rawTxs {
		to := common.HexToAddress(rawTx.toAddress)
		gasLimit := uint64(21000)
		gasPrice := big.NewInt(0)

		if _, ok := nonceMap[rawTx.privateKeyHex]; ok {
			nonceMap[rawTx.privateKeyHex]++
		}

		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonceMap[rawTx.privateKeyHex],
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &to,
			Value:    big.NewInt(int64(rawTx.amount)),
			Data:     rawTx.data,
		})

		privateKey, err := crypto.HexToECDSA(rawTx.privateKeyHex)
		if err != nil {
			log.Fatal(err)
		}

		chainID := m.cfg.ChainConfig.ChainID
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal(err)
		}
		rawTxBytes, err := rlp.EncodeToBytes(signedTx)
		if err != nil {
			log.Fatal(err)
		}
		batchTx.TxsBytes = append(batchTx.TxsBytes, rawTxBytes)
	}

	batchTxBytes, err := json.Marshal(batchTx)
	if err != nil {
		return err
	}

	requestBody := fmt.Sprintf(
		`	{
		"jsonrpc": "2.0",
		"id": 0,
		"method": "eth_sendBatchRawTransactions",
		"params": ["0x%x"] 
	}`, batchTxBytes)
	_, err = sendRequest(m.hostAddress, requestBody)
	return err
}
