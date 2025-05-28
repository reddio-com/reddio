package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/reddio-com/reddio/cmd/node/app"
	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/erc20"
)

var (
	evmConfigPath        string
	yuConfigPath         string
	poaConfigPath        string
	evmProcessorSelector string
	useSql               bool
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.StringVar(&evmProcessorSelector, "evmProcessorSelector", "serial", "")
	flag.BoolVar(&useSql, "use-sql", false, "")
}

func main() {
	flag.Parse()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	config := config2.GetGlobalConfig()
	config.IsBenchmarkMode = true
	config.EvmProcessorSelector = evmProcessorSelector
	config.AsyncCommit = true
	go func() {
		app.Start(evmConfigPath, yuConfigPath, poaConfigPath, "")
	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")
	if err := assertErc20Transfer(context.Background(), evmConfig); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("assert success")
	os.Exit(0)
}

func assertErc20Transfer(ctx context.Context, evmCfg *evm.GethConfig) error {
	log.Println("start asserting transfer eth")
	ethManager := &erc20.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		erc20.NewRandomTest("[rand_test 2 account, 1 transfer]", 2, cfg.InitialEthCount, 1),
	)
	return ethManager.Run(ctx)
}
