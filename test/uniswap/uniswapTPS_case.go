package uniswap

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/test/contracts"
	"github.com/reddio-com/reddio/test/pkg"
)

const (
	numTestUsers            = 100
	accountInitialFunds     = 1e18
	GasLimit                = 21000
	maxRetries              = 300
	waitForConfirmationTime = 1 * time.Second
)

type ERC20DeployedContracts struct {
	tokenAddress     common.Address
	tokenTransaction *types.Transaction
	tokenInstance    *contracts.Token
}

type UniswapV2DeployedContracts struct {
	weth9Address                 common.Address
	uniswapV2FactoryAddress      common.Address
	uniswapV2Router01Address     common.Address
	weth9Transaction             *types.Transaction
	uniswapV2FactoryTransaction  *types.Transaction
	uniswapV2Router01Transaction *types.Transaction
	weth9Instance                *contracts.WETH9
	uniswapV2FactoryInstance     *contracts.UniswapV2Factory
	uniswapV2RouterInstance      *contracts.UniswapV2Router01
}

type UniswapV2TPSStatisticsTestCase struct {
	CaseName     string
	walletCount  int
	initialCount uint64
}

func (cd *UniswapV2TPSStatisticsTestCase) Name() string {
	return cd.CaseName
}

func NewUniswapV2TPSStatisticsTestCase(name string, count int, initial uint64) *UniswapV2TPSStatisticsTestCase {
	return &UniswapV2TPSStatisticsTestCase{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
	}
}

