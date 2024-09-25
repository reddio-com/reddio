package watcher

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/reddio-com/reddio/watcher/contract"
)

type EthSubscriber struct {
	ethClient *ethclient.Client
	client    *rpc.Client
	filterer  *contract.ParentBridgeCoreFacetFilterer
}

func NewEthSubscriber(l1ClientAddress string, coreContractAddress common.Address) (*EthSubscriber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// TODO replace with our own client once we have one.
	// Geth pulls in a lot of dependencies that we don't use.
	client, err := rpc.DialContext(ctx, l1ClientAddress)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)
	filterer, err := contract.NewParentBridgeCoreFacetFilterer(coreContractAddress, ethClient)
	if err != nil {
		return nil, err
	}
	return &EthSubscriber{
		ethClient: ethClient,
		client:    client,
		filterer:  filterer,
	}, nil
}

func (s *EthSubscriber) WatchDownwardMessage(
	ctx context.Context,
	sink chan<- *contract.ParentBridgeCoreFacetDownwardMessage,
	sequence []*big.Int,
) (event.Subscription, error) {
	return s.filterer.WatchDownwardMessage(&bind.WatchOpts{Context: ctx}, sink, sequence)
}

func (s *EthSubscriber) ChainID(ctx context.Context) (*big.Int, error) {
	return s.ethClient.ChainID(ctx)
}

func (s *EthSubscriber) Close() {
	s.ethClient.Close()
}
