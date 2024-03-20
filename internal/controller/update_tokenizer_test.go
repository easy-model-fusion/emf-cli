package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/spf13/viper"
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

func TestRunTokenizerUpdate(t *testing.T) {
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

	// Process update
	TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing update")
	_, err = config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")
}

func TestNoModuleTransformersUpdate(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name: "model1",
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

	// Process update
	TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing Update")
	_, err = config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")
}

func TestWrongModelNameUpdate(t *testing.T) {
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
	args = append(args, "modelX")
	args = append(args, "tokenizer1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process update
	TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing update")
}

func TestNoArgsUpdate(t *testing.T) {
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

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process remove
	TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

func TestNoTokenizersArgsUpdate(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	var expectedSelections []string
	expectedSelections = append(expectedSelections, "tokenizer1")

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Initialize selected models list
	var args []string
	args = append(args, "model1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	// Process update
	TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}
