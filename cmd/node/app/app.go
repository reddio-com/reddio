package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	yuConfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"

	watcher "github.com/reddio-com/reddio/bridge/controller"
	"github.com/reddio-com/reddio/bridge/controller/api"
	"github.com/reddio-com/reddio/bridge/controller/route"
	"github.com/reddio-com/reddio/bridge/utils/database"
	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	"github.com/reddio-com/reddio/parallel"

	"github.com/reddio-com/reddio/bridge/relayer"
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

	StartupL1Watcher(chain, evmCfg)

	chain.Startup()
}

func InitReddio(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	parallelTri := parallel.NewParallelEVM()
	watcherTri := watcher.NewL2EventsWatcher(evmCfg)

	chain := startup.InitDefaultKernel(yuCfg).WithTripods(poaTri, solidityTri, parallelTri, watcherTri)
	// chain.WithExecuteFn(chain.OrderedExecute)
	chain.WithExecuteFn(parallelTri.Execute)
	return chain
}

func startPromServer() {
	// Export Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func StartupL1Watcher(chain *kernel.Kernel, cfg *evm.GethConfig) {
	ctx := context.Background()
	if !cfg.EnableBridge {
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
	l1ToL2Relayer, err := relayer.NewL1ToL2Relayer(ctx, cfg, l1Client, l2Client, chain)
	if err != nil {
		logrus.Fatal("init bridge relayer failed: ", err)
	}

	if cfg.EnableBridge {
		l1Watcher, err := watcher.NewL1EventsWatcher(ctx, cfg, l1Client, l1ToL2Relayer)
		if err != nil {
			logrus.Fatal("init L1 client failed: ", err)
		}
		err = l1Watcher.Run(ctx)
		if err != nil {
			logrus.Fatal("l1 client run failed: ", err)
		}
		//fmt.Println("CFG.BridgeDBConfig: ", cfg.BridgeDBConfig)
		db, err := database.InitDB(cfg.BridgeDBConfig)
		if err != nil {
			log.Fatal("failed to init db", "err", err)
		}

		// defer func() {
		// 	fmt.Println("closing rpc db")
		// 	if deferErr := database.CloseDB(db); deferErr != nil {
		// 		log.Fatal("failed to close db", "err", err)
		// 	}
		// }()
		api.InitController(db)

		router := gin.Default()
		route.Route(router)

		go func() {
			port := cfg.BridgePort
			log.Println("Starting Bridge API server on", port)
			if runServerErr := router.Run(fmt.Sprintf(":%s", port)); runServerErr != nil {
				log.Fatal("run http server failure", "error", runServerErr)
			}
		}()
	}

}
