package config

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/model"
	"os"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/viper"
)

func init() {
	app.Init("", "")
}

func setupConfigDirStrings(t *testing.T, initialModels []string) string {
	// Create a temporary directory for the test
	confDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}

	// Set up a temporary config file with some initial models
	initialConfigFile := filepath.Join(confDir, "config.yaml")

	err = createConfigFileStrings(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	return confDir
}

func createConfigFileStrings(filePath string, models []string) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	if len(models) > 0 {
		// Write models to the config file
		_, err = file.WriteString("models:\n")
		for _, model := range models {
			_, err := file.WriteString("  - " + model + "\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAddModel(t *testing.T) {
	// Use the setup function
	initialModels := []string{"model1", "model2"}
	confDir := setupConfigDirStrings(t, initialModels)

	// Call the AddModel function to add new models
	newModels := []string{"model3", "model4"}
	err := Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModel(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels := viper.GetStringSlice("models")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")
	for i := 0; i < len(initialModels); i++ {
		test.AssertEqual(t, updatedModels[i], expectedModels[i], "Models not updated as expected.")
	}

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

func TestAddModelOnEmptyConfFile(t *testing.T) {
	// Use the setup function
	initialModels := []string{}
	confDir := setupConfigDirStrings(t, initialModels)

	// Call the AddModel function to add new models
	newModels := []string{"model1", "model2"}
	err := Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")
	err = AddModel(newModels)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Assert that the models have been updated correctly
	updatedModels := viper.GetStringSlice("models")
	expectedModels := append(initialModels, newModels...)
	test.AssertEqual(t, len(updatedModels), len(expectedModels), "Models list length is not as expected.")
	for i := 0; i < len(initialModels); i++ {
		test.AssertEqual(t, updatedModels[i], expectedModels[i], "Models not updated as expected.")
	}

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

func TestErrorOnAddModelWithEmptyViper(t *testing.T) {
	// Use the setup function
	initialModels := []string{}
	confDir := setupConfigDirStrings(t, initialModels)

	// Call the AddModel function to add new models
	newModels := []string{"model1", "model2"}
	err := AddModel(newModels)
	test.AssertNotEqual(t, err, nil, "Should get error while updating configuration file.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}

func TestGetModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{Name: "model1", PipeLine: "pipeline1", DirectoryPath: "/path/to/model1", AddToBinary: true},
		{Name: "model2", PipeLine: "pipeline2", DirectoryPath: "/path/to/model2", AddToBinary: false},
	}
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

// TestIsModelsEmpty_EmptyModels tests the IsModelsEmpty function with an empty models slice.
func TestIsModelsEmpty_EmptyModels(t *testing.T) {
	// Init
	var models []model.Model

	// Execute
	isEmpty := IsModelsEmpty(models)

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestIsModelsEmpty_NonEmptyModels tests the IsModelsEmpty function with a non-empty models slice.
func TestIsModelsEmpty_NonEmptyModels(t *testing.T) {
	// Init
	models := []model.Model{
		{Name: "Model1", PipeLine: "Pipeline1", DirectoryPath: "/path/to/directory1", AddToBinary: true},
		{Name: "Model2", PipeLine: "Pipeline2", DirectoryPath: "/path/to/directory2", AddToBinary: false},
	}

	// Execute
	isEmpty := IsModelsEmpty(models)

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestRemoveModels_Success tests the RemoveModels function for successful removal of specified models.
func TestRemoveModels_Success(t *testing.T) {
	// Setup directory
	confDir, initialConfigFile := setupConfigDir(t)

	// Setup file
	initialModels := []model.Model{
		{Name: "model1", PipeLine: "pipeline1", DirectoryPath: "/path/to/model1", AddToBinary: true},
		{Name: "model2", PipeLine: "pipeline2", DirectoryPath: "/path/to/model2", AddToBinary: false},
	}
	err := setupConfigFile(initialConfigFile, initialModels)
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
		{Name: "model1", PipeLine: "pipeline1", DirectoryPath: "/path/to/model1", AddToBinary: true},
		{Name: "model2", PipeLine: "pipeline2", DirectoryPath: "/path/to/model2", AddToBinary: false},
	}
	err := setupConfigFile(initialConfigFile, initialModels)
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

// setupConfigFile creates a configuration file with specified models.
func setupConfigFile(filePath string, models []model.Model) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		return err
	}

	if len(models) > 0 {
		// Write models to the config file
		_, err = file.WriteString("models:\n")
		for _, m := range models {
			_, err := file.WriteString(fmt.Sprintf("  - name: %s\n    pipeline: %s\n    directorypath: %s\n    addtobinary: %t\n",
				m.Name, m.PipeLine, m.DirectoryPath, m.AddToBinary))
			if err != nil {
				return err
			}
		}
	} else {
		_, err = file.WriteString("other:\n")
		_, _ = file.WriteString(fmt.Sprintf("  - name: error\n    pipeline: error\n    directorypath: error\n"))
		_, _ = file.WriteString(fmt.Sprintf("  - name: error\n    pipeline: error\n    directorypath: error\n"))
		_, _ = file.WriteString(fmt.Sprintf("  - name: error\n    pipeline: error\n    directorypath: error\n"))
	}

	return nil
}

// cleanConfDir removes the temporary directory created during testing.
func cleanConfDir(t *testing.T, confDir string) {
	if err := os.RemoveAll(confDir); err != nil {
		t.Errorf("Error cleaning up temporary directory: %v", err)
	}
}
