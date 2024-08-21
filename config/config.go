package config

type Config struct {
	IsParallel bool
}

func defaultConfig() *Config {
	return &Config{
		IsParallel: true,
	}
}

var GlobalConfig *Config

func init() {
	GlobalConfig = defaultConfig()
}

func GetGlobalConfig() *Config {
	return GlobalConfig
}
