package config

import (
	"github.com/easy-model-fusion/client/internal/app"
	"os"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/viper"
)

func init() {
	app.Init("", "")
}

func TestAddModel(t *testing.T) {
	// Create a temporary directory for the test
	confDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(confDir)

	// Set up a temporary config file with some initial models
	initialModels := []string{"model1", "model2"}
	initialConfigFile := filepath.Join(confDir, "config.yaml")

	err = createConfigFile(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the AddModel function to add new models
	newModels := []string{"model3", "model4"}
	err = Load(confDir)
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
}

func TestAddModelOnEmptyConfFile(t *testing.T) {
	// Create a temporary directory for the test
	confDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(confDir)

	// Set up a temporary config file with some initial models
	initialModels := []string{}
	initialConfigFile := filepath.Join(confDir, "config.yaml")

	err = createConfigFile(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the AddModel function to add new models
	newModels := []string{"model1", "model2"}
	err = Load(confDir)
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
}

func TestErrorOnAddModelWithEmptyViper(t *testing.T) {
	// Create a temporary directory for the test
	confDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(confDir)

	// Set up a temporary config file with some initial models
	initialModels := []string{}
	initialConfigFile := filepath.Join(confDir, "config.yaml")

	err = createConfigFile(initialConfigFile, initialModels)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")

	// Call the AddModel function to add new models
	newModels := []string{"model1", "model2"}
	err = AddModel(newModels)
	test.AssertNotEqual(t, err, nil, "Should get error while updating configuration file.")
}

func createConfigFile(filePath string, models []string) error {
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
