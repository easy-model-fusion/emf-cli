package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.Init("", "")
	app.InitGit("", "")
	os.Exit(m.Run())
}

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

func Test_runModelRemove(t *testing.T) {
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
	runModelRemove(nil, args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels), 1, "Only one model should be left.")
	test.AssertEqual(t, newModels[0].Name, "model1", "Model1 shouldn't be deleted")

}
