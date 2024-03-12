package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
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

// getTokenizer initiates a basic tokenizer with an id as suffix
func getTokenizer(suffix int) Tokenizer {
	idStr := fmt.Sprint(suffix)
	return Tokenizer{
		Class: "tokenizer" + idStr,
		Path:  "path" + idStr,
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

// TestTokenizersToMap_Success tests the TokenizersToMap function to return a map from a slice of tokenizers.
func TestTokenizersToMap_Success(t *testing.T) {
	// Init
	model := getModel(0)
	model.Tokenizers = []Tokenizer{getTokenizer(0), getTokenizer(1), getTokenizer(2)}
	expected := map[string]Tokenizer{
		model.Tokenizers[0].Class: model.Tokenizers[0],
		model.Tokenizers[1].Class: model.Tokenizers[1],
		model.Tokenizers[2].Class: model.Tokenizers[2],
	}

	// Execute
	result := TokenizersToMap(model)

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
	models := []Model{getModel(0), getModel(1)}

	// Execute
	names := GetNames(models)

	// Assert
	test.AssertEqual(t, len(models), len(names), "Lengths should be equal.")
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

// TestGetModelsWithAddToBinaryFileTrue_Success tests the GetModelsWithAddToBinaryFileTrue to return the sub-slice.
func TestGetModelsWithAddToBinaryFileTrue_Success(t *testing.T) {
	// Init
	models := []Model{getModel(0), getModel(1)}
	models[0].AddToBinaryFile = false
	expected := []Model{models[1]}

	// Execute
	result := GetModelsWithAddToBinaryFileTrue(models)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths should be equal.")
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
		Path:    "",
		Module:  "",
		Class:   "",
		Options: map[string]string{},
		Tokenizer: downloader.Tokenizer{
			Path:    "",
			Class:   "",
			Options: map[string]string{},
		},
	}
	expected := Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
		Tokenizers: []Tokenizer{
			{
				Path:  "/path/to/tokenizer",
				Class: "tokenizer_class",
				Options: map[string]string{
					"option1": "true",
					"option2": "'text'",
				},
			},
		},
	}

	// Execute
	result := MapToModelFromDownloaderModel(expected, downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Options), len(result.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(result.Tokenizers[0].Options))
}

// TestMapToModelFromDownloaderModel_Fill tests the MapToModelFromDownloaderModel to return the correct Model.
func TestMapToModelFromDownloaderModel_Fill(t *testing.T) {
	// Init
	downloaderModel := downloader.Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
		Tokenizer: downloader.Tokenizer{
			Path:  "/path/to/tokenizer",
			Class: "tokenizer_class",
			Options: map[string]string{
				"option1": "true",
				"option2": "'text'",
			},
		},
	}
	expected := Model{
		Path:   filepath.Clean("/path/to/model"),
		Module: "module_name",
		Class:  "class_name",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
		Tokenizers: []Tokenizer{
			{
				Path:  filepath.Clean("/path/to/tokenizer"),
				Class: "tokenizer_class",
				Options: map[string]string{
					"option1": "true",
					"option2": "'text'",
				},
			},
		},
	}

	// Execute
	result := MapToModelFromDownloaderModel(Model{}, downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Options), len(result.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(result.Tokenizers[0].Options))
}

// TestMapToModelFromDownloaderModel_ReplaceTokenizer tests the MapToModelFromDownloaderModel to return the correct Model.
func TestMapToModelFromDownloaderModel_ReplaceTokenizer(t *testing.T) {
	// Init
	downloaderModel := downloader.Model{
		Path:    "/path/to/model",
		Module:  "module_name",
		Class:   "class_name",
		Options: map[string]string{},
		Tokenizer: downloader.Tokenizer{
			Path:  "/path/to/tokenizer",
			Class: "tokenizer_class",
			Options: map[string]string{
				"option1": "true",
				"option2": "'text'",
			},
		},
	}
	input := Model{
		Path:    filepath.Clean("/path/to/model"),
		Module:  "module_name",
		Class:   "class_name",
		Options: map[string]string{},
		Tokenizers: []Tokenizer{
			{
				Path:    "",
				Class:   "",
				Options: map[string]string{},
			},
		},
	}
	expected := input
	expected.Tokenizers[0].Path = downloaderModel.Tokenizer.Path
	expected.Tokenizers[0].Class = downloaderModel.Tokenizer.Class
	expected.Tokenizers[0].Options = downloaderModel.Tokenizer.Options

	// Execute
	result := MapToModelFromDownloaderModel(input, downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Options), len(result.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(result.Tokenizers[0].Options))
}

