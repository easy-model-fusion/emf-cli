package confgen

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
