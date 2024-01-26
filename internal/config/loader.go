package config

import (
	"github.com/spf13/viper"
)

// Load loads the current configuration file
func Load(confDirPath string) error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDirPath)

	// Attempt to read the configuration file
	return viper.ReadInConfig()
}
