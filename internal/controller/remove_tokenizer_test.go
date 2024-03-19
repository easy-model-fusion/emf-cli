package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"

	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

func TestRemoveTokenizer(t *testing.T) {
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
	TokenizerRemoveCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 0, "Only one model should be left.")
}

func TestRemoveTokenizer_withNoModule(t *testing.T) {
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

	// Process remove
	TokenizerRemoveCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 1, "Only one model should be left.")
}

func TestRemoveTokenizer_withWrongModule(t *testing.T) {
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

	// Process remove
	TokenizerRemoveCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

func TestRemoveTokenizer_withNoArgs(t *testing.T) {
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
	TokenizerRemoveCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

func TestRemoveTokenizer_withNoTokenizerArgs(t *testing.T) {
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
	ui := test.MockUI{MultiselectResult: expectedSelections}
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

	// Process remove
	TokenizerRemoveCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}
