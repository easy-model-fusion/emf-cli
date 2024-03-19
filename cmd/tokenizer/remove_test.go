package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
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
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})
	// Initialize selected models list
	var args []string
	args = append(args, "model1")
	args = append(args, "tokenizer1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	runTokenizerRemove(nil, args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 0, "tokenizer deleted")
}
