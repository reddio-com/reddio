package uniswap

import (
	"context"
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

type UniswapV2DeployedAccuracyContracts struct {
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

const swapTimes = 200

func (ca *UniswapV2AccuracyTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	//create a wallet for contract deployment
	wallets, err := m.GenerateRandomWallet(1, 1e18)
	if err != nil {
		return err
	}
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
	uniswapV2Contract, err := deployUniswapV2AccuracyContracts(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}

	//Arrange :
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

	isConfirmed, err = waitForConfirmation(client, tokenAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}

	WethAmountApproved := big.NewInt(1e18)
	WethAApproveTx, err := uniswapV2Contract.weth9Instance.Approve(auth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	isConfirmed, err = waitForConfirmation(client, WethAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}

	WethAApproveTx, err = uniswapV2Contract.weth9Instance.Approve(testUserAuth, common.HexToAddress(uniswapV2Contract.uniswapV2Router01Address.Hex()), WethAmountApproved)
	if err != nil {
		log.Fatalf("Failed to create approve transaction: %v", err)
	}

	isConfirmed, err = waitForConfirmation(client, WethAApproveTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Approve transaction was not confirmed")
	}

	//add ETH liquidity
	amountADesired := big.NewInt(1e18)
	auth.Value = big.NewInt(1e18)

	addLiquidityETHTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidityETH(auth, uniswapV2Contract.tokenAAddress, amountADesired, big.NewInt(0), big.NewInt(0), common.HexToAddress(wallets[0].Address), big.NewInt(time.Now().Unix()+1000))
	if err != nil {
		log.Fatalf("Failed to create add liquidity transaction: %v", err)
	}

	isConfirmed, err = waitForConfirmation(client, addLiquidityETHTx.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm add liquidity transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Add liquidity transaction was not confirmed")
	}

	// Perform  swaps
	const maxRetries = 300
	const retryDelay = 10 * time.Millisecond
	var retryErrors []struct {
		Nonce int
		Err   error
	}
	//Act
	for i := 0; i < swapTimes; i++ {
		// Set swap parameters
		testUserAuth.Value = big.NewInt(100)

		//Get output amount
		amounts, err := uniswapV2Contract.uniswapV2RouterInstance.GetAmountsOut(nil, big.NewInt(100), []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress})
		if err != nil {
			log.Fatalf("Failed to get amounts out: %v", err)
		}
		amountOut := amounts[len(amounts)-1]
		// Execute swap operation
		nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(testUser[0].Address))
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}

		for j := 0; j < maxRetries; j++ {
			swapETHForExactTokensTx, err := uniswapV2Contract.uniswapV2RouterInstance.SwapETHForExactTokens(testUserAuth, amountOut, []common.Address{uniswapV2Contract.weth9Address, uniswapV2Contract.tokenAAddress}, common.HexToAddress(testUser[0].Address), big.NewInt(time.Now().Unix()+1000))
			if err == nil {
				// Wait for transaction confirmation
				if i == (swapTimes - 1) {
					isConfirmed, err := waitForConfirmation(client, swapETHForExactTokensTx.Hash())
					if err != nil {
						log.Fatalf("Failed to confirm swapETHForExactTokensTx transaction: %v", err)
					}
					if !isConfirmed {
						log.Fatalf("SwapETHForExactTokens transaction was not confirmed")
					}
				}
				break
			}
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

	// Get account balance
	tokenABalance, err := uniswapV2Contract.tokenAInstance.BalanceOf(nil, common.HexToAddress(testUser[0].Address))
	if err != nil {
		log.Fatalf("Failed to get TokenA balance: %v", err)
	}

	ethBalance, err := client.BalanceAt(context.Background(), common.HexToAddress(testUser[0].Address), nil)
	if err != nil {
		log.Fatalf("Failed to get ETH balance: %v", err)
	}
	// Expect results
	// Initial state
	tokenAReserve := big.NewInt(1e18)
	ethReserve := big.NewInt(1e18)
	k := new(big.Int).Mul(tokenAReserve, ethReserve)

	// User initial state
	expectedEthBalance := big.NewInt(1e18)
	expectedTokenABalance := big.NewInt(0)

	swapEth := big.NewInt(100)

	for i := 0; i < swapTimes; i++ {
		eth2tokenAPrice := 99
		tokenAReceived := big.NewInt(int64(eth2tokenAPrice))

		// Update reserves
		ethReserve.Add(ethReserve, swapEth)
		tokenAReserve.Sub(k, ethReserve)

		// Update user balance
		expectedEthBalance.Sub(expectedEthBalance, swapEth)
		expectedTokenABalance.Add(expectedTokenABalance, tokenAReceived)
	}

	if tokenABalance.Cmp(expectedTokenABalance) != 0 {
		log.Fatalf("Expected user TokenA balance to be %s, but got %s", expectedTokenABalance.String(), tokenABalance.String())
	}

	if ethBalance.Cmp(expectedEthBalance) != 0 {
		log.Fatalf("Expected user ETH balance to be %s, but got %s", expectedEthBalance.String(), ethBalance.String())
	}
	return err
}

// deploy UniswapV2 Contracts
func deployUniswapV2AccuracyContracts(auth *bind.TransactOpts, client *ethclient.Client) (*UniswapV2DeployedAccuracyContracts, error) {
	var err error
	deployed := &UniswapV2DeployedAccuracyContracts{}

	// Deploy token A
	deployed.tokenAAddress, deployed.tokenATransaction, deployed.tokenAInstance, err = contracts.DeployToken(auth, client)
	if err != nil {
		return nil, err
	}

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

	isConfirmed, err = waitForConfirmation(client, deployed.uniswapV2Router01Transaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Transaction was not confirmed")
	}

	return deployed, nil
}
