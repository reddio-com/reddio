package config

import (
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	EvmProcessorSelector string `yaml:"evmProcessorSelector"`
	MaxConcurrency       int    `yaml:"maxConcurrency"`
	IsBenchmarkMode      bool   `yaml:"isBenchmarkMode"`
	AsyncCommit          bool   `yaml:"asyncCommit"`
}

func defaultConfig() *Config {
	return &Config{
		EvmProcessorSelector: "serial",
		AsyncCommit:          false,
		MaxConcurrency:       runtime.NumCPU(),
	}
}

var GlobalConfig *Config

func init() {
	GlobalConfig = defaultConfig()
}

func GetGlobalConfig() *Config {
	return GlobalConfig
}

func LoadConfig(path string) error {
	if path == "" {
		return nil
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	c := &Config{}
	err = toml.Unmarshal(content, c)
	if err != nil {
		return err
	}
	GlobalConfig = c
	return nil
}
