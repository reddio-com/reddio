package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reddio-com/reddio/bridge/test/bindings"
	"github.com/reddio-com/reddio/bridge/utils"
)

const (
	ZeroAddress = "0x0"
)

var (
	// Testnet
	//use your own l1 and l2 endpoint
	sepoliaHelpConfig = helpConfig{
		testAdmin:       "",
		L1ClientAddress: "",
		//L2ClientAddress: "https://reddio-dev.reddio.com/",
		//L2ClientAddress: "https://reddio-evm-bridge.reddio.com/",
		L2ClientAddress:            "http://localhost:9092",
		ParentlayerContractAddress: "0x9F7e49fcAB7eD379451e8422D20908bF439011A5",
		ChildlayerContractAddress:  "0xeC054c6ee2DbbeBC9EbCA50CdBF94A94B02B2E40",
		//ChildlayerContractAddress: "0xeC054c6ee2DbbeBC9EbCA50CdBF94A94B02B2E40",
		//testPublicKey1:            "0x0CC0cD4A9024A2d15BbEdd348Fbf7Cd69B5489bA",
		testPublicKey1:          "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645",
		testPublicKey2:          "0x66eb032B3a74d85C8b6965a4df788f3C31678b1a",
		adminPublicKey:          "0x7Bd36074b61Cfe75a53e1B9DF7678C96E6463b02",
		maxRetries:              300,
		waitForConfirmationTime: 12 * time.Second,
		L1ETHAddress:            "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		//L1ERC20Address:          "0xF1E77FF9A4d4fc09CD955EfC44cB843617C73F23",
		L1ERC20Address:   "0x9627E313C18be25fC03100bbD3bf48743B4dee70",
		L1ERC721Address:  "0xA399AA7a6b2f4b36E36f2518FeE7C2AEC48dfD10",
		L1ERC1155Address: "0x3713cC896e86AA63Ec97088fB5894E3c985792e7",
		L1REDAddress:     "0xB878927d79975BDb288ab53271f171534A49eb7D",
		l2gaslimit:       *big.NewInt(0),
	}
)

type helpConfig struct {
	testAdmin                  string
	L1ClientAddress            string
	L2ClientAddress            string
	ParentlayerContractAddress string
	ChildlayerContractAddress  string
	testPublicKey1             string
	testPublicKey2             string
	adminPublicKey             string
	maxRetries                 int
	waitForConfirmationTime    time.Duration
	L1ETHAddress               string
	L1ERC20Address             string
	L1ERC721Address            string
	L1ERC1155Address           string
	L1REDAddress               string
	l2gaslimit                 big.Int
}

// Deposit Tests
func SetupForkedChain() error {
	return nil
}
func TestPrepareSendNativeToken(t *testing.T) {
	t.Run("PrepareSendNativeToken", func(t *testing.T) {
		fmt.Println("PrepareSendNativeToken")

		contractAddress := common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress)

		sendAmount := big.NewInt(9e18)

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		require.True(t, ok)
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		nonce, err := l2Client.PendingNonceAt(context.Background(), fromAddress)
		require.NoError(t, err)

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)

		chainID, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)

		tx := types.NewTransaction(
			nonce,
			contractAddress,
			sendAmount,
			21000, // gas limit
			gasPrice,
			nil, // data
		)

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		require.NoError(t, err)

		err = l2Client.SendTransaction(context.Background(), signedTx)
		require.NoError(t, err)

		fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, signedTx.Hash())
		require.NoError(t, err)
		assert.True(t, success, "Transaction was not confirmed")

		fmt.Println("Transaction confirmed: ", signedTx.Hash().Hex())
	})
}

