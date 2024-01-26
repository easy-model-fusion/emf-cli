package config

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/spf13/viper"
)

func Load(confDirPath string) error {
	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confDirPath)

	// Attempt to read the configuration file
	return viper.ReadInConfig()
}

// LoadConfig Config loaded and return an error upon failure
func LoadConfig() error {
	logger := app.L().WithTime(false)
	if err := Load("."); err != nil {
		logger.Error("Error loading config file:" + err.Error())
		return err
	}
	return nil
}
