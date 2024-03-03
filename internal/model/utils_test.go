package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"path"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/emf-cli/test"
)

// getModel initiates a basic model with an id as suffix
func getModel(suffix int) Model {
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

// TestEmpty_True tests the Empty function with an empty models slice.
func TestEmpty_True(t *testing.T) {
	// Init
	var models []Model

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, true, "Expected true.")
}

// TestEmpty_False tests the Empty function with a non-empty models slice.
func TestEmpty_False(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	isEmpty := Empty(models)

	// Assert
	test.AssertEqual(t, isEmpty, false, "Expected false.")
}

// TestContainsByName_True tests the ContainsByName function with an element's name contained by the slice.
func TestContainsByName_True(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := ContainsByName(models, models[0].Name)

	// Assert
	test.AssertEqual(t, contains, true, "Expected true.")
}

// TestContainsByName_False tests the ContainsByName function with an element's name not contained by the slice.
func TestContainsByName_False(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	contains := ContainsByName(models, getModel(2).Name)

	// Assert
	test.AssertEqual(t, contains, false, "Expected false.")
}

// TestDifference tests the Difference function to return the correct difference.
func TestDifference(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1), getModel(2), getModel(3), getModel(4)}
	index := 2
	sub := models[:index]
	expected := models[index:]

	// Execute
	difference := Difference(models, sub)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

// TestUnion tests the Union function to return the correct union.
func TestUnion(t *testing.T) {
	// Init
	index := 2
	models1 := []Model{getModel(0), getModel(1), getModel(2), getModel(3), getModel(4)}
	models2 := models1[:index]
	expected := models2

	// Execute
	union := Union(models1, models2)

	// Assert
	test.AssertEqual(t, len(expected), len(union), "Lengths should be equal.")
}

// TestGetModelsByNames tests the GetModelsByNames function to return the correct models.
func TestGetModelsByNames(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}
	names := []string{models[0].Name, models[1].Name}

	// Execute
	result := GetModelsByNames(models, names)

	// Assert
	test.AssertEqual(t, len(models), len(result), "Lengths should be equal.")
}

// TestGetNames tests the GetNames function to return the correct names.
func TestGetNames(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}

	// Execute
	names := GetNames(models)

	// Assert
	test.AssertEqual(t, len(models), len(names), "Lengths should be equal.")
}

// TestConstructConfigPaths_Default tests the ConstructConfigPaths for a default model.
func TestConstructConfigPaths_Default(t *testing.T) {
	// Init
	model := getModel(0)

	// Execute
	model = ConstructConfigPaths(model)

	// Assert
	test.AssertEqual(t, model.Path, path.Join(app.DownloadDirectoryPath, model.Name))
}

// TestConstructConfigPaths_Transformers tests the ConstructConfigPaths for a transformers model.
func TestConstructConfigPaths_Transformers(t *testing.T) {
	// Init
	model := getModel(0)
	model.Module = huggingface.TRANSFORMERS

	// Execute
	model = ConstructConfigPaths(model)

	// Assert
	test.AssertEqual(t, model.Path, path.Join(app.DownloadDirectoryPath, model.Name, "model"))
}

// TestConstructConfigPaths_Transformers tests the ConstructConfigPaths for a transformers model.
func TestConstructConfigPaths_TransformersTokenizers(t *testing.T) {
	// Init
	model := getModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = []Tokenizer{{Class: "tokenizer"}}

	// Execute
	model = ConstructConfigPaths(model)

	// Assert
	test.AssertEqual(t, model.Tokenizers[0].Path, path.Join(app.DownloadDirectoryPath, model.Name, "tokenizer"))
}

// TestMapToModelFromDownloaderModel_Empty tests the MapToModelFromDownloaderModel to return the correct Model.
func TestMapToModelFromDownloaderModel_Empty(t *testing.T) {
	// Init
	downloaderModel := downloader.Model{
		Path:   "",
		Module: "",
		Class:  "",
		Tokenizer: downloader.Tokenizer{
			Path:  "",
			Class: "",
		},
	}
	expected := Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Tokenizers: []Tokenizer{
			{Path: "/path/to/tokenizer", Class: "tokenizer_class"},
		},
	}

	// Execute
	result := MapToModelFromDownloaderModel(expected, downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
}

// TestMapToModelFromDownloaderModel_Fill tests the MapToModelFromDownloaderModel to return the correct Config.
func TestMapToModelFromDownloaderModel_Fill(t *testing.T) {
	// Init
	downloaderModel := downloader.Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Tokenizer: downloader.Tokenizer{
			Path:  "/path/to/tokenizer",
			Class: "tokenizer_class",
		},
	}
	expected := Model{
		Path:   filepath.Clean("/path/to/model"),
		Module: "module_name",
		Class:  "class_name",
		Tokenizers: []Tokenizer{
			{Path: filepath.Clean("/path/to/tokenizer"), Class: "tokenizer_class"},
		},
	}

	// Execute
	result := MapToModelFromDownloaderModel(Model{}, downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
}

// TestMapToTokenizerFromDownloaderTokenizer_Success tests the MapToTokenizerFromDownloaderTokenizer to return the correct Tokenizer.
func TestMapToTokenizerFromDownloaderTokenizer_Success(t *testing.T) {
	// Init
	downloaderTokenizer := downloader.Tokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}
	expected := Tokenizer{Path: filepath.Clean("/path/to/tokenizer"), Class: "tokenizer_class"}

	// Execute
	result := MapToTokenizerFromDownloaderTokenizer(downloaderTokenizer)

	// Assert
	test.AssertEqual(t, expected, result)
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
	model := MapToModelFromHuggingfaceModel(huggingfaceModel)

	pterm.Info.Println(model.Module)
	// Assert
	test.AssertEqual(t, model.Name, huggingfaceModel.Name)
	test.AssertEqual(t, model.PipelineTag, huggingfaceModel.PipelineTag)
	test.AssertEqual(t, model.Module, huggingfaceModel.LibraryName)
	test.AssertEqual(t, model.Source, HUGGING_FACE)
}

// TestGetModelsWithSourceHuggingface_Success tests the GetModelsWithSourceHuggingface to return the sub-slice.
func TestGetModelsWithSourceHuggingface_Success(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}
	models[0].Source = ""
	expected := []Model{models[1]}

	// Execute
	result := GetModelsWithSourceHuggingface(models)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestGetModelsWithIsDownloadedTrue_Success tests the GetModelsWithIsDownloadedTrue to return the sub-slice.
func TestGetModelsWithIsDownloadedTrue_Success(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}
	models[0].IsDownloaded = false
	expected := []Model{models[1]}

	// Execute
	result := GetModelsWithIsDownloadedTrue(models)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
}

// TestModelsToMap_Success tests the ModelsToMap function to return a map from a slice of models.
func TestModelsToMap_Success(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1), getModel(2)}
	expected := map[string]Model{
		models[0].Name: models[0],
		models[1].Name: models[1],
		models[2].Name: models[2],
	}

	// Execute
	result := ModelsToMap(models)

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestGetTokenizerNames_Success tests the GetNames function to return the correct names.
func TestGetTokenizerNames_Success(t *testing.T) {
	// Init
	input := getModel(0)
	input.Tokenizers = []Tokenizer{{Class: "tokenizer1"}, {Class: "tokenizer2"}, {Class: "tokenizer3"}}
	expected := []string{
		input.Tokenizers[0].Class,
		input.Tokenizers[1].Class,
		input.Tokenizers[2].Class,
	}

	// Execute
	names := GetTokenizerNames(input)

	// Assert
	test.AssertEqual(t, len(expected), len(names), "Lengths should be equal.")
}