// testCaseinfo:
// 1. deposit ETH to L2
// 2. check the balance of testPublicKey2 in l2 is increased by depositAmount
func TestDepositETH(t *testing.T) {
	t.Run("DepositETH", func(t *testing.T) {
		fmt.Println("DepositETH")
		depositAmount := big.NewInt(10000)
		//Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		if err != nil {
			logrus.Fatal("failed to connect to L1 geth", "endpoint", sepoliaHelpConfig.L1ClientAddress, "err", err)
		}
		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		if err != nil {
			logrus.Fatal("failed to connect to L2 geth", "endpoint", sepoliaHelpConfig.L2ClientAddress, "err", err)
		}
		defer l1Client.Close()
		defer l2Client.Close()
		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		if err != nil {
			logrus.Fatalf("failed to create ChildTokenMessageTransmitterFacet contract: %v", err)
		}
		//if this L2BridgeTokenAddress is not exist,need to register it at previous step
		l2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		if err != nil {
			logrus.Fatalf("failed to get bridged token address: %v", err)
		}
		var startBalance *big.Int
		if l2BridgeTokenAddress == (common.Address{}) {
			startBalance = big.NewInt(0)
		} else {
			l2BridgeERC20Token, err := bindings.NewERC20Token(l2BridgeTokenAddress, l2Client)
			if err != nil {
				logrus.Fatalf("failed to create ERC20Token contract: %v", err)
			}

			startBalance, err = l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
			if err != nil {
				logrus.Fatalf("failed to get balance of testPublicKey: %v", err)
			}
		}
		// Action
		// get gas price
		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		if err != nil {
			logrus.Fatal("failed to get gas price", "err", err)
		}
		chainid, err := l1Client.ChainID(context.Background())
		if err != nil {
			logrus.Fatal("failed to get chain id", "err", err)
		}
		t.Log("gas price", "price", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		if err != nil {
			logrus.Fatalf("Error loading private key: %v", err)
		}
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		if err != nil {
			logrus.Fatalf("failed to create private key %v", err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		if err != nil {
			logrus.Fatalf("failed to create authorized transactor: %v", err)
		}
		auth.GasPrice = gasPrice

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		if err != nil {
			logrus.Fatalf("failed to create ParentTokenMessageTransmitterFacet contract: %v", err)
		}
		auth.Value = depositAmount
		tx, err := ParentTokenMessageTransmitterFacet.DepositETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), depositAmount, big.NewInt(0))
		if err != nil {
			logrus.Fatalf("failed to deposit eth: %v", err)
		}
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		if err != nil {
			logrus.Fatalf("failed to wait for confirmation: %v", err)
		}
		assert.True(t, success)

		// wait for the L2 confirmation
		time.Sleep(5 * time.Second)

		//Check the balance of the testPublicKey\\
		l2BridgeTokenAddress, err = ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		if err != nil {
			logrus.Fatalf("failed to get bridged token address: %v", err)
		}
		l2BridgeERC20Token, err := bindings.NewERC20Token(l2BridgeTokenAddress, l2Client)
		if err != nil {
			logrus.Fatalf("failed to create ERC20Token contract: %v", err)
		}
		fmt.Println("L2BridgeTokenAddress: ", l2BridgeTokenAddress)
		balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		if err != nil {
			logrus.Fatalf("failed to get balance of testPublicKey: %v", err)
		}
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Add(depositAmount, startBalance)

		assert.Equal(t, expectedBalance, balance)
	})
}
func TestTransferETHToZeroAddress(t *testing.T) {
	// Arrange
	l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
	require.NoError(t, err)
	defer l2Client.Close()
	startBalance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress("0x0000000000000000000000000000000000000000"), nil)
	if err != nil {
		logrus.Fatalf("Failed to get balance: %v", err)
	}
	fmt.Println("Start Balance: ", startBalance)
	assert.Equal(t, 1, 2)

	gasPrice, err := l2Client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	chainID, err := l2Client.ChainID(context.Background())
	require.NoError(t, err)

	privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
	require.NoError(t, err)
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	auth.GasPrice = gasPrice

	// 获取当前账户的nonce
	nonce, err := l2Client.PendingNonceAt(context.Background(), auth.From)
	require.NoError(t, err)

	// 设置auth的nonce
	auth.Nonce = big.NewInt(int64(nonce))

	transferAmount := big.NewInt(1e18) // 1 ETH

	tx := types.NewTransaction(
		nonce,
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		transferAmount,
		auth.GasLimit,
		auth.GasPrice,
		nil,
	)

	signedTx, err := auth.Signer(auth.From, tx)
	require.NoError(t, err)

	err = l2Client.SendTransaction(context.Background(), signedTx)

	require.NoError(t, err)
	fmt.Println("Transaction sent: ", signedTx.Hash().Hex())

	success, err := waitForConfirmation(l2Client, signedTx.Hash())
	require.NoError(t, err)
	assert.True(t, success)
	fmt.Println("Transaction confirmed: ", signedTx.Hash().Hex())
}

