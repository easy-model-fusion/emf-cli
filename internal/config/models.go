package config

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/spf13/viper"
)

// GetModels retrieves models from the configuration.
func GetModels() ([]model.Model, error) {
	// Define a slice for models
	var models []model.Model

	// Retrieve models using the generic function
	if err := utils.GetViperItem("models", &models); err != nil {
		return nil, err
	}
	return models, nil
}

// IsModelsEmpty checks if the models slice is empty.
func IsModelsEmpty(models []model.Model) bool {
	logger := app.L().WithTime(false)

	// No models currently downloaded
	if len(models) == 0 {
		logger.Error("Models list is empty. Nothing to remove.")
		return true
	}

	return false
}

func AddModel(models []string) error {
	originalModelsList := viper.GetStringSlice("models")
	updatedModelsList := append(originalModelsList, models...)
	viper.Set("models", updatedModelsList)

	// Attempt to write the configuration file
	return viper.WriteConfig()
}

// RemoveModels filters out specified models and writes to the configuration file.
func RemoveModels(models []model.Model, modelsToRemove []string) error {

	// Create a map for faster lookup
	modelsMap := utils.MapFromArrayString(modelsToRemove)

	// Filter out the models to be removed
	var updatedModels []model.Model
	for _, existingModel := range models {
		if _, exists := modelsMap[existingModel.Name]; !exists {
			// Keep the model if it's not in the modelsToRemove list
			updatedModels = append(updatedModels, existingModel)
		}
	}

	// Attempt to write the configuration file
	return utils.WriteViperConfig()
}

// RemoveAllModels empties the models list and writes to the configuration file.
func RemoveAllModels() error {
	// Empty the models
	viper.Set("models", []string{})

	// Attempt to write the configuration file
	return utils.WriteViperConfig()
}
