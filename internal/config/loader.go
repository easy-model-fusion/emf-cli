package config

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/spf13/viper"
)

var FilePath string

// Load loads the current configuration file
func Load(confDirPath string) error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDirPath)

	// Attempt to read the configuration file
	return viper.ReadInConfig()
}

// UpdateConfigFilePath updates configuration file path
func UpdateConfigFilePath() string {
	FilePath = app.UI().AskForUsersInput("Enter the configuration file path")
	return FilePath
}

func init() {
	// Default configuration file path
	FilePath = "."
}
