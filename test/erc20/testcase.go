package erc20

import (
	"context"
	"errors"
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
	nodeUrl                 = "http://localhost:9092"
	accountInitialFunds     = 1e18
	chainID                 = 50341
	waitForConfirmationTime = 1 * time.Second
	maxRetries              = 300
)

type TestContract struct {
	UniswapV2Router common.Address      `json:"uniswapV2Router"`
	TokenPairs      [][2]common.Address `json:"tokenPairs"`
}

type ERC20DeployedContract struct {
	tokenAddress     common.Address
	tokenTransaction *types.Transaction
	tokenInstance    *contracts.Token
}

type TestCase interface {
	Prepare(ctx context.Context, m *pkg.WalletManager) (common.Address, error)
	Run(ctx context.Context, m *pkg.WalletManager) error
	Name() string
}

type TestData struct {
	TestContracts common.Address
}

type RandomTransferTestCase struct {
	CaseName     string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.Erc20TransferManager

	wallets      []*pkg.CaseEthWallet
	transCase    *pkg.Erc20TransferCase
	erc20Wallets []*pkg.CaseERC20Wallet
}

func NewRandomTest(name string, count int, initial uint64, steps int, contractAddr common.Address) *RandomTransferTestCase {
	return &RandomTransferTestCase{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewErc20TransferManager(contractAddr),
	}
}

func (tc *RandomTransferTestCase) Name() string {
	return tc.CaseName
}

func (tc *RandomTransferTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	var wallets []*pkg.EthWallet
	var err error
	wallets, err = m.GenerateRandomWallets(tc.walletCount, tc.initialCount)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("%s create wallets finish", tc.CaseName))
	tc.wallets = pkg.GenerateCaseWallets(tc.initialCount, wallets)

	contractAddress, err := tc.Prepare(ctx, m)
	fmt.Println(contractAddress)
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

	var caseERC20Wallets []*pkg.CaseERC20Wallet

	for _, ethWallet := range tc.wallets {
		caseERC20Wallets = append(caseERC20Wallets, convertEthWalletToERC20Wallet(ethWallet))
	}

	tc.transCase = tc.tm.GenerateRandomErc20TransferSteps(tc.steps, caseERC20Wallets)
	return runAndAssert(tc.transCase, m, caseERC20Wallets)
}

func (tc *RandomTransferTestCase) Prepare(ctx context.Context, m *pkg.WalletManager) (common.Address, error) {
	deployerUsers, err := m.GenerateRandomWallets(1, accountInitialFunds)
	fmt.Println(deployerUsers) // 临时打印变量，避免未使用报错
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to generate deployer user: %v", err.Error())
	}

	client, err := ethclient.Dial(nodeUrl)

	if err != nil {
		return common.Address{}, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	// get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	contractAddress, err := tc.prepareDeployerContract(deployerUsers[0], gasPrice, client)
	fmt.Println(contractAddress) // 临时打印变量，避免未使用报错
	if err != nil {
		return common.Address{}, fmt.Errorf("prepare contract failed, err:%v", err)
	}

	return contractAddress, nil
}

func (tc *RandomTransferTestCase) prepareDeployerContract(deployerUser *pkg.EthWallet, gasPrice *big.Int, client *ethclient.Client) (contractAddress common.Address, err error) {
	privateKey, err := crypto.HexToECDSA(deployerUser.PK)
	if err != nil {
		return common.Address{}, nil
	}

	depolyerAuth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	depolyerAuth.GasPrice = gasPrice
	depolyerAuth.GasLimit = uint64(6e7)
	depolyerNonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(deployerUser.Address))

	if err != nil {
		return common.Address{}, nil
	}

	depolyerAuth.Nonce = big.NewInt(int64(depolyerNonce))

	contractAddress, err = deployERC20Contracts(depolyerAuth, client)

	return contractAddress, nil
}

func runAndAssert(transferCase *pkg.Erc20TransferCase, m *pkg.WalletManager, wallets []*pkg.CaseERC20Wallet) error {
	if err := transferCase.Run(m); err != nil {
		return err
	}
	log.Println("wait transfer transaction done")
	time.Sleep(5 * time.Second)
	success, err := assert(transferCase, m, wallets)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("transfer manager assert failed")
	}

	bm := pkg.GetDefaultBlockManager()
	block, err := bm.GetCurrentBlock()
	if err != nil {
		return err
	}
	log.Printf("Block(%d) StateRoot: %s", block.Height, block.StateRoot.String())
	return nil
}

func assert(transferCase *pkg.Erc20TransferCase, walletsManager *pkg.WalletManager, wallets []*pkg.ERC20Wallet) (bool, error) {
	var got map[string]*pkg.CaseERC20Wallet
	var success bool
	var err error
	for i := 0; i < 20; i++ {
		got, success, err = transferCase.AssertExpect(walletsManager, wallets)
		if err != nil {
			return false, err
		}
		if success {
			return true, nil
		} else {
			// wait block
			time.Sleep(4 * time.Second)
			continue
		}
	}

	printChange(got, transferCase.Expect, transferCase)
	return false, nil
}

func printChange(got, expect map[string]*pkg.CaseEthWallet, transferCase *pkg.Erc20TransferCase) {
	for _, step := range transferCase.Steps {
		log.Println(fmt.Sprintf("%v transfer %v eth to %v", step.From.Address, step.Count, step.To.Address))
	}
	for k, v := range got {
		ev, ok := expect[k]
		if ok {
			if v.EthCount != ev.EthCount {
				log.Println(fmt.Sprintf("%v got:%v expect:%v", k, v.EthCount, ev.EthCount))
			}
		}
	}
}

// deploy Erc20 token contracts
func deployERC20Contracts(auth *bind.TransactOpts, client *ethclient.Client) (common.Address, error) {
	var err error

	deployedToken := &ERC20DeployedContract{}
	deployedToken.tokenAddress, deployedToken.tokenTransaction, deployedToken.tokenInstance, err = contracts.DeployToken(auth, client)

	if err != nil {
		return common.Address{}, err
	}

	return deployedToken.tokenAddress, nil
}

func waitForConfirmation(client *ethclient.Client, txHash common.Hash) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				return true, nil
			}
			return false, fmt.Errorf("transaction failed with status: %v", receipt.Status)
		}
		time.Sleep(waitForConfirmationTime)
	}
	return false, fmt.Errorf("transaction was not confirmed after %d retries", maxRetries)
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

func convertEthWalletToERC20Wallet(ethWallet *pkg.CaseEthWallet) *pkg.CaseERC20Wallet {
	return &pkg.CaseERC20Wallet{
		ERC20Wallet: &pkg.ERC20Wallet{
			Address: common.HexToAddress(ethWallet.Address),
			Balance: 0,
		},
		TokenCount: 0,
	}
}

func convertERC20WalletToCaseERC20Wallet(erc20Wallet *pkg.ERC20Wallet, tokenCount uint64) *pkg.CaseERC20Wallet {
	return &pkg.CaseERC20Wallet{
		ERC20Wallet: erc20Wallet,
		TokenCount:  tokenCount,
	}
}
