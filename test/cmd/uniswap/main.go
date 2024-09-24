package main

import (
	"context"
	"flag"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/reddio-com/reddio/cmd/node/app"
	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/transfer"
	"github.com/reddio-com/reddio/test/uniswap"
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
		log.Printf("Number of goroutines after app.Start: %d", runtime.NumGoroutine())
		if config.IsParallel {
			log.Println("start uniswap test in parallel")
		} else {
			log.Println("start uniswap test in serial")
		}
		app.Start(evmConfigPath, yuConfigPath, poaConfigPath, "")
	}()
	time.Sleep(5 * time.Second)
	log.Println("finish start reddio")
	yuCfg := startup.InitKernelConfigFromPath(yuConfigPath)
	if err := assertUniswapV2(context.Background(), evmConfig, yuCfg); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println("assert success")
	os.Exit(0)

}

func assertUniswapV2(ctx context.Context, evmCfg *evm.GethConfig, yuCfg *config.KernelConf) error {
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg, yuCfg)
	ethManager.AddTestCase(
		uniswap.NewUniswapV2AccuracyTestCase("UniswapV2 Accuracy TestCase", 2, cfg.InitialEthCount),
	)
	return ethManager.Run(ctx)
}
