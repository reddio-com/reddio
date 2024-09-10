package uniswap

import (
	"context"
	"fmt"
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
	waitForConfirmationTime = 1 * time.Second
	maxRetries              = 300
)

type ERC20DeployedContract struct {
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

// deploy Erc20 token contracts
func deployERC20Contracts(auth *bind.TransactOpts, client *ethclient.Client, deployNum int) ([]*ERC20DeployedContract, error) {

	var err error
	deployedTokens := make([]*ERC20DeployedContract, 0)

	for i := 0; i < deployNum; i++ {
		deployedToken := &ERC20DeployedContract{}
		deployedToken.tokenAddress, deployedToken.tokenTransaction, deployedToken.tokenInstance, err = contracts.DeployToken(auth, client)
		if err != nil {
			return nil, err
		}

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
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	// Deploy UniswapV2Factory
	deployed.uniswapV2FactoryAddress, deployed.uniswapV2FactoryTransaction, deployed.uniswapV2FactoryInstance, err = contracts.DeployUniswapV2Factory(auth, client, auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	// Deploy UniswapV2Router01
	deployed.uniswapV2Router01Address, deployed.uniswapV2Router01Transaction, deployed.uniswapV2RouterInstance, err = contracts.DeployUniswapV2Router01(auth, client, deployed.uniswapV2FactoryAddress, deployed.weth9Address)
	if err != nil {
		return nil, err
	}
	auth.Nonce.Add(auth.Nonce, big.NewInt(1))

	return deployed, nil
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

func generateTestAuth(client *ethclient.Client, user *pkg.EthWallet, chainID int64, gasPrice *big.Int, gasLimit uint64) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(user.PK)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		return nil, fmt.Errorf("failed to create authorized transactor: %v", err)
	}

	auth.GasPrice = gasPrice
	auth.GasLimit = gasLimit

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(user.Address))
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}
