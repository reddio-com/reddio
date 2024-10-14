package main

import (
	"flag"

	"github.com/reddio-com/reddio/cmd/node/app"
)

var (
	evmConfigPath    string
	yuConfigPath     string
	PoaConfigPath    string
	ReddioConfigPath string
)

func init() {
	flag.StringVar(&evmConfigPath, "evm-config", "./conf/evm.toml", "path to evm-config file")
	flag.StringVar(&yuConfigPath, "yu-config", "./conf/yu.toml", "path to yu-config file")
	flag.StringVar(&PoaConfigPath, "poa-config", "./conf/poa.toml", "path to poa-config file")
	flag.StringVar(&ReddioConfigPath, "reddio-config", "./conf/config.toml", "path to reddio-config file")
}

func main() {
	flag.Parse()
	app.Start(evmConfigPath, yuConfigPath, PoaConfigPath, ReddioConfigPath)
}
