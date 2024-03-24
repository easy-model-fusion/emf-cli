package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/spf13/viper"
	"testing"

	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

// Sets the configuration file with the given models
func setupConfigFile(models model.Models) error {
	config.FilePath = "."
	// Load configuration file
	err := config.GetViperConfig(".")
	if err != nil {
		return err
	}
	// Write models to the config file
	viper.Set("models", models)
	return config.WriteViperConfig()
}

// TestRemoveTokenizer_Success tests the RunTokenizerRemove function
func TestRemoveTokenizer_Success(t *testing.T) {
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

	// Create ui mock
	ui := mock.MockUI{UserConfirmationResult: true}
	app.SetUI(ui)
	ic := RemoveTokenizerController{}
	// Process remove
	if err := ic.RunTokenizerRemove(args); err != nil {
		test.AssertEqual(t, err, nil, "Error on update")
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 0, "Only one model should be left.")
}

// TestRemoveTokenizer_WithModuleNotTransformers tests the RunTokenizerRemove function with no transformers module
func TestRemoveTokenizer_WithModuleNotTransformers(t *testing.T) {
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

	ic := RemoveTokenizerController{}
	// Process remove
	if err := ic.RunTokenizerRemove(args); err != nil {
		expectedErrMsg := "only transformers models have tokenizers"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 1, "Only one model should be left.")
}

// TestRemoveTokenizer_WithWrongModel tests the RunTokenizerRemove function with wrong model
func TestRemoveTokenizer_WithWrongModel(t *testing.T) {
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
	ic := RemoveTokenizerController{}
	// Process remove
	if err := ic.RunTokenizerRemove(args); err != nil {
		test.AssertEqual(t, err, nil, "Error on update")
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

// TestRemoveTokenizer_WithNoTokenizerArgs_Success tests the RunTokenizerRemove function with no tokenizers args
func TestRemoveTokenizer_WithNoTokenizerArgs_Success(t *testing.T) {
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
	ic := RemoveTokenizerController{}
	// Process remove
	if err := ic.RunTokenizerRemove(args); err != nil {
		test.AssertEqual(t, err, nil, "Error on update")
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

// TestRemoveTokenizer_WithWrongTokenizerArgs tests the RunTokenizerRemove function with wrong tokenizers args
func TestRemoveTokenizer_WithWrongTokenizerArgs(t *testing.T) {
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
	args = append(args, "X")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := RemoveTokenizerController{}
	// Process remove
	if err := ic.RunTokenizerRemove(args); err != nil {
		test.AssertEqual(t, err, nil, "Error on update")
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}
