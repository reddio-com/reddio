package uniswap

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/test/contracts"
	"github.com/reddio-com/reddio/test/pkg"
)

const (
	accountInitialFunds      = 1e18
	gasLimit                 = 6e7
	accountInitialERC20Token = 1e18
	approveAmount            = 1e18
	amountADesired           = 1e15
	amountBDesired           = 1e15
	stepCount                = 50000
	tokenContractNum         = 2
	maxErrorThreshold        = 10 // Maximum number of errors before stopping the test
)

type UniswapV2TPSStatisticsTestCase struct {
	nodeUrl       string
	ChainID       int64
	MaxUsers      int
	NonConflict   bool
	TestUsers     int
	DeployedUsers int
	rm            *rate.Limiter
	CaseName      string
	loadTestData  TestData
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

func NewUniswapV2TPSStatisticsTestCase(name, nodeURL string, t, d, maxUser int, rm *rate.Limiter, needLoad, nonConflict bool, chainID int64) *UniswapV2TPSStatisticsTestCase {
	tc := &UniswapV2TPSStatisticsTestCase{
		nodeUrl:       nodeURL,
		MaxUsers:      maxUser,
		NonConflict:   nonConflict,
		DeployedUsers: t,
		TestUsers:     d,
		CaseName:      name,
		rm:            rm,
		ChainID:       chainID,
	}
	if needLoad {
		loadedTestData, err := loadTestDataFromFile("test/tmp/prepared_test_data.json")
		if err != nil {
			logrus.Fatalf("Failed to load test data: %v", err)
			return nil
		}
		tc.loadTestData = loadedTestData
	}
	return tc
}

// Run TestUniswapTPS is a test case to measure the transactions per second (TPS) of Uniswap.
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
	err := cd.executeTest(cd.nodeUrl, cd.ChainID, gasLimit, stepCount)
	if err != nil {
		logrus.Fatalf("Failed to execute test and calculate TPS: %v", err)
	}
	return err
}

