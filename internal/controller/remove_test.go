package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func init() {
	app.Init("", "")
	app.InitGit("", "")
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
	_ = ts.CreateConfigurationFileFullTestSuite(t, models)
	defer ts.CleanTestSuite(t)

	// Process remove
	_, _, err := processRemove(args, true)
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
	_ = ts.CreateConfigurationFileFullTestSuite(t, models)
	defer ts.CleanTestSuite(t)

	// Process remove
	_, _, err := processRemove(args, false)
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 1, "Only one model should be left.")
	test.AssertEqual(t, newModels[0].Name, "model1", "Model1 shouldn't be deleted")
}
