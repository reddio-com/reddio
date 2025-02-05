package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/reddio-com/reddio/cmd/node/app"
	"github.com/reddio-com/reddio/evm"
	"github.com/reddio-com/reddio/test/conf"
	"github.com/reddio-com/reddio/test/testx"
	"github.com/reddio-com/reddio/test/transfer"
)

var (
	evmConfigPath string
	yuConfigPath  string
	poaConfigPath string
	isParallel    bool
	useSql        bool
)

func init() {
	flag.StringVar(&evmConfigPath, "evmConfigPath", "./conf/evm.toml", "")
	flag.StringVar(&yuConfigPath, "yuConfigPath", "./conf/yu.toml", "")
	flag.StringVar(&poaConfigPath, "poaConfigPath", "./conf/poa.toml", "")
	flag.BoolVar(&isParallel, "parallel", true, "")
	flag.BoolVar(&useSql, "use-sql", false, "")
}

func main() {
	flag.Parse()
	yuCfg, poaCfg, evmConfig, config := testx.GenerateConfig(yuConfigPath, evmConfigPath, poaConfigPath, useSql, isParallel)
	yuCfg.SqliteDBConf.Path = "sqlite.db"
	yuCfg.TxnConf.EnableSqliteStorage = true
	go func() {
		if config.IsParallel {
			logrus.Info("start transfer test in parallel")
		} else {
			logrus.Info("start transfer test in serial")
		}
		app.StartByConfig(yuCfg, poaCfg, evmConfig)
	}()
	time.Sleep(5 * time.Second)
	logrus.Info("finish start reddio")
	if err := assertEthTransfer(context.Background(), evmConfig); err != nil {
		logrus.Info(err)
		os.Exit(1)
	}
	logrus.Info("assert success")
	os.Exit(0)
}

func assertEthTransfer(ctx context.Context, evmCfg *evm.GethConfig) error {
	logrus.Info("start asserting transfer eth")
	ethManager := &transfer.EthManager{}
	cfg := conf.Config.EthCaseConf
	ethManager.Configure(cfg, evmCfg)
	ethManager.AddTestCase(
		transfer.NewRandomTest("[rand_test 2 account, 1 transfer]", 2, cfg.InitialEthCount, 1),
		transfer.NewRandomTest("[rand_test 20 account, 100 transfer]", 20, cfg.InitialEthCount, 100),
		transfer.NewConflictTest("[conflict_test 20 account, 50 transfer]", 20, cfg.InitialEthCount, 50),
	)
	return ethManager.Run(ctx)
}
