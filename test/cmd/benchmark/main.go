package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/pkg"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	configPath    string
	evmConfigPath string
	maxBlock      int
)

func init() {
	flag.StringVar(&configPath, "configPath", "", "")
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm_cfg.toml", "")
	flag.IntVar(&maxBlock, "maxBlock", 10, "")
}

func main() {
	flag.Parse()
	if err := conf.LoadConfig(configPath); err != nil {
		panic(err)
	}
	yuCfg := startup.InitDefaultKernelConfig()
	yuCfg.IsAdmin = true
	evmConfig := evm.LoadEvmConfig(evmConfigPath)
	go func() {
		log.Println("start reddio")
		app.Start(evmConfigPath, yuCfg)
		log.Println("exit reddio")
	}()
	totalCount, err := blockBenchmark(evmConfig, maxBlock)
	if err != nil {
		os.Exit(1)
	}
	log.Println(fmt.Sprintf("totalTxn Count %v", totalCount))
	os.Exit(0)
}

func blockBenchmark(evmCfg *evm.GethConfig, target int) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bm := pkg.GetDefaultBlockManager()
	go runBenchmark(ctx, evmCfg)
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

func runBenchmark(ctx context.Context, evmCfg *evm.GethConfig) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		ethManager := &transfer.EthManager{}
		cfg := conf.Config.EthCaseConf
		ethManager.Configure(cfg, evmCfg)
		ethManager.AddTestCase(transfer.NewRandomTest("[rand_test 100 account, 5000 transfer]", 100, cfg.InitialEthCount, 5000, false))
		ethManager.Run()
		time.Sleep(5 * time.Second)
	}
}
