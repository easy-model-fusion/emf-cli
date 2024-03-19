package model

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func SetupDownloaderForFailure() {
	app.Downloader().(*mock.MockDownloader).DownloaderModel = downloadermodel.Model{}
	app.Downloader().(*mock.MockDownloader).DownloaderError = errors.New("")
}

func SetupDownloaderForSuccess(model downloadermodel.Model) {
	// Mock python script to succeed
	app.Downloader().(*mock.MockDownloader).DownloaderModel = model
	app.Downloader().(*mock.MockDownloader).DownloaderError = nil
}

func TestMain(m *testing.M) {
	app.SetDownloader(&mock.MockDownloader{})
	app.SetUI(&mock.MockUI{})
	app.SetPython(&mock.MockPython{})
	os.Exit(m.Run())
}

// TestFromDownloaderModel_Empty tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_Empty(t *testing.T) {
	// Init
	downloaderModel := downloadermodel.Model{
		Path:    "",
		Module:  "",
		Class:   "",
		Options: map[string]string{},
		Tokenizer: downloadermodel.Tokenizer{
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

// TestFromDownloaderModel_Fill tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_Fill(t *testing.T) {
	// Init
	downloaderModel := downloadermodel.Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Options: map[string]string{
			"option1": "true",
			"option2": "'text'",
		},
		Tokenizer: downloadermodel.Tokenizer{
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

// TestFromDownloaderModel_ReplaceTokenizer tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_ReplaceTokenizer(t *testing.T) {
	// Init
	downloaderModel := downloadermodel.Model{
		Path:    "/path/to/model",
		Module:  "module_name",
		Class:   "class_name",
		Options: map[string]string{},
		Tokenizer: downloadermodel.Tokenizer{
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

// TestGetConfig_Failure tests the Model.GetConfig upon failure.
func TestGetConfig_Failure(t *testing.T) {
	// Mock python script to fail
	SetupDownloaderForFailure()

	// Init
	input := GetModel(0)
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expected := input

	// Execute
	success := input.GetConfig(downloaderArgs)

	// Assert
	test.AssertEqual(t, success, false)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}

// TestGetConfig_Success tests the Model.GetConfig upon success.
func TestGetConfig_Success(t *testing.T) {
	// Init
	input := GetModel(0)
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expectedDownloaderResult := downloadermodel.Model{Class: "test"}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)

	// Mock downloader to succeed
	SetupDownloaderForSuccess(expectedDownloaderResult)

	// Execute
	success := input.GetConfig(downloaderArgs)

	// Assert
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}

// TestModel_Download_Failure tests the Model.Download upon failure.
func TestModel_Download_Failure(t *testing.T) {
	// Mock python script to fail
	SetupDownloaderForFailure()

	// Init
	input := GetModel(0)
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expected := input

	// Execute
	success := input.Download(downloaderArgs)

	// Assert
	test.AssertEqual(t, success, false)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}

// TestModel_Download_Success tests the Model.Download upon success.
func TestModel_Download_Success(t *testing.T) {
	// Init
	input := GetModel(0)
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expectedDownloaderResult := downloadermodel.Model{Path: "test/Path"}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)
	expected.AddToBinaryFile = true
	expected.IsDownloaded = true

	// Mock python script to succeed
	SetupDownloaderForSuccess(expectedDownloaderResult)

	// Execute
	success := input.Download(downloaderArgs)

	// Assert
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}

// TestTokenizer_Download_Failure tests the Model.DownloadTokenizer upon failure.
func TestTokenizer_Download_Failure(t *testing.T) {
	// Mock python script to fail
	SetupDownloaderForFailure()

	// Init
	input := GetModel(0)
	tokenizer := GetTokenizer(0)
	input.Tokenizers = Tokenizers{tokenizer}
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expected := input

	// Execute
	success := input.DownloadTokenizer(tokenizer, downloaderArgs)

	// Assert
	test.AssertEqual(t, success, false)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}

// TestTokenizer_Download_Success tests the Model.DownloadTokenizer upon success.
func TestTokenizer_Download_Success(t *testing.T) {
	// Init
	input := GetModel(0)
	tokenizer := GetTokenizer(0)
	input.Tokenizers = Tokenizers{tokenizer}
	downloaderArgs := downloadermodel.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	expectedDownloaderResult := downloadermodel.Model{Tokenizer: downloadermodel.Tokenizer{Path: "test/Path"}}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)
	expected.AddToBinaryFile = true
	expected.IsDownloaded = true

	// Mock python script to succeed
	SetupDownloaderForSuccess(expectedDownloaderResult)

	// Execute
	success := input.DownloadTokenizer(tokenizer, downloaderArgs)

	// Assert
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}
