package app

import (
	"net/http"

	"github.com/common-nighthawk/go-figure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yu-org/yu/apps/poa"
	yuConfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	"github.com/reddio-com/reddio/parallel"
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

	chain.Startup()

}

func InitReddio(yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	parallelTri := parallel.NewParallelEVM()

	chain := startup.InitDefaultKernel(yuCfg).WithTripods(poaTri, solidityTri, parallelTri)
	//chain.WithExecuteFn(chain.OrderedExecute)
	chain.WithExecuteFn(parallelTri.Execute)
	return chain
}

func startPromServer() {
	// Export Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
