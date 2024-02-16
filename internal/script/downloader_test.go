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
