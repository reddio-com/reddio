package transfer

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
)

type EthManager struct {
	evmCfg *evm.GethConfig
	config *conf.EthCaseConf
	wm     *pkg.WalletManager
	tm     *pkg.TransferManager
}

func (m *EthManager) Configure(cfg *conf.EthCaseConf, evmCfg *evm.GethConfig) {
	m.config = cfg
	m.evmCfg = evmCfg
	m.wm = pkg.NewWalletManager(m.evmCfg, m.config.HostUrl)
	m.tm = pkg.NewTransferManager()
}

func (m *EthManager) Run() error {
	wallets, err := m.wm.GenerateRandomWallet(m.config.GenWalletCount, m.config.InitialEthCount)
	if err != nil {
		return err
	}
	log.Println("create wallets finish")
	tc := m.tm.GenerateTransferSteps(m.config.TestSteps, pkg.GenerateCaseWallets(m.config.InitialEthCount, wallets))
	err = tc.Run(m.wm)
	if err != nil {
		return err
	}
	success, err := m.assert(tc, m.wm, wallets)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("transfer manager assert failed")
	}
	return nil
}

func (m *EthManager) assert(tc *pkg.TransferCase, walletsManager *pkg.WalletManager, wallets []*pkg.EthWallet) (bool, error) {
	var got map[string]*pkg.CaseEthWallet
	var success bool
	var err error
	for i := 0; i < m.config.RetryCount; i++ {
		_, success, err = tc.AssertExpect(walletsManager, wallets)
		if err != nil {
			return false, err
		}
		if success {
			return true, nil
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}
	printChange(got, tc.Expect)
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
