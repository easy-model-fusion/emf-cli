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

// Empty checks if the models slice is empty.
func Empty(models []model.Model) bool {
	// No models currently downloaded
	if len(models) == 0 {
		pterm.Info.Println("Models list is empty.")
		return true
	}
	return false
}

// Exists verifies if a model exists already in the configuration file or not
func Exists(models []model.Model, name string) (bool, error) {
	for _, currentModel := range models {
		println(currentModel.Name)
		if currentModel.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// Contains checks if a model exists in a slice
func Contains(models []model.Model, model model.Model) bool {
	for _, item := range models {
		if model == item {
			return true
		}
	}
	return false
}

// Difference returns the models in `parentSlice` that are not present in `subSlice`
func Difference(parentSlice, subSlice []model.Model) []model.Model {
	var difference []model.Model
	for _, s1 := range parentSlice {
		if !Contains(subSlice, s1) {
			difference = append(difference, s1)
		}
	}
	return difference
}

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

// GetAllModelNames retrieves model names from the configuration.
func GetAllModelNames() ([]string, error) {
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

// GetModelsByNames retrieves the models by their names given an input slice.
func GetModelsByNames(models []model.Model, namesSlice []string) []model.Model {
	// Create a map for faster lookup
	namesMap := utils.MapFromArrayString(namesSlice)

	// Slice of all the models that were found
	var namesModels []model.Model

	// Find the requested models
	for _, existingModel := range models {
		// Check if this model exists and adds it to the result
		if _, exists := namesMap[existingModel.Name]; exists {
			namesModels = append(namesModels, existingModel)
		}
	}

	return namesModels
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

// RemoveModelPhysically only removes the model from the project's downloaded models
func RemoveModelPhysically(model model.Model) error {

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

// RemoveAllModels removes all the models and updates the configuration file.
func RemoveAllModels() error {

	// Get the models from the configuration file
	models, err := GetModels()
	if err != nil {
		return err
	}

	// Trying to remove every model
	for _, item := range models {
		_ = RemoveModelPhysically(item)
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

// RemoveModelsByNames filters out specified models, removes them and updates the configuration file.
func RemoveModelsByNames(models []model.Model, modelsNamesToRemove []string) error {
	// Find all the models that should be removed
	modelsToRemove := GetModelsByNames(models, modelsNamesToRemove)

	// Trying to remove the models
	for _, item := range modelsToRemove {
		_ = RemoveModelPhysically(item)
	}

	// Find all the remaining models
	remainingModels := Difference(models, modelsToRemove)

	// Update the models
	viper.Set("models", remainingModels)

	// Attempt to write the configuration file
	err := WriteViperConfig()
	if err != nil {
		return err
	}

	return nil

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

		models, err := GetModels()
		if err != nil {
			return err
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

			exist, err := Exists(models, name)
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
