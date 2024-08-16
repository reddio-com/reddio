package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	configPath    string
	evmConfigPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm_cfg.toml", "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	yuCfg := startup.InitDefaultKernelConfig()
	go func() {
		app.Start(evmConfigPath, yuCfg)
	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")
	if err := assertEthTransfer(evmConfig); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("assert success")
	os.Exit(0)
}

func assertEthTransfer(evmCfg *evm.GethConfig) error {
	log.Println("start asserting transfer eth")
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		transfer.NewRandomTest("[rand_test 2 account, 1 transfer]", 2, cfg.InitialEthCount, 1, true),
		transfer.NewRandomTest("[rand_test 20 account, 100 transfer]", 20, cfg.InitialEthCount, 100, true),
		transfer.NewConflictTest("[conflict_test 20 account, 50 transfer]", 20, cfg.InitialEthCount, 50),
	)
	return ethManager.Run()
}
