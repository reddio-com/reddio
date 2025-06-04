package main

import (
	"flag"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/utils/s3"
)

var (
	evmConfigPath    string
	yuConfigPath     string
	PoaConfigPath    string
	ReddioConfigPath string
	useS3Config      bool
	s3Bucket         string
)

func init() {
	flag.StringVar(&evmConfigPath, "evm-config", "./conf/evm.toml", "path to evm-config file")
	flag.StringVar(&yuConfigPath, "yu-config", "./conf/yu.toml", "path to yu-config file")
	flag.StringVar(&PoaConfigPath, "poa-config", "./conf/poa.toml", "path to poa-config file")
	flag.StringVar(&ReddioConfigPath, "reddio-config", "./conf/config.toml", "path to reddio-config file")
	flag.BoolVar(&useS3Config, "use-s3-config", false, "use s3 config file")
	flag.StringVar(&s3Bucket, "s3-bucket", "", "s3 bucket name")
}

func main() {
	flag.Parse()
	if useS3Config {
		s3Config, err := s3.InitS3Config(s3Bucket)
		if err != nil {
			panic(err)
		}
		app.StartByCfgData(s3Config.GetConfig())
	} else {
		app.Start(evmConfigPath, yuConfigPath, PoaConfigPath, ReddioConfigPath)
	}
}
