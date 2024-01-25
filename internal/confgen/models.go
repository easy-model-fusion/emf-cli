package confgen

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/spf13/viper"
)

func AddModel(models []string) { // Get the current working directory
	err := app.LoadConfFile(".")
	if err != nil {
		app.L().Fatal("Error while loading configuration file." + err.Error())
	}

	// Access and print the original models list
	originalModelsList := viper.GetStringSlice("models")

	// Add additional items to the models list
	updatedModelsList := append(originalModelsList, models...)

	// Update the models list in the configuration
	viper.Set("models", updatedModelsList)

	// Attempt to write the configuration file
	if err := viper.WriteConfig(); err != nil {
		app.L().Fatal("Error while adding new model(s) to the config file")
	}

	app.L().Info("New model(s) successfully added.")
}