func (cd *UniswapV2TPSStatisticsTestCase) prepareDeployerContract(deployerUser *pkg.EthWallet, testUsers []*pkg.EthWallet, gasPrice *big.Int, client *ethclient.Client) (UniswapV2Router common.Address, TokenPairs [][2]common.Address, err error) {
	// set tx auth
	privateKey, err := crypto.HexToECDSA(deployerUser.PK)
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	depolyerAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(cd.ChainID))
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("failed to create authorized transactor: %v", err)
	}
	depolyerAuth.GasPrice = gasPrice
	depolyerAuth.GasLimit = uint64(6e7)
	depolyerNonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(deployerUser.Address))
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	depolyerAuth.Nonce = big.NewInt(int64(depolyerNonce))
	// deploy contracts
	uniswapV2Contract, err := deployUniswapV2Contracts(depolyerAuth, client)
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("Failed to deploy contract: %v", err)
	}
	ERC20DeployedContracts, err := deployERC20Contracts(depolyerAuth, client, tokenContractNum)
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("Failed to deploy ERC20 contracts: %v", err)
	}
	lastIndex := len(ERC20DeployedContracts) - 1
	isConfirmed, err := waitForConfirmation(client, ERC20DeployedContracts[lastIndex].tokenTransaction.Hash())
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("Failed to confirm approve transaction: %v", err)
	}
	if !isConfirmed {
		return [20]byte{}, nil, fmt.Errorf("transaction was not confirmed")
	}
	logrus.Info("deploy contracts done")
	err = dispatchTestToken(client, depolyerAuth, ERC20DeployedContracts, testUsers, big.NewInt(accountInitialERC20Token))
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("failed to dispatch test tokens: %v", err)
	}
	var lastTxHash common.Hash
	for _, contract := range ERC20DeployedContracts {
		_, err := contract.tokenInstance.Approve(depolyerAuth, uniswapV2Contract.uniswapV2Router01Address, big.NewInt(approveAmount))
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to create approve transaction for user %s: %v", deployerUser.Address, err)
		}

		depolyerAuth.Nonce = depolyerAuth.Nonce.Add(depolyerAuth.Nonce, big.NewInt(1))
		for _, user := range testUsers {
			testAuth, err := generateTestAuth(client, user, cd.ChainID, gasPrice, gasLimit)
			if err != nil {
				return [20]byte{}, nil, fmt.Errorf("failed to generate test auth for user %s: %v", user.Address, err)
			}
			tx, err := contract.tokenInstance.Approve(testAuth, uniswapV2Contract.uniswapV2Router01Address, big.NewInt(approveAmount))
			if err != nil {
				return [20]byte{}, nil, fmt.Errorf("failed to create approve transaction for user %s: %v", user.Address, err)
			}
			lastTxHash = tx.Hash()
			// logrus.Infof("Approve transaction hash for user %s: %s", user.Address, tx.Hash().Hex())
			testAuth.Nonce = testAuth.Nonce.Add(testAuth.Nonce, big.NewInt(1))
		}
	}
	isConfirmed, err = waitForConfirmation(client, lastTxHash)
	if err != nil {
		return [20]byte{}, nil, err
	}
	if !isConfirmed {
		return [20]byte{}, nil, fmt.Errorf("transaction %s was not confirmed", lastTxHash.Hex())
	}
	logrus.Info("dispatchTestToken done")
	tokenPairs := generateTokenPairs(ERC20DeployedContracts)
	// add liquidity
	for _, pair := range tokenPairs {
		addLiquidityTx, err := uniswapV2Contract.uniswapV2RouterInstance.AddLiquidity(
			depolyerAuth,
			pair[0],
			pair[1],
			big.NewInt(amountADesired),
			big.NewInt(amountBDesired),
			big.NewInt(0),
			big.NewInt(0),
			common.HexToAddress(deployerUser.Address),
			big.NewInt(time.Now().Unix()+1000),
		)
		if err != nil {
			return [20]byte{}, nil, fmt.Errorf("failed to create add liquidity transaction for pair %s - %s: %v", pair[0].Hex(), pair[1].Hex(), err)
		}
		depolyerAuth.Nonce = depolyerAuth.Nonce.Add(depolyerAuth.Nonce, big.NewInt(1))
		lastTxHash = addLiquidityTx.Hash()
	}
	isConfirmed, err = waitForConfirmation(client, lastTxHash)
	if err != nil {
		return [20]byte{}, nil, fmt.Errorf("failed to confirm add liquidity transaction: %v", err)
	}
	if !isConfirmed {
		return [20]byte{}, nil, errors.New("add liquidity transaction was not confirmed")
	}
	logrus.Info("generateTokenPairs done")
	return uniswapV2Contract.uniswapV2Router01Address, tokenPairs, nil
}

