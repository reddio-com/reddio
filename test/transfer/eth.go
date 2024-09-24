package transfer

import (
	"context"
	"fmt"
	"github.com/yu-org/yu/config"
	"log"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
)

type EthManager struct {
	evmCfg *evm.GethConfig
	yuCfg  *config.KernelConf
	config *conf.EthCaseConf
	wm     *pkg.WalletManager
	//tm     *pkg.TransferManager
	testcases []TestCase
}

func (m *EthManager) Configure(cfg *conf.EthCaseConf, evmCfg *evm.GethConfig, yuCfg *config.KernelConf, e2eMode bool) {
	m.config = cfg
	m.evmCfg = evmCfg
	m.yuCfg = yuCfg
	m.wm = pkg.NewWalletManager(m.evmCfg, m.yuCfg, m.config.HostUrl, e2eMode)
	m.testcases = []TestCase{}
}

func (m *EthManager) PreCreateWallets(walletCount int, initCount uint64) ([]*pkg.EthWallet, error) {
	wallets, err := m.wm.BatchGenerateRandomWallets(walletCount, initCount)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (m *EthManager) AddTestCase(tc ...TestCase) {
	m.testcases = append(m.testcases, tc...)
}

func (m *EthManager) Run(ctx context.Context) error {
	for _, tc := range m.testcases {
		log.Println(fmt.Sprintf("start to test %v", tc.Name()))
		if err := tc.Run(ctx, m.wm); err != nil {
			return fmt.Errorf("%s failed, err:%v", tc.Name(), err)
		}
		log.Println(fmt.Sprintf("test %v success", tc.Name()))
	}
	return nil
}

func (m *EthManager) BatchRun(ctx context.Context) error {
	for _, tc := range m.testcases {
		if err := tc.BatchRun(ctx, m.wm); err != nil {
			return fmt.Errorf("%s failed, err:%v", tc.Name(), err)
		}
		log.Println(fmt.Sprintf("test %v success", tc.Name()))
	}
	return nil
}
