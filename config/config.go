package config

import "runtime"

type Config struct {
	IsParallel     bool
	MaxConcurrency int
}

func defaultConfig() *Config {
	return &Config{
		IsParallel:     true,
		MaxConcurrency: runtime.NumCPU(),
	}
}

var GlobalConfig *Config

func init() {
	GlobalConfig = defaultConfig()
}

func GetGlobalConfig() *Config {
	return GlobalConfig
}
