package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
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

// RemoveItemPhysically only removes the item from the project's downloaded item
func RemoveItemPhysically(itemPath string) error {

	// Check if the item_path exists
	if exists, err := fileutil.IsExistingPath(itemPath); err != nil {
		// Skipping item : an error occurred
		return err
	} else if exists {

		// Split the path into a slice of strings
		directories := stringutil.SplitPath(itemPath)

		// Removing item
		err = os.RemoveAll(itemPath)
		if err != nil {
			return err
		}

		// Excluding the tail since it has already been removed
		directories = directories[:len(directories)-1]

		// Cleaning up : removing every empty directory on the way to the item (from tail to head)
		for i := len(directories) - 1; i >= 0; i-- {
			// Build path to parent directory
			path := filepath.Join(directories[:i+1]...)

			// Delete directory if empty
			err = fileutil.DeleteDirectoryIfEmpty(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveAllModels removes all the models and updates the configuration file.
func RemoveAllModels() (info string, err error) {

	// Get the models from the configuration file
	models, err := GetModels()
	if err != nil {
		return info, err
	}

	// User did not add any model yet
	if len(models) == 0 {
		info = "There is no models to be removed."
		return info, err
	}

	// Trying to remove every model
	for _, item := range models {
		modelPath := filepath.Join(app.DownloadDirectoryPath, item.Name)
		spinner := app.UI().StartSpinner(fmt.Sprintf("Removing item %s...", item.Name))
		err = RemoveItemPhysically(modelPath)
		if err != nil {
			spinner.Fail("failed to remove item")
			continue
		} else {
			spinner.Success()
		}
	}

	// Empty the models
	viper.Set("models", []string{})

	// Attempt to write the configuration file
	err = WriteViperConfig()
	return info, err
}

// RemoveModelsByNames filters out specified models, removes them and updates the configuration file.
func RemoveModelsByNames(models model.Models, modelsNamesToRemove []string) (warning string, info string, err error) {
	// Find all the models that should be removed
	modelsToRemove := models.FilterWithNames(modelsNamesToRemove)

	// Indicate the models that were not found in the configuration file
	notFoundModels := stringutil.SliceDifference(modelsNamesToRemove, modelsToRemove.GetNames())
	if len(notFoundModels) != 0 {
		warning = fmt.Sprintf("The following models were not found in the configuration file : %s", notFoundModels)
	}

	// User did not provide any input
	if len(modelsToRemove) == 0 {
		info = "No valid models were inputted."
		return warning, info, err
	}

	// Trying to remove the models
	for _, item := range modelsToRemove {
		modelPath := filepath.Join(app.DownloadDirectoryPath, item.Name)
		spinner := app.UI().StartSpinner(fmt.Sprintf("Removing item %s...", item.Name))
		err = RemoveItemPhysically(modelPath)
		if err != nil {
			spinner.Fail("failed to remove item")
			continue
		} else {
			spinner.Success()
		}
	}

	// Find all the remaining models
	remainingModels := models.Difference(modelsToRemove)

	// Update the models
	viper.Set("models", remainingModels)

	// Attempt to write the configuration file
	err = WriteViperConfig()

	return warning, info, err
}

// Validate to validate a model before adding it.
func Validate(current model.Model, yes bool) (warning string, success bool, err error) {

	// Check if model is already configured
	models, err := GetModels()
	if err != nil {
		return warning, false, err
	}

	if models.ContainsByName(current.Name) {
		warning = fmt.Sprintf("Model '%s' is already configured", current.Name)
		return warning, false, err
	}

	// Build path for validation
	current.UpdatePaths()

	// Validate the model : if model is already downloaded
	downloaded, err := current.DownloadedOnDevice(true)
	if err != nil {
		return warning, false, err
	} else if downloaded && !current.AddToBinaryFile {
		// Model won't be downloaded but a version is already downloaded
		message := fmt.Sprintf("Model '%s' is already downloaded. Do you wish to delete it?", current.Name)
		overwrite := yes || app.UI().AskForUsersConfirmation(message)
		if !overwrite {
			warning = fmt.Sprintf("This model is already downloaded and should be checked manually %s", current.Name)
			return warning, false, err
		}

		// Removing model
		modelPath := filepath.Join(app.DownloadDirectoryPath, current.Name)
		spinner := app.UI().StartSpinner(fmt.Sprintf("Removing item %s...", current.Name))
		err = RemoveItemPhysically(modelPath)
		if err != nil {
			return warning, false, err
		} else {
			spinner.Success()
			return warning, true, err
		}
	} else if downloaded {
		// A version of the model is already downloaded
		message := fmt.Sprintf("Model '%s' is already downloaded. Do you wish to overwrite it?", current.Name)
		overwrite := yes || app.UI().AskForUsersConfirmation(message)
		if !overwrite {
			warning = fmt.Sprintf("This model is already downloaded and should be checked manually %s", current.Name)
			return warning, false, err
		}
	}

	return warning, true, err
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