// testCaseinfo:
// 1. approve ERC20 token
// 2. deposit ERC20 token
// 3. check the balance of testPublicKey2 in l2 is increased by depositAmount
func TestDepositERC20(t *testing.T) {
	t.Run("DepositERC20", func(t *testing.T) {
		fmt.Println("DepositERC2020")
		depositAmount := big.NewInt(100)
		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		l2TokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC20 Token Address: ", l2TokenAddress)

		var startBalance *big.Int
		if l2TokenAddress == (common.Address{}) {
			startBalance = big.NewInt(0)
		} else {
			l2RC20Token, err := bindings.NewERC20Token(l2TokenAddress, l2Client)
			require.NoError(t, err)

			startBalance, err = l2RC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
			require.NoError(t, err)
		}

		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l1Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		erc20Token, err := bindings.NewERC20Token(common.HexToAddress(sepoliaHelpConfig.L1ERC20Address), l1Client)
		require.NoError(t, err)
		tx, err := erc20Token.Approve(auth, common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), depositAmount)
		require.NoError(t, err)
		fmt.Println("Approve transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		require.NoError(t, err)
		tx, err = ParentTokenMessageTransmitterFacet.DepositERC20Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), depositAmount, big.NewInt(0))
		require.NoError(t, err)
		fmt.Println("DepositERC20 transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)
		time.Sleep(5 * time.Second)
		l2TokenAddress, err = ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC20 Token Address2: ", l2TokenAddress)
		l2RC20Token, err := bindings.NewERC20Token(l2TokenAddress, l2Client)
		require.NoError(t, err)
		balance, err := l2RC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Add(depositAmount, startBalance)

		assert.Equal(t, expectedBalance, balance)
	})
}

// testCaseinfo:
// 1. approve RED token
// 2. deposit RED token
// 3. check the balance of testPublicKey2 in l2 is increased by depositAmount
func TestDepositRED(t *testing.T) {
	t.Run("DepositRED", func(t *testing.T) {
		fmt.Println("DepositRED")
		depositAmount := big.NewInt(1e14)
		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		startBalance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), nil)
		if err != nil {
			logrus.Fatalf("Failed to get balance: %v", err)
		}
		fmt.Println("Start Balance: ", startBalance)

		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l1Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		erc20Token, err := bindings.NewERC20Token(common.HexToAddress(sepoliaHelpConfig.L1REDAddress), l1Client)
		require.NoError(t, err)
		tx, err := erc20Token.Approve(auth, common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), depositAmount)
		require.NoError(t, err)
		fmt.Println("Approve transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		require.NoError(t, err)
		tx, err = ParentTokenMessageTransmitterFacet.DepositRED(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), depositAmount, big.NewInt(0))
		require.NoError(t, err)
		fmt.Println("DepositRED transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)
		time.Sleep(5 * time.Second)
		require.NoError(t, err)
		balance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), nil)
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Add(depositAmount, startBalance)

		assert.Equal(t, expectedBalance, balance)
	})
}

// testCaseinfo:
// 1. mint ERC721 token
// 2. approve ERC721 token
// 3. deposit ERC721 token
// 4. check the owner of tokenID is testPublicKey2
func TestDepositERC721(t *testing.T) {
	t.Run("DepositERC721", func(t *testing.T) {
		fmt.Println("DepositERC721")

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l1Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		erc721Token, err := bindings.NewERC721Token(common.HexToAddress(sepoliaHelpConfig.L1ERC721Address), l1Client)
		require.NoError(t, err)

		tokenID := big.NewInt(4)

		owner, err := erc721Token.OwnerOf(callOpts, tokenID)
		require.NoError(t, err)
		fmt.Println("Owner of tokenID: ", owner)
		assert.Equal(t, common.HexToAddress(sepoliaHelpConfig.adminPublicKey), owner)

		tx, err := erc721Token.SetApprovalForAll(auth, common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), true)
		require.NoError(t, err)
		fmt.Println("SetApprovalForAll transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		require.NoError(t, err)
		tx, err = ParentTokenMessageTransmitterFacet.DepositERC721Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC721Address), common.HexToAddress(sepoliaHelpConfig.adminPublicKey), tokenID, big.NewInt(0))
		require.NoError(t, err)
		fmt.Println("DepositERC721 transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		time.Sleep(5 * time.Second)

		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC721TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC721Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC721 Token Address2: ", L2BridgeTokenAddress)
		l2ERC721Token, err := bindings.NewERC721Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)
		owner, err = l2ERC721Token.OwnerOf(callOpts, tokenID)
		require.NoError(t, err)
		fmt.Println("Owner of tokenID: ", owner)

		assert.Equal(t, common.HexToAddress(sepoliaHelpConfig.adminPublicKey), owner)
	})
}

