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
	configPath       string
	evmConfigPath    string
	qps              int
	duration         time.Duration
	action           string
	preCreateWallets int
)

const benchmarkDataPath = "./bin/eth_benchmark_data.json"

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.IntVar(&qps, "qps", 10000, "")
	flag.DurationVar(&duration, "duration", 5*time.Minute, "")
	flag.StringVar(&action, "action", "run", "")
	flag.IntVar(&preCreateWallets, "preCreateWallets", 100, "")
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
	wallets, err := ethManager.PreCreateWallets(preCreateWallets, cfg.InitialEthCount)
	if err != nil {
		return err
	}
	_, err = os.Stat(benchmarkDataPath)
	if err == nil {
		os.Remove(benchmarkDataPath)
	}
	file, err := os.Create(benchmarkDataPath)
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
	d, err := os.ReadFile(benchmarkDataPath)
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
	ethManager.AddTestCase(transfer.NewRandomBenchmarkTest("[rand_test 1000 transfer]", cfg.InitialEthCount, wallets, limiter))
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
