package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/emf-cli/test"
)

// TestConstructConfigPaths_Default tests the ConstructConfigPaths for a default model.
func TestConstructConfigPaths_Default(t *testing.T) {
	// Init
	model := GetModel(0)

	// Execute
	model.ConstructConfigPaths()

	// Assert
	test.AssertEqual(t, model.Path, path.Join(app.DownloadDirectoryPath, model.Name))
}

// TestConstructConfigPaths_Transformers tests the ConstructConfigPaths for a transformers model.
func TestConstructConfigPaths_Transformers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS

	// Execute
	model.ConstructConfigPaths()

	// Assert
	test.AssertEqual(t, model.Path, path.Join(app.DownloadDirectoryPath, model.Name, "model"))
}

// TestConstructConfigPaths_Transformers tests the ConstructConfigPaths for a transformers model.
func TestConstructConfigPaths_TransformersTokenizers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = []Tokenizer{{Class: "tokenizer"}}

	// Execute
	model.ConstructConfigPaths()

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
		Tokenizers: Tokenizers{
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
	input := expected

	// Execute
	input.FromDownloaderModel(downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, input.Path)
	test.AssertEqual(t, expected.Module, input.Module)
	test.AssertEqual(t, expected.Class, input.Class)
	test.AssertEqual(t, len(expected.Options), len(input.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(input.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, input.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, input.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(input.Tokenizers[0].Options))
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
	input := Model{}

	// Execute
	input.FromDownloaderModel(downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, input.Path)
	test.AssertEqual(t, expected.Module, input.Module)
	test.AssertEqual(t, expected.Class, input.Class)
	test.AssertEqual(t, len(expected.Options), len(input.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(input.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, input.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, input.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(input.Tokenizers[0].Options))
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
	input.FromDownloaderModel(downloaderModel)

	// Assert
	test.AssertEqual(t, expected.Path, input.Path)
	test.AssertEqual(t, expected.Module, input.Module)
	test.AssertEqual(t, expected.Class, input.Class)
	test.AssertEqual(t, len(expected.Options), len(input.Options))
	test.AssertEqual(t, len(expected.Tokenizers), len(input.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, input.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, input.Tokenizers[0].Class)
	test.AssertEqual(t, len(expected.Tokenizers[0].Options), len(input.Tokenizers[0].Options))
}

// TestMapToTokenizerFromDownloaderTokenizer_Success tests the MapToTokenizerFromDownloaderTokenizer to return the correct Tokenizer.
/*func TestMapToTokenizerFromDownloaderTokenizer_Success(t *testing.T) {
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
}*/

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
