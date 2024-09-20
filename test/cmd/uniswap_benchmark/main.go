package main

import (
	"context"
	"flag"
	"time"

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
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.IntVar(&maxBlock, "maxBlock", 500, "")
	flag.IntVar(&qps, "qps", 1500, "")
	flag.StringVar(&action, "action", "prepare", "")
	flag.DurationVar(&duration, "duration", time.Minute*5, "")
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
		uniswap.NewUniswapV2TPSStatisticsTestCase("UniswapV2 TPS StatisticsTestCase", limiter))
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
	manager.Prepare(ctx)
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
