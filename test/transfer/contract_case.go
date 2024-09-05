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

type UniswapV2DeployedContracts struct {
	tokenAAddress                common.Address
	tokenBAddress                common.Address
	weth9Address                 common.Address
	uniswapV2FactoryAddress      common.Address
	uniswapV2Router01Address     common.Address
	tokenATransaction            *types.Transaction
	tokenBTransaction            *types.Transaction
	weth9Transaction             *types.Transaction
	uniswapV2FactoryTransaction  *types.Transaction
	uniswapV2Router01Transaction *types.Transaction
	tokenAInstance               *contracts.Token
	tokenBInstance               *contracts.ERC20T
	weth9Instance                *contracts.WETH9
	uniswapV2FactoryInstance     *contracts.UniswapV2Factory
	uniswapV2RouterInstance      *contracts.UniswapV2Router01
}

type UniswapV2AccuracyTestCase struct {
	CaseName     string
	walletCount  int
	initialCount uint64
}

func (cd *UniswapV2AccuracyTestCase) Name() string {
	return cd.CaseName
}

func NewUniswapV2AccuracyTestCase(name string, count int, initial uint64) *UniswapV2AccuracyTestCase {
	return &UniswapV2AccuracyTestCase{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
	}
}

func (ca *UniswapV2AccuracyTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	log.Printf("Running %s", ca.CaseName)
	//create a wallet for contract deployment
	wallets, err := m.GenerateRandomWallet(1, 1e18)
	if err != nil {
		return err
	}
	//deployer wallet
	log.Printf("deployer wallet address: %s", wallets[0].Address)
	log.Printf("deployer test wallet pk: %s", wallets[0].PK) //Note: private key, should be kept secret on production environment

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

	// set tx auth
	privateKey, err := crypto.HexToECDSA(wallets[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(60000000)

	// deploy contracts
	uniswapV2Contract, err := deployUniswapV2Contracts(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	//Arrange :
	///
	testUser, err := m.GenerateRandomWallet(1, 1e18)
	if err != nil {
		return err
	}
	testUserPK, err := crypto.HexToECDSA(testUser[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	testUserAuth, err := bind.NewKeyedTransactorWithChainID(testUserPK, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	testUserAuth.GasPrice = gasPrice
	testUserAuth.GasLimit = uint64(60000000)

	amountApproved := big.NewInt(1e18)
	tokenAApproveTx, err := uniswapV2Contract.tokenAInstance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), amountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	log.Printf("tokenAApproveTx transaction hash: %s", tokenAApproveTx.Hash().Hex())

	isConfirmed, err := waitForConfirmation(client, tokenAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}
	tokenAApproveTx, err = uniswapV2Contract.tokenAInstance.Approve(testUserAuth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), amountApproved)
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
	log.Println("testUser tokenAApproveTx transaction confirmed")

	WethAmountApproved := big.NewInt(1e18)
	WethAApproveTx, err := uniswapV2Contract.weth9Instance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
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

	WethAApproveTx, err = uniswapV2Contract.weth9Instance.Approve(testUserAuth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
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
	log.Println("testUse WethAApproveTx transaction confirmed")

	//add ETH liquidity
	amountADesired := big.NewInt(1e18)
	auth.Value = big.NewInt(1e18)

	addLiquidityETHTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidityETH(auth, uniswapV2Contract.tokenAAddress, amountADesired, big.NewInt(0), big.NewInt(0), common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
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

	//amountOut := big.NewInt(99)
	// Perform  swaps
	const maxRetries = 300
	const retryDelay = 10 * time.Millisecond
	var retryErrors []struct {
		Nonce int
		Err   error
	}
	//Act
	for i := 0; i < 200; i++ {
		log.Printf("Swap %d", i)
		// Set swap parameters
		testUserAuth.Value = big.NewInt(100)

		//Get output amount
		amounts, err := uniswapV2Contract.uniswapV2RouterInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress})
		if err != nil {
			log.Fatalf("Failed to get amounts out: %v", err)
		}
		amountOut := amounts[len(amounts)-1]
		log.Printf("TokenA Amount out: %v", amountOut)
		// Execute swap operation
		nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(testUser[0].Address))
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}
		log.Printf("Nonce: %v", nonce)

		for j := 0; j < maxRetries; j++ {
			swapETHForExactTokensTx, err := uniswapV2Contract.uniswapV2RouterInstance.SwapETHForExactTokens(testUserAuth, amountOut, []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress}, common.HexToAddress(testUser[0].Address), big.NewInt(time.Now().Unix()+1000))
			if err == nil {
				log.Printf("SwapETHForExactTokens transaction hash: %s", swapETHForExactTokensTx.Hash().Hex())
				// Wait for transaction confirmation
				if i == 199 {
					isConfirmed, err := waitForConfirmation(client, swapETHForExactTokensTx.Hash())
					if err != nil {
						log.Fatalf("Failed to confirm swapETHForExactTokensTx transaction: %v", err)
					}
					if !isConfirmed {
						log.Fatalf("SwapETHForExactTokens transaction was not confirmed")
					}
					// log.Println("SwapETHForExactTokens transaction confirmed")
					//time.Sleep(5 * time.Second)
				}
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

	}
	// 获取 swap 以后的 tokenA 兑换 tokenB 的价格
	amounts, err := uniswapV2Contract.uniswapV2RouterInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress})
	if err != nil {
		log.Fatalf("Failed to get amounts out: %v", err)
	}
	tokenAPriceInTokenB := amounts[len(amounts)-1]
	log.Printf("TokenA price in TokenB: %v", tokenAPriceInTokenB)

	// 获取账户余额
	tokenABalance, err := uniswapV2Contract.tokenAInstance.BalanceOf(nil, common.HexToAddress(testUser[0].Address))
	if err != nil {
		log.Fatalf("Failed to get TokenA balance: %v", err)
	}
	log.Printf("Account TokenA balance: %v", tokenABalance)

	ethBalance, err := client.BalanceAt(context.Background(), common.HexToAddress(testUser[0].Address), nil)
	if err != nil {
		log.Fatalf("Failed to get ETH balance: %v", err)
	}
	log.Printf("Account ETH balance: %v", ethBalance)

	// expect results
	// 初始状态
	tokenAReserve := big.NewInt(1000000)
	ethReserve := big.NewInt(1000)
	k := new(big.Int).Mul(tokenAReserve, ethReserve)

	// 用户初始状态
	userEth := big.NewInt(0)
	userTokenA := big.NewInt(0)

	// 交易次数
	numSwaps := 200
	swapEth := big.NewInt(100)
	feeMultiplier := big.NewFloat(0.997)

	for i := 0; i < numSwaps; i++ {
		// 计算有效的 ETH 输入
		swapEthEffective := new(big.Float).Mul(new(big.Float).SetInt(swapEth), feeMultiplier)
		swapEthEffectiveInt, _ := swapEthEffective.Int(nil)

		// 更新 ETH 储备量
		newEthReserve := new(big.Int).Add(ethReserve, swapEthEffectiveInt)

		// 计算新的 TokenA 储备量
		newTokenAReserve := new(big.Int).Div(k, newEthReserve)

		// 计算用户获得的 TokenA 数量
		tokenAReceived := new(big.Int).Sub(tokenAReserve, newTokenAReserve)

		// 更新储备量
		tokenAReserve = newTokenAReserve
		ethReserve = newEthReserve

		// 更新用户余额
		userEth.Add(userEth, swapEth)
		userTokenA.Add(userTokenA, tokenAReceived)
	}

	// 输出结果
	fmt.Printf("用户账户中的 ETH: %s\n", userEth.String())
	fmt.Printf("用户账户中的 TokenA: %s\n", userTokenA.String())
	fmt.Printf("流动性池中的 TokenA 储备量: %s\n", tokenAReserve.String())
	fmt.Printf("流动性池中的 ETH 储备量: %s\n", ethReserve.String())

	// 计算最终的 TokenA 价格
	tokenAPrice := new(big.Float).Quo(new(big.Float).SetInt(ethReserve), new(big.Float).SetInt(tokenAReserve))
	fmt.Printf("流动性池中的 TokenA 价格: %f ETH\n", tokenAPrice)

	return err
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

func (cd *UniswapV2TPSStatisticsTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	log.Printf("Running %s", cd.CaseName)
	//create a wallet for contract deployment
	wallets, err := m.GenerateRandomWallet(1, 100000000000000000)
	if err != nil {
		return err
	}
	//deployer wallet
	log.Printf("deployer wallet address: %s", wallets[0].Address)
	log.Printf("deployer test wallet pk: %s", wallets[0].PK) //Note: private key, should be kept secret on production environment

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

	// set tx auth
	privateKey, err := crypto.HexToECDSA(wallets[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(50341))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(60000000)

	// deploy contracts
	uniswapV2Contract, err := deployUniswapV2Contracts(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	// interact with the contract
	// ETH <-> TOKEN A
	amountApproved := big.NewInt(1000000000000000000)
	tokenAApproveTx, err := uniswapV2Contract.tokenAInstance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), amountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	log.Printf("tokenAApproveTx transaction hash: %s", tokenAApproveTx.Hash().Hex())

	isConfirmed, err := waitForConfirmation(client, tokenAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}
	log.Println("tokenAApproveTx transaction confirmed")

	WethAmountApproved := big.NewInt(1000000000000000000)
	WethAApproveTx, err := uniswapV2Contract.weth9Instance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
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
	addLiquidityETHTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidityETH(auth, uniswapV2Contract.tokenAAddress, amountADesired, big.NewInt(0), big.NewInt(0), common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
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
		amounts, err := uniswapV2Contract.uniswapV2RouterInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress})
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

		for j := 0; j < maxRetries; j++ {
			swapETHForExactTokensTx, err := uniswapV2Contract.uniswapV2RouterInstance.SwapETHForExactTokens(auth, amountOut, []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress}, common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
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
	for j, retryErr := range retryErrors {
		log.Printf("Attempt %d: Nonce %d, Error: %v", j+1, retryErr.Nonce, retryErr.Err)
	}
	duration := endTime.Sub(startTime).Seconds()
	tps := float64(2000) / duration
	log.Printf("TPS: %.2f", tps)

	return err
}

// deploy UniswapV2 Contracts
/*
   Deploy token A ("AAAToken", "AAA")
   Deploy token B ("BBBToken", "BBB")
   Deploy WETH
   Deploy UniswapV2Factory (FeeToSetter)
   Deploy UniswapV2Router01 (WETH addresse, factory addresse)
*/
func deployUniswapV2Contracts(auth *bind.TransactOpts, client *ethclient.Client) (*UniswapV2DeployedContracts, error) {
	var err error
	deployed := &UniswapV2DeployedContracts{}

	// Deploy token A
	deployed.tokenAAddress, deployed.tokenATransaction, deployed.tokenAInstance, err = contracts.DeployToken(auth, client)
	if err != nil {
		return nil, err
	}
	log.Printf("TokenA deployed at address: %s", deployed.tokenAAddress.Hex())

	isConfirmed, err := waitForConfirmation(client, deployed.tokenATransaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// Deploy token B
	deployed.tokenBAddress, deployed.tokenBTransaction, deployed.tokenBInstance, err = contracts.DeployERC20T(auth, client, "BBBToken", "BBB")
	if err != nil {
		return nil, err
	}
	log.Printf("TokenB deployed at address: %s", deployed.tokenBAddress.Hex())

	isConfirmed, err = waitForConfirmation(client, deployed.tokenBTransaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// Deploy WETH9
	deployed.weth9Address, deployed.weth9Transaction, deployed.weth9Instance, err = contracts.DeployWETH9(auth, client)
	if err != nil {
		return nil, err
	}
	log.Printf("WETH deployed at address: %s", deployed.weth9Address.Hex())

	isConfirmed, err = waitForConfirmation(client, deployed.weth9Transaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// Deploy UniswapV2Factory
	deployed.uniswapV2FactoryAddress, deployed.uniswapV2FactoryTransaction, deployed.uniswapV2FactoryInstance, err = contracts.DeployUniswapV2Factory(auth, client, auth.From)
	if err != nil {
		return nil, err
	}
	log.Printf("UniswapV2Factory deployed at address: %s", deployed.uniswapV2FactoryAddress.Hex())

	isConfirmed, err = waitForConfirmation(client, deployed.uniswapV2FactoryTransaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	// Deploy UniswapV2Router01
	deployed.uniswapV2Router01Address, deployed.uniswapV2Router01Transaction, deployed.uniswapV2RouterInstance, err = contracts.DeployUniswapV2Router01(auth, client, deployed.uniswapV2FactoryAddress, deployed.weth9Address)
	if err != nil {
		return nil, err
	}
	log.Printf("UniswapV2Router01 deployed at address: %s", deployed.uniswapV2Router01Address.Hex())

	isConfirmed, err = waitForConfirmation(client, deployed.uniswapV2Router01Transaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	return deployed, nil
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
