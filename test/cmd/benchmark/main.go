package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	configPath    string
	evmConfigPath string
	qps           int
	duration      time.Duration
	action        string
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.IntVar(&qps, "qps", 10000, "")
	flag.DurationVar(&duration, "duration", 5*time.Minute, "")
	flag.StringVar(&action, "action", "run", "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	switch action {
	case "prepare":
		prepareBenchmark(evmConfig)
	case "run":
		blockBenchmark(evmConfig, qps)
	}
}

func prepareBenchmark(evmCfg *evm.GethConfig) error {
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	wallets, err := ethManager.PreCreateWallets(100, cfg.InitialEthCount)
	if err != nil {
		return err
	}
	_, err = os.Stat("eth_benchmark_data.json")
	if err == nil {
		os.Remove("eth_benchmark_data.json")
	}
	file, err := os.Create("eth_benchmark_data.json")
	if err != nil {
		return err
	}
	defer file.Close()
	d, err := json.Marshal(wallets)
	if err != nil {
		return err
	}
	_, err = file.Write(d)
	return err
}

func loadWallets() ([]*pkg.EthWallet, error) {
	d, err := os.ReadFile("eth_benchmark_data.json")
	if err != nil {
		return nil, err
	}
	exp := make([]*pkg.EthWallet, 0)
	if err := json.Unmarshal(d, &exp); err != nil {
		return nil, err
	}
	return exp, nil
}

func blockBenchmark(evmCfg *evm.GethConfig, qps int) error {
	wallets, err := loadWallets()
	if err != nil {
		return err
	}
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	ethManager.AddTestCase(transfer.NewRandomBenchmarkTest("[rand_test 100 account, 1000 transfer]", 100, cfg.InitialEthCount, 50, wallets, limiter))
	runBenchmark(ethManager)
	return nil
}

func runBenchmark(manager *transfer.EthManager) {
	after := time.After(duration)
	for {
		select {
		case <-after:
			return
		default:
		}
		manager.Run(context.Background())
	}
}
