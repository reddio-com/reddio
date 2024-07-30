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
	cfg := startup.InitDefaultKernelConfig()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	go func() {
		app.Start(cfg, evmConfig)
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
	ethManager.Configure(conf.Config.EthCaseConf, evmCfg)
	return ethManager.Run()
}
