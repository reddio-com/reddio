package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yu-org/yu/core/startup"
	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/cmd/node/app"
	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	configPath    string
	evmConfigPath string
	maxBlock      int
	qps           int
	isParallel    bool
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm_cfg.toml", "")
	flag.IntVar(&maxBlock, "maxBlock", 50, "")
	flag.IntVar(&qps, "qps", 1500, "")
	flag.BoolVar(&isParallel, "parallel", true, "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	yuCfg := startup.InitDefaultKernelConfig()
	yuCfg.IsAdmin = true
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	config := config2.GetGlobalConfig()
	config.IsParallel = isParallel
	go func() {
		if config.IsParallel {
			log.Println("start reddio in parallel")
		} else {
			log.Println("start reddio in serial")
		}
		app.Start(evmConfigPath, yuCfg)
		log.Println("exit reddio")
	}()
	time.Sleep(3 * time.Second)
	totalCount, err := blockBenchmark(evmConfig, maxBlock, qps)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	log.Println(fmt.Sprintf("totalTxn Count %v", totalCount))
	os.Exit(0)
}

func blockBenchmark(evmCfg *evm.GethConfig, target int, qps int) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bm := pkg.GetDefaultBlockManager()
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	wallets, err := ethManager.PreCreateWallets(100, cfg.InitialEthCount)
	if err != nil {
		return 0, err
	}
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	ethManager.AddTestCase(transfer.NewRandomBenchmarkTest("[rand_test 100 account, 5000 transfer]", 100, cfg.InitialEthCount, 5000, wallets, limiter))
	go runBenchmark(ctx, ethManager)
	totalCount := 0
	for i := 1; i <= target; {
		finish, txnCount, err := bm.GetBlockTxnCountByIndex(i)
		if err != nil {
			return 0, err
		}
		if finish {
			i++
			totalCount += txnCount
			continue
		}
		time.Sleep(3 * time.Second)
	}
	bm.StopBlockChain()
	return totalCount, nil
}

func runBenchmark(ctx context.Context, manager *transfer.EthManager) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		manager.Run(ctx)
	}
}
