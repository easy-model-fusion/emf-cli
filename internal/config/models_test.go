package config

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/viper"
)

type AlternativeStructure struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
}

func init() {
	app.Init("", "")
}

// TestErrorOnAddModelWithEmptyViper tests the AddModel function
func TestAddModel(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}

	// Call the AddModel function to add new models
	newModels := []model.Model{
		{
			Name:          "model3",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline3"}},
			DirectoryPath: "/path/to/model3",
			AddToBinary:   true},
		{
			Name:          "model4",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline3"}},
			DirectoryPath: "/path/to/model4",
			AddToBinary:   false},
	}
	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModel(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels, err := GetModels()
	test.AssertEqual(t, err, nil, "Error while getting updated models.")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestErrorOnAddModelWithEmptyViper tests the AddModel function with an empty config file
func TestAddModelOnEmptyConfFile(t *testing.T) {
	// Use the setup function
	initialModels := []model.Model{}
	confDir, initialConfigFile := setupConfigDir(t)

	// Call the AddModel function to add new models
	newModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}

	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModel(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels, err := GetModels()
	test.AssertEqual(t, err, nil, "Error while getting updated models.")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestErrorOnAddModelWithEmptyViper tests the AddModel function with a missing config file
func TestErrorOnAddModelWithEmptyViper(t *testing.T) {
	viper.Reset()
	// Call the AddModel function to add new models
	newModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}

	err := AddModel(newModels)
	test.AssertNotEqual(t, err, nil, "Should get error while updating configuration file.")
}

// TestGetModels_MissingConfig tests the GetModels function
func TestGetModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}
	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the GetModels function
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	retrievedModels, err := GetModels()
	test.AssertEqual(t, err, nil, "Error while retrieving models from configuration.")

	// Assert that the models have been retrieved correctly
	test.AssertEqual(t, len(retrievedModels), len(initialModels), "Retrieved models do not match initial models.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestGetModels_MissingConfig tests the GetModels function with a missing config file
func TestGetModels_MissingConfig(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	var initialModels []model.Model
	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the GetModels function
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")

	// Assert that the models have been retrieved correctly
	retrievedModels, err := GetModels()
	test.AssertEqual(t, len(retrievedModels), 0, "Retrieved models should be empty.")
	test.AssertEqual(t, err, nil, "Retrieving models should not have failed.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestIsModelsEmpty_EmptyModels tests the Empty function with an empty models slice.
func TestIsModelsEmpty_EmptyModels(t *testing.T) {
	// Init
	var models []model.Model

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestIsModelsEmpty_NonEmptyModels tests the Empty function with a non-empty models slice.
func TestIsModelsEmpty_NonEmptyModels(t *testing.T) {
	// Init
	models := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestRemoveModels_Success tests the RemoveModels function for successful removal of specified models.
func TestRemoveModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}
	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the RemoveModels function
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = RemoveModels(initialModels, []string{"model1"})
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Get the newly stored data
	var updatedModels []model.Model
	err = viper.UnmarshalKey("models", &updatedModels)
	test.AssertEqual(t, err, nil, "Error while unmarshalling models from configuration file.")

	// Assert that the models have been removed correctly
	expectedModels := initialModels[1:]
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "The selected models were not removed correctly.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestRemoveAllModels_Success tests the RemoveAllModels function for successful removal of all models.
func TestRemoveAllModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{
			Name:          "model1",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline1"}},
			DirectoryPath: "/path/to/model1",
			AddToBinary:   true},
		{
			Name:          "model2",
			Config:        model.Config{Diffusers: model.Diffusers{PipeLine: "pipeline2"}},
			DirectoryPath: "/path/to/model2",
			AddToBinary:   false},
	}
	err := setupConfigFile(t, initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the RemoveAllModels function
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = RemoveAllModels()
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Get the newly stored data
	var updatedModels []model.Model
	err = viper.UnmarshalKey("models", &updatedModels)
	test.AssertEqual(t, err, nil, "Error while unmarshalling models from configuration file.")

	// Assert that the models have been removed correctly
	var expectedModels []model.Model
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Not all models were removed.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// setupConfigDir creates a temporary directory.
func setupConfigDir(t *testing.T) (string, string) {
	// Create a temporary directory for the test
	confDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}

	// Set up a temporary config file with some initial models
	initialConfigFile := filepath.Join(confDir, "config.yaml")

	return confDir, initialConfigFile
}

// setupConfigFile creates a configuration file.
func setupConfigFile(t *testing.T, filePath string, models []model.Model) error {
	file, err := os.Create(filePath)
	defer utils.CloseFile(file)
	if err != nil {
		return err
	}

	if len(models) > 0 {

		// Write models to the config file
		err = writeToConfigFile(file, "models", models)
		if err != nil {
			return err
		}

	} else {

		// Setting up alternative config data
		var alternative = []AlternativeStructure{
			{Field1: "data1", Field2: "data2", Field3: "data3"},
		}

		// Writing alternative data to the config file
		err = writeToConfigFile(file, "alternative", alternative)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeToConfigFile writes the specified item to the configuration file.
func writeToConfigFile(file *os.File, itemName string, itemData interface{}) error {
	// Marshal models to JSON
	jsonData, err := json.Marshal(itemData)
	if err != nil {
		return err
	}

	// Unmarshal JSON to YAML
	var yamlData interface{}
	err = yaml.Unmarshal(jsonData, &yamlData)
	if err != nil {
		return err
	}

	// Convert YAML data to []byte
	yamlBytes, err := yaml.Marshal(yamlData)
	if err != nil {
		return err
	}

	// Write models to the config file
	_, err = file.WriteString(fmt.Sprintf("%s:\n", itemName))
	if err != nil {
		return err
	}

	// Write YAML data to config file
	_, err = file.Write(yamlBytes)
	if err != nil {
		return err
	}

	return nil
}

// cleanConfDir removes the temporary directory created during testing.
func cleanConfDir(t *testing.T, confDir string) {
	if err := os.RemoveAll(confDir); err != nil {
		t.Errorf("Error cleaning up temporary directory: %v", err)
	}
}

func TestModelExists_OnExistentModel(t *testing.T) {
	// TODO implement this after fixing writeToConfigFile
	// bug ==> writing name and pipeline
	// cause of bug ==> name of variables key values
}

func TestModelExists_OnNotExistentModel(t *testing.T) {
	// TODO implement this after fixing writeToConfigFile
	// bug ==> writing name and pipeline
	// cause of bug ==> name of variables key values
}
