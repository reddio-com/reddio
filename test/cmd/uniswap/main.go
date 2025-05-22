package main

import (
	"context"
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/testx"
	"github.com/reddio-com/reddio/test/transfer"
	"github.com/reddio-com/reddio/test/uniswap"
)

var (
	evmConfigPath string
	yuConfigPath  string
	poaConfigPath string
	isParallel    bool
	useSql        bool
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.BoolVar(&isParallel, "parallel", true, "")
	flag.BoolVar(&useSql, "use-sql", false, "")
}

func main() {
	flag.Parse()
	yuCfg, poaCfg, evmConfig, config := testx.GenerateConfig(yuConfigPath, evmConfigPath, poaConfigPath, useSql, isParallel)
	go func() {
		logrus.Infof("Number of goroutines after app.Start: %d", runtime.NumGoroutine())
		if config.IsParallel {
			logrus.Info("start uniswap test in parallel")
		} else {
			logrus.Info("start uniswap test in serial")
		}
		app.StartByConfig(yuCfg, poaCfg, evmConfig)
	}()
	time.Sleep(5 * time.Second)
	logrus.Info("finish start reddio")
	if err := assertUniswapV2(context.Background(), evmConfig); err != nil {
		logrus.Info(err)
		os.Exit(1)
	}
	logrus.Info("assert success")
	os.Exit(0)
}

func assertUniswapV2(ctx context.Context, evmCfg *evm.GethConfig) error {
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		uniswap.NewUniswapV2AccuracyTestCase("UniswapV2 Accuracy TestCase", 2, cfg.InitialEthCount, evmCfg.ChainID),
	)
	return ethManager.Run(ctx)
}
