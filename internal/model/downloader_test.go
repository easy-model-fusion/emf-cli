package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/test"
	"path/filepath"
	"testing"
)

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
