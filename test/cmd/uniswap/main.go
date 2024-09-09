package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/transfer"
	"github.com/reddio-com/reddio/test/uniswap"
)

var (
	configPath    string
	evmConfigPath string
	isParallel    bool
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm_cfg.toml", "")
	flag.BoolVar(&isParallel, "parallel", true, "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	yuCfg := startup.InitDefaultKernelConfig()
	yuCfg.Txpool.PoolSize = 10000000
	config := config2.GetGlobalConfig()
	config.IsParallel = isParallel

	go func() {
		log.Printf("Number of goroutines after app.Start: %d", runtime.NumGoroutine())
		if config.IsParallel {
			log.Println("start transfer test in parallel")
		} else {
			log.Println("start transfer test in serial")
		}
		app.Start(evmConfigPath, yuCfg)

	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")
	if err := assertUniswapV2(context.Background(), evmConfig); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("assert success")
	os.Exit(0)

}

func assertUniswapV2(ctx context.Context, evmCfg *evm.GethConfig) error {
	log.Println("start asserting transfer eth")
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		uniswap.NewUniswapV2AccuracyTestCase("UniswapV2 Accuracy TestCase", 2, cfg.InitialEthCount),
		uniswap.NewUniswapV2TPSStatisticsTestCase("UniswapV2 TPS StatisticsTestCase", 2, cfg.InitialEthCount),
	)
	return ethManager.Run(ctx)
}
