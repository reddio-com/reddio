package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/common-nighthawk/go-figure"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/yu-org/yu/apps/poa"
	yuConfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"
	"gorm.io/gorm"

	"github.com/reddio-com/reddio/bridge/checker"
	watcher "github.com/reddio-com/reddio/bridge/controller"
	"github.com/reddio-com/reddio/bridge/controller/api"
	"github.com/reddio-com/reddio/bridge/controller/route"
	"github.com/reddio-com/reddio/bridge/relayer"
	"github.com/reddio-com/reddio/bridge/utils/database"
	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	"github.com/reddio-com/reddio/parallel"
	"github.com/reddio-com/reddio/utils"
)

func StartByConfig(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) {
	utils.IniLimiter()
	go startPromServer()
	StartUpChain(yuCfg, poaCfg, evmCfg)
}

func Start(evmPath, yuPath, poaPath, configPath string) {
	yuCfg := startup.InitKernelConfigFromPath(yuPath)
	poaCfg := poa.LoadCfgFromPath(poaPath)
	evmCfg := evm.LoadEvmConfig(evmPath)
	err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	StartByConfig(yuCfg, poaCfg, evmCfg)
}

func StartUpChain(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) {
	figure.NewColorFigure("Reddio", "big", "green", false).Print()
	logrus.Info("--- Start the Reddio Chain ---")
	var db *gorm.DB
	var err error
	if evmCfg.EnableBridge {
		db, err = database.InitDB(evmCfg.BridgeDBConfig)
		if err != nil {
			logrus.Fatal("failed to init db", "err", err)
		}
	}
	chain := InitReddio(yuCfg, poaCfg, evmCfg, db)

	ethrpc.StartupEthRPC(chain, evmCfg)
	if evmCfg.EnableBridge {
		StartupL1Watcher(evmCfg, db)
		StartupRelayer(chain, evmCfg, db)
		StartupBridgeRpc(evmCfg, db)
	}
	if evmCfg.EnableBridgeChecker {
		StartupChecker(evmCfg, db)
	}
	chain.Startup()
	logrus.Info("start the server")
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sigint:
		logrus.Info("stop the server")
		chain.Stop()
		chain.WaitExit()
	}

}

func InitReddio(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig, db *gorm.DB) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	parallelTri := parallel.NewParallelEVM()
	watcherTri := watcher.NewL2EventsWatcher(evmCfg, db)

	chain := startup.InitDefaultKernel(yuCfg).WithTripods(poaTri, solidityTri, parallelTri, watcherTri)
	// chain.WithExecuteFn(chain.OrderedExecute)
	chain.WithExecuteFn(parallelTri.Execute)
	return chain
}

func startPromServer() {
	// Export Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Fatal("start prometheus server failed", "error", err)
	}
}

func StartupL1Watcher(cfg *evm.GethConfig, db *gorm.DB) {
	ctx := context.Background()

	l1Client, err := ethclient.Dial(cfg.L1ClientAddress)
	if err != nil {
		logrus.Fatal("failed to connect to L1 geth", "endpoint", cfg.L1ClientAddress, "err", err)
	}

	l1Watcher, err := watcher.NewL1EventsWatcher(ctx, cfg, l1Client, db)
	if err != nil {
		logrus.Fatal("init L1 client failed: ", err)
	}
	l1Watcher.Start()

}

func StartupRelayer(chain *kernel.Kernel, cfg *evm.GethConfig, db *gorm.DB) {
	ctx := context.Background()

	l1Client, err := ethclient.Dial(cfg.L1ClientAddress)
	if err != nil {
		logrus.Fatal("failed to connect to L1 geth", "endpoint", cfg.L1ClientAddress, "err", err)
	}

	l1Relayer, err := relayer.NewL1Relayer(ctx, cfg, l1Client, chain, db)
	if err != nil {
		logrus.Fatal("init bridge relayer failed: ", err)
	}

	go l1Relayer.StartPolling()
	l2Relayer, err := relayer.NewL2Relayer(ctx, cfg, db)
	if err != nil {
		logrus.Fatal("init bridge relayer failed: ", err)
	}

	go l2Relayer.StartPolling()
}

func StartupChecker(cfg *evm.GethConfig, db *gorm.DB) {
	ctx := context.Background()
	checker := checker.NewChecker(ctx, cfg, db)
	go checker.StartChecking()

}

func StartupBridgeRpc(cfg *evm.GethConfig, db *gorm.DB) {
	api.InitController(db)

	router := gin.Default()
	route.Route(router)

	go func() {
		port := cfg.BridgePort
		if runServerErr := router.Run(fmt.Sprintf(":%s", port)); runServerErr != nil {
			logrus.Fatal("run http server failure", "error", runServerErr)
		}
	}()
}