// TestUniswapTPS is a test case to measure the transactions per second (TPS) of Uniswap.
// The test case follows these steps:
// 1. Arrange
//   - Create a deployer user with sufficient balance
//   - Create multiple test users with initial balance
//   - Deploy Uniswap core contracts (WETH9, UniswapV2Factory, UniswapV2Router01)
//   - Deploy a set of ERC20 token contracts
//   - Approve the router contract to spend tokens on behalf of test users
//   - Distribute the test tokens to each test user
//   - Generate all possible token pairs from the deployed ERC20 tokens
//   - Add liquidity to the token pairs on Uniswap
//
// 2. Act
//   - Each test user performs a series of token swaps on Uniswap
//   - The swaps are performed concurrently to simulate real-world usage
//
// 3. Assert
//   - Calculate and report the transactions per second (TPS) achieved during the test
func (cd *UniswapV2TPSStatisticsTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	log.Printf("Running %s", cd.CaseName)
	depolyerUser, err := m.GenerateRandomWallet(1, accountInitialFunds)
	if err != nil {
		log.Fatalf("Failed to generate deployer user: %v", err)
		return err
	}
	//Fixme: should be removed in production environment
	log.Printf("deployer wallet address: %s", depolyerUser[0].Address)
	log.Printf("deployer test wallet pk: %s", depolyerUser[0].PK) //Note: private key, should be kept secret on production environment

	testUsersWallets, err := m.GenerateRandomWallet(100, accountInitialFunds)
	if err != nil {
		log.Fatalf("Failed to generate test users: %v", err)
		return err
	}
	log.Printf("testUsersWallets length: %d", len(testUsersWallets))
	//Fixme: add this to the config file
	client, err := ethclient.Dial("http://localhost:9092")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	defer client.Close()
	if err != nil {
		log.Fatalf("Failed to Close  the Ethereum client: %v", err)
	}

	// get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())

	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	log.Printf("Gas price: %v", gasPrice)

	// set tx auth
	privateKey, err := crypto.HexToECDSA(depolyerUser[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	depolyerAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	depolyerAuth.GasPrice = gasPrice
	depolyerAuth.GasLimit = uint64(6e7)
	depolyerNonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(depolyerUser[0].Address))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	depolyerAuth.Nonce = big.NewInt(int64(depolyerNonce))

	// deploy contracts
	uniswapV2Contract, err := deployUniswapV2Contracts(depolyerAuth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	log.Printf("UniswapV2 contracts deployed :%s", uniswapV2Contract.uniswapV2Router01Address)
	//Fixme: add depployNum to the config file
	ERC20DeployedContracts, err := deployERC20Contracts(depolyerAuth, client, 10)
	if err != nil {
		log.Fatalf("Failed to deploy ERC20 contracts: %v", err)
	}
	lastIndex := len(ERC20DeployedContracts) - 1
	log.Printf("ERC20 contracts deployed, the last tokenAddress: %s", ERC20DeployedContracts[lastIndex].tokenAddress.Hex())
	isConfirmed, err := waitForConfirmation(client, ERC20DeployedContracts[lastIndex].tokenTransaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf(" transaction was not confirmed")
	}
	log.Printf("wait for ERC20 contract deployment done")
	// interact with the contract
	// ETH <-> TOKEN A
	// amountApproved := big.NewInt(1e18)
	// tokenAApproveTx, err := uniswapV2Contract.tokenAInstance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), amountApproved)
	// if err != nil {
	// 	log.Fatalf("Failed to create approve transaction: %v", err)
	// }

	// log.Printf("tokenAApproveTx transaction hash: %s", tokenAApproveTx.Hash().Hex())

	// isConfirmed, err := waitForConfirmation(client, tokenAApproveTx.Hash())
	// if err != nil {
	// 	log.Fatalf("Failed to confirm approve transaction: %v", err)
	// }
	// if !isConfirmed {
	// 	log.Fatalf("Approve transaction was not confirmed")
	// }
	// log.Println("tokenAApproveTx transaction confirmed")

	// WethAmountApproved := big.NewInt(1e18)
	// WethAApproveTx, err := uniswapV2Contract.weth9Instance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
	// if err != nil {
	// 	log.Fatalf("Failed to create approve transaction: %v", err)
	// }

	// log.Printf("WethAApproveTx transaction hash: %s", WethAApproveTx.Hash().Hex())

	// isConfirmed, err = waitForConfirmation(client, WethAApproveTx.Hash())
	// if err != nil {
	// 	log.Fatalf("Failed to confirm approve transaction: %v", err)
	// }
	// if !isConfirmed {
	// 	log.Fatalf("Approve transaction was not confirmed")
	// }
	// log.Println("WethAApproveTx transaction confirmed")

	// //add ETH liquidity
	// amountADesired := big.NewInt(1e18)
	// auth.Value = big.NewInt(1e18)
	// addLiquidityETHTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidityETH(auth, uniswapV2Contract.tokenAAddress, amountADesired, big.NewInt(0), big.NewInt(0), common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
	// if err != nil {
	// 	log.Fatalf("Failed to create add liquidity transaction: %v", err)
	// }
	// log.Printf("Add liquidity transaction hash: %s", addLiquidityETHTx.Hash().Hex())

	// isConfirmed, err = waitForConfirmation(client, addLiquidityETHTx.Hash())
	// if err != nil {
	// 	log.Fatalf("Failed to confirm add liquidity transaction: %v", err)
	// }
	// if !isConfirmed {
	// 	log.Fatalf("Add liquidity transaction was not confirmed")
	// }
	// log.Println("Add liquidity transaction confirmed")

	// //swap from ETH to token A
	// startTime := time.Now()
	// log.Printf("Start time: %s", startTime.Format(time.RFC3339))
	// //amountOut := big.NewInt(99)
	// // Perform  swaps
	// const maxRetries = 300
	// const retryDelay = 10 * time.Millisecond
	// var retryErrors []struct {
	// 	Nonce int
	// 	Err   error
	// }
	// for i := 0; i < 2000; i++ {
	// 	log.Printf("Swap %d", i)
	// 	// Set swap parameters
	// 	auth.Value = big.NewInt(100)

	// 	//Get output amount
	// 	amounts, err := uniswapV2Contract.uniswapV2RouterInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress})
	// 	if err != nil {
	// 		log.Fatalf("Failed to get amounts out: %v", err)
	// 	}
	// 	amountOut := amounts[len(amounts)-1]
	// 	log.Printf("TokenA Amount out: %v", amountOut)
	// 	// Execute swap operation
	// 	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(wallets[0].Address))
	// 	if err != nil {
	// 		log.Fatalf("Failed to get nonce: %v", err)
	// 	}
	// 	log.Printf("Nonce: %v", nonce)

	// 	for j := 0; j < maxRetries; j++ {
	// 		swapETHForExactTokensTx, err := uniswapV2Contract.uniswapV2RouterInstance.SwapETHForExactTokens(auth, amountOut, []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress}, common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
	// 		if err == nil {
	// 			log.Printf("SwapETHForExactTokens transaction hash: %s", swapETHForExactTokensTx.Hash().Hex())
	// 			break
	// 		}
	// 		log.Printf("Attempt %d: Failed to swapETHForExactTokensTx transaction: %v", i+1, err)
	// 		retryErrors = append(retryErrors, struct {
	// 			Nonce int
	// 			Err   error
	// 		}{Nonce: int(nonce), Err: err}) // record和 nonce
	// 		time.Sleep(retryDelay)
	// 	}

	// 	if err != nil {
	// 		log.Fatalf("Failed to swapETHForExactTokensTx transaction after %d attempts: %v", maxRetries, err)
	// 	}

	// 	// if (i+1)%540 == 0 {
	// 	// 	amountOut.Sub(amountOut, big.NewInt(1))
	// 	// 	log.Printf("Decreased amountOut to: %v", amountOut)
	// 	// }
	// 	// Wait for transaction confirmation
	// 	// isConfirmed, err := waitForConfirmation(client, swapETHForExactTokensTx.Hash())
	// 	// if err != nil {
	// 	// 	log.Fatalf("Failed to confirm swapETHForExactTokensTx transaction: %v", err)
	// 	// }
	// 	// if !isConfirmed {
	// 	// 	log.Fatalf("SwapETHForExactTokens transaction was not confirmed")
	// 	// }
	// 	// log.Println("SwapETHForExactTokens transaction confirmed")
	// }
	// endTime := time.Now()
	// log.Printf("End time: %s", endTime.Format(time.RFC3339))

	// // Calculate TPS
	// for j, retryErr := range retryErrors {
	// 	log.Printf("Attempt %d: Nonce %d, Error: %v", j+1, retryErr.Nonce, retryErr.Err)
	// }
	// duration := endTime.Sub(startTime).Seconds()
	// tps := float64(2000) / duration
	// log.Printf("TPS: %.2f", tps)

	return err
}

// deploy Erc20 token contracts
func deployERC20Contracts(auth *bind.TransactOpts, client *ethclient.Client, deployNum int) ([]*ERC20DeployedContracts, error) {

	var err error
	deployedTokens := make([]*ERC20DeployedContracts, 0)

	for i := 0; i < deployNum; i++ {
		deployedToken := &ERC20DeployedContracts{}
		deployedToken.tokenAddress, deployedToken.tokenTransaction, deployedToken.tokenInstance, err = contracts.DeployToken(auth, client)
		if err != nil {
			return nil, err
		}
		log.Printf("Token deployed at address: %s", deployedToken.tokenAddress.Hex())
		log.Printf("Token deployed txHash: %s", deployedToken.tokenTransaction.Hash().Hex())

		deployedTokens = append(deployedTokens, deployedToken)
		auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	}

	return deployedTokens, nil

}

// deploy UniswapV2 Contracts
/*
   Deploy WETH
   Deploy UniswapV2Factory (FeeToSetter)
   Deploy UniswapV2Router01 (WETH addresse, factory addresse)
*/
func deployUniswapV2Contracts(auth *bind.TransactOpts, client *ethclient.Client) (*UniswapV2DeployedContracts, error) {
	var err error
	deployed := &UniswapV2DeployedContracts{}

	// Deploy WETH9
	deployed.weth9Address, deployed.weth9Transaction, deployed.weth9Instance, err = contracts.DeployWETH9(auth, client)
	if err != nil {
		return nil, err
	}
	log.Printf("WETH deployed at address: %s", deployed.weth9Address.Hex())
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	// Deploy UniswapV2Factory
	deployed.uniswapV2FactoryAddress, deployed.uniswapV2FactoryTransaction, deployed.uniswapV2FactoryInstance, err = contracts.DeployUniswapV2Factory(auth, client, auth.From)
	if err != nil {
		return nil, err
	}
	log.Printf("UniswapV2Factory deployed at address: %s", deployed.uniswapV2FactoryAddress.Hex())
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	// Deploy UniswapV2Router01
	deployed.uniswapV2Router01Address, deployed.uniswapV2Router01Transaction, deployed.uniswapV2RouterInstance, err = contracts.DeployUniswapV2Router01(auth, client, deployed.uniswapV2FactoryAddress, deployed.weth9Address)
	if err != nil {
		return nil, err
	}
	log.Printf("UniswapV2Router01 deployed at address: %s", deployed.uniswapV2Router01Address.Hex())
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	return deployed, nil
}

func waitForConfirmation(client *ethclient.Client, txHash common.Hash) (bool, error) {
	log.Printf("Waiting for transaction to be confirmed: %s", txHash.Hex())
	for i := 0; i < maxRetries; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				log.Printf("Transaction confirmed: %s", txHash.Hex())
				return true, nil
			}
			log.Printf("Transaction failed: %s", txHash.Hex())
			return false, fmt.Errorf("transaction failed with status: %v", receipt.Status)
		}
		time.Sleep(waitForConfirmationTime)
	}
	return false, fmt.Errorf("transaction was not confirmed after %d retries", maxRetries)
}