func TestDepositERC1155Batch(t *testing.T) {
	t.Run("DepositERC1155Batch", func(t *testing.T) {
		fmt.Println("DepositERC1155Batch")

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC1155TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC1155 Token Address: ", L2BridgeTokenAddress)

		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l1Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		erc1155Token, err := bindings.NewERC1155Token(common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address), l1Client)
		require.NoError(t, err)
		tokenIDs := []*big.Int{big.NewInt(1), big.NewInt(2)}
		amounts := []*big.Int{big.NewInt(100), big.NewInt(200)}
		data := []byte{}

		tx, err := erc1155Token.MintBatch(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), tokenIDs, amounts, data)
		require.NoError(t, err)
		fmt.Println("MintBatch transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		for _, tokenID := range tokenIDs {
			balance, err := erc1155Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), tokenID)
			require.NoError(t, err)
			fmt.Printf("Balance of tokenID %s: %s\n", tokenID.String(), balance.String())
			assert.Equal(t, amounts[tokenID.Int64()-1], balance)
		}

		tx, err = erc1155Token.SetApprovalForAll(auth, common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), true)
		require.NoError(t, err)
		fmt.Println("SetApprovalForAll transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		require.NoError(t, err)
		tx, err = ParentTokenMessageTransmitterFacet.DepositERC1155Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), tokenIDs, amounts, &sepoliaHelpConfig.l2gaslimit)
		require.NoError(t, err)
		fmt.Println("DepositERC1155Batch transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		time.Sleep(5 * time.Second)

		L2BridgeTokenAddress, err = ChildBridgeCoreFacet.GetBridgedERC1155TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC1155 Token Address2: ", L2BridgeTokenAddress)
		l2ERC1155Token, err := bindings.NewERC1155Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)
		for _, tokenID := range tokenIDs {
			balance, err := l2ERC1155Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey2), tokenID)
			require.NoError(t, err)
			fmt.Printf("Balance of tokenID %s: %s\n", tokenID.String(), balance.String())
			assert.Equal(t, amounts[tokenID.Int64()-1], balance)
		}
	})
}

// Withdraw Tests

func TestWithdrawETH(t *testing.T) {
	t.Run("WithdrawETH", func(t *testing.T) {
		fmt.Println("WithdrawETH1")

		withdrawAmount := big.NewInt(5000)

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		// get L2BridgeTokenAddress
		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		require.NoError(t, err)

		l2BridgeERC20Token, err := bindings.NewERC20Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)

		startBalance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)

		//TransferL2ETH(t, l2Client, sepoliaHelpConfig.testAdmin, common.HexToAddress(sepoliaHelpConfig.testPublicKey), big.NewInt(1e18))
		// Action
		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice
		auth.GasLimit = 3000000
		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), withdrawAmount)
		require.NoError(t, err)
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		// check testPublicKey L2eth balance
		balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Sub(startBalance, withdrawAmount)

		assert.Equal(t, expectedBalance, balance)
	})

}
func TestWithdrawRED(t *testing.T) {
	t.Run("WithdrawRED", func(t *testing.T) {
		fmt.Println("WithdrawRED")
		withdrawAmount := big.NewInt(2000)

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		startBalance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), nil)
		if err != nil {
			logrus.Fatalf("Failed to get balance: %v", err)
		}
		fmt.Println("Start Balance: ", startBalance)

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice
		auth.GasLimit = 3000000
		auth.Value = withdrawAmount
		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)
		nonce, err := l2Client.PendingNonceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)
		fmt.Println("Nonce: ", nonce)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawRED(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)
		fmt.Println("WithdrawRED transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		time.Sleep(5 * time.Second)
		require.NoError(t, err)
		balance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), nil)
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Sub(startBalance, withdrawAmount)

		assert.Equal(t, expectedBalance, balance)
	})
}
func TestWithdrawERC20(t *testing.T) {
	t.Run("WithdrawERC20", func(t *testing.T) {
		fmt.Println("WithdrawERC20")

		withdrawAmount := big.NewInt(100)

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address))
		require.NoError(t, err)

		l2BridgeERC20Token, err := bindings.NewERC20Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)

		startBalance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawErc20Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey1), withdrawAmount)
		require.NoError(t, err)
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Sub(startBalance, withdrawAmount)

		assert.Equal(t, expectedBalance, balance)
	})
}

