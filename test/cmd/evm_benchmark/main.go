package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/reddio-com/reddio/test/pkg"
	"github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"
	"os"
	"time"

	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	configPath    string
	evmConfigPath string
	yuConfigPath  string
	qps           int
	duration      time.Duration
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.IntVar(&qps, "qps", 10000, "")
	flag.DurationVar(&duration, "duration", 30*time.Second, "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	yuConfig := startup.InitKernelConfigFromPath(yuConfigPath)

	prepareAndRun(evmConfig, yuConfig, qps)
}

func prepareAndRun(evmCfg *evm.GethConfig, yuCfg *config.KernelConf, qps int) error {
	// ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	wm := pkg.NewWalletManagerForEVM(evmCfg, yuCfg, cfg.HostUrl)
	//ethManager.Configure(wm)
	wallets, err := wm.GenerateRandomWalletsForEVM(100, cfg.InitialEthCount)
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
	if err != nil {
		return err
	}

	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	testcase := transfer.NewRandomEVMBenchmarkTest("[rand_test 100 account, 1000 transfer]", 100, cfg.InitialEthCount, 50, wallets, limiter)

	runBenchmark(testcase, wm)
	return nil
}

func runBenchmark(testcase *transfer.RandomEVMBenchmarkTest, wm *pkg.WalletManager) {
	after := time.After(duration)
	for {
		select {
		case <-after:
			return
		default:
		}
		testcase.Run(context.Background(), wm)
	}
}
