package app

import (
	"github.com/spf13/viper"
)

func LoadConfFile(confDirPath string) error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDirPath)

	// Attempt to read the configuration file
	return viper.ReadInConfig()
}
