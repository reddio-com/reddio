package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"itachi/evm"
	"itachi/test/trasnfer/pkg"
)

var (
	hostUrl         = "localhost:9092"
	cfgPath         = "./conf/evm_cfg.toml"
	genWalletCount  = 2
	initialEthCount = uint64(100 * 100)
	testSteps       = 1
	retryCount      = 3
)

func init() {
	flag.StringVar(&hostUrl, "hostUrl", "localhost:9092", "")
	flag.StringVar(&cfgPath, "cfgPath", "./conf/evm_cfg.toml", "")
	flag.IntVar(&genWalletCount, "genWalletCount", 2, "")
	flag.Uint64Var(&initialEthCount, "initialEthCount", uint64(100*100), "")
	flag.IntVar(&testSteps, "testSteps", 1, "")
	flag.IntVar(&retryCount, "retryCount", 3, "")
}

func main() {
	flag.Parse()
	if initialEthCount < 2 {
		panic(fmt.Sprintf("initialEthCount should be larget than 1"))
	}
	cfg := evm.LoadEvmConfig(cfgPath)
	walletsManager := pkg.NewWalletManager(cfg, hostUrl)
	wallets, err := walletsManager.GenerateRandomWallet(genWalletCount, initialEthCount)
	if err != nil {
		panic(fmt.Errorf("generate wallets error:%v", err))
	}
	time.Sleep(5 * time.Second)
	log.Println("create wallets success")
	testStepsManager := pkg.NewTransferManager()
	tc := testStepsManager.GenerateTransferSteps(testSteps, pkg.GenerateCaseWallets(initialEthCount, wallets))
	err = tc.Run(walletsManager)
	if err != nil {
		panic(fmt.Errorf("run testcase failed, err:%v", err))
	}
	success, err := assertWithRetry(tc, walletsManager, wallets)
	if err != nil {
		panic(fmt.Sprintf("assert result err: %v", err))
	}
	if success {
		log.Println("success")
	} else {
		log.Println("failed")
	}
}

func assertWithRetry(tc *pkg.TransferCase, walletsManager *pkg.WalletManager, wallets []*pkg.EthWallet) (bool, error) {
	var got map[string]*pkg.CaseEthWallet
	var success bool
	var err error
	for i := 0; i < retryCount; i++ {
		got, success, err = tc.AssertExpect(walletsManager, wallets)
		if err != nil {
			return false, err
		}
		if success {
			return true, nil
		} else {
			time.Sleep(5 * time.Second)
			continue
		}
	}
	printChange(got, tc.Expect)
	return false, nil
}

func printChange(got, expect map[string]*pkg.CaseEthWallet) {
	for k, v := range got {
		ev, ok := expect[k]
		if ok {
			if v.EthCount != ev.EthCount {
				log.Println(fmt.Sprintf("address:%v got:%v expect:%v", k, v.EthCount, ev.EthCount))
			}
		}
	}
}