// TestMapToTokenizerFromDownloaderTokenizer_Success tests the MapToTokenizerFromDownloaderTokenizer to return the correct Tokenizer.
func TestMapToTokenizerFromDownloaderTokenizer_Success(t *testing.T) {
	// Init
	downloaderTokenizer := downloader.Tokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
	}
	expected := Tokenizer{
		Path:  filepath.Clean("/path/to/tokenizer"),
		Class: "tokenizer_class",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
	}

	// Execute
	result := MapToTokenizerFromDownloaderTokenizer(downloaderTokenizer)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Options), len(result.Options))
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

// TestModelDownloadedOnDevice_FalseMissing tests the ModelDownloadedOnDevice function to return false upon missing.
func TestModelDownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	model := getModel(0)
	model.Path = ""

	// Execute
	exists, err := ModelDownloadedOnDevice(model)

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
	model := getModel(0)
	model.Path = modelDirectory

	// Execute
	exists, err := ModelDownloadedOnDevice(model)

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
	model := getModel(0)
	model.Path = modelDirectory

	// Execute
	exists, err := ModelDownloadedOnDevice(model)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}

// TestTokenizerDownloadedOnDevice_FalseMissing tests the TokenizerDownloadedOnDevice function to return false upon missing.
func TestTokenizerDownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	tokenizer := Tokenizer{Path: ""}

	// Execute
	exists, err := TokenizerDownloadedOnDevice(tokenizer)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestTokenizerDownloadedOnDevice_FalseEmpty tests the TokenizerDownloadedOnDevice function to return false upon empty.
