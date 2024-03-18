package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestTokenizerUpdateCmd tests the TokenizerUpdateCmd function
func TestTokenizerUpdateCmd(t *testing.T) {
	// Test with valid arguments

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
	var tokenizer model.Tokenizers
	tokenizer = append(tokenizer, model.Tokenizer{
		Path:    "path/to/tokenizer1",
		Class:   "tokenizer1",
		Options: nil,
	})

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	t.Run("ValidArguments", func(t *testing.T) {
		args := []string{"model1", "tokenizer1"}
		TokenizerUpdateCmd(args)
		// Assert the output here based on your implementation
		test.AssertEqual(t, true, "Tokenizer updated.")
	})

	// Test with missing arguments
	t.Run("MissingArguments", func(t *testing.T) {
		var args []string
		TokenizerUpdateCmd(args)
		// Assert the output here based on your implementation
		test.AssertEqual(t, err, false, "Tokenizer not updated.")
	})

	// Test with non-existent model
	t.Run("NonExistentModel", func(t *testing.T) {
		args := []string{"nonexistent_model"}
		TokenizerUpdateCmd(args)
		// Assert the output here based on your implementation
		test.AssertEqual(t, err, false, "Tokenizer not updated.")
	})

}
