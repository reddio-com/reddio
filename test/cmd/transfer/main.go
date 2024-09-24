package main

import (
	"context"
	"flag"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"log"
	"os"
	"time"

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
	isParallel    bool
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.BoolVar(&isParallel, "parallel", true, "")
}

func main() {
	flag.Parse()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	config := config2.GetGlobalConfig()
	config.IsParallel = isParallel
	go func() {
		if config.IsParallel {
			log.Println("start transfer test in parallel")
		} else {
			log.Println("start transfer test in serial")
		}
		app.Start(evmConfigPath, yuConfigPath, poaConfigPath, "")
	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")

	yuCfg := startup.InitKernelConfigFromPath(yuConfigPath)

	if err := assertEthTransfer(context.Background(), evmConfig, yuCfg); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("assert success")
	os.Exit(0)
}

func assertEthTransfer(ctx context.Context, evmCfg *evm.GethConfig, yuCfg *config.KernelConf) error {
	log.Println("start asserting transfer eth")
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg, yuCfg)
	ethManager.AddTestCase(
		transfer.NewRandomTest("[rand_test 2 account, 1 transfer]", 2, cfg.InitialEthCount, 1),
		transfer.NewRandomTest("[rand_test 20 account, 100 transfer]", 20, cfg.InitialEthCount, 100),
		transfer.NewConflictTest("[conflict_test 20 account, 50 transfer]", 20, cfg.InitialEthCount, 50),
	)
	return ethManager.Run(ctx)
}
