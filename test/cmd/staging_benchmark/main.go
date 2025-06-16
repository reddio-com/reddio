package main

import (
	"context"
	"flag"
	"time"

	"golang.org/x/time/rate"

	"github.com/reddio-com/reddio/test/transfer"
)

var (
	qps               int
	duration          time.Duration
	preCreateWallets  int
	nodeUrl           string
	genesisPrivateKey string
	chainID           int64
)

func init() {
	flag.IntVar(&qps, "qps", 10, "")
	flag.DurationVar(&duration, "duration", 5*time.Minute, "")
	flag.IntVar(&preCreateWallets, "pre-create-wallets", 100, "")
	flag.StringVar(&nodeUrl, "nodeUrl", "http://localhost:9092", "")
	flag.StringVar(&genesisPrivateKey, "key", "32e3b56c9f2763d2332e6e4188e4755815ac96441e899de121969845e343c2ff", "")
	flag.Int64Var(&chainID, "chainId", 50341, "")
}

func main() {
	flag.Parse()
	ethManager := &transfer.EthManager{}
	ethManager.Configure(nil, nodeUrl, genesisPrivateKey, chainID)
	wallets, err := ethManager.PreCreateWallets(preCreateWallets, 5)
	if err != nil {
		panic(err)
	}
	limiter := rate.NewLimiter(rate.Limit(qps), qps)
	ethManager.AddTestCase(transfer.NewStagingBenchmark(wallets, limiter))
	runBenchmark(ethManager)
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