func TestWithdrawERC721(t *testing.T) {
	// Implement ERC721 withdraw test
	t.Run("WithdrawERC721", func(t *testing.T) {
		fmt.Println("WithdrawERC721")

		tokenID := big.NewInt(1)

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC721TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC721Address))
		require.NoError(t, err)

		l2ERC721Token, err := bindings.NewERC721Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)

		owner, err := l2ERC721Token.OwnerOf(callOpts, tokenID)
		require.NoError(t, err)
		assert.Equal(t, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), owner)

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawErc721Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC721Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey2), tokenID)
		require.NoError(t, err)
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		owner, err = l2ERC721Token.OwnerOf(callOpts, tokenID)
		require.NoError(t, err)
		fmt.Println("Owner of tokenID after withdrawal: ", owner)
		assert.NotEqual(t, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), owner)
	})
}

func TestWithdrawERC1155Batch(t *testing.T) {
	// Implement ERC1155 batch withdraw test
	t.Run("WithdrawERC1155Batch", func(t *testing.T) {
		fmt.Println("WithdrawERC1155Batch")

		tokenIDs := []*big.Int{big.NewInt(1), big.NewInt(2)}
		amounts := []*big.Int{big.NewInt(10), big.NewInt(20)}

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC1155TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address))
		require.NoError(t, err)

		l2ERC1155Token, err := bindings.NewERC1155Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)

		balances, err := l2ERC1155Token.BalanceOfBatch(callOpts, []common.Address{common.HexToAddress(sepoliaHelpConfig.testPublicKey1), common.HexToAddress(sepoliaHelpConfig.testPublicKey1)}, tokenIDs)
		require.NoError(t, err)
		assert.Equal(t, amounts[0], balances[0])
		assert.Equal(t, amounts[1], balances[1])

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawErc1155BatchToken(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC1155Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey2), tokenIDs, amounts)
		require.NoError(t, err)
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		balances, err = l2ERC1155Token.BalanceOfBatch(callOpts, []common.Address{common.HexToAddress(sepoliaHelpConfig.testPublicKey1), common.HexToAddress(sepoliaHelpConfig.testPublicKey1)}, tokenIDs)
		require.NoError(t, err)
		fmt.Println("Balances of tokenIDs after withdrawal: ", balances)
		assert.Equal(t, big.NewInt(0), balances[0])
		assert.Equal(t, big.NewInt(0), balances[1])
	})
}

