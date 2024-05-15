package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"path/filepath"
	"testing"
)

// GetModel initiates a basic model with an id as suffix
func GetModel(id int) Model {
	idStr := fmt.Sprint(id)
	return Model{
		Name:            "model" + idStr,
		Module:          huggingface.Module("module" + idStr),
		Class:           "class" + idStr,
		Source:          HUGGING_FACE,
		AddToBinaryFile: true,
		IsDownloaded:    true,
	}
}

// GetModels initiates a list of basic models starting with id 0
func GetModels(length int) Models {
	var models Models
	for i := 1; i <= length; i++ {
		models = append(models, GetModel(i-1))
	}
	return models
}

// GetTokenizer initiates a basic tokenizer with an id as suffix
func GetTokenizer(id int) Tokenizer {
	idStr := fmt.Sprint(id)
	return Tokenizer{
		Class: "tokenizer" + idStr,
		Path:  "path" + idStr,
	}
}

// GetTokenizers initiates a list of basic tokenizers starting with id 0
func GetTokenizers(length int) Tokenizers {
	var tokenizers Tokenizers
	for i := 1; i <= length; i++ {
		tokenizers = append(tokenizers, GetTokenizer(i-1))
	}
	return tokenizers
}

// TestEmpty_True tests the Models.Empty function with an empty models slice.
func TestEmpty_True(t *testing.T) {
	// Init
	var models Models

	// Execute
	isEmpty := models.Empty()

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestEmpty_False tests the Models.Empty function with a non-empty models slice.
func TestEmpty_False(t *testing.T) {
	// Init
	models := GetModels(1)

	// Execute
	isEmpty := models.Empty()

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestContainsByName_True tests the Models.ContainsByName function with an element's name contained by the slice.
func TestContainsByName_True(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	contains := models.ContainsByName(models[0].Name)

	// Assert
	test.AssertEqual(t, contains, true, "Expected true.")
}

// TestContainsByName_False tests the Models.ContainsByName function with an element's name not contained by the slice.
func TestContainsByName_False(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	contains := models.ContainsByName(GetModel(3).Name)

	// Assert
	test.AssertEqual(t, contains, false, "Expected false.")
}

// TestDifference tests the Models.Difference function to return the correct difference.
func TestDifference(t *testing.T) {
	// Init
	models := GetModels(5)
	index := 2
	sub := models[:index]
	expected := models[index:]

	// Execute
	difference := models.Difference(sub)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

// TestUnion tests the Models.Union function to return the correct union.
func TestUnion(t *testing.T) {
	// Init
	index := 2
	models1 := GetModels(5)
	models2 := models1[:index]
	expected := models2

	// Execute
	union := models1.Union(models2)

	// Assert
	test.AssertEqual(t, len(expected), len(union), "Lengths should be equal.")
}

// TestMap_Success tests the Models.Map function to return a map from a slice of models.
func TestMap_Success(t *testing.T) {
	// Init
	models := GetModels(3)
	expected := map[string]Model{
		models[0].Name: models[0],
		models[1].Name: models[1],
		models[2].Name: models[2],
	}

	// Execute
	result := models.Map()

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestTokenizer_Map_Success tests the Tokenizers.Map function to return a map from a slice of tokenizers.
func TestTokenizer_Map_Success(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Tokenizers = GetTokenizers(3)
	expected := map[string]Tokenizer{
		model.Tokenizers[0].Class: model.Tokenizers[0],
		model.Tokenizers[1].Class: model.Tokenizers[1],
		model.Tokenizers[2].Class: model.Tokenizers[2],
	}

	// Execute
	result := model.Tokenizers.Map()

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestGetNames_Success tests the Models.GetNames function to return the correct model names.
func TestGetNames_Success(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	names := models.GetNames()

	// Assert
	test.AssertEqual(t, len(models), len(names), "Lengths should be equal.")
}

// TestTokenizer_GetNames_Success tests the Tokenizers.GetNames function to return the correct names.
func TestTokenizer_GetNames_Success(t *testing.T) {
	// Init
	input := GetModel(0)
	input.Tokenizers = Tokenizers{{Class: "tokenizer1"}, {Class: "tokenizer2"}, {Class: "tokenizer3"}}
	expected := []string{
		input.Tokenizers[0].Class,
		input.Tokenizers[1].Class,
		input.Tokenizers[2].Class,
	}

	// Execute
	names := input.Tokenizers.GetNames()

	// Assert
	test.AssertEqual(t, len(expected), len(names), "Lengths should be equal.")
}

// TestFilterWithNames_Success tests the Models.FilterWithNames function to return the correct models.
func TestFilterWithNames_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	names := []string{models[0].Name, models[1].Name}

	// Execute
	result := models.FilterWithNames(names)

	// Assert
	test.AssertEqual(t, len(models), len(result), "Lengths should be equal.")
}

// TestFilterWithSourceHuggingface_Success tests the Models.FilterWithSourceHuggingface to return the sub-slice.
func TestFilterWithSourceHuggingface_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].Source = ""
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithSourceHuggingface()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestFilterWithIsDownloadedTrue_Success tests the Models.FilterWithIsDownloadedTrue to return the sub-slice.
func TestFilterWithIsDownloadedTrue_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].IsDownloaded = false
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithIsDownloadedTrue()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestFilterWithIsDownloadedTrue_Success tests the Models.FilterWithIsDownloadedTrue to return the sub-slice.
func TestFilterWithIsDownloadedOrAddToBinaryTrue_Success(t *testing.T) {
	// Init
	models := GetModels(4)
	models[0].IsDownloaded = false
	models[0].AddToBinaryFile = false
	models[2].IsDownloaded = false
	models[3].AddToBinaryFile = false
	expected := Models{models[1], models[2], models[3]}

	// Execute
	result := models.FilterWithIsDownloadedOrAddToBinaryFileTrue()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
	for i, currentModel := range expected {
		test.AssertEqual(t, result[i].Name, currentModel.Name, "returned models should be equal to expected model.")
	}
}

// TestFilterWithAddToBinaryFileTrue_Success tests the Models.FilterWithAddToBinaryFileTrue to return the sub-slice.
func TestFilterWithAddToBinaryFileTrue_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].AddToBinaryFile = false
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithAddToBinaryFileTrue()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestGetBasePath tests the GetBasePath to return the correct base path to the model.
func TestGetBasePath(t *testing.T) {
	// Init
	modelName := "name"
	model := Model{Name: modelName}

	// Execute
	basePath := model.GetBasePath()

	// Assert
	test.AssertEqual(t, basePath, filepath.Join(app.DownloadDirectoryPath, model.Name))
}

// TestUpdatePaths_Default tests the Model.UpdatePaths for a default model.
func TestUpdatePaths_Default(t *testing.T) {
	// Init
	model := GetModel(0)

	// Execute
	model.UpdatePaths()

	// Assert
	test.AssertEqual(t, model.Path, filepath.Join(app.DownloadDirectoryPath, model.Name))
}

// TestUpdatePaths_Transformers tests the Model.UpdatePaths for a transformers model.
func TestUpdatePaths_Transformers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS

	// Execute
	model.UpdatePaths()

	// Assert
	test.AssertEqual(t, model.Path, filepath.Join(app.DownloadDirectoryPath, model.Name, "model"))
}

