package transfer

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

type ContractDeploymentTestCase struct {
	CaseName     string
	walletCount  int
	initialCount uint64
}

func (cd *ContractDeploymentTestCase) Name() string {
	return cd.CaseName
}

func (cd *ContractDeploymentTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	//create a wallet for contract deployment
	wallets, err := m.GenerateRandomWallet(1, 100000000000000000)

	log.Println("initialCount:", cd.initialCount)
	if err != nil {
		return err
	}
	log.Printf("%s create wallets finish", cd.CaseName)

	//deployer wallet
	log.Printf("deployer wallet address: %s", wallets[0].Address)
	log.Printf("deployer wallet pk: %s", wallets[0].PK)

	client, err := ethclient.Dial("http://localhost:9092")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	defer client.Close()
	if err != nil {
		log.Fatalf("Failed to Close  the Ethereum client: %v", err)
	}

	balance, err := client.BalanceAt(ctx, common.HexToAddress(wallets[0].Address), nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	log.Printf("deployer wallet balance: %s", balance.String())
	// get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas price: %v", err)
	}
	log.Printf("Gas price: %v", gasPrice)

	// get the current nonce
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(wallets[0].Address))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	log.Printf("Nonce: %v", nonce)
	// set tx auth
	privateKey, err := crypto.HexToECDSA(wallets[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(50341)) // 使用你的私钥和链ID
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(60000000)

	// deploy contract
	/*
	   Deploy token A ("AAAToken", "AAA")
	   Deploy token B ("BBBToken", "BBB")
	   Deploy WETH
	   Deploy UniswapV2Factory (FeeToSetter)
	   Deploy UniswapV2Router01 (WETH addresse, factory addresse)
	*/

	// deploy token A
	tokenAAddress, TokenATx, tokenAInstance, err := contracts.DeployToken(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	log.Printf("TokenA deployed at address: %s", tokenAAddress.Hex())
	log.Printf("Transaction hash: %s", TokenATx.Hash().Hex())
	log.Printf("Contract instance: %v", tokenAInstance)

	isConfirmed, err := waitForConfirmation(client, TokenATx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// deploy token B ("BBB Token", "BBB")
	TokenBAddress, TokenBTx, tokenBInstance, err := contracts.DeployERC20T(auth, client, "BBBToken", "BBB")
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	log.Printf("TokenB deployed at address: %s", TokenBAddress.Hex())
	log.Printf("Transaction hash: %s", TokenBTx.Hash().Hex())
	log.Printf("Contract instance: %v", tokenBInstance)

	isConfirmed, err = waitForConfirmation(client, TokenBTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// deploy WETH
	WETHAddress, WETHTx, WETHInstance, err := contracts.DeployWETH9(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	log.Printf("WETH deployed at address: %s", WETHAddress.Hex())
	log.Printf("Transaction hash: %s", WETHTx.Hash().Hex())
	log.Printf("Contract instance: %v", WETHInstance)

	isConfirmed, err = waitForConfirmation(client, WETHTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// deploy UniswapV2Factory
	uniswapV2FactoryAddress, uniswapV2FactoryTx, uniswapV2FactoryInstance, err := contracts.DeployUniswapV2Factory(auth, client, common.HexToAddress(wallets[0].Address))
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Printf("UniswapV2Factory deployed at address: %s", uniswapV2FactoryAddress.Hex())
	log.Printf("Transaction hash: %s", uniswapV2FactoryTx.Hash().Hex())
	log.Printf("Contract instance: %v", uniswapV2FactoryInstance)

	isConfirmed, err = waitForConfirmation(client, uniswapV2FactoryTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// deploy UniswapV2Router01
	uniswapV2Router01Address, uniswapV2RouterTx, routerInstance, err := contracts.DeployUniswapV2Router01(auth, client, uniswapV2FactoryAddress, WETHAddress)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	log.Printf("UniswapV2Router01 deployed at address: %s", uniswapV2Router01Address.Hex())
	log.Printf("Transaction hash: %s", uniswapV2RouterTx.Hash().Hex())
	log.Printf("Contract instance: %v", routerInstance)

	// wait for confirmation
	isConfirmed, err = waitForConfirmation(client, uniswapV2RouterTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// interact with the contract
	// ETH <-> TOKEN A
	//approve 1000000000000000000 token A to uniswapV2Router01Address
	amountApproved := big.NewInt(1000000000000000000)
	tokenAApproveTx, err := tokenAInstance.Approve(auth, common.HexToAddress(uniswapV2Router01Address.Hex()), amountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	log.Printf("tokenAApproveTx transaction hash: %s", tokenAApproveTx.Hash().Hex())

	isConfirmed, err = waitForConfirmation(client, tokenAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}
	log.Println("tokenAApproveTx transaction confirmed")

	WethAmountApproved := big.NewInt(1000000000000000000)
	WethAApproveTx, err := WETHInstance.Approve(auth, common.HexToAddress(uniswapV2Router01Address.Hex()), WethAmountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	log.Printf("WethAApproveTx transaction hash: %s", WethAApproveTx.Hash().Hex())

	isConfirmed, err = waitForConfirmation(client, WethAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}
	log.Println("WethAApproveTx transaction confirmed")
	//add ETH liquidity
	amountADesired := big.NewInt(1000000)
	auth.Value = big.NewInt(1000)
	addLiquidityETHTx, err := routerInstance.AddLiquidityETH(auth, tokenAAddress, amountADesired, big.NewInt(0), big.NewInt(0), common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
	if err != nil {
		log.Fatalf("Failed to create add liquidity transaction: %v", err)
	}
	log.Printf("Add liquidity transaction hash: %s", addLiquidityETHTx.Hash().Hex())

	isConfirmed, err = waitForConfirmation(client, addLiquidityETHTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm add liquidity transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Add liquidity transaction was not confirmed")
	}
	log.Println("Add liquidity transaction confirmed")

	//swap from ETH to token A
	startTime := time.Now()
	log.Printf("Start time: %s", startTime.Format(time.RFC3339))
	//amountOut := big.NewInt(99)
	// Perform  swaps
	const maxRetries = 300
	const retryDelay = 10 * time.Millisecond
	var retryErrors []struct {
		Nonce int
		Err   error
	}
	for i := 0; i < 2000; i++ {
		log.Printf("Swap %d", i)
		// Set swap parameters
		auth.Value = big.NewInt(100)

		//Get output amount
		amounts, err := routerInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{WETHAddress, tokenAAddress})
		if err != nil {
			log.Fatalf("Failed to get amounts out: %v", err)
		}
		amountOut := amounts[len(amounts)-1]
		log.Printf("TokenA Amount out: %v", amountOut)
		// Execute swap operation
		nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(wallets[0].Address))
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}
		log.Printf("Nonce: %v", nonce)

		for i := 0; i < maxRetries; i++ {
			swapETHForExactTokensTx, err := routerInstance.SwapETHForExactTokens(auth, amountOut, []common.Address{WETHAddress, tokenAAddress}, common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
			if err == nil {
				log.Printf("SwapETHForExactTokens transaction hash: %s", swapETHForExactTokensTx.Hash().Hex())
				break
			}
			log.Printf("Attempt %d: Failed to swapETHForExactTokensTx transaction: %v", i+1, err)
			retryErrors = append(retryErrors, struct {
				Nonce int
				Err   error
			}{Nonce: int(nonce), Err: err}) // record和 nonce
			time.Sleep(retryDelay)
		}

		if err != nil {
			log.Fatalf("Failed to swapETHForExactTokensTx transaction after %d attempts: %v", maxRetries, err)
		}

		// if (i+1)%540 == 0 {
		// 	amountOut.Sub(amountOut, big.NewInt(1))
		// 	log.Printf("Decreased amountOut to: %v", amountOut)
		// }
		// Wait for transaction confirmation
		// isConfirmed, err := waitForConfirmation(client, swapETHForExactTokensTx.Hash())
		// if err != nil {
		// 	log.Fatalf("Failed to confirm swapETHForExactTokensTx transaction: %v", err)
		// }
		// if !isConfirmed {
		// 	log.Fatalf("SwapETHForExactTokens transaction was not confirmed")
		// }
		// log.Println("SwapETHForExactTokens transaction confirmed")
	}
	endTime := time.Now()
	log.Printf("End time: %s", endTime.Format(time.RFC3339))

	// Calculate TPS
	duration := endTime.Sub(startTime).Seconds()
	tps := float64(2000) / duration
	log.Printf("TPS: %.2f", tps)

	log.Println("All retry attempts failed. Error reasons and nonces:")
	for j, retryErr := range retryErrors {
		log.Printf("Attempt %d: Nonce %d, Error: %v", j+1, retryErr.Nonce, retryErr.Err)
	}
	// Swap from TokenA to ETH
	// Perform 10 swaps
	// for i := 0; i < 10; i++ {
	// 	log.Printf("Swap %d: TokenA to ETH", i)
	// 	// Set swap parameters
	// 	amountIn := big.NewInt(100)

	// 	// Get output amount
	// 	amounts, err := routerInstance.GetAmountsOut(nil, amountIn, []common.Address{tokenAAddress, WETHAddress})
	// 	if err != nil {
	// 		log.Fatalf("Failed to get amounts out: %v", err)
	// 	}
	// 	amountOut := amounts[len(amounts)-1]
	// 	log.Printf("ETH Amount out: %v", amountOut)
	// 	// Execute swap operation
	// 	swapTokensForExactETHTx, err := routerInstance.SwapExactTokensForETH(auth, amountIn, amountOut, []common.Address{tokenAAddress, WETHAddress}, common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
	// 	if err != nil {
	// 		log.Fatalf("Failed to swapTokensForExactETHTx transaction: %v", err)
	// 	}
	// 	log.Printf("SwapTokensForExactETH transaction hash: %s", swapTokensForExactETHTx.Hash().Hex())

	// 	// Wait for transaction confirmation
	// 	isConfirmed, err := waitForConfirmation(client, swapTokensForExactETHTx.Hash())
	// 	if err != nil {
	// 		log.Fatalf("Failed to confirm swapTokensForExactETHTx transaction: %v", err)
	// 	}
	// 	if !isConfirmed {
	// 		log.Fatalf("SwapTokensForExactETH transaction was not confirmed")
	// 	}
	// 	log.Println("SwapTokensForExactETH transaction confirmed")

	return err
}

func waitForConfirmation(client *ethclient.Client, txHash common.Hash) (bool, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				log.Printf("Transaction confirmed: %s", txHash.Hex())
				return true, nil
			}
			log.Printf("Transaction failed: %s", txHash.Hex())
			return false, fmt.Errorf("transaction failed with status: %v", receipt.Status)
		}
		time.Sleep(1 * time.Second)
	}
}

func NewContractDeploymentTest(name string, count int, initial uint64) *ContractDeploymentTestCase {
	return &ContractDeploymentTestCase{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
	}
}
