package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/emf-cli/test"
)

// TestDownloadedOnDevice_FalseMissing tests the Model.DownloadedOnDevice function to return false upon missing.
func TestDownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Path = ""

	// Execute
	exists, err := model.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestDownloadedOnDevice_FalseEmpty tests the Model.DownloadedOnDevice function to return false upon empty.
func TestDownloadedOnDevice_FalseEmpty(t *testing.T) {
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

// TestDownloadedOnDevice_True tests the Model.DownloadedOnDevice function to return true.
func TestDownloadedOnDevice_True(t *testing.T) {
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

// TestTokenizer_DownloadedOnDevice_FalseMissing tests the Tokenizer.DownloadedOnDevice function to return false upon missing.
func TestTokenizer_DownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	tokenizer := Tokenizer{Path: ""}

	// Execute
	exists, err := tokenizer.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestTokenizer_DownloadedOnDevice_FalseEmpty tests the Tokenizer.DownloadedOnDevice function to return false upon empty.
func TestTokenizer_DownloadedOnDevice_FalseEmpty(t *testing.T) {
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
	exists, err := tokenizer.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, false)
}

// TestTokenizer_DownloadedOnDevice_True tests the Tokenizer.DownloadedOnDevice function to return true.
func TestTokenizer_DownloadedOnDevice_True(t *testing.T) {
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
	exists, err := tokenizer.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}

// TestGetTokenizersNotDownloadedOnDevice_NotTransformers tests the Model.GetTokenizersNotDownloadedOnDevice function while not using a transformer model.
func TestGetTokenizersNotDownloadedOnDevice_NotTransformers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.DIFFUSERS

	// Execute
	var expected Tokenizers
	result := model.GetTokenizersNotDownloadedOnDevice()

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestGetTokenizersNotDownloadedOnDevice_Missing tests the Model.GetTokenizersNotDownloadedOnDevice function with no missing tokenizer.
func TestGetTokenizersNotDownloadedOnDevice_Missing(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = Tokenizers{{Path: "tokenizer"}}

	// Execute
	expected := model.Tokenizers
	result := model.GetTokenizersNotDownloadedOnDevice()

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestGetTokenizersNotDownloadedOnDevice_NotMissing tests the Model.GetTokenizersNotDownloadedOnDevice function with missing tokenizers.
func TestGetTokenizersNotDownloadedOnDevice_NotMissing(t *testing.T) {
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
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = Tokenizers{{Path: tokenizerDirectory}}

	// Execute
	var expected Tokenizers
	result := model.GetTokenizersNotDownloadedOnDevice()

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

// TestFromHuggingfaceModel_Success tests the FromHuggingfaceModel to return the correct Model.
func TestFromHuggingfaceModel_Success(t *testing.T) {
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
