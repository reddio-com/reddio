package uniswap

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/reddio-com/reddio/test/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/reddio-com/reddio/test/pkg"
)

const (
	nodeUrl                      = "http://localhost:9092"
	numTestUsers                 = 100
	accountInitialFunds          = 1e18
	gasLimit                     = 6e7
	chainID                      = 50341
	accountInitialERC20Token     = 1e18
	approveAmount                = 1e18
	amountADesired               = 1e15
	amountBDesired               = 1e15
	maxSwapAmount                = 1e9
	maxBlocks                    = 20
	allowFailedTransactionsCount = 10
	stepCount                    = 5000
	retriesInterval              = 3 * time.Second
	tokenContractNum             = 100
)

type UniswapV2TPSStatisticsTestCase struct {
	CaseName string
}

type TPSStats struct {
	TPS              float64
	BlockCount       int
	TimeInterval     uint64
	TransactionCount int
	StartBlockNumber *big.Int
}

func (cd *UniswapV2TPSStatisticsTestCase) Name() string {
	return cd.CaseName
}

func NewUniswapV2TPSStatisticsTestCase(name string) *UniswapV2TPSStatisticsTestCase {

	return &UniswapV2TPSStatisticsTestCase{
		CaseName: name,
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
	if _, err := os.Stat("test/tmp"); os.IsNotExist(err) {
		Prepare(m)
	}
	err := executeTestAndCalculateTPS(nodeUrl, chainID, gasLimit, stepCount, allowFailedTransactionsCount)
	if err != nil {
		log.Fatalf("Failed to execute test and calculate TPS: %v", err)
	}
	return err
}
func Prepare(m *pkg.WalletManager) error {

	deployerUser, err := m.GenerateRandomWallet(1, accountInitialFunds)
	if err != nil {
		log.Fatalf("Failed to generate deployer user: %v", err)
		return err
	}
	testUsers, err := m.GenerateRandomWallet(100, accountInitialFunds)
	if err != nil {
		log.Fatalf("Failed to generate test users: %v", err)
		return err
	}
	var lastTxHash common.Hash
	client, err := ethclient.Dial(nodeUrl)
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
	privateKey, err := crypto.HexToECDSA(deployerUser[0].PK)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	depolyerAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	depolyerAuth.GasPrice = gasPrice
	depolyerAuth.GasLimit = uint64(6e7)
	depolyerNonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(deployerUser[0].Address))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	depolyerAuth.Nonce = big.NewInt(int64(depolyerNonce))

	// deploy contracts
	uniswapV2Contract, err := deployUniswapV2Contracts(depolyerAuth, client)
	if err != nil {
		log.Fatalf("Failed to deploy contract: %v", err)
	}
	ERC20DeployedContracts, err := deployERC20Contracts(depolyerAuth, client, tokenContractNum)
	if err != nil {
		log.Fatalf("Failed to deploy ERC20 contracts: %v", err)
	}
	lastIndex := len(ERC20DeployedContracts) - 1
	isConfirmed, err := waitForConfirmation(client, ERC20DeployedContracts[lastIndex].tokenTransaction.Hash())
	if err != nil {
		log.Fatalf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf(" transaction was not confirmed")
	}
	//interact with the contract
	///dispatchTestToken
	// dispatchTestToken([] TokenAddresses,testUsers)
	err = dispatchTestToken(client, depolyerAuth, ERC20DeployedContracts, testUsers, big.NewInt(accountInitialERC20Token))
	if err != nil {
		log.Fatalf("Failed to dispatch test tokens: %v", err)
	}

	///approve
	// approve (weth9Address,testUsers)
	// approve ([] TokenAddresses,testUsers)

	for _, contract := range ERC20DeployedContracts {
		_, err := contract.tokenInstance.Approve(depolyerAuth, uniswapV2Contract.uniswapV2Router01Address, big.NewInt(approveAmount))
		if err != nil {
			log.Fatalf("Failed to create approve transaction for user %s: %v", deployerUser[0].Address, err)
		}

		depolyerAuth.Nonce = depolyerAuth.Nonce.Add(depolyerAuth.Nonce, big.NewInt(1))
		if err != nil {
			return fmt.Errorf("failed to generate test auth for user %s: %v", deployerUser[0].Address, err)
		}
		for _, user := range testUsers {
			testAuth, err := generateTestAuth(client, user, chainID, gasPrice, gasLimit)
			if err != nil {
				return fmt.Errorf("failed to generate test auth for user %s: %v", user.Address, err)
			}
			tx, err := contract.tokenInstance.Approve(testAuth, uniswapV2Contract.uniswapV2Router01Address, big.NewInt(approveAmount))
			lastTxHash = tx.Hash()
			if err != nil {
				return fmt.Errorf("failed to create approve transaction for user %s: %v", user.Address, err)
			}
			//log.Printf("Approve transaction hash for user %s: %s", user.Address, tx.Hash().Hex())
			testAuth.Nonce = testAuth.Nonce.Add(testAuth.Nonce, big.NewInt(1))
		}
	}
	isConfirmed, err = waitForConfirmation(client, lastTxHash)
	if err != nil {
		return err
	}
	if !isConfirmed {
		return fmt.Errorf("transaction %s was not confirmed", lastTxHash.Hex())
	}
	///generateTokenPairs
	//// C(TokenAddresses.size,2)
	//generateTokenPairs([]TokenAddresses) return Pairs
	tokenPairs := generateTokenPairs(ERC20DeployedContracts)

	//add liquidity
	for i, pair := range tokenPairs {
		addLiquidityTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidity(
			depolyerAuth,
			pair[0],
			pair[1],
			big.NewInt(amountADesired),
			big.NewInt(amountBDesired),
			big.NewInt(0),
			big.NewInt(0),
			common.HexToAddress(deployerUser[0].Address),
			big.NewInt(time.Now().Unix()+1000),
		)
		if err != nil {
			log.Fatalf("Failed to create add liquidity transaction for pair %s - %s: %v", pair[0].Hex(), pair[1].Hex(), err)
		}
		//log.Printf("Add liquidity transaction hash for pair %s - %s: %s", pair[0].Hex(), pair[1].Hex(), addLiquidityTx.Hash().Hex())
		depolyerAuth.Nonce = depolyerAuth.Nonce.Add(depolyerAuth.Nonce, big.NewInt(1))
		lastTxHash = addLiquidityTx.Hash()
		if (i+1)%500 == 0 {
			isConfirmed, err := waitForConfirmation(client, lastTxHash)
			if err != nil {
				return err
			}
			if !isConfirmed {
				return fmt.Errorf("transaction %s was not confirmed", lastTxHash.Hex())
			}
		}

	}

	isConfirmed, err = waitForConfirmation(client, lastTxHash)
	if err != nil {
		log.Fatalf("Failed to confirm add liquidity transaction: %v", err)
	}
	if !isConfirmed {
		log.Fatalf("Add liquidity transaction was not confirmed")
	}

	preparedTestData := TestData{
		TestUsers:       testUsers,
		UniswapV2Router: uniswapV2Contract.uniswapV2Router01Address,
		TokenPairs:      tokenPairs,
	}

	saveTestDataToFile("test/tmp/prepared_test_data.json", preparedTestData)
	return err

}