// TestUpdatePaths_TransformersTokenizers tests the Model.UpdatePaths for a transformers model.
func TestUpdatePaths_TransformersTokenizers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = []Tokenizer{{Class: "tokenizer"}}

	// Execute
	model.UpdatePaths()

	// Assert
	test.AssertEqual(t, model.Tokenizers[0].Path, filepath.Join(app.DownloadDirectoryPath, model.Name, "tokenizer"))
}

// TestFilterWithClass_Success tests the Tokenizers.FilterWithClass function to return the correct models.
func TestFilterWithClass_Success(t *testing.T) {
	//Init
	tokenizers := GetTokenizers(2)
	names := []string{tokenizers[0].Class, tokenizers[1].Class}

	// Execute
	result := tokenizers.FilterWithClass(names)

	// Assert
	test.AssertEqual(t, len(tokenizers), len(result), "Lengths should be equal.")

}

func TestTokenizers_Difference(t *testing.T) {
	// Init
	tokenizers := GetTokenizers(3)
	sub := tokenizers[:2]
	expected := tokenizers[2:]

	// Execute
	difference := tokenizers.Difference(sub)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

func TestSetAccessTokenKey(t *testing.T) {
	//Init
	model := GetModel(1)
	model.Name = "1model/name1.6-test_escape"
	expected := "ACCESS_TOKEN_1MODEL_NAME1_6_TEST_ESCAPE"

	// Set access token
	err := model.setAccessTokenKey()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, expected, model.AccessToken)
}

func TestModelSaveAndGetAccessToken(t *testing.T) {
	//Init
	model := GetModel(1)
	model.Name = "1model/name1.6-test_escape"

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Set access token
	err := model.SaveAccessToken("token")
	test.AssertEqual(t, err, nil)
	savedToken, err := model.GetAccessToken()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, savedToken, "token")
}
