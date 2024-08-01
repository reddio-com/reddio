package transfer

import (
	"fmt"
	"log"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
)

type EthManager struct {
	evmCfg *evm.GethConfig
	config *conf.EthCaseConf
	wm     *pkg.WalletManager
	//tm     *pkg.TransferManager
	testcases []TestCase
}

func (m *EthManager) Configure(cfg *conf.EthCaseConf, evmCfg *evm.GethConfig) {
	m.config = cfg
	m.evmCfg = evmCfg
	m.wm = pkg.NewWalletManager(m.evmCfg, m.config.HostUrl)
	m.testcases = []TestCase{
		NewRandomTest("[2 account, 1 transfer]", 2, cfg.InitialEthCount, 1),
		NewRandomTest("[20 account, 100 transfer]", 20, cfg.InitialEthCount, 100),
		NewConflictTest("[20 account, 50 transfer]", 20, cfg.InitialEthCount, 50),
	}
}

func (m *EthManager) Run() error {
	for _, tc := range m.testcases {
		log.Println(fmt.Sprintf("start to test %v", tc.Name()))
		if err := tc.Run(m.wm); err != nil {
			return fmt.Errorf("%s failed, err:%v", tc.Name(), err)
		}
		log.Println(fmt.Sprintf("test %v success", tc.Name()))
	}
	return nil
}
