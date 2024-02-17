package config

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/script"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
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

// AddModel adds models to configuration file
func AddModel(updatedModels []model.Model) error {
	// Get existent models
	configModels, err := GetModels()
	if err != nil {
		return err
	}

	// Keeping those that haven't changed
	unchangedModels := model.Difference(configModels, updatedModels)

	// Combining the unchanged models with the updated models
	models := append(unchangedModels, updatedModels...)

	// Update the models
	viper.Set("models", models)

	// Attempt to write the configuration file
	err = WriteViperConfig()

	if err != nil {
		return err
	}

	return nil
}

// RemoveModelPhysically only removes the model from the project's downloaded models
func RemoveModelPhysically(model model.Model) error {

	// Path to the model
	modelPath := filepath.Join(app.ModelsDownloadPath, model.Name)

	// Starting client spinner animation
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing model %s...", model.Name))

	// Check if the model_path exists
	if exists, err := utils.IsExistingPath(modelPath); err != nil {
		// Skipping model : an error occurred
		spinner.Fail(err)
		return err
	} else if exists {
		// Model path is in the current project

		// Split the path into a slice of strings
		directories := utils.SplitPath(modelPath)

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

			// Delete directory if empty
			err = utils.DeleteDirectoryIfEmpty(path)
			if err != nil {
				spinner.Fail(err)
			}
		}
		spinner.Success(fmt.Sprintf("Removed model %s", model.Name))
	} else {
		// Model path is not in the current project
		spinner.Warning(fmt.Sprintf("Model '%s' was not found in the project directory. The model will be removed from this project's configuration file.", model.Name))
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

	// User did not add any model yet
	if len(models) == 0 {
		pterm.Info.Printfln("There is no models to be removed.")
		return nil
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
	modelsToRemove := model.GetModelsByNames(models, modelsNamesToRemove)

	// Indicate the models that were not found in the configuration file
	notFoundModels := utils.SliceDifference(modelsNamesToRemove, model.GetNames(modelsToRemove))
	if len(notFoundModels) != 0 {
		pterm.Warning.Println(fmt.Sprintf("The following models were not found in the configuration file : %s", notFoundModels))
	}

	// User did not provide any input
	if len(modelsToRemove) == 0 {
		pterm.Info.Printfln("No valid models were inputted.")
		return nil
	}

	// Trying to remove the models
	for _, item := range modelsToRemove {
		_ = RemoveModelPhysically(item)
	}

	// Find all the remaining models
	remainingModels := model.Difference(models, modelsToRemove)

	// Update the models
	viper.Set("models", remainingModels)

	// Attempt to write the configuration file
	err := WriteViperConfig()
	if err != nil {
		return err
	}

	return nil
}

// DownloadModel downloads physically a model.
func DownloadModel(modelObj model.Model) model.Model {

	// Exclude from download if not requested
	if !modelObj.AddToBinary {
		return modelObj
	}

	// Reset in case the download fails
	modelObj.AddToBinary = false
	overwrite := false

	// Get mandatory model data for the download script
	modelName := modelObj.Name
	moduleName := modelObj.Config.Module
	className := modelObj.Config.Class

	// Local path where the model will be downloaded
	downloadPath := app.ModelsDownloadPath
	modelPath := filepath.Join(downloadPath, modelName)

	// Check if the model_path already exists
	if exists, err := utils.IsExistingPath(modelPath); err != nil {
		// Skipping model : an error occurred
		return modelObj
	} else if exists {
		// Model path already exists : ask the user if he would like to overwrite it
		overwrite = utils.AskForUsersConfirmation(fmt.Sprintf("Model '%s' already downloaded at '%s'. Do you want to overwrite it?", modelName, modelPath))

		// User does not want to overwrite : skipping to the next model
		if !overwrite {
			return modelObj
		}
	}

	// Prepare the script arguments
	downloaderArgs := script.DownloaderArgs{
		DownloadPath: downloadPath,
		ModelName:    modelName,
		ModelModule:  moduleName,
		ModelClass:   className,
		Overwrite:    overwrite,
	}
	args := script.ProcessArgsForDownload(downloaderArgs)

	// Run the script to download the model
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading model '%s'...", modelName))
	var scriptModel, err, exitCode = utils.ExecuteScript(".venv", script.DownloaderName, args)

	// An error occurred while running the script
	if err != nil {
		spinner.Fail(err)
		switch exitCode {
		case 2:
			pterm.Info.Println("Run the 'add custom' command to manually add the model.")
		}
		return modelObj
	}

	// No data was returned by the script
	if scriptModel == nil {
		spinner.Fail(fmt.Sprintf("The script didn't return any data for '%s'", modelName))
		return modelObj
	}

	// Unmarshall JSON response
	var dsm script.DownloaderModel
	err = json.Unmarshal(scriptModel, &dsm)
	if err != nil {
		return modelObj
	}

	// Download was successful
	spinner.Success(fmt.Sprintf("Successfully downloaded model '%s'", modelName))

	// Update the model for the configuration file
	modelObj.Config = model.MapToConfigFromScriptDownloadModel(modelObj.Config, dsm)
	modelObj.AddToBinary = true

	return modelObj
}
