package uniswap

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
)

type EthManager struct {
	evmCfg *evm.GethConfig
	config *conf.EthCaseConf
	wm     *pkg.WalletManager
	// tm     *pkg.TransferManager
	testcases []TestCase
}

func (m *EthManager) Configure(cfg *conf.EthCaseConf, evmCfg *evm.GethConfig) {
	m.config = cfg
	m.evmCfg = evmCfg
	m.wm = pkg.NewWalletManager(m.evmCfg, m.config.HostUrl)
	m.testcases = []TestCase{}
}

func (m *EthManager) PreCreateWallets(walletCount int, initCount uint64) ([]*pkg.EthWallet, error) {
	wallets, err := m.wm.GenerateRandomWallets(walletCount, initCount)
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
		logrus.Infof("start to test %v", tc.Name())
		if err := tc.Run(ctx, m.wm); err != nil {
			return fmt.Errorf("%s failed, err:%v", tc.Name(), err)
		}
		logrus.Infof("test %v success", tc.Name())
	}
	return nil
}

func (m *EthManager) Prepare(ctx context.Context) error {
	for _, tc := range m.testcases {
		logrus.Infof("start to prepare %v", tc.Name())
		if err := tc.Prepare(ctx, m.wm); err != nil {
			return fmt.Errorf("%s failed, err:%v", tc.Name(), err)
		}
		logrus.Infof("prepare %v success", tc.Name())
	}
	return nil
}
