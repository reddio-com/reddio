package main

import (
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/evm"
)

func main() {
	cfg := startup.InitDefaultKernelConfig()
	gethCfg := evm.LoadEvmConfig("./conf/evm_cfg.toml")
	app.Start(cfg, gethCfg)
}