func executeTestAndCalculateTPS(nodeUrl string, chainID int64, gasLimit uint64, stepCount int, allowFailedTransactionsCount int) error {
	loadedTestData, err := loadTestDataFromFile("test/tmp/prepared_test_data.json")
	if err != nil {
		log.Fatalf("Failed to load test data: %v", err)
		return err
	}
	// randomswap from token A to token A
	steps := generateRandomSwapSteps(loadedTestData.TestUsers, loadedTestData.TokenPairs, stepCount)

	time.Sleep(5 * time.Second)
	//get recent block number
	//lastBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
		return err

	}
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return err

	}
	// get gas price
	// gasPrice, err := client.SuggestGasPrice(context.Background())
	// if err != nil {
	// 	log.Fatalf("Failed to suggest gas price: %v", err)
	// }

	//FixME: should use gasPrice from the chain
	gasPrice := new(big.Int).SetUint64(2000000000)
	resultChan := make(chan *TPSStats)
	errorChan := make(chan error)

	uniswapV2RouterInstance, err := contracts.NewUniswapV2Router01(loadedTestData.UniswapV2Router, client)
	if err != nil {
		log.Fatalf("Failed to create Uniswap V2 Router instance: %v", err)
		return err

	}

	err = executeSwapSteps(client, uniswapV2RouterInstance, steps, chainID, gasPrice, gasLimit)
	if err != nil {
		log.Fatalf("Failed to perform swap steps: %v", err)
		return err

	}
	go calculateTPSByTransactionsCount(client, stepCount-allowFailedTransactionsCount, resultChan, errorChan)

	select {
	case stats := <-resultChan:
		if stats != nil {
			log.Printf("Statistics for the last %d transactions starting from the current block :", stats.TransactionCount)
			log.Printf("Starting block number: %s", stats.StartBlockNumber.String())
			log.Printf("TPS: %.2f", stats.TPS)
			log.Printf("Time interval: %d seconds", stats.TimeInterval)
			log.Printf("Number of blocks: %d", stats.BlockCount)
		}
	case err := <-errorChan:
		if err != nil {
			log.Fatalf("Failed to calculate TPS: %v", err)
		}
	}
	return err

}