func TestMultipleWithdrawETH(t *testing.T) {
	t.Run("MultipleWithdrawETH", func(t *testing.T) {
		fmt.Println("MultipleWithdrawETH")

		withdrawAmount := big.NewInt(1)
		numWithdrawals := 3

		// Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		require.NoError(t, err)
		defer l1Client.Close()

		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		// get L2BridgeTokenAddress
		L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		require.NoError(t, err)

		l2BridgeERC20Token, err := bindings.NewERC20Token(L2BridgeTokenAddress, l2Client)
		require.NoError(t, err)

		startBalance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
		require.NoError(t, err)

		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		for i := 0; i < numWithdrawals; i++ {
			// Action
			gasPrice, err := l2Client.SuggestGasPrice(context.Background())
			require.NoError(t, err)
			chainid, err := l2Client.ChainID(context.Background())
			require.NoError(t, err)
			t.Log("gas price", "price", gasPrice)
			fmt.Println("gasPrice", gasPrice)

			auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
			require.NoError(t, err)
			auth.GasPrice = gasPrice
			auth.GasLimit = 3000000

			ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
			require.NoError(t, err)

			tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), withdrawAmount)
			require.NoError(t, err)
			fmt.Println("Transaction sent: ", tx.Hash().Hex())

			success, err := waitForConfirmation(l2Client, tx.Hash())
			require.NoError(t, err)
			assert.True(t, success)

			// check testPublicKey L2eth balance
			balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey1))
			require.NoError(t, err)
			fmt.Println("Balance of testPublicKey after withdrawal", i+1, ":", balance)
			expectedBalance := new(big.Int).Sub(startBalance, new(big.Int).Mul(withdrawAmount, big.NewInt(int64(i+1))))

			assert.Equal(t, expectedBalance, balance)
		}
	})
}
func TestDistributeETHAndWithdraw(t *testing.T) {
	t.Run("DistributeETHAndWithdraw", func(t *testing.T) {
		fmt.Println("DistributeETHAndWithdraw")
		//check the balance of the adminPrivateKey

		distributeL2EthAmount := big.NewInt(1e18)
		distributeBridgeTokenAmount := big.NewInt(3)

		testWithdrawBridgeTokenAmount := big.NewInt(1)

		// random 3 addresses
		privateKeys := make([]*ecdsa.PrivateKey, 3)
		addresses := make([]common.Address, 3)
		for i := 0; i < 3; i++ {
			privateKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			privateKeys[i] = privateKey
			addresses[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
			fmt.Printf("test Address %d: %s\n", i+1, addresses[i].String())
		}

		// Arrange
		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		adminPrivateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		// adminPrivateKey, err := crypto.HexToECDSA(adminPrivateKeyStr)
		// require.NoError(t, err)
		adminRdoBalance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.adminPublicKey), nil)
		fmt.Printf("Admin l2ETH Balance: %s\n", adminRdoBalance)
		L2BridgeTokenAddress, err := GetL2BridgeTokenAddress(l2Client, sepoliaHelpConfig.ChildlayerContractAddress, sepoliaHelpConfig.L1ETHAddress)
		require.NoError(t, err)
		fmt.Println("L2BridgeTokenAddress: ", L2BridgeTokenAddress)
		L2BridgeToken, err := bindings.NewERC20Token(L2BridgeTokenAddress, l2Client)
		L2BridgeTokenBalance, err := L2BridgeToken.BalanceOf(nil, common.HexToAddress(sepoliaHelpConfig.adminPublicKey))
		fmt.Printf("Admin l2ETH Balance: %s\n", L2BridgeTokenBalance)

		for _, address := range addresses {
			TransferL2ETH(t, l2Client, adminPrivateKeyStr, address, distributeL2EthAmount)
			//ApproveERC20(t, l2Client, adminPrivateKeyStr, L2BridgeTokenAddress, address, distributeL1EthAmount)
			TransferERC20(t, l2Client, adminPrivateKeyStr, L2BridgeTokenAddress, address, distributeBridgeTokenAmount)
			l2EthBalance, err := l2Client.BalanceAt(context.Background(), address, nil)
			require.NoError(t, err)
			fmt.Printf("l2Eth Balance of address %s after distribution: %s\n", address, l2EthBalance.String())
			assert.Equal(t, distributeL2EthAmount, l2EthBalance)
			l1EthBalance, err := L2BridgeToken.BalanceOf(nil, address)
			require.NoError(t, err)
			fmt.Printf("l1Eth Balance of address %s after distribution: %s\n", address, l1EthBalance.String())
			assert.Equal(t, distributeBridgeTokenAmount, l1EthBalance)

		}

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}

		var wg sync.WaitGroup
		for i := 0; i < len(privateKeys); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				// Action
				gasPrice, err := l2Client.SuggestGasPrice(context.Background())
				require.NoError(t, err)
				chainid, err := l2Client.ChainID(context.Background())
				require.NoError(t, err)
				fmt.Println("gasPrice", gasPrice)

				privateKey := privateKeys[i]
				require.NoError(t, err)

				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
				require.NoError(t, err)
				auth.GasPrice = gasPrice
				auth.GasLimit = 3000000

				ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
				require.NoError(t, err)

				tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, addresses[i], testWithdrawBridgeTokenAmount)
				require.NoError(t, err)
				fmt.Println("Transaction sent: ", tx.Hash().Hex())

				success, err := waitForConfirmation(l2Client, tx.Hash())
				require.NoError(t, err)
				assert.True(t, success)

				balance, err := L2BridgeToken.BalanceOf(callOpts, addresses[i])
				require.NoError(t, err)
				fmt.Println("Balance of address after withdrawal", i+1, ":", balance)
				expectedBalance := new(big.Int).Sub(distributeBridgeTokenAmount, testWithdrawBridgeTokenAmount)

				assert.Equal(t, expectedBalance, balance)
			}(i)
		}
		wg.Wait()
	})
}

