package config

import (
	"github.com/spf13/viper"
)

func AddModel(models []string) error {
	originalModelsList := viper.GetStringSlice("models")
	updatedModelsList := append(originalModelsList, models...)
	viper.Set("models", updatedModelsList)

	// Attempt to write the configuration file
	return viper.WriteConfig()
}

func RemoveModel(models []string) error {
	return viper.WriteConfig()
}

func RemoveAllModels() error {
	viper.Set("models", []string{})
	return viper.WriteConfig()
}
