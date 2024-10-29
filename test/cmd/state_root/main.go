package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	evmConfigPath string
	yuConfigPath  string
	poaConfigPath string
	action        string
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.StringVar(&action, "action", "gen", "")
}

func main() {
	flag.Parse()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	config := config2.GlobalConfig
	config.IsBenchmarkMode = true
	go func() {
		yuCfg := startup.InitKernelConfigFromPath(yuConfigPath)
		poaCfg := poa.LoadCfgFromPath(poaConfigPath)
		evmCfg := evm.LoadEvmConfig(evmConfigPath)
		poaCfg.BlockInterval = 30 * 1000
		poaCfg.PrettyLog = true
		app.StartUpChain(yuCfg, poaCfg, evmCfg)
	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")
	switch action {
	case "gen":
		if err := assertEthTransfer(context.Background(), evmConfig); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.Println("assert success")
		os.Exit(0)
	}

}

func assertEthTransfer(ctx context.Context, evmCfg *evm.GethConfig) error {
	log.Println("start asserting transfer eth")
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		transfer.NewStateRootTestCase("[state_root_test 50 account, 5000 transfer]", 50, 90*100*100*100, 5000),
	)
	return ethManager.Run(ctx)
}
