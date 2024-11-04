package test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/bridge/test/bindings"
	"github.com/reddio-com/reddio/bridge/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Testnet
	sepoliaHelpConfig = helpConfig{
		testAdmin:                  "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff",
		L1ClientAddress:            "wss://sepolia.infura.io/ws/v3/80b72ad34e16495595abeb6ccc30255a",
		L2ClientAddress:            "http://localhost:9092",
		ParentlayerContractAddress: "0xBac6aE08c64D389555A3E64D7C8167339327b77e",
		ChildlayerContractAddress:  "0xeC054c6ee2DbbeBC9EbCA50CdBF94A94B02B2E40",
		testPublicKey:              "0x7888b7B844B4B16c03F8daCACef7dDa0F5188645",
		maxRetries:                 300,
		waitForConfirmationTime:    12 * time.Second,
		L1ETHAddress:               "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		L1ERC20Address:             "0xF1E77FF9A4d4fc09CD955EfC44cB843617C73F23",
		L1ERC721Address:            "0xA399AA7a6b2f4b36E36f2518FeE7C2AEC48dfD10",
		L1ERC1155Address:           "0x3713cC896e86AA63Ec97088fB5894E3c985792e7",
	}
)

type helpConfig struct {
	testAdmin                  string
	L1ClientAddress            string
	L2ClientAddress            string
	ParentlayerContractAddress string
	ChildlayerContractAddress  string
	testPublicKey              string
	maxRetries                 int
	waitForConfirmationTime    time.Duration
	L1ETHAddress               string
	L1ERC20Address             string
	L1ERC721Address            string
	L1ERC1155Address           string
}

// Deposit Tests
func SetupForkedChain() error {
	return nil
}

func prepareTest() {
	// tranfer test Ether to testPublicKey

}
func TestDepositETH(t *testing.T) {
	t.Run("DepositETH", func(t *testing.T) {
		fmt.Println("DepositETH1")
		depositAmount := big.NewInt(100)
		//Arrange
		l1Client, err := ethclient.Dial(sepoliaHelpConfig.L1ClientAddress)
		if err != nil {
			log.Fatal("failed to connect to L1 geth", "endpoint", sepoliaHelpConfig.L1ClientAddress, "err", err)
		}
		l2Client, err := ethclient.Dial(sepoliaHelpConfig.L2ClientAddress)
		if err != nil {
			log.Fatal("failed to connect to L2 geth", "endpoint", sepoliaHelpConfig.L2ClientAddress, "err", err)
		}
		defer l1Client.Close()
		defer l2Client.Close()
		callOpts := &bind.CallOpts{
			Context: context.Background(),
		}
		ChildBridgeCoreFacet, err := bindings.NewChildBridgeCoreFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		if err != nil {
			log.Fatalf("failed to create ChildTokenMessageTransmitterFacet contract: %v", err)
		}
		//if this L2BridgeTokenAddress is not exist,need to register it at previous step
		l2BridgeTokenAddress, err := ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		if err != nil {
			log.Fatalf("failed to get bridged token address: %v", err)
		}
		var startBalance *big.Int
		if l2BridgeTokenAddress == (common.Address{}) {
			startBalance = big.NewInt(0)
		} else {
			l2BridgeERC20Token, err := bindings.NewERC20Token(l2BridgeTokenAddress, l2Client)
			if err != nil {
				log.Fatalf("failed to create ERC20Token contract: %v", err)
			}

			startBalance, err = l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
			if err != nil {
				log.Fatalf("failed to get balance of testPublicKey: %v", err)
			}
		}
		// Action
		// get gas price
		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal("failed to get gas price", "err", err)
		}
		chainid, err := l1Client.ChainID(context.Background())
		if err != nil {
			log.Fatal("failed to get chain id", "err", err)
		}
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env")
		if err != nil {
			log.Fatalf("Error loading private key: %v", err)
		}
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		if err != nil {
			log.Fatalf("failed to create private key %v", err)
		}

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		if err != nil {
			log.Fatalf("failed to create authorized transactor: %v", err)
		}
		auth.GasPrice = gasPrice

		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		if err != nil {
			log.Fatalf("failed to create ParentTokenMessageTransmitterFacet contract: %v", err)
		}
		auth.Value = depositAmount
		tx, err := ParentTokenMessageTransmitterFacet.DepositETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey), depositAmount, big.NewInt(0))
		if err != nil {
			log.Fatalf("failed to deposit eth: %v", err)
		}
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		if err != nil {
			log.Fatalf("failed to wait for confirmation: %v", err)
		}
		assert.True(t, success)

		// wait for the L2 confirmation
		time.Sleep(5 * time.Second)

		//Check the balance of the testPublicKey\\
		l2BridgeTokenAddress, err = ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ETHAddress))
		if err != nil {
			log.Fatalf("failed to get bridged token address: %v", err)
		}
		l2BridgeERC20Token, err := bindings.NewERC20Token(l2BridgeTokenAddress, l2Client)
		if err != nil {
			log.Fatalf("failed to create ERC20Token contract: %v", err)
		}
		fmt.Println("L2BridgeTokenAddress: ", l2BridgeTokenAddress)
		balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
		if err != nil {
			log.Fatalf("failed to get balance of testPublicKey: %v", err)
		}
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Add(depositAmount, startBalance)

		assert.Equal(t, expectedBalance, balance)
	})
}

