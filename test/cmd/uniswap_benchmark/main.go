package main

import (
	"context"
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/core/startup"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/uniswap"
)

var (
	configPath    string
	evmConfigPath string
	maxBlock      int
	qps           int
	action        string
	duration      time.Duration
	deployUsers   int
	testUsers     int
	nonConflict   bool
	maxUsers      int
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.IntVar(&maxBlock, "maxBlock", 500, "")
	flag.IntVar(&qps, "qps", 1500, "")
	flag.StringVar(&action, "action", "run", "")
	flag.DurationVar(&duration, "duration", time.Minute*5, "")
	flag.IntVar(&deployUsers, "deployUsers", 10, "")
	flag.IntVar(&testUsers, "testUsers", 100, "")
	flag.BoolVar(&nonConflict, "nonConflict", false, "")
	flag.IntVar(&maxUsers, "maxUsers", 0, "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	yuCfg := startup.InitDefaultKernelConfig()
	yuCfg.IsAdmin = true
	yuCfg.Txpool.PoolSize = 10000000
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	ethManager := &uniswap.EthManager{}
	cfg := conf.Config.EthCaseConf
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	ethManager.Configure(cfg, evmConfig)
	ethManager.AddTestCase(
		uniswap.NewUniswapV2TPSStatisticsTestCase("UniswapV2 TPS StatisticsTestCase", deployUsers, testUsers, maxUsers, limiter, action == "run", nonConflict, evmConfig.ChainID))
	switch action {
	case "prepare":
		prepareBenchmark(context.Background(), ethManager)
	case "run":
		blockBenchmark(ethManager)
	}
}

func blockBenchmark(ethManager *uniswap.EthManager) {
	after := time.After(duration)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runBenchmark(ctx, ethManager)
	for {
		select {
		case <-after:
			return
		}
	}
}

func prepareBenchmark(ctx context.Context, manager *uniswap.EthManager) {
	err := manager.Prepare(ctx)
	if err != nil {
		logrus.Infof("err:%v", err.Error())
	}
}

func runBenchmark(ctx context.Context, manager *uniswap.EthManager) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		manager.Run(ctx)
	}
}
