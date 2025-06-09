package testx

import (
	"github.com/yu-org/yu/apps/poa"
	yuConfig "github.com/yu-org/yu/config"
	"github.com/yu-org/yu/core/startup"

	config2 "github.com/reddio-com/reddio/config"
	"github.com/reddio-com/reddio/evm"
)

func GenerateConfig(yuConfigPath, evmConfigPath, poaConfigPath string, isParallel bool) (yuCfg *yuConfig.KernelConf, poaCfg *poa.PoaConfig, evmConfig *evm.GethConfig, config *config2.Config) {
	yuCfg = startup.InitKernelConfigFromPath(yuConfigPath)
	evmConfig = evm.LoadEvmConfig(evmConfigPath)
	config = config2.GetGlobalConfig()
	config.IsBenchmarkMode = true
	config.IsParallel = isParallel
	config.AsyncCommit = false
	config.RateLimitConfig.GetReceipt = 0
	poaCfg = poa.LoadCfgFromPath(poaConfigPath)
	return yuCfg, poaCfg, evmConfig, config
}
