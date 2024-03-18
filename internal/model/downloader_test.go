package model

import (
	"encoding/json"
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/test"
	mock "github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func GetDownloaderResultForTransformers() string {
	return "{\n" +
		"\"module\": \"transformers\",\n" +
		"\"class\": \"Class\",\n" +
		"\"tokenizer\": {\n" +
		"\"class\": \"Tokenizer\",\n" +
		"\"path\": \"./models/microsoft/phi-2\\\\Tokenizer\",\n" +
		"\"options\": {\n" +
		"\"trust_remote_code\": \"True\"\n" +
		"}\n" +
		"},\n" +
		"\"path\": \"./models/provider/name\\\\model\",\n" +
		"\"options\": {\n" +
		"\"torch_dtype\": \"\\\"auto\\\"\",\n" +
		"\"trust_remote_code\": \"True\"\n" +
		"}\n}"
}

func GetDownloaderResultForTransformersTokenizer() string {
	return "{\n" +
		"\"module\": \"transformers\",\n" +
		"\"tokenizer\": {\n" +
		"\"class\": \"AutoTokenizer\",\n" +
		"\"path\": \"./models/microsoft/phi-2\\\\AutoTokenizer\",\n" +
		"\"options\": {\n" +
		"\"trust_remote_code\": \"True\"\n" +
		"}\n" +
		"}\n}"
}

func SetupDownloaderForFailure() {
	app.Python().(*mock.MockPython).ScriptResult = []byte{}
	app.Python().(*mock.MockPython).Error = errors.New("")
}

func SetupDownloaderForSuccess(result string) {
	// Mock python script to succeed
	app.Python().(*mock.MockPython).ScriptResult = []byte(result)
	app.Python().(*mock.MockPython).Error = nil
}

func TestMain(m *testing.M) {
	app.SetUI(&mock.MockUI{})
	app.SetPython(&mock.MockPython{})
	os.Exit(m.Run())
}

// TestFromDownloaderModel_Empty tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_Empty(t *testing.T) {
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

// TestFromDownloaderModel_Fill tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_Fill(t *testing.T) {
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

// TestFromDownloaderModel_ReplaceTokenizer tests the Model.FromDownloaderModel to return the correct Model.
func TestFromDownloaderModel_ReplaceTokenizer(t *testing.T) {
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

// TestGetConfig_Failure tests the Model.GetConfig upon failure.
func TestGetConfig_Failure(t *testing.T) {
	// Mock python script to fail
	SetupDownloaderForFailure()

	// Init
	input := GetModel(0)
	downloaderArgs := downloader.Args{
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
	// Mock python script to succeed
	output := GetDownloaderResultForTransformers()
	SetupDownloaderForSuccess(output)

	// Init
	input := GetModel(0)
	downloaderArgs := downloader.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	var expectedDownloaderResult downloader.Model
	err := json.Unmarshal([]byte(output), &expectedDownloaderResult)
	if err != nil {
		t.FailNow()
	}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)

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
	downloaderArgs := downloader.Args{
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
	// Mock python script to succeed
	output := GetDownloaderResultForTransformers()
	SetupDownloaderForSuccess(output)

	// Init
	input := GetModel(0)
	downloaderArgs := downloader.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	var expectedDownloaderResult downloader.Model
	err := json.Unmarshal([]byte(output), &expectedDownloaderResult)
	if err != nil {
		t.FailNow()
	}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)
	expected.AddToBinaryFile = true
	expected.IsDownloaded = true

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
	downloaderArgs := downloader.Args{
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
	// Mock python script to succeed
	output := GetDownloaderResultForTransformersTokenizer()
	SetupDownloaderForSuccess(output)

	// Init
	input := GetModel(0)
	tokenizer := GetTokenizer(0)
	input.Tokenizers = Tokenizers{tokenizer}
	downloaderArgs := downloader.Args{
		ModelName:     input.Name,
		ModelModule:   string(input.Module),
		DirectoryPath: app.DownloadDirectoryPath,
	}

	// Expected
	var expectedDownloaderResult downloader.Model
	err := json.Unmarshal([]byte(output), &expectedDownloaderResult)
	if err != nil {
		t.FailNow()
	}
	expected := input
	expected.FromDownloaderModel(expectedDownloaderResult)
	expected.AddToBinaryFile = true
	expected.IsDownloaded = true

	// Execute
	success := input.DownloadTokenizer(tokenizer, downloaderArgs)

	// Assert
	test.AssertEqual(t, success, true)
	test.AssertEqual(t, reflect.DeepEqual(expected, input), true)
}
