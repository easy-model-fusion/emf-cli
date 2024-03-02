package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/codegen"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
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

// AddModels adds models to configuration file
func AddModels(updatedModels []model.Model) error {
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
func RemoveModelPhysically(modelName string) error {

	// Path to the model
	modelPath := filepath.Join(downloader.DirectoryPath, modelName)

	// Starting client spinner animation
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing model %s...", modelName))

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
func RemoveModelsByNames(models []model.Model, modelsNamesToRemove []string) error {
	// Find all the models that should be removed
	modelsToRemove := model.GetModelsByNames(models, modelsNamesToRemove)

	// Indicate the models that were not found in the configuration file
	notFoundModels := stringutil.SliceDifference(modelsNamesToRemove, model.GetNames(modelsToRemove))
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

// GenerateModelsPythonCode generates the python code for the given models
func GenerateModelsPythonCode(models []model.Model) error {
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
