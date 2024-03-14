package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// GetModels retrieves models from the configuration.
func GetModels() (model.Models, error) {
	// Define a slice for models
	var models model.Models

	// Retrieve models using the generic function
	if err := GetViperItem("models", &models); err != nil {
		return nil, err
	}
	return models, nil
}

// AddModels adds models to configuration file
func AddModels(updatedModels model.Models) error {
	// Get existent models
	configModels, err := GetModels()
	if err != nil {
		return err
	}

	// Keeping those that haven't changed
	unchangedModels := configModels.Difference(updatedModels)

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
func RemoveModelPhysically(modelName string) error {

	// Path to the model
	modelPath := filepath.Join(app.DownloadDirectoryPath, modelName)

	// Starting client spinner animation
	spinner := app.UI().StartSpinner(fmt.Sprintf("Removing model %s...", modelName))

	// Check if the model_path exists
	if exists, err := fileutil.IsExistingPath(modelPath); err != nil {
		// Skipping model : an error occurred
		spinner.Fail(err)
		return err
	} else if exists {
		// Model path is in the current project

		// Split the path into a slice of strings
		directories := stringutil.SplitPath(modelPath)

		// Removing model
		err = os.RemoveAll(modelPath)
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
			err = fileutil.DeleteDirectoryIfEmpty(path)
			if err != nil {
				spinner.Fail(err)
			}
		}
		spinner.Success(fmt.Sprintf("Removed model %s", modelName))
	} else {
		// Model path is not in the current project
		spinner.Warning(fmt.Sprintf("Model '%s' was not found in the project directory. The model will be removed from this project's configuration file.", modelName))
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
		_ = RemoveModelPhysically(item.Name)
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
func RemoveModelsByNames(models model.Models, modelsNamesToRemove []string) error {
	// Find all the models that should be removed
	modelsToRemove := models.FilterWithNames(modelsNamesToRemove)

	// Indicate the models that were not found in the configuration file
	notFoundModels := stringutil.SliceDifference(modelsNamesToRemove, modelsToRemove.GetNames())
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
		_ = RemoveModelPhysically(item.Name)
	}

	// Find all the remaining models
	remainingModels := models.Difference(modelsToRemove)

	// Update the models
	viper.Set("models", remainingModels)

	// Attempt to write the configuration file
	err := WriteViperConfig()
	if err != nil {
		return err
	}

	return nil
}

// Validate to validate a model before adding it.
func Validate(current model.Model) bool {

	// Check if model is already configured
	models, err := GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return false
	}

	if model.ContainsByName(models, current.Name) {
		pterm.Warning.Printfln("Model '%s' is already configured", current.Name)
		return false
	}

	// Build path for validation
	current = model.ConstructConfigPaths(current)

	// Validate the model : if model is already downloaded
	downloaded, err := model.ModelDownloadedOnDevice(current, true)
	if err != nil {
		pterm.Error.Println(err)
		return false
	} else if downloaded && !current.AddToBinaryFile {
		// Model won't be downloaded but a version is already downloaded
		message := fmt.Sprintf("Model '%s' is already downloaded. Do you wish to delete it?", current.Name)
		overwrite := app.UI().AskForUsersConfirmation(message)
		if !overwrite {
			pterm.Warning.Println("This model is already downloaded and should be checked manually", current.Name)
			return false
		}

		// Removing model
		err = RemoveModelPhysically(current.Name)
		if err != nil {
			return false
		}
	} else if downloaded {
		// A version of the model is already downloaded
		message := fmt.Sprintf("Model '%s' is already downloaded. Do you wish to overwrite it?", current.Name)
		overwrite := app.UI().AskForUsersConfirmation(message)
		if !overwrite {
			pterm.Warning.Println("This model is already downloaded and should be checked manually", current.Name)
			return false
		}
	}

	return true
}

// GenerateExistingModelsPythonCode generates the python code for all the configured models
func GenerateExistingModelsPythonCode() error {
	// Get existing models
	models, err := GetModels()
	if err != nil {
		return err
	}

	// Generating code for these models
	return GenerateModelsPythonCode(models)
}

// GenerateModelsPythonCode generates the python code for the given models
func GenerateModelsPythonCode(models model.Models) error {
	genFile := &codegen.File{
		Name: "generated_models.py",
		HeaderComments: []string{
			"Code generated by EMF",
			"DO NOT EDIT!",
		},
		Classes: []*codegen.Class{},
	}

	// loop through the models and add the imports and classes to the generated file
	for _, currentModel := range models {

		// remove duplicates and add the imports to the generated file
		for _, imp := range currentModel.GenImports() {
			found := false
			// search for duplicates
			for _, alreadyImp := range genFile.Imports {
				if alreadyImp.Equals(&imp) {
					found = true
					break
				}
			}
			if !found {
				// no duplicates found => add the import
				genFile.Imports = append(genFile.Imports, imp)
			}
		}

		genFile.Classes = append(genFile.Classes, currentModel.GenClass())

	}

	cg := codegen.NewPythonCodeGenerator(true)
	result, err := cg.Generate(genFile)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join("sdk", "generated_models.py"), []byte(result), 0644)
}
