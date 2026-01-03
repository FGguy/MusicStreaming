package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = "musicstreaming"
	ConfigFileType = "yaml"
)

type Config struct {
	MusicDirectories []string `mapstructure:"music-directories"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	ConfigPath, ok := os.LookupEnv("CONFIG_PATH")
	if ok {
		viper.AddConfigPath(ConfigPath)
	} else {
		viper.AddConfigPath(".") // Default to current directory
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
