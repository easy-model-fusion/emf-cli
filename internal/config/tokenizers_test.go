package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	app.Init("", "")
}

// getModel initiates a basic model with an id as suffix
func getModelWithTokenizer(suffix int) model.Model {
	idStr := fmt.Sprint(suffix)
	return model.Model{
		Name:   "model" + idStr,
		Module: huggingface.TRANSFORMERS,
		Class:  "class" + idStr,
		Tokenizers: model.Tokenizers{
			{Path: "test/path" + idStr, Class: "tokenizer" + idStr, Options: map[string]string{"option1": "value1"}},
		},
		AddToBinaryFile: true,
	}
}

// TestRemoveTokenizer_Success tests the RemoveTokenizersByName function for successful removal of specified tokenizer.
func TestRemoveTokenizer_Success(t *testing.T) {
	// Init the models
	models := []model.Model{getModelWithTokenizer(0)}

	modelToUse := models[0]

	configTokenizerMap := modelToUse.Tokenizers.Map()
	var tokenizersToRemove model.Tokenizers

	tokenizer, _ := configTokenizerMap["tokenizer0"]
	tokenizersToRemove = append(tokenizersToRemove, tokenizer)

	// Setup configuration directory and file
	confDir, initialConfigFile := setupConfigDir(t)
	err := setupConfigFile(initialConfigFile, models, false)
	test.AssertEqual(t, err, nil, "Error while creating temporary configuration file.")
	err = Load(confDir)
	test.AssertEqual(t, err, nil, "Error while loading configuration file.")

	// Call the RemoveModels function
	_, err = RemoveTokenizersByName(modelToUse, tokenizersToRemove)
	test.AssertEqual(t, err, nil, "Error while updating configuration file.")

	// Get the newly stored data
	var updatedModels []model.Model
	err = viper.UnmarshalKey("models", &updatedModels)
	test.AssertEqual(t, err, nil, "Error while unmarshalling models from configuration file.")

	// Assert that the tokenizer have been removed correctly
	test.AssertEqual(t, len(updatedModels[0].Tokenizers), 0, "The selected tokenizer were not removed correctly.")

	// Clean up directory afterward
	cleanConfDir(t, confDir)
}
