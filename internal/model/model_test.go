package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/pterm/pterm"
	"os"
	"testing"
)

// GetModel initiates a basic model with an id as suffix
func GetModel(suffix int) Model {
	idStr := fmt.Sprint(suffix)
	return Model{
		Name:            "model" + idStr,
		Module:          huggingface.Module("module" + idStr),
		Class:           "class" + idStr,
		Source:          HUGGING_FACE,
		AddToBinaryFile: true,
		IsDownloaded:    true,
	}
}

func GetModels(length int) Models {
	var models Models
	for i := 1; i <= length; i++ {
		models = append(models, GetModel(i-1))
	}
	return models
}

// TestEmpty_True tests the Empty function with an empty models slice.
func TestEmpty_True(t *testing.T) {
	// Init
	var models Models

	// Execute
	isEmpty := models.Empty()

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestEmpty_False tests the Empty function with a non-empty models slice.
func TestEmpty_False(t *testing.T) {
	// Init
	models := GetModels(1)

	// Execute
	isEmpty := models.Empty()

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestContainsByName_True tests the ContainsByName function with an element's name contained by the slice.
func TestContainsByName_True(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	contains := models.ContainsByName(models[0].Name)

	// Assert
	test.AssertEqual(t, contains, true, "Expected true.")
}

// TestContainsByName_False tests the ContainsByName function with an element's name not contained by the slice.
func TestContainsByName_False(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	contains := models.ContainsByName(GetModel(3).Name)

	// Assert
	test.AssertEqual(t, contains, false, "Expected false.")
}

// TestDifference tests the Difference function to return the correct difference.
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

// TestUnion tests the Union function to return the correct union.
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

// TestModelsToMap_Success tests the ModelsToMap function to return a map from a slice of models.
func TestModelsToMap_Success(t *testing.T) {
	// Init
	models := GetModels(3)
	expected := map[string]Model{
		models[0].Name: models[0],
		models[1].Name: models[1],
		models[2].Name: models[2],
	}

	// Execute
	result := models.ToMap()

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestGetNames_Success tests the GetNames function to return the correct model names.
func TestGetNames_Success(t *testing.T) {
	// Init
	models := GetModels(2)

	// Execute
	names := models.GetNames()

	// Assert
	test.AssertEqual(t, len(models), len(names), "Lengths should be equal.")
}

// TestGetModelsByNames tests the GetModelsByNames function to return the correct models.
func TestGetModelsByNames(t *testing.T) {
	// Init
	models := GetModels(2)
	names := []string{models[0].Name, models[1].Name}

	// Execute
	result := models.GetByNames(names)

	// Assert
	test.AssertEqual(t, len(models), len(result), "Lengths should be equal.")
}

// TestGetModelsWithSourceHuggingface_Success tests the GetModelsWithSourceHuggingface to return the sub-slice.
func TestGetModelsWithSourceHuggingface_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].Source = ""
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithSourceHuggingface()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestGetModelsWithIsDownloadedTrue_Success tests the GetModelsWithIsDownloadedTrue to return the sub-slice.
func TestGetModelsWithIsDownloadedTrue_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].IsDownloaded = false
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithIsDownloadedTrue()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestGetModelsWithAddToBinaryFileTrue_Success tests the GetModelsWithAddToBinaryFileTrue to return the sub-slice.
func TestGetModelsWithAddToBinaryFileTrue_Success(t *testing.T) {
	// Init
	models := GetModels(2)
	models[0].AddToBinaryFile = false
	expected := Models{models[1]}

	// Execute
	result := models.FilterWithAddToBinaryFileTrue()

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestMapToModelFromHuggingfaceModel_Success tests the MapToModelFromHuggingfaceModel to return the correct Model.
func TestMapToModelFromHuggingfaceModel_Success(t *testing.T) {
	// Init
	huggingfaceModel := huggingface.Model{
		Name:        "name",
		PipelineTag: "pipeline",
		LibraryName: "library",
	}

	// Execute
	model := FromHuggingfaceModel(huggingfaceModel)

	pterm.Info.Println(model.Module)
	// Assert
	test.AssertEqual(t, model.Name, huggingfaceModel.Name)
	test.AssertEqual(t, model.PipelineTag, huggingfaceModel.PipelineTag)
	test.AssertEqual(t, model.Module, huggingfaceModel.LibraryName)
	test.AssertEqual(t, model.Source, HUGGING_FACE)
}

// TestModelDownloadedOnDevice_FalseMissing tests the ModelDownloadedOnDevice function to return false upon missing.
func TestModelDownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Path = ""

	// Execute
	exists, err := model.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestModelDownloadedOnDevice_FalseEmpty tests the ModelDownloadedOnDevice function to return false upon empty.
func TestModelDownloadedOnDevice_FalseEmpty(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelDirectory, err := os.MkdirTemp("", "modelDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(modelDirectory)

	// Init
	model := GetModel(0)
	model.Path = modelDirectory

	// Execute
	exists, err := model.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestModelDownloadedOnDevice_True tests the ModelDownloadedOnDevice function to return true.
func TestModelDownloadedOnDevice_True(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelDirectory, err := os.MkdirTemp("", "modelDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(modelDirectory)

	// Create temporary file inside the model base path
	file, err := os.CreateTemp(modelDirectory, "")
	if err != nil {
		t.Fatal(err)
	}
	fileutil.CloseFile(file)

	// Init
	model := GetModel(0)
	model.Path = modelDirectory

	// Execute
	exists, err := model.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}
