package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test/dmock"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/pterm/pterm"
	"os"
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
	exists, err := model.DownloadedOnDevice(false)

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
	exists, err := model.DownloadedOnDevice(false)

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
	exists, err := model.DownloadedOnDevice(false)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}

// TestModelDownloadedOnDevice_UseBasePath_True tests the ModelDownloadedOnDevice function to return true.
func TestModelDownloadedOnDevice_UseBasePath_True(t *testing.T) {
	// Create a temporary directory representing the model base path
	modelName := fileutil.PathJoin("microsoft", "phi-2")
	modelDirectory := fileutil.PathJoin(app.DownloadDirectoryPath, modelName)
	modelPath := fileutil.PathJoin(modelDirectory, "model")
	err := os.MkdirAll(modelPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Init
	model := GetModel(0)
	model.Name = modelName

	// Execute
	exists, err := model.DownloadedOnDevice(true)

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
	modelPath := fileutil.PathJoin(app.DownloadDirectoryPath, "custom-provider", "custom-model")
	modelPath = filepath.ToSlash(modelPath)
	err := os.MkdirAll(modelPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice("")

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
	modelDirectoryPath := fileutil.PathJoin(app.DownloadDirectoryPath, "stabilityai", "sdxl-turbo")
	err := os.MkdirAll(modelDirectoryPath, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice("")

	// Assert
	test.AssertEqual(t, len(models), 0)
}

// TestBuildModelsFromDevice_HuggingfaceDiffusers tests the BuildModelsFromDevice function to work for huggingface diffusers models.
func TestBuildModelsFromDevice_HuggingfaceDiffusers(t *testing.T) {
	// Create a temporary directory representing the path to the diffusers model which is not empty
	modelName := fileutil.PathJoin("stabilityai", "sdxl-turbo")
	modelDirectory := fileutil.PathJoin(app.DownloadDirectoryPath, modelName)
	err := os.MkdirAll(fileutil.PathJoin(modelDirectory, "not-empty"), 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(app.DownloadDirectoryPath)

	// Execute
	app.InitHuggingFace(huggingface.BaseUrl, "")
	models := BuildModelsFromDevice("")

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
	modelName := fileutil.PathJoin("microsoft", "phi-2")
	modelDirectory := fileutil.PathJoin(app.DownloadDirectoryPath, modelName)
	modelPath := fileutil.PathJoin(modelDirectory, "model")
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
	models := BuildModelsFromDevice("")

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

// Tests TidyConfiguredModel on clean model
func TestTidyConfiguredModel_CleanModel(t *testing.T) {
	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Init
	model := Model{
		Name:   "model1",
		Path:   "models",
		Module: huggingface.DIFFUSERS,
	}

	// Synchronize model
	warnings, success, clean, err := model.TidyConfiguredModel("")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, clean, true)
}

// Tests TidyConfiguredModel
func TestTidyConfiguredModel_Success(t *testing.T) {
	// Init
	model := Model{
		Name:   "Test",
		Path:   "invalid/path",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: Tokenizers{
			Tokenizer{
				Class: "class1",
				Path:  "invalid/Path",
			},
			Tokenizer{
				Class: "class2",
				Path:  "invalid/Path",
			},
		},
	}
	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Synchronize model
	warnings, success, clean, err := model.TidyConfiguredModel("")

	// Assertions
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, clean, false)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
}

// Tests TidyConfiguredModel
func TestTidyConfiguredModel_Fail(t *testing.T) {
	// Init
	model := Model{
		Name:   "Test",
		Path:   "invalid/path",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: Tokenizers{
			Tokenizer{
				Class: "class1",
				Path:  "invalid/Path",
			},
			Tokenizer{
				Class: "class2",
				Path:  "invalid/Path",
			},
		},
	}
	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Synchronize model
	warnings, success, clean, err := model.TidyConfiguredModel("")

	// Assertions
	test.AssertEqual(t, success, false)
	test.AssertEqual(t, clean, false)
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
}

// Tests TidyConfiguredModel
func TestTidyConfiguredModel_FailTokenizersTidy(t *testing.T) {
	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Init
	model := Model{
		Name:   "model4/name",
		Path:   "models/model4/name/model",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: Tokenizers{
			Tokenizer{
				Class: "tokenizer",
				Path:  "models/model4/name/tokenizer",
			},
			Tokenizer{
				Class: "tokenizer2",
				Path:  "invalid/Path",
			},
		},
	}
	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Synchronize model
	warnings, success, clean, err := model.TidyConfiguredModel("")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, clean, false)
	test.AssertEqual(t, warnings[0], "The following tokenizer(s) couldn't be downloaded for 'model4/name': [tokenizer2]")
}

// Tests successful update on diffusers model
func TestModelUpdate_Diffusers(t *testing.T) {
	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Init
	model := Model{
		Name:   "model1",
		Path:   "models",
		Module: huggingface.DIFFUSERS,
	}

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Update model
	_, success, err := model.Update(true, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, success, true)
}

// Tests update on diffusers model with no user confirmation
func TestModelUpdate_WithNoConfirmation(t *testing.T) {
	// Init
	model := Model{
		Name:   "Test",
		Path:   "invalid/path",
		Module: huggingface.DIFFUSERS,
	}

	// Create ui mock
	ui := mock.MockUI{UserConfirmationResult: false}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Update model
	_, success, err := model.Update(false, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, success, false)
}

// Tests failed update
func TestModelUpdate_Failed(t *testing.T) {
	// Init
	model := Model{
		Name:   "Test",
		Path:   "invalid/path",
		Module: huggingface.DIFFUSERS,
	}

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Update model
	_, success, err := model.Update(true, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, success, false)
}

// Tests successful update on transformers model
func TestModelUpdate_Transformers(t *testing.T) {
	// Init
	model := Model{
		Name:   "Test",
		Path:   "invalid/path",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: Tokenizers{
			Tokenizer{
				Class: "class1",
				Path:  "invalid/Path",
			},
			Tokenizer{
				Class: "class2",
				Path:  "invalid/Path",
			},
		},
	}
	expectedSelections := []string{"class1"}

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Update model
	_, success, err := model.Update(true, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, success, true)
}

// TestModel_GetModelDirectorySuccess test the success case of GetModelDirectory
func TestModel_GetModelDirectorySuccess(t *testing.T) {
	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Init
	model := Model{
		Name:   "model4/name",
		Path:   "folder/folder2/models/model4/name/model",
		Module: huggingface.TRANSFORMERS,
	}

	resultPath, err := model.GetModelDirectory()
	expectedPath := "folder/folder2/models"
	test.AssertEqual(t, err, nil, "No error message")
	test.AssertEqual(t, resultPath, expectedPath, "Path is as expected")
}

// TestModel_GetModelDirectoryFail test the fail case of GetModelDirectory
func TestModel_GetModelDirectoryFail(t *testing.T) {
	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Init
	model := Model{
		Name:   "model4/name",
		Path:   "",
		Module: huggingface.TRANSFORMERS,
	}

	_, err := model.GetModelDirectory()
	expectedMessage := "directory invalid ."
	test.AssertEqual(t, err.Error(), expectedMessage, "Directory error message")
}
