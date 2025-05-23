package config

import (
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
)

type Config struct {
	IsParallel      bool `yaml:"isParallel"`
	MaxConcurrency  int  `yaml:"maxConcurrency"`
	IsBenchmarkMode bool `yaml:"isBenchmarkMode"`
	AsyncCommit     bool `yaml:"asyncCommit"`
}

func defaultConfig() *Config {
	return &Config{
		AsyncCommit:    false,
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