func (cd *UniswapV2TPSStatisticsTestCase) Prepare(ctx context.Context, m *pkg.WalletManager) error {
	deployerUsers, err := m.GenerateRandomWallets(cd.DeployedUsers, accountInitialFunds)
	if err != nil {
		return fmt.Errorf("failed to generate deployer user: %v", err.Error())
	}
	testUsers, err := m.GenerateRandomWallets(cd.TestUsers, accountInitialFunds)
	if err != nil {
		return fmt.Errorf("failed to generate test users: %v", err)
	}
	client, err := ethclient.Dial(cd.nodeUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()
	fmt.Println("create testUsers done")

	// get gas price
	gasPrice := big.NewInt(1)
	preparedTestData := TestData{
		TestUsers:     testUsers,
		TestContracts: make([]TestContract, 0),
	}
	for index, deployerUser := range deployerUsers {
		logrus.Infof("start to deploy %v contract", index)
		router, tokenPairs, err := cd.prepareDeployerContract(deployerUser, testUsers, gasPrice, client)
		if err != nil {
			return fmt.Errorf("prepare contract failed, err:%v", err)
		}
		preparedTestData.TestContracts = append(preparedTestData.TestContracts, TestContract{router, tokenPairs})
		logrus.Infof("create %v deploy contract done", index+1)
	}
	saveTestDataToFile("test/tmp/prepared_test_data.json", preparedTestData)
	return err
}

func (cd *UniswapV2TPSStatisticsTestCase) executeTest(nodeUrl string, chainID int64, gasLimit uint64, stepCount int) error {
	var steps []SwapStep
	if cd.NonConflict {
		steps = cd.generateNoConflictSwapSteps(cd.loadTestData)
	} else {
		steps = generateRandomSwapSteps(cd.loadTestData, stepCount)
	}
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		logrus.Infof("Failed to connect to the Ethereum client: %v", err)
		return err
	}
	gasPrice := new(big.Int).SetUint64(1)
	err = cd.executeSwapSteps(client, steps, chainID, gasPrice, gasLimit)
	if err != nil {
		logrus.Infof("Failed to perform swap steps: %v", err)
		return err
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
	Router   common.Address
}

func (cd *UniswapV2TPSStatisticsTestCase) generateNoConflictSwapSteps(testData TestData) []SwapStep {
	var steps []SwapStep
	maxUsers := len(testData.TestUsers)
	if cd.MaxUsers < maxUsers && cd.MaxUsers > 0 {
		maxUsers = cd.MaxUsers
	}
	testUsers := testData.TestUsers
	for i := 0; i < maxUsers; i++ {
		user := testUsers[i]
		contract := testData.TestContracts[i]
		pair := contract.TokenPairs[0]
		tokenIn := pair[0]
		tokenOut := pair[1]
		amountIn := big.NewInt(rand.Int63n(1e5))
		step := SwapStep{
			User:     user,
			TokenIn:  tokenIn,
			TokenOut: tokenOut,
			AmountIn: amountIn,
			Router:   contract.UniswapV2Router,
		}
		steps = append(steps, step)
	}
	return steps
}

func generateRandomSwapSteps(testData TestData, stepCount int) []SwapStep {
	var steps []SwapStep
	testUsers := testData.TestUsers
	for i := 0; i < stepCount; i++ {
		user := testUsers[rand.Intn(len(testUsers))]
		contract := testData.TestContracts[rand.Intn(len(testData.TestContracts))]
		pair := contract.TokenPairs[rand.Intn(len(contract.TokenPairs))]
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
			Router:   contract.UniswapV2Router,
		}
		steps = append(steps, step)
	}
	return steps
}

func (cd *UniswapV2TPSStatisticsTestCase) executeSwapSteps(client *ethclient.Client, steps []SwapStep, chainID int64, gasPrice *big.Int, gasLimit uint64) error {
	var errors []error
	for i, step := range steps {
		// Check if we've hit the error threshold
		if len(errors) >= maxErrorThreshold {
			return fmt.Errorf("stopped after %d errors: %v", len(errors), errors)
		}

		if err := cd.rm.Wait(context.Background()); err != nil {
			errors = append(errors, fmt.Errorf("rate limiter error at step %d: %w", i, err))
			continue
		}
		
		if err := executeSwapStep(client, step, chainID, gasPrice, gasLimit); err != nil {
			errors = append(errors, fmt.Errorf("swap step %d failed: %w", i, err))
			logrus.WithFields(logrus.Fields{
				"step":     i,
				"tokenIn":  step.TokenIn.Hex(),
				"tokenOut": step.TokenOut.Hex(),
				"user":     step.User.Address,
				"error":    err,
			}).Error("Failed to execute swap step")
			
			// Check if this is a critical error that should stop the test
			if isCriticalError(err) {
				return fmt.Errorf("critical error at step %d: %w", i, err)
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("multiple swap steps failed (%d errors): %v", len(errors), errors)
	}
	return nil
}

// isCriticalError determines if an error is critical enough to stop the test
func isCriticalError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for specific critical errors
	criticalErrors := []string{
		"insufficient balance",
		"execution reverted",
		"gas required exceeds allowance",
		"nonce too low",
	}
	
	errStr := err.Error()
	for _, critical := range criticalErrors {
		if strings.Contains(errStr, critical) {
			return true
		}
	}
	return false
}

func executeSwapStep(client *ethclient.Client, step SwapStep, chainID int64, gasPrice *big.Int, gasLimit uint64) error {
	auth, err := generateTestAuth(client, step.User, chainID, gasPrice, gasLimit)
	if err != nil {
		return fmt.Errorf("failed to generate auth for user %s: %v", step.User.Address, err)
	}
	uniswapV2RouterInstance, err := contracts.NewUniswapV2Router01(step.Router, client)
	if err != nil {
		return fmt.Errorf("failed to create Uniswap V2 Router instance: %v", err)
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
		return fmt.Errorf("failed to create swap transaction for user %s: %v", step.User.Address, err)
	}
	return nil
}
