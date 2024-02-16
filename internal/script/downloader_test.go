package script

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

// TestIsDownloaderScriptModelEmpty_True tests the IsDownloaderScriptModelEmpty to return true.
func TestIsDownloaderScriptModelEmpty_True(t *testing.T) {
	// Init
	sm := DownloaderModel{}

	// Execute
	result := IsDownloaderScriptModelEmpty(sm)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestIsDownloaderScriptModelEmpty_False tests the IsDownloaderScriptModelEmpty to return true.
func TestIsDownloaderScriptModelEmpty_False(t *testing.T) {
	// Init
	sm := DownloaderModel{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
	}

	// Execute
	result := IsDownloaderScriptModelEmpty(sm)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestIsScriptTokenizerEmpty_True tests the IsDownloaderScriptTokenizer to return true.
func TestIsScriptTokenizerEmpty_True(t *testing.T) {
	// Init
	st := DownloaderTokenizer{}

	// Execute
	result := IsDownloaderScriptTokenizer(st)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestIsScriptTokenizerEmpty_False tests the IsDownloaderScriptTokenizer to return true.
func TestIsScriptTokenizerEmpty_False(t *testing.T) {
	// Init
	st := DownloaderTokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}

	// Execute
	result := IsDownloaderScriptTokenizer(st)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestProcessArgsForDownload tests the ProcessArgsForDownload
func TestProcessArgsForDownload(t *testing.T) {
	// Init
	args := DownloadArgs{
		DownloadPath:     "/path/to/download",
		ModelName:        "model",
		ModelModule:      "module",
		ModelClass:       "class",
		ModelOptions:     []string{"opt1=val1", "opt2=val2"},
		TokenizerClass:   "tokenizer",
		TokenizerOptions: []string{"tok_opt1=val1"},
		Skip:             "model",
		Overwrite:        true,
	}
	expected := []string{
		"--emf-client", "/path/to/download", "model", "module",
		"--model-class", "class",
		"--model-options", "opt1=val1", "opt2=val2",
		"--tokenizer-class", "tokenizer",
		"--tokenizer-options", "tok_opt1=val1",
		"--skip", "model",
		"--overwrite",
	}

	// Execute
	result := ProcessArgsForDownload(args)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}
