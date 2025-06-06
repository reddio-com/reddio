package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/utils/s3"
)

var (
	evmConfigPath    string
	yuConfigPath     string
	PoaConfigPath    string
	ReddioConfigPath string

	loadConfigType string
	folder         string
	Bucket         string
)

func init() {
	flag.StringVar(&evmConfigPath, "evm-config", "./conf/evm.toml", "path to evm-config file")
	flag.StringVar(&yuConfigPath, "yu-config", "./conf/yu.toml", "path to yu-config file")
	flag.StringVar(&PoaConfigPath, "poa-config", "./conf/poa.toml", "path to poa-config file")
	flag.StringVar(&ReddioConfigPath, "reddio-config", "./conf/config.toml", "path to reddio-config file")
	flag.StringVar(&loadConfigType, "load-config-type", "file", "load-config-type json")
	flag.StringVar(&folder, "folder", "", "path to bucket folder")
	flag.StringVar(&Bucket, "bucket", "", "s3 bucket name")
}

func main() {
	flag.Parse()
	switch loadConfigType {
	case "s3":
		logrus.Info("load config from s3")
		if len(Bucket) < 1 {
			panic(fmt.Errorf("s3 bucket name is required"))
		}
		s3Config, err := s3.InitS3Config(folder, Bucket)
		if err != nil {
			panic(err)
		}
		app.StartByCfgData(s3Config.GetConfig())
	case "file":
		logrus.Info("load config from file")
		app.Start(evmConfigPath, yuConfigPath, PoaConfigPath, ReddioConfigPath)
	default:
		logrus.Info("load config from file")
		app.Start(evmConfigPath, yuConfigPath, PoaConfigPath, ReddioConfigPath)

	}
}