func TestTokenizerDownloadedOnDevice_FalseEmpty(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelDirectory, err := os.MkdirTemp("", "modelDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(modelDirectory)

	// Create a temporary directory representing the tokenizer path
	tokenizerDirectory, err := os.MkdirTemp(modelDirectory, "tokenizerDirectory")
	if err != nil {
		t.Fatal(err)
	}

	// Init
	tokenizer := Tokenizer{Path: tokenizerDirectory}

	// Execute
	exists, err := TokenizerDownloadedOnDevice(tokenizer)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestTokenizersNotDownloadedOnDevice_True tests the TokenizersNotDownloadedOnDevice function to return true.
func TestTokenizersNotDownloadedOnDevice_True(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelDirectory, err := os.MkdirTemp("", "modelDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(modelDirectory)

	// Create a temporary directory representing the tokenizer path
	tokenizerDirectory, err := os.MkdirTemp(modelDirectory, "tokenizerDirectory")
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary file inside the tokenizer path
	file, err := os.CreateTemp(tokenizerDirectory, "")
	if err != nil {
		t.Fatal(err)
	}
	fileutil.CloseFile(file)

	// Init
	tokenizer := Tokenizer{Path: tokenizerDirectory}

	// Execute
	exists, err := TokenizerDownloadedOnDevice(tokenizer)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}

// TestTokenizersNotDownloadedOnDevice_NotTransformers tests the TokenizersNotDownloadedOnDevice function while not using a transformer model.
func TestTokenizersNotDownloadedOnDevice_NotTransformers(t *testing.T) {
	// Init
	model := getModel(0)
	model.Module = huggingface.DIFFUSERS

	// Execute
	var expected []Tokenizer
	result := TokenizersNotDownloadedOnDevice(model)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestTokenizersNotDownloadedOnDevice_Missing tests the TokenizersNotDownloadedOnDevice function with no missing tokenizer.
func TestTokenizersNotDownloadedOnDevice_Missing(t *testing.T) {
	// Init
	model := getModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = []Tokenizer{{Path: "tokenizer"}}

	// Execute
	expected := model.Tokenizers
	result := TokenizersNotDownloadedOnDevice(model)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestTokenizersNotDownloadedOnDevice_NotMissing tests the TokenizersNotDownloadedOnDevice function with missing tokenizers.
func TestTokenizersNotDownloadedOnDevice_NotMissing(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelDirectory, err := os.MkdirTemp("", "modelDirectory")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(modelDirectory)

	// Create a temporary directory representing the tokenizer path
	tokenizerDirectory, err := os.MkdirTemp(modelDirectory, "tokenizerDirectory")
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary file inside the tokenizer path
	file, err := os.CreateTemp(tokenizerDirectory, "")
	if err != nil {
		t.Fatal(err)
	}
	fileutil.CloseFile(file)

	// Init
	model := getModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = []Tokenizer{{Path: tokenizerDirectory}}

	// Execute
	var expected []Tokenizer
	result := TokenizersNotDownloadedOnDevice(model)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestBuildModelsFromDevice_Custom tests the BuildModelsFromDevice function to work for custom configured models.
func TestBuildModelsFromDevice_Custom(t *testing.T) {
	// Create a temporary directory representing the path to the custom model
	modelPath := path.Join(app.DownloadDirectoryPath, "custom-provider", "custom-model")
	modelPath = filepath.ToSlash(modelPath)
	err := os.MkdirAll(modelPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice()

	// Assert
	test.AssertEqual(t, len(models), 1)
	test.AssertEqual(t, models[0].Source, CUSTOM)
	test.AssertEqual(t, models[0].Path, modelPath)
	test.AssertEqual(t, models[0].AddToBinaryFile, true)
	test.AssertEqual(t, models[0].IsDownloaded, true)

}

// TestBuildModelsFromDevice_HuggingfaceEmpty tests the BuildModelsFromDevice function to work for huggingface empty models.
func TestBuildModelsFromDevice_HuggingfaceEmpty(t *testing.T) {
	// Create a temporary directory representing the path to the model which is empty
	modelDirectoryPath := path.Join(app.DownloadDirectoryPath, "stabilityai", "sdxl-turbo")
	err := os.MkdirAll(modelDirectoryPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice()

	// Assert
	test.AssertEqual(t, len(models), 0)
}

// TestBuildModelsFromDevice_HuggingfaceDiffusers tests the BuildModelsFromDevice function to work for huggingface diffusers models.
func TestBuildModelsFromDevice_HuggingfaceDiffusers(t *testing.T) {
	// Create a temporary directory representing the path to the diffusers model which is not empty
	modelName := path.Join("stabilityai", "sdxl-turbo")
	modelDirectory := path.Join(app.DownloadDirectoryPath, modelName)
	err := os.MkdirAll(path.Join(modelDirectory, "not-empty"), 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice()

	// Assert
	test.AssertEqual(t, len(models), 1)
	test.AssertEqual(t, models[0].Name, modelName)
	test.AssertEqual(t, models[0].Module, huggingface.DIFFUSERS)
	test.AssertEqual(t, models[0].Source, HUGGING_FACE)
	test.AssertEqual(t, models[0].Path, modelDirectory)
	test.AssertEqual(t, models[0].AddToBinaryFile, true)
	test.AssertEqual(t, models[0].IsDownloaded, true)
	test.AssertEqual(t, models[0].Version, "")
}

// TestBuildModelsFromDevice_HuggingfaceTransformers tests the BuildModelsFromDevice function to work for huggingface transformers models.
func TestBuildModelsFromDevice_HuggingfaceTransformers(t *testing.T) {
	// Create a temporary directory representing the path to the transformers model
	modelName := path.Join("microsoft", "phi-2")
	modelDirectory := path.Join(app.DownloadDirectoryPath, modelName)
	modelPath := path.Join(modelDirectory, "model")
	err := os.MkdirAll(modelPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Create a temporary directory representing the path to the tokenizer for the model
	tokenizerDirectory, err := os.MkdirTemp(modelDirectory, "tokenizer")
	tokenizerDirectory = filepath.ToSlash(tokenizerDirectory)
	if err != nil {
		t.Fatal(err)
	}

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice()

	// Assert
	test.AssertEqual(t, len(models), 1)
	test.AssertEqual(t, models[0].Name, modelName)
	test.AssertEqual(t, models[0].Module, huggingface.TRANSFORMERS)
	test.AssertEqual(t, models[0].Source, HUGGING_FACE)
	test.AssertEqual(t, models[0].Path, modelPath)
	test.AssertEqual(t, models[0].AddToBinaryFile, true)
	test.AssertEqual(t, models[0].IsDownloaded, true)
	test.AssertEqual(t, models[0].Version, "")

	test.AssertEqual(t, len(models[0].Tokenizers), 1)
	test.AssertEqual(t, models[0].Tokenizers[0].Path, tokenizerDirectory)

}
