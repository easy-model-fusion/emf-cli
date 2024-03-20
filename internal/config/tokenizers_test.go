package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/viper"
	"os"
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

// TestRemoveTokenizerPhysically_AddToBinaryFalse tests the RemoveTokenizerPhysically with the property addToBinary to false.
func TestRemoveTokenizerPhysically_AddToBinaryFalse(t *testing.T) {
	// Init
	modelToRemove := getModel(0)
	modelToRemove.AddToBinaryFile = false

	// Execute
	err := RemoveTokenizerPhysically(modelToRemove.Name)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")
}

// TestRemoveTokenizerPhysically_NotPhysical tests the RemoveTokenizerPhysically with a non-physically present tokenizer.
func TestRemoveTokenizerPhysically_NotPhysical(t *testing.T) {
	// Init
	modelToRemove := getModel(0)

	// Execute
	err := RemoveTokenizerPhysically(modelToRemove.Name)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")
}

// TestRemoveTokenizerPhysically_Success tests the RemoveTokenizerPhysically with a physically present tokenizer.
func TestRemoveTokenizerPhysically_Success(t *testing.T) {
	// Init
	modelToUse := getModelWithTokenizer(0)

	configTokenizerMap := modelToUse.Tokenizers.Map()

	tokenizer, _ := configTokenizerMap["tokenizer0"]

	// Create temporary tokenizer
	setupModelDirectory(t, tokenizer.Path)
	defer os.RemoveAll(tokenizer.Path)

	// Execute
	err := RemoveTokenizerPhysically(tokenizer.Path)
	test.AssertEqual(t, nil, err, "Removal should not have failed since it's not physically downloaded.")

	// Assert that the model was physically removed
	exists, err := fileutil.IsExistingPath(tokenizer.Path)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "Model should have been removed.")
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
