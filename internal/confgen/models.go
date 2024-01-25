package confgen

import (
	"github.com/spf13/viper"
)

func AddModel(models []string) error { // Get the current working directory
	// Access and print the original models list
	originalModelsList := viper.GetStringSlice("models")

	// Add additional items to the models list
	updatedModelsList := append(originalModelsList, models...)

	// Update the models list in the configuration
	viper.Set("models", updatedModelsList)

	// Attempt to write the configuration file
	return viper.WriteConfig()
}