func dispatchTestToken(client *ethclient.Client, ownerAuth *bind.TransactOpts, ERC20DeployedContracts []*ERC20DeployedContract, testUsers []*pkg.EthWallet, accountInitialERC20Token *big.Int) error {
	var lastTxHash common.Hash

	for _, contract := range ERC20DeployedContracts {
		for _, user := range testUsers {
			amount := accountInitialERC20Token
			tx, err := contract.tokenInstance.Transfer(ownerAuth, common.HexToAddress(user.Address), amount)
			if err != nil {
				return err
			}
			lastTxHash = tx.Hash()
			ownerAuth.Nonce = ownerAuth.Nonce.Add(ownerAuth.Nonce, big.NewInt(1))

		}
	}

	isConfirmed, err := waitForConfirmation(client, lastTxHash)
	if err != nil {
		return err
	}
	if !isConfirmed {
		return fmt.Errorf("transaction %s was not confirmed", lastTxHash.Hex())
	}

	return nil
}

func generateTokenPairs(contracts []*ERC20DeployedContract) [][2]common.Address {
	var pairs [][2]common.Address
	for i := 0; i < len(contracts); i++ {
		for j := i + 1; j < len(contracts); j++ {
			pair := [2]common.Address{contracts[i].tokenAddress, contracts[j].tokenAddress}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

type SwapStep struct {
	User     *pkg.EthWallet
	TokenIn  common.Address
	TokenOut common.Address
	AmountIn *big.Int
}

func generateRandomSwapSteps(testUsers []*pkg.EthWallet, tokenPairs [][2]common.Address, stepCount int) []SwapStep {
	var steps []SwapStep
	for i := 0; i < stepCount; i++ {
		user := testUsers[rand.Intn(len(testUsers))]

		pair := tokenPairs[rand.Intn(len(tokenPairs))]

		//random swap direction
		swapDirection := rand.Intn(2)
		var tokenIn, tokenOut common.Address
		if swapDirection == 0 {
			tokenIn = pair[0]
			tokenOut = pair[1]
		} else {
			tokenIn = pair[1]
			tokenOut = pair[0]
		}

		amountIn := big.NewInt(rand.Int63n(1e5))

		step := SwapStep{
			User:     user,
			TokenIn:  tokenIn,
			TokenOut: tokenOut,
			AmountIn: amountIn,
		}
		//fmt.Printf("Step %d: User %s swaps %s to %s with amount %s\n", i+1, user.Address, tokenIn.Hex(), tokenOut.Hex(), amountIn.String())

		steps = append(steps, step)
	}
	return steps
}

func executeSwapSteps(client *ethclient.Client, uniswapV2RouterInstance *contracts.UniswapV2Router01, steps []SwapStep, chainID int64, gasPrice *big.Int, gasLimit uint64) error {
	var wg sync.WaitGroup
	results := make(chan error, len(steps))
	for _, step := range steps {
		wg.Add(1)
		go executeSwapStep(client, uniswapV2RouterInstance, step, chainID, gasPrice, gasLimit, &wg, results)
	}

	wg.Wait()
	close(results)

	return nil
}

func executeSwapStep(client *ethclient.Client, uniswapV2RouterInstance *contracts.UniswapV2Router01, step SwapStep, chainID int64, gasPrice *big.Int, gasLimit uint64, wg *sync.WaitGroup, results chan<- error) {
	defer wg.Done()

	auth, err := generateTestAuth(client, step.User, chainID, gasPrice, gasLimit)
	if err != nil {
		results <- fmt.Errorf("failed to generate auth for user %s: %v", step.User.Address, err)
		return
	}

	_, err = uniswapV2RouterInstance.SwapExactTokensForTokens(
		auth,
		step.AmountIn,
		big.NewInt(0),
		[]common.Address{step.TokenIn, step.TokenOut},
		common.HexToAddress(step.User.Address),
		big.NewInt(time.Now().Unix()+1000),
	)
	if err != nil {
		results <- fmt.Errorf("failed to create swap transaction for user %s: %v", step.User.Address, err)
		return
	}

	results <- nil
}

// calculateTPSByTransactionsCount calculates the transactions per second (TPS) for the last `transactionCount` transactions,
// starting from the current block minus `startOffset`.
func calculateTPSByTransactionsCount(client *ethclient.Client, transactionCount int, resultChan chan<- *TPSStats, errorChan chan<- error) {
	defer close(resultChan)
	defer close(errorChan)

	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		errorChan <- fmt.Errorf("failed to get latest block: %v", err)
		return
	}
	//count from the latest block -1
	blockNumber := new(big.Int).SetUint64(latestBlockNumber)

	var blocks []*types.Block
	totalTransactions := 0
	blockCount := 0
	for totalTransactions < transactionCount || len(blocks) < 2 {
		var block *types.Block
		retries := 0
		for {
			block, err = client.BlockByNumber(context.Background(), blockNumber)
			if err == nil {
				break
			}
			if retries >= maxRetries {
				errorChan <- fmt.Errorf("failed to get block %d after %d retries: %v", blockNumber, retries, err)
				return
			}
			log.Printf("Block %d not found, retrying ... (attempt %d/%d)", latestBlockNumber, retries+1, maxRetries)

			time.Sleep(retriesInterval)
			retries++
		}
		log.Printf("block.Time(): %d,block.Number(): %d", block.Time(), block.Number())
		blocks = append(blocks, block)
		totalTransactions += len(block.Transactions())
		log.Printf("totalTransactions: %d", totalTransactions)
		blockNumber.Add(blockNumber, big.NewInt(1))
		log.Printf("latestBlockNumber: %d", blockNumber)
		blockCount++
		if blockCount >= maxBlocks && totalTransactions < transactionCount {
			log.Fatalf("Reached maximum block count of %d with less than %d transactions. Stopping.", maxBlocks, transactionCount)
			break
		}
	}

	timeInterval := blocks[len(blocks)-1].Time() - blocks[0].Time()

	if timeInterval == 0 {
		errorChan <- fmt.Errorf("time interval is zero, cannot calculate TPS")
		return
	}
	tps := float64(totalTransactions) / float64(timeInterval)

	stats := &TPSStats{
		TPS:              tps,
		BlockCount:       len(blocks),
		TimeInterval:     timeInterval,
		TransactionCount: totalTransactions,
		StartBlockNumber: blocks[len(blocks)-1].Number(),
	}
	resultChan <- stats
}
