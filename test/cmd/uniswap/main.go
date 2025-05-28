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
	evmConfigPath        string
	yuConfigPath         string
	poaConfigPath        string
	evmProcessorSelector string
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.StringVar(&evmProcessorSelector, "evmProcessorSelector", "serial", "")
}

func main() {
	flag.Parse()
	yuCfg, poaCfg, evmConfig, _ := testx.GenerateConfig(yuConfigPath, evmConfigPath, poaConfigPath, evmProcessorSelector)
	go func() {
		logrus.Infof("Number of goroutines after app.Start: %d", runtime.NumGoroutine())
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
		uniswap.NewUniswapV2AccuracyTestCase("UniswapV2 Accuracy TestCase", 2, cfg.InitialEthCount),
	)
	return ethManager.Run(ctx)
}
