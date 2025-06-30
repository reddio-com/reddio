package main

import (
	"context"
	"flag"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/uniswap"
)

var (
	dataPath      string
	evmConfigPath string
	qps           int
	action        string
	duration      time.Duration
	deployUsers   int
	testUsers     int
	nonConflict   bool
	maxUsers      int
)

func init() {
	flag.StringVar(&dataPath, "data-path", "./bin/prepared_test_data.json", "Path to uniswap data")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.IntVar(&qps, "qps", 5, "")
	flag.StringVar(&action, "action", "prepare", "")
	flag.DurationVar(&duration, "duration", time.Minute*3, "")
	flag.IntVar(&deployUsers, "deployUsers", 1, "")
	flag.IntVar(&testUsers, "testUsers", 2, "")
	flag.BoolVar(&nonConflict, "nonConflict", false, "")
	flag.IntVar(&maxUsers, "maxUsers", 0, "")
}

func main() {
	flag.Parse()
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	ethManager := &uniswap.EthManager{}
	cfg := conf.Config.EthCaseConf
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	ethManager.Configure(cfg, evmConfig)
	ethManager.AddTestCase(
		uniswap.NewUniswapV2TPSStatisticsTestCase("UniswapV2 TPS StatisticsTestCase", deployUsers, testUsers, maxUsers, limiter, action == "run", nonConflict, dataPath))
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
