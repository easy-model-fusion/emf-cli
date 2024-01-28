package config

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetModels retrieves models from the configuration.
func GetModels() ([]model.Model, error) {
	// Define a slice for models
	var models []model.Model

	// Retrieve models using the generic function
	if err := GetViperItem("models", &models); err != nil {
		return nil, err
	}
	return models, nil
}

// GetModelsNames retrieves models from the configuration.
func GetModelsNames() ([]string, error) {
	models, err := GetModels()
	if err != nil {
		return nil, err
	}

	var modelsNames []string
	for _, item := range models {
		modelsNames = append(modelsNames, item.Name)
	}
	return modelsNames, nil
}

// IsModelsEmpty checks if the models slice is empty.
func IsModelsEmpty(models []model.Model) bool {
	// No models currently downloaded
	if len(models) == 0 {
		pterm.Info.Println("Models list is empty.")
		return true
	}
	return false
}

// AddModel adds models to configuration file
func AddModel(models []model.Model) error {
	// get existent models
	originalModelsList, err := GetModels()
	if err != nil {
		return err
	}
	// add new models
	updatedModels := append(originalModelsList, models...)
	// Update the models
	viper.Set("models", updatedModels)

	// Attempt to write the configuration file
	err = WriteViperConfig()

	if err != nil {
		return err
	}

	return nil
}

// RemoveModels filters out specified models and writes to the configuration file.
func RemoveModels(models []model.Model, modelsToRemove []string) error {

	// Create a map for faster lookup
	modelsMap := utils.MapFromArrayString(modelsToRemove)

	// Filter out the models to be removed
	var updatedModels []model.Model
	var removedModels []string
	for _, existingModel := range models {
		if _, exists := modelsMap[existingModel.Name]; !exists {
			// Keep the model if it's not in the modelsToRemove list
			updatedModels = append(updatedModels, existingModel)
		} else {
			// Indicate which model was effectively removed
			removedModels = append(removedModels, existingModel.Name)
		}
	}

	// Create a map for faster lookup
	removedModelsMap := utils.MapFromArrayString(removedModels)

	// Displaying if the selected models where removed or not
	for _, modelToRemove := range modelsToRemove {
		if _, exists := removedModelsMap[modelToRemove]; exists {
			pterm.Success.Printfln("Model '%s' was removed.", modelToRemove)
		} else {
			pterm.Warning.Printfln("Model '%s' is not installed. Nothing to remove.", modelToRemove)
		}
	}

	// Update the models
	viper.Set("models", updatedModels)

	// Attempt to write the configuration file
	err := WriteViperConfig()
	if err != nil {
		return err
	}

	// TODO : remove the downloaded models : Issue #21 => removedModels

	return nil

}

// RemoveAllModels empties the models list and writes to the configuration file.
func RemoveAllModels() error {

	// Empty the models
	viper.Set("models", []string{})

	// Attempt to write the configuration file
	err := WriteViperConfig()
	if err != nil {
		return err
	}

	// TODO : remove the downloaded models : Issue #21

	return nil
}

func RemoveModelsFromList(currentModels []model.Model, modelsToRemove []string) []model.Model {
	// Create a map for faster lookup
	modelsMap := utils.MapFromArrayString(modelsToRemove)

	// Filter out the models to be removed
	var updatedModels []model.Model
	for _, existingModel := range currentModels {
		if _, exists := modelsMap[existingModel.Name]; !exists {
			// Keep the model if it's not in the modelsToRemove list
			updatedModels = append(updatedModels, existingModel)
		}
	}

	return updatedModels
}

func ModelExists(name string) (bool, error) {
	models, err := GetModels()
	if err != nil {
		return false, err
	}
	for _, currentModel := range models {
		println(currentModel.Name)
		if currentModel.Name == name {
			return true, nil
		}
	}
	println(models)

	return false, nil
}

// ValidModelName returns an error if the which arg is not a valid model name.
func ValidModelName() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}

		for _, name := range args {
			valid, err := huggingface.ValidModel(name)
			if err != nil {
				return err
			}
			if !valid {
				return fmt.Errorf("'%s' is not a valid model name", name)
			}

			// Load the configuration file
			err = Load(".")
			if err != nil {
				return err
			}

			exist, err := ModelExists(name)
			if err != nil {
				return err
			}
			if exist {
				return fmt.Errorf("'%s' model is already included in the project", name)
			}
		}

		return nil
	}
}
