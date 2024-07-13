package app

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/yu-org/yu/apps/poa"
	"github.com/yu-org/yu/core/kernel"
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/evm/ethrpc"
	reddioKernel "github.com/reddio-com/reddio/kernel"
)

func StartUpChain(poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) {
	figure.NewColorFigure("Reddio", "big", "green", false).Print()

	chain := InitReddio(poaCfg, evmCfg)

	ethrpc.StartupEthRPC(chain, evmCfg)

	chain.Startup()

}

func InitReddio(poaCfg *poa.PoaConfig, evmCfg *evm.GethConfig) *kernel.Kernel {
	poaTri := poa.NewPoa(poaCfg)
	solidityTri := evm.NewSolidity(evmCfg)
	chain := startup.InitDefaultKernel(
		poaTri, solidityTri,
	)
	//chain.WithExecuteFn(chain.OrderedExecute)
	rk := reddioKernel.NewReddioKernel(chain, solidityTri)
	chain.WithExecuteFn(rk.Execute)
	return chain
}
