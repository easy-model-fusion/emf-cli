package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"os"
	"testing"
)

// Sets the configuration file with the given models
func setupConfigFile(models model.Models) error {
	// Load configuration file
	err := config.GetViperConfig(".")
	if err != nil {
		return err
	}
	// Write models to the config file
	viper.Set("models", models)
	return config.WriteViperConfig()
}

// Tests removeModels with valid model from models list
func TestRemoveModels(t *testing.T) {
	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "model1",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	models = append(models, model.Model{
		Name:            "model2",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize selected models list
	var selectedModels []string
	selectedModels = append(selectedModels, "model2")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Load configuration file
	err := config.GetViperConfig(".")
	test.AssertEqual(t, err, nil, "No error expected")

	// Remove selected models
	warning, info, err := removeModels(models, selectedModels)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warning, "", "No warning expected")
	test.AssertEqual(t, info, "", "No info expected")
}

// Tests removeModels with no models selected
func TestRemoveModels_WithNoModelsSelected(t *testing.T) {
	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "test",
		Path:            "path/to/model",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize empty selected models list
	var selectedModels []string

	// Remove selected models
	warning, info, err := removeModels(models, selectedModels)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warning, "", "No warning expected")
	test.AssertEqual(t, info, "There is no models to be removed.", "Information message expected")
}

// Tests removeModels with no models configured
func TestRemoveModels_WithNoConfiguredModels(t *testing.T) {
	// initialize empty models list
	var models model.Models
	// Initialize selected models list
	var selectedModels []string
	selectedModels = append(selectedModels, "testModel")

	// Remove selected models
	warning, info, err := removeModels(models, selectedModels)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warning, "", "No warning expected")
	test.AssertEqual(t, info, "There is no models to be removed.", "Information message expected")
}

// Tests selectModelsToDelete
func TestSelectModelsToDelete(t *testing.T) {
	// Initialize model names list
	var modelNames []string
	modelNames = append(modelNames, "model1")
	modelNames = append(modelNames, "model2")
	modelNames = append(modelNames, "model3")
	var expectedSelections []string
	expectedSelections = append(expectedSelections, "model1")
	expectedSelections = append(expectedSelections, "model3")

	// Create ui mock
	ui := test.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Select models
	selectedModels := selectModelsToDelete(modelNames, false)

	// Assertions
	test.AssertEqual(t, len(selectedModels), 2, "2 models should be returned")
	test.AssertEqual(t, selectedModels[0], expectedSelections[0])
	test.AssertEqual(t, selectedModels[1], expectedSelections[1])
}

// Tests selectModelsToDelete with all models selected
func TestSelectModelsToDelete_WithSelectedAll(t *testing.T) {
	// Initialize model names list
	var modelNames []string
	modelNames = append(modelNames, "model1")
	modelNames = append(modelNames, "model2")

	// Select models
	selectedModels := selectModelsToDelete(modelNames, true)

	// Assertions
	test.AssertEqual(t, len(selectedModels), len(modelNames), "All models should be returned")
	test.AssertEqual(t, selectedModels[0], modelNames[0])
	test.AssertEqual(t, selectedModels[1], modelNames[1])
}

// Tests selectModelsToDelete with empty model names list
func TestSelectModelsToDelete_WithEmptyModelNames(t *testing.T) {
	// Initialize model names list
	var modelNames []string

	// Select models
	selectedModels := selectModelsToDelete(modelNames, false)

	// Assertions
	test.AssertEqual(t, len(selectedModels), 0, "Empty models list should be returned")
}

// Tests processRemove method without any model entered in args
func TestProcessRemove_WithoutArgs(t *testing.T) {
	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "model1",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	models = append(models, model.Model{
		Name:            "model2",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize selected models list
	var args []string

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	_, _, err = processRemove(args, true)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 0, "All models should be removed.")
}

// Tests processRemove method with some models entered in args
func TestProcessRemove_WithArgs(t *testing.T) {
	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "model1",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	models = append(models, model.Model{
		Name:            "model2",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize selected models list
	var args []string
	args = append(args, "model2")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	_, _, err = processRemove(args, false)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 1, "Only one model should be left.")
	test.AssertEqual(t, newModels[0].Name, "model1", "Model1 shouldn't be deleted")
}

// Tests RunModelRemove
func TestRunModelRemove_WithNoModels(t *testing.T) {
	// initialize models list
	var models model.Models
	// Initialize selected models list
	var args []string

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	RunModelRemove(args, true)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 0, "Only one model should be left.")
}

// Tests RunModelRemove
func TestRunModelRemove(t *testing.T) {
	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "model1",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	models = append(models, model.Model{
		Name:            "model2",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize selected models list
	var args []string
	args = append(args, "model2")
	args = append(args, "invalidModel")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	RunModelRemove(args, false)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 1, "Only one model should be left.")
	test.AssertEqual(t, newModels[0].Name, "model1", "Model1 shouldn't be deleted")
}

// Tests processRemove method should throw error on invalid configuration file path
func TestProcessRemove_WithErrorOnLoadingConfigurationFile(t *testing.T) {
	// Initialize selected models list
	var args []string

	//create mock UI
	ui := test.MockUI{UserInputResult: "path/test"}
	app.SetUI(ui)

	// Process remove
	_, _, err := processRemove(args, false)
	test.AssertNotEqual(t, err, nil, "Error expected while loading configuration file")
}

// Tests RunModelRemove with invalid configuration file path
func TestRunModelRemove_WitInvalidConfigPath(t *testing.T) {
	// Get current Directory
	currentDir, err := os.Getwd()
	test.AssertEqual(t, err, nil, "No error expected while getting current directory")

	// initialize models list
	var models model.Models
	models = append(models, model.Model{
		Name:            "model1",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	models = append(models, model.Model{
		Name:            "model2",
		Path:            "path/to/model1",
		Source:          "CUSTOM",
		AddToBinaryFile: true,
		IsDownloaded:    true,
	})
	// Initialize selected models list
	var args []string
	args = append(args, "invalidModel")

	// Create temporary configuration file
	ts := test.TestSuite{}
	confDir := ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err = setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Change current Directory
	os.Chdir(currentDir)

	//create mock UI
	ui := test.MockUI{UserInputResult: "path/test"}
	app.SetUI(ui)

	// Process remove
	RunModelRemove(args, false)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	os.Chdir(confDir)
	err = config.Load(".")
	test.AssertEqual(t, err, nil, "No error expected while loading configuration file")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 2, "Only one model should be left.")
	test.AssertEqual(t, newModels[0].Name, "model1", "Model1 shouldn't be deleted")
	test.AssertEqual(t, newModels[1].Name, "model2", "Model2 shouldn't be deleted")
}
