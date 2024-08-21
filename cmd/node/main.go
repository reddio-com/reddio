package main

import (
	"github.com/yu-org/yu/core/startup"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/config"
)

func main() {
	yuCfg := startup.InitDefaultKernelConfig()
	app.Start("./conf/evm_cfg.toml", config.GetGlobalConfig(), yuCfg)
}
