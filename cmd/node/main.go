package main

import (
	"github.com/reddio-com/reddio/cmd/node/app"
)

func main() {
	app.Start("./conf/evm_cfg.toml")
}
