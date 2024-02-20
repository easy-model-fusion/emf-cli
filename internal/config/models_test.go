package config

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/script"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/emf-cli/test"
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
func setupConfigFile(filePath string, models []model.Model) error {
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

// setupModelDirectory creates a temporary directory for the model.
func setupModelDirectory(t *testing.T, modelPath string) {
	// Create a temporary path to the model for the test
	err := os.MkdirAll(modelPath, 0750)
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary data inside the model for the test
	file, err := os.CreateTemp(modelPath, "")
	if err != nil {
		t.Fatal(err)
	}
	utils.CloseFile(file)
}

// getModel initiates a basic model with an id as suffix
func getModel(suffix int) model.Model {
	idStr := fmt.Sprint(suffix)
	return model.Model{
		Name:               "model" + idStr,
		Config:             model.Config{Module: "module" + idStr, Class: "class" + idStr},
		ShouldBeDownloaded: true,
	}
}

// TestGetModels_MissingConfig tests the GetModels function.
func TestGetModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{getModel(0), getModel(1)}
	err := setupConfigFile(initialConfigFile, initialModels)
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

// TestGetModels_MissingConfig tests the GetModels function with a missing config file.
func TestGetModels_MissingConfig(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	var initialModels []model.Model
	err := setupConfigFile(initialConfigFile, initialModels)
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

// TestErrorOnAddModelWithEmptyViper tests the AddModels function
func TestAddModel(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{getModel(0), getModel(1)}

	// Call the AddModels function to add new models
	newModels := []model.Model{getModel(2), getModel(3)}
	err := setupConfigFile(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModels(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels, err := GetModels()
	test.AssertEqual(t, err, nil, "Error while getting updated models.")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestErrorOnAddModelWithEmptyViper tests the AddModels function with an empty config file
func TestAddModelOnEmptyConfFile(t *testing.T) {
	// Use the setup function
	var initialModels []model.Model
	confDir, initialConfigFile := setupConfigDir(t)

	// Call the AddModels function to add new models
	newModels := []model.Model{getModel(0), getModel(1)}

	err := setupConfigFile(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModels(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels, err := GetModels()
	test.AssertEqual(t, err, nil, "Error while getting updated models.")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestErrorOnAddModelWithEmptyViper tests the AddModels function with a missing config file
func TestErrorOnAddModelWithEmptyViper(t *testing.T) {
	viper.Reset()
	// Call the AddModels function to add new models
	newModels := []model.Model{getModel(0), getModel(1)}

	err := AddModels(newModels)
	test.AssertNotEqual(t, err, nil, "Should get error while updating configuration file.")
}

// TestRemoveModelPhysically_AddToBinaryFalse tests the RemoveModelPhysically with the property addToBinary to false.
func TestRemoveModelPhysically_AddToBinaryFalse(t *testing.T) {
	// Init
	modelToRemove := getModel(0)
	modelToRemove.ShouldBeDownloaded = false

	// Execute
	err := RemoveModelPhysically(modelToRemove.Name)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")
}

// TestRemoveModelPhysically_NonPhysical tests the RemoveModelPhysically with a non-physically present model.
func TestRemoveModelPhysically_NotPhysical(t *testing.T) {
	// Init
	modelToRemove := getModel(0)

	// Execute
	err := RemoveModelPhysically(modelToRemove.Name)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")
}

// TestRemoveModelPhysically_Success tests the RemoveModelPhysically with a physically present model.
func TestRemoveModelPhysically_Success(t *testing.T) {
	// Init
	modelToRemove := getModel(0)
	modelPath := filepath.Join(script.DownloadModelsPath, modelToRemove.Name)

	// Create temporary model
	setupModelDirectory(t, modelPath)
	defer os.RemoveAll(modelPath)

	// Execute
	err := RemoveModelPhysically(modelToRemove.Name)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")

	// Assert that the model was physically removed
	exists, err := utils.IsExistingPath(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "Model should have been removed.")
}

// TestRemoveAllModels_Success tests the RemoveAllModels function for successful removal of all models.
func TestRemoveAllModels_Success(t *testing.T) {
	// Init the models
	models := []model.Model{getModel(0), getModel(1), getModel(2)}

	// Create temporary models
	modelPath0 := filepath.Join(script.DownloadModelsPath, models[0].Name)
	setupModelDirectory(t, modelPath0)
	defer os.RemoveAll(modelPath0)
	modelPath1 := filepath.Join(script.DownloadModelsPath, models[1].Name)
	setupModelDirectory(t, modelPath1)
	defer os.RemoveAll(modelPath1)
	modelPath2 := filepath.Join(script.DownloadModelsPath, models[2].Name)
	setupModelDirectory(t, modelPath2)
	defer os.RemoveAll(modelPath2)

	// Setup configuration directory and file
	confDir, initialConfigFile := setupConfigDir(t)
	err := setupConfigFile(initialConfigFile, models)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")

	// Call the RemoveAllModels function
	err = RemoveAllModels()
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that all models were physically removed
	exists, err := utils.IsExistingPath(script.DownloadModelsPath)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "All models should have been removed.")

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

// TestRemoveModels_Success tests the RemoveModelsByNames function for successful removal of specified models.
func TestRemoveModels_Success(t *testing.T) {
	// Init the models
	models := []model.Model{getModel(0), getModel(1), getModel(2)}

	// Create temporary models
	modelPath0 := filepath.Join(script.DownloadModelsPath, models[0].Name)
	setupModelDirectory(t, modelPath0)
	modelPath1 := filepath.Join(script.DownloadModelsPath, models[1].Name)
	setupModelDirectory(t, modelPath1)
	modelPath2 := filepath.Join(script.DownloadModelsPath, models[2].Name)
	setupModelDirectory(t, modelPath2)
	defer os.RemoveAll(script.DownloadModelsPath)

	// Models to remove
	removeStartIndex := 1
	remainingModelsExpected := models[:removeStartIndex]
	var names []string
	for i := removeStartIndex; i < len(models); i++ {
		names = append(names, models[i].Name)
	}

	// Setup configuration directory and file
	confDir, initialConfigFile := setupConfigDir(t)
	err := setupConfigFile(initialConfigFile, models)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")

	// Call the RemoveModels function
	err = RemoveModelsByNames(models, names)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that all models were not physically removed
	exists, err := utils.IsExistingPath(script.DownloadModelsPath)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, true, exists, "All models should not have been removed.")

	// Assert that the request models were physically removed
	exists, err = utils.IsExistingPath(modelPath1)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "Model 1 should not have been removed.")
	exists, err = utils.IsExistingPath(modelPath2)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "Model 2 should not have been removed.")

	// Get the newly stored data
	var updatedModels []model.Model
	err = viper.UnmarshalKey("models", &updatedModels)
	test.AssertEqual(t, err, nil, "Error while unmarshalling models from configuration file.")

	// Assert that the models have been removed correctly
	test.AssertEqual(t, len(updatedModels), len(remainingModelsExpected), "The selected models were not removed correctly.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

// TestDownloadModel_BinaryFalse tests the DownloadModel function for skipping when download not needed.
func TestDownloadModel_BinaryFalse(t *testing.T) {
	// Init
	modelToDownload := getModel(0)
	modelToDownload.ShouldBeDownloaded = false

	// Execute
	result, ok := DownloadModel(modelToDownload)

	// Assert
	test.AssertEqual(t, ok, true)
	test.AssertEqual(t, result.ShouldBeDownloaded, false)
}

// TestDownloadModel_Fail tests the DownloadModel function where the script should fail.
func TestDownloadModel_Fail(t *testing.T) {
	// Init
	modelToDownload := getModel(0)

	// Execute
	result, ok := DownloadModel(modelToDownload)

	// Assert
	test.AssertEqual(t, ok, false)
	test.AssertEqual(t, result.ShouldBeDownloaded, false)
}

// TestDownloadModel_Fail tests the DownloadModel function where the script should succeed.
func TestDownloadModel_Success(t *testing.T) {
	t.Skip() // TODO : mock the script.DownloaderExecute call to return (sdm, nil)

	// Init
	modelToDownload := getModel(0)
	sdm := script.DownloaderModel{Module: "moduleAsATest"}

	// Execute
	result, ok := DownloadModel(modelToDownload)

	// Assert
	test.AssertEqual(t, ok, true)
	test.AssertEqual(t, result.ShouldBeDownloaded, true)
	test.AssertEqual(t, result.Config.Module, sdm.Module)
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