func TestDepositERC20(t *testing.T) {
	t.Run("DepositERC20", func(t *testing.T) {
		fmt.Println("DepositERC20")
		depositAmount := big.NewInt(100) // 设置存款金额
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

			startBalance, err = l2RC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
			require.NoError(t, err)
		}

		// 获取 gas price
		gasPrice, err := l1Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l1Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		// 批准 ERC20 代币转移
		erc20Token, err := bindings.NewERC20Token(common.HexToAddress(sepoliaHelpConfig.L1ERC20Address), l1Client)
		require.NoError(t, err)
		tx, err := erc20Token.Approve(auth, common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), depositAmount)
		require.NoError(t, err)
		fmt.Println("Approve transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		// 存款 ERC20 代币
		ParentTokenMessageTransmitterFacet, err := bindings.NewParentTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ParentlayerContractAddress), l1Client)
		require.NoError(t, err)
		tx, err = ParentTokenMessageTransmitterFacet.DepositERC20Token(auth, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address), common.HexToAddress(sepoliaHelpConfig.testPublicKey), depositAmount, big.NewInt(0))
		require.NoError(t, err)
		fmt.Println("DepositERC20 transaction sent: ", tx.Hash().Hex())

		success, err = waitForConfirmation(l1Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		// 等待 L2 交易确认
		time.Sleep(5 * time.Second)

		// 检查 testPublicKey 的余额
		l2TokenAddress, err = ChildBridgeCoreFacet.GetBridgedERC20TokenChild(callOpts, common.HexToAddress(sepoliaHelpConfig.L1ERC20Address))
		require.NoError(t, err)
		fmt.Println("L2 ERC20 Token Address2: ", l2TokenAddress)
		l2RC20Token, err := bindings.NewERC20Token(l2TokenAddress, l2Client)
		require.NoError(t, err)
		balance, err := l2RC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Add(depositAmount, startBalance)

		assert.Equal(t, expectedBalance, balance)
	})
}

func TestDepositERC721(t *testing.T) {
	// Implement ERC721 deposit test
}

func TestDepositERC1155Batch(t *testing.T) {
	// Implement ERC1155 batch deposit test
}

// Withdraw Tests

func TestWithdrawETH(t *testing.T) {
	t.Run("WithdrawETH", func(t *testing.T) {
		fmt.Println("WithdrawETH")

		withdrawAmount := big.NewInt(50)

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

		startBalance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
		require.NoError(t, err)

		//TransferL2ETH(t, l2Client, sepoliaHelpConfig.testAdmin, common.HexToAddress(sepoliaHelpConfig.testPublicKey), big.NewInt(1e18))
		// Action
		gasPrice, err := l2Client.SuggestGasPrice(context.Background())
		require.NoError(t, err)
		chainid, err := l2Client.ChainID(context.Background())
		require.NoError(t, err)
		t.Log("gas price", "price", gasPrice)
		fmt.Println("gasPrice", gasPrice)
		privateKeyStr, err := utils.LoadPrivateKey("../test/.sepolia.env")
		require.NoError(t, err)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		require.NoError(t, err)

		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainid)
		require.NoError(t, err)
		auth.GasPrice = gasPrice

		ChildTokenMessageTransmitterFacet, err := bindings.NewChildTokenMessageTransmitterFacet(common.HexToAddress(sepoliaHelpConfig.ChildlayerContractAddress), l2Client)
		require.NoError(t, err)

		tx, err := ChildTokenMessageTransmitterFacet.WithdrawETH(auth, common.HexToAddress(sepoliaHelpConfig.testPublicKey), withdrawAmount)
		require.NoError(t, err)
		fmt.Println("Transaction sent: ", tx.Hash().Hex())

		success, err := waitForConfirmation(l2Client, tx.Hash())
		require.NoError(t, err)
		assert.True(t, success)

		// check testPublicKey L2eth balance
		balance, err := l2BridgeERC20Token.BalanceOf(callOpts, common.HexToAddress(sepoliaHelpConfig.testPublicKey))
		require.NoError(t, err)
		fmt.Println("Balance of testPublicKey: ", balance)
		expectedBalance := new(big.Int).Sub(startBalance, withdrawAmount)

		assert.Equal(t, expectedBalance, balance)
	})

}

func TestWithdrawERC20(t *testing.T) {
	// Implement ERC20 withdraw test
}

func TestWithdrawERC721(t *testing.T) {
	// Implement ERC721 withdraw test
}

func TestWithdrawERC1155Batch(t *testing.T) {
	// Implement ERC1155 batch withdraw test
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

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

	success, err := waitForConfirmation(l2Client, signedTx.Hash())
	require.NoError(t, err)
	assert.True(t, success, "L2 transaction was not confirmed")
}