func TestErrorAndCorrectTransactions(t *testing.T) {
	t.Run("DistributeETHAndWithdraw", func(t *testing.T) {
		fmt.Println("DistributeETHAndWithdraw")
		//check the balance of the adminPrivateKey

		distributeL2EthAmount := big.NewInt(1e18)
		distributeBridgeTokenAmount := big.NewInt(3)

		testWithdrawBridgeTokenAmount := big.NewInt(1)

		// random 3 addresses
		privateKeys := make([]*ecdsa.PrivateKey, 3)
		addresses := make([]common.Address, 3)
		for i := 0; i < 3; i++ {
			privateKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			privateKeys[i] = privateKey
			addresses[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
			fmt.Printf("test Address %d: %s\n", i+1, addresses[i].String())
		}

		// Arrange
		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		require.NoError(t, err)
		defer l2Client.Close()

		adminPrivateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env", "PRIVATE_KEY")
		require.NoError(t, err)
		// adminPrivateKey, err := crypto.HexToECDSA(adminPrivateKeyStr)
		// require.NoError(t, err)
		adminRdoBalance, err := l2Client.BalanceAt(context.Background(), common.HexToAddress(sepoliaHelpConfig.adminPublicKey), nil)
		fmt.Printf("Admin l2ETH Balance: %s\n", adminRdoBalance)
		L2BridgeTokenAddress, err := GetL2BridgeTokenAddress(l2Client, sepoliaHelpConfig.ChildlayerContractAddress, sepoliaHelpConfig.L1ETHAddress)
		require.NoError(t, err)
		fmt.Println("L2BridgeTokenAddress: ", L2BridgeTokenAddress)
		L2BridgeToken, err := bindings.NewERC20Token(L2BridgeTokenAddress, l2Client)
		L2BridgeTokenBalance, err := L2BridgeToken.BalanceOf(nil, common.HexToAddress(sepoliaHelpConfig.adminPublicKey))
		fmt.Printf("Admin l2ETH Balance: %s\n", L2BridgeTokenBalance)

		for _, address := range addresses {
			TransferL2ETH(t, l2Client, adminPrivateKeyStr, address, distributeL2EthAmount)
			//ApproveERC20(t, l2Client, adminPrivateKeyStr, L2BridgeTokenAddress, address, distributeL1EthAmount)
			TransferERC20(t, l2Client, adminPrivateKeyStr, L2BridgeTokenAddress, address, distributeBridgeTokenAmount)
			l2EthBalance, err := l2Client.BalanceAt(context.Background(), address, nil)
			require.NoError(t, err)
			fmt.Printf("l2Eth Balance of address %s after distribution: %s\n", address, l2EthBalance.String())
			assert.Equal(t, distributeL2EthAmount, l2EthBalance)
			l1EthBalance, err := L2BridgeToken.BalanceOf(nil, address)
			require.NoError(t, err)
			fmt.Printf("l1Eth Balance of address %s after distribution: %s\n", address, l1EthBalance.String())
			assert.Equal(t, distributeBridgeTokenAmount, l1EthBalance)

		}

		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}

		var wg1 sync.WaitGroup
		for i := 0; i < len(privateKeys); i++ {
			wg1.Add(1)
			go func(i int) {
				defer wg1.Done()
				// Action
				gasPrice, err := l2Client.SuggestGasPrice(context.Background())
				require.NoError(t, err)
				chainid, err := l2Client.ChainID(context.Background())
				require.NoError(t, err)
				fmt.Println("gasPrice", gasPrice)

				privateKey := privateKeys[i]
				require.NoError(t, err)

				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
				require.NoError(t, err)
				auth.GasPrice = gasPrice
				auth.GasLimit = 3000000

				ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
				require.NoError(t, err)

				tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, addresses[i], testWithdrawBridgeTokenAmount)
				require.NoError(t, err)
				fmt.Println("Transaction sent: ", tx.Hash().Hex())

				success, err := waitForConfirmation(l2Client, tx.Hash())
				require.NoError(t, err)
				assert.True(t, success)

				balance, err := L2BridgeToken.BalanceOf(callOpts, addresses[i])
				require.NoError(t, err)
				fmt.Println("Balance of address after withdrawal", i+1, ":", balance)
				expectedBalance := new(big.Int).Sub(distributeBridgeTokenAmount, testWithdrawBridgeTokenAmount)

				assert.Equal(t, expectedBalance, balance)
			}(i)
		}
		wg1.Wait()

		//test error transaction
		privateKey, err := crypto.HexToECDSA(adminPrivateKeyStr)
		require.NoError(t, err)

		chainID, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)

		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
		require.NoError(t, err)
		auth.GasPrice = gasPrice
		auth.GasLimit = 3000000

		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		for i := 0; i < 100; i++ {
			tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey1), big.NewInt(9e18))
			require.NoError(t, err)
			fmt.Printf("Error Transaction %d sent: %s\n", i+1, tx.Hash().Hex())
		}
		time.Sleep(5 * time.Second)
		//test error transaction end
		var wg sync.WaitGroup
		for i := 0; i < len(privateKeys); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				// Action
				gasPrice, err := l2Client.SuggestGasPrice(context.Background())
				require.NoError(t, err)
				chainid, err := l2Client.ChainID(context.Background())
				require.NoError(t, err)
				fmt.Println("gasPrice", gasPrice)

				privateKey := privateKeys[i]
				require.NoError(t, err)

				auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
				require.NoError(t, err)
				auth.GasPrice = gasPrice
				auth.GasLimit = 3000000

				ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
				require.NoError(t, err)

				tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, addresses[i], testWithdrawBridgeTokenAmount)
				require.NoError(t, err)
				fmt.Println("Transaction sent: ", tx.Hash().Hex())

				// success, err := waitForConfirmation(l2Client, tx.Hash())
				// require.NoError(t, err)
				// assert.True(t, success)

				balance, err := L2BridgeToken.BalanceOf(callOpts, addresses[i])
				require.NoError(t, err)
				fmt.Println("Balance of address after withdrawal", i+1, ":", balance)
				expectedBalance := new(big.Int).Sub(distributeBridgeTokenAmount, testWithdrawBridgeTokenAmount)

				assert.Equal(t, expectedBalance, balance)
			}(i)
		}
		wg.Wait()
	})
}

