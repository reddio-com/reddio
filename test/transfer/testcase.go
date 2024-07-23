package transfer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/reddio-com/reddio/test/pkg"
)

type TestCase struct {
	Name         string
	walletCount  int
	initialCount uint64
	steps        int
	tm           *pkg.TransferManager
}

func NewTestcase(name string, count int, initial uint64, steps int) *TestCase {
	return &TestCase{
		Name:         name,
		walletCount:  count,
		initialCount: initial,
		steps:        steps,
		tm:           pkg.NewTransferManager(),
	}
}

func (tc *TestCase) Run(m *pkg.WalletManager) error {
	wallets, err := m.GenerateRandomWallet(tc.walletCount, tc.initialCount)
	if err != nil {
		return err
	}
	log.Println("create wallets finish")
	transferCase := tc.tm.GenerateTransferSteps(tc.steps, pkg.GenerateCaseWallets(tc.initialCount, wallets))
	if err := transferCase.Run(m); err != nil {
		return err
	}
	success, err := tc.assert(transferCase, m, wallets)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("transfer manager assert failed")
	}
	return nil
}

func (tc *TestCase) assert(transferCase *pkg.TransferCase, walletsManager *pkg.WalletManager, wallets []*pkg.EthWallet) (bool, error) {
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
