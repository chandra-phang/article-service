package configloader

import (
	"errors"

	"github.com/spf13/viper"
)

var root = RootConfig{}

func GetRootConfig() RootConfig {
	return root
}

func LoadConfigFromFile(filePath string) error {
	cfg, err := loadConfig(filePath)
	if err != nil {
		return err
	}
	root = cfg
	return nil
}

func loadConfig(filePath string) (RootConfig, error) {
	cfg := RootConfig{}

	// setting config file on viper
	viper.SetConfigFile(filePath)

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	// updating config on global variable
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	isEmptyConfig := cfg.AppConfig == AppConfig{} &&
		cfg.DbConfig == DbConfig{}

	if isEmptyConfig {
		return cfg, errors.New("empty config file loaded")
	}

	return cfg, nil
}
