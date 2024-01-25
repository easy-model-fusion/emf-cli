package app

import (
	"github.com/spf13/viper"
)

func LoadConfFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		L().Fatal("Error while loading the config file")
	}
	L().Info("Config file loaded")
}
