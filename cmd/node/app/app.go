package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	yuConfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	"github.com/reddio-com/reddio/parallel"
	"github.com/reddio-com/reddio/relayer"
	watcher "github.com/reddio-com/reddio/watcher/controller"
)

func Start(evmPath, yuPath, poaPath, configPath string) {
	yuCfg := startup.InitKernelConfigFromPath(yuPath)
	poaCfg := poa.LoadCfgFromPath(poaPath)
	evmCfg := evm.LoadEvmConfig(evmPath)
	err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	go startPromServer()
	StartUpChain(yuCfg, poaCfg, evmCfg)
}

func StartUpChain(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) {
	figure.NewColorFigure("Reddio", "big", "green", false).Print()

	chain := InitReddio(yuCfg, poaCfg, evmCfg)

	ethrpc.StartupEthRPC(chain, evmCfg)

	StartupEventsWatcher(chain, evmCfg)

	chain.Startup()

}

func InitReddio(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	parallelTri := parallel.NewParallelEVM()

	chain := startup.InitDefaultKernel(
		yuCfg, poaTri, solidityTri, parallelTri,
	)
	//chain.WithExecuteFn(chain.OrderedExecute)
	chain.WithExecuteFn(parallelTri.Execute)
	return chain
}

func startPromServer() {
	// Export Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func StartupEventsWatcher(chain *kernel.Kernel, cfg *evm.GethConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if !cfg.EnableL1Client || !cfg.EnableL2Client {
		logrus.Info("no client enabled, stop init watcher")
		return
	}

	l1Client, err := ethclient.Dial(cfg.L1ClientAddress)
	if err != nil {
		log.Fatal("failed to connect to L1 geth", "endpoint", cfg.L1ClientAddress, "err", err)
	}

	l2Client, err := ethclient.Dial(cfg.L2ClientAddress)
	if err != nil {
		log.Fatal("failed to connect to L2 geth", "endpoint", cfg.L2ClientAddress, "err", err)
	}
	// set up the bridge relayer
	bridgeRelayer, err := relayer.NewBridgeRelayer(ctx, cfg, l1Client, l2Client, chain)
	if err != nil {
		logrus.Fatal("init bridge relayer failed: ", err)
	}

	// set up the L1 and L2 event watchers
	if cfg.EnableL1Client {
		l1Watcher, err := watcher.NewL1EventsWatcher(ctx, cfg, l1Client, bridgeRelayer)
		if err != nil {
			logrus.Fatal("init L1 client failed: ", err)
		}
		err = l1Watcher.Run(cfg, context.Background())
		if err != nil {
			logrus.Fatal("l1 client run failed: ", err)
		}
	}
	if cfg.EnableL2Client {
		l2Watcher, err := watcher.NewL2EventsWatcher(ctx, cfg, l2Client, bridgeRelayer)
		if err != nil {
			logrus.Fatal("init L1 client failed: ", err)
		}
		err = l2Watcher.Run(cfg, context.Background())
		if err != nil {
			logrus.Fatal("l1 client run failed: ", err)
		}
	}

}
