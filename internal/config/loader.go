package config

import (
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/spf13/viper"
)

var FilePath string

// Load loads the current configuration file
func Load() error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(FilePath)

	// Attempt to read the configuration file
	return viper.ReadInConfig()
}

func UpdateConfigFilePath() {
	FilePath = utils.AskForUsersInput("Enter the configuration file path")
}

func init() {
	// Default configuration file path
	FilePath = "."
}
