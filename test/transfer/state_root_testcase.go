package transfer

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/reddio-com/reddio/test/pkg"
)

var (
	resultJson = "stateRootTestResult.json"
)

type StateRootTestCase struct {
	*RandomTransferTestCase
}

func NewStateRootTestCase(name string, count int, initial uint64, steps int) *StateRootTestCase {
	return &StateRootTestCase{
		RandomTransferTestCase: NewRandomTest(name, count, initial, steps),
	}
}

func (st *StateRootTestCase) Name() string {
	return "StateRootTestCase"
}

func (st *StateRootTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	if err := st.RandomTransferTestCase.Run(ctx, m); err != nil {
		return err
	}
	result := StateRootTestResult{
		Wallets:      st.wallets,
		TransferCase: st.transCase,
		StateRoot:    getStateRoot(),
	}
	content, _ := json.Marshal(result)
	os.Remove("stateRootTestResult.json")
	file, err := os.Create(resultJson)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

type StateRootTestResult struct {
	Wallets      []*pkg.CaseEthWallet
	TransferCase *pkg.TransferCase
	StateRoot    common.Hash
}

func getStateRoot() common.Hash {
	return [32]byte{}
}

type StateRootAssertTestCase struct {
}

func (s *StateRootAssertTestCase) Run(ctx context.Context, m *pkg.WalletManager) error {
	content, err := os.ReadFile(resultJson)
	if err != nil {
		return err
	}
	result := &StateRootTestResult{}
	if err := json.Unmarshal(content, result); err != nil {
		return err
	}
	if err := runAndAssert(result.TransferCase, m, getWallets(result.Wallets)); err != nil {
		return err
	}
	stateRoot := getStateRoot()
	return result.StateRoot == stateRoot
}

func (s *StateRootAssertTestCase) Name() string {
	return "StateRootAssertTestCase"
}

func getWallets(ws []*pkg.CaseEthWallet) []*pkg.EthWallet {
	got := make([]*pkg.EthWallet, 0)
	for _, w := range ws {
		got = append(got, w.EthWallet)
	}
	return got
}
