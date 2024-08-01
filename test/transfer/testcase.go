package transfer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/reddio-com/reddio/test/pkg"
)

type TestCase interface {
	Run(m *pkg.WalletManager) error
	Name() string
}

type RandomTransferTestCase struct {
	CaseName     string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.TransferManager
}

func NewRandomTest(name string, count int, initial uint64, steps int) *RandomTransferTestCase {
	return &RandomTransferTestCase{
		CaseName:     name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewTransferManager(),
	}
}

func (tc *RandomTransferTestCase) Name() string {
	return tc.CaseName
}

func (tc *RandomTransferTestCase) Run(m *pkg.WalletManager) error {
	wallets, err := m.GenerateRandomWallet(tc.walletCount, tc.initialCount)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("%s create wallets finish", tc.CaseName))
	transferCase := tc.tm.GenerateRandomTransferSteps(tc.steps, pkg.GenerateCaseWallets(tc.initialCount, wallets))
	return runAndAssert(transferCase, m, wallets)
}

func runAndAssert(transferCase *pkg.TransferCase, m *pkg.WalletManager, wallets []*pkg.EthWallet) error {
	if err := transferCase.Run(m); err != nil {
		return err
	}
	success, err := assert(transferCase, m, wallets)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("transfer manager assert failed")
	}
	return nil
}

func assert(transferCase *pkg.TransferCase, walletsManager *pkg.WalletManager, wallets []*pkg.EthWallet) (bool, error) {
	var got map[string]*pkg.CaseEthWallet
	var success bool
	var err error
	for i := 0; i < 3; i++ {
		_, success, err = transferCase.AssertExpect(walletsManager, wallets)
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
	printChange(got, transferCase.Expect)
	return false, nil
}

func printChange(got, expect map[string]*pkg.CaseEthWallet) {
	for k, v := range got {
		ev, ok := expect[k]
		if ok {
			if v.EthCount != ev.EthCount {
				log.Println(fmt.Sprintf("address:%v got:%v expect:%v", k, v.EthCount, ev.EthCount))
			}
		}
	}
}