func waitForConfirmation(client *ethclient.Client, txHash common.Hash) (bool, error) {
	for i := 0; i < sepoliaHelpConfig.maxRetries; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return true, nil
			}
			return false, fmt.Errorf("transaction failed with status: %v", receipt.Status)
		}
		time.Sleep(sepoliaHelpConfig.waitForConfirmationTime)
	}
	return false, fmt.Errorf("transaction was not confirmed after %d retries", sepoliaHelpConfig.maxRetries)
}

func TransferL2ETH(t *testing.T, l2Client *ethclient.Client, fromPrivateKey string, toAddress common.Address, amount *big.Int) {
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := l2Client.PendingNonceAt(context.Background(), fromAddress)
	require.NoError(t, err)

	gasPrice, err := l2Client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	chainID, err := l2Client.ChainID(context.Background())
	require.NoError(t, err)

	tx := types.NewTransaction(nonce, toAddress, amount, 21000, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	require.NoError(t, err)

	err = l2Client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	fmt.Printf("TransferL2ETH Transaction sent: %s\n", signedTx.Hash().Hex())

	success, err := waitForConfirmation(l2Client, signedTx.Hash())
	require.NoError(t, err)
	assert.True(t, success, "L2 transaction was not confirmed")
}

func TransferERC20(t *testing.T, client *ethclient.Client, fromPrivateKey string, tokenAddress common.Address, toAddress common.Address, amount *big.Int) {
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress: ", fromAddress)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000)

	erc20Token, err := bindings.NewERC20Token(tokenAddress, client)
	require.NoError(t, err)

	tx, err := erc20Token.Transfer(auth, toAddress, amount)
	require.NoError(t, err)

	fmt.Printf("TransferERC20 Transaction sent: %s\n", tx.Hash().Hex())

	success, err := waitForConfirmation(client, tx.Hash())
	require.NoError(t, err)
	assert.True(t, success, "ERC20 transaction was not confirmed")

}

func GetL2BridgeTokenAddress(l2Client *ethclient.Client, childLayerContractAddress, l1TokenAddress string) (common.Address, error) {
	callOpts := &bind.CallOpts{
		Context: context.Background(),
	}
	ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(childLayerContractAddress), l2Client)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create ChildBridgeCoreFacet contract: %v", err)
	}

	L2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(l1TokenAddress))
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get bridged token address: %v", err)
	}

	return L2BridgeTokenAddress, nil
}
func ApproveERC20(t *testing.T, client *ethclient.Client, fromPrivateKey string, tokenAddress common.Address, spenderAddress common.Address, amount *big.Int) {
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	require.NoError(t, err)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000)

	erc20Token, err := bindings.NewERC20Token(tokenAddress, client)
	require.NoError(t, err)

	tx, err := erc20Token.Approve(auth, spenderAddress, amount)
	require.NoError(t, err)

	fmt.Printf("Approve transaction sent: %s\n", tx.Hash().Hex())

	success, err := waitForConfirmation(client, tx.Hash())
	require.NoError(t, err)
	assert.True(t, success, "ERC20 approve transaction was not confirmed")
}
