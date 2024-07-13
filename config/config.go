package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ConfigRedis ConfigRedis `yaml:"redis"`
	ConfigHTTP  ConfigHTTP  `yaml:"http"`
}

type ConfigRedis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type ConfigHTTP struct {
	Host string `yaml:"host"`
}

// New return new config
func New() (*Config, error) {
	configFileData, err := os.ReadFile("config.yaml")
	if err != nil {
		slog.Error("Error while reading config file", "err", err)

		return nil, err
	}

	config := new(Config)

	err = yaml.Unmarshal(configFileData, config)
	if err != nil {
		slog.Error("Error while marshalling config file", "err", err)

		return nil, err
	}

	return config, nil
}
