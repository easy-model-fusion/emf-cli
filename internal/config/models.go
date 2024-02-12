package config

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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

// GetModelNames retrieves models from the configuration.
func GetModelNames() ([]string, error) {
	models, err := GetModels()
	if err != nil {
		return nil, err
	}

	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames, nil
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

func RemoveModel(model model.Model) error {

	// Nothing to remove if not downloaded
	if !model.AddToBinary {
		return nil
	}

	// Path to the model
	modelPath := filepath.Join(app.ModelsDownloadPath, model.Name)

	// Starting client spinner animation
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing model %s...", model.Name))

	// Check if the model_path already exists
	if exists, err := utils.IsExistingPath(modelPath); err != nil {
		// Skipping model : an error occurred
		spinner.Fail(err)
		return err
	} else if exists {
		// Model path is in the current project

		// Split the path into a slice of strings
		directories := utils.ArrayFromPath(modelPath)

		// Removing model
		err := os.RemoveAll(modelPath)
		if err != nil {
			spinner.Fail(err)
			return err
		}

		// Excluding the tail since it has already been removed
		directories = directories[:len(directories)-1]

		// Cleaning up : removing every empty directory on the way to the model (from tail to head)
		for i := len(directories) - 1; i >= 0; i-- {
			// Build path to parent directory
			path := filepath.Join(directories[:i+1]...)

			// Check if the directory is empty
			if empty, err := utils.IsDirectoryEmpty(path); err != nil {
				spinner.Fail(err)
				return err
			} else if empty {
				// Current directory is empty : removing it
				err := os.Remove(path)
				if err != nil {
					spinner.Fail(err)
					return err
				}
			}
		}
		spinner.Success(fmt.Sprintf("Removed model %s", model.Name))
	} else {
		// Model path is not in the current project
		spinner.Info(fmt.Sprintf("Model '%s' was not found in the project directory. It might have been removed manually or belongs to another project. The model will be removed from this project's configuration file.", model.Name))
	}
	return nil
}

// RemoveModels filters out specified models and writes to the configuration file.
func RemoveModels(models []model.Model, modelsToRemove []string) error {
	// Filter out the models to be removed
	updatedModels, removedModels := RemoveModelsFromList(models, modelsToRemove)

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

	// Get the models from the configuration file
	models, err := GetModels()
	if err != nil {
		return err
	}

	// Trying to remove every model
	for _, item := range models {
		_ = RemoveModel(item)
	}

	// Empty the models
	viper.Set("models", []string{})

	// Attempt to write the configuration file
	err = WriteViperConfig()
	if err != nil {
		return err
	}

	return nil
}

func RemoveModelsFromList(currentModels []model.Model, modelsToRemove []string) ([]model.Model, []string) {
	// Create a map for faster lookup
	modelsMap := utils.MapFromArrayString(modelsToRemove)

	// Filter out the models to be removed
	var updatedModels []model.Model
	var removedModels []string
	for _, existingModel := range currentModels {
		if _, exists := modelsMap[existingModel.Name]; !exists {
			// Keep the model if it's not in the modelsToRemove list
			updatedModels = append(updatedModels, existingModel)
		} else {
			// Indicate which model was effectively removed
			removedModels = append(removedModels, existingModel.Name)
		}
	}

	return updatedModels, removedModels
}

// ModelExists verifies if a model exists already in the configuration file or not
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

	return false, nil
}

// ValidModelName returns an error if the which arg is not a valid model name.
func ValidModelName() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		if app.H() == nil {
			return fmt.Errorf("hugging face api is not initialized")
		}

		for _, name := range args {
			valid, err := app.H().ValidModel(name)
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
