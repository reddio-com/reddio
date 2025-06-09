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
	evmConfigPath     string
	yuConfigPath      string
	poaConfigPath     string
	isParallel        bool
	nodeUrl           string
	genesisPrivateKey string
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.BoolVar(&isParallel, "parallel", true, "")
	flag.StringVar(&nodeUrl, "nodeUrl", "http://localhost:9092", "")
	flag.StringVar(&genesisPrivateKey, "key", "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff", "")

}

func main() {
	flag.Parse()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	config := config2.GetGlobalConfig()
	config.IsBenchmarkMode = true
	config.IsParallel = isParallel
	config.AsyncCommit = true
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
	ethManager.Configure(cfg, nodeUrl, genesisPrivateKey, evmCfg.ChainConfig.ChainID.Int64())
	ethManager.AddTestCase(
		erc20.NewRandomTest("[rand_test 2 account, 1 transfer]", nodeUrl, 2, cfg.InitialEthCount, 1, evmCfg.ChainID),
	)
	return ethManager.Run(ctx)
}
