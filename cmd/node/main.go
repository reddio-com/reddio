package main

import (
	"github.com/reddio-com/reddio/cmd/node/app"
)

func main() {
	app.Start("./conf/evm.toml", "./conf/yu.toml", "./conf/poa.toml", "./conf/config.toml")
}
