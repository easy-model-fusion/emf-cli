package model

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/script"
	"path/filepath"
	"testing"

	"github.com/easy-model-fusion/client/test"
)

// getModel initiates a basic model with an id as suffix
func getModel(suffix int) Model {
	idStr := fmt.Sprint(suffix)
	return Model{
		Name:        "model" + idStr,
		Config:      Config{Module: "module" + idStr, Class: "class" + idStr},
		AddToBinary: true,
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

// TestMapToConfigFromDownloadScriptModel_Empty tests the MapToConfigFromDownloadScriptModel to return the correct Config.
func TestMapToConfigFromDownloadScriptModel_Empty(t *testing.T) {
	// Init
	dsm := script.DownloaderModel{
		Path:   "",
		Module: "",
		Class:  "",
		Tokenizer: script.DownloaderTokenizer{
			Path:  "",
			Class: "",
		},
	}
	expected := Config{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Tokenizers: []Tokenizer{
			{Path: "/path/to/tokenizer", Class: "tokenizer_class"},
		},
	}

	// Execute
	result := MapToConfigFromScriptDownloadModel(expected, dsm)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
}

// TestMapToConfigFromDownloadScriptModel_Fill tests the MapToConfigFromDownloadScriptModel to return the correct Config.
func TestMapToConfigFromDownloadScriptModel_Fill(t *testing.T) {
	// Init
	sm := script.DownloaderModel{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
		Tokenizer: script.DownloaderTokenizer{
			Path:  "/path/to/tokenizer",
			Class: "tokenizer_class",
		},
	}
	expected := Config{
		Path:   filepath.Clean("/path/to/model"),
		Module: "module_name",
		Class:  "class_name",
		Tokenizers: []Tokenizer{
			{Path: filepath.Clean("/path/to/tokenizer"), Class: "tokenizer_class"},
		},
	}

	// Execute
	result := MapToConfigFromScriptDownloadModel(Config{}, sm)

	// Assert
	test.AssertEqual(t, expected.Path, result.Path)
	test.AssertEqual(t, expected.Module, result.Module)
	test.AssertEqual(t, expected.Class, result.Class)
	test.AssertEqual(t, len(expected.Tokenizers), len(result.Tokenizers))
	test.AssertEqual(t, expected.Tokenizers[0].Path, result.Tokenizers[0].Path)
	test.AssertEqual(t, expected.Tokenizers[0].Class, result.Tokenizers[0].Class)
}

// TestMapToTokenizerFromDownloaderScriptTokenizer tests the MapToTokenizerFromDownloaderScriptTokenizer to return the correct Tokenizer.
func TestMapToTokenizerFromDownloaderScriptTokenizer(t *testing.T) {
	// Init
	st := script.DownloaderTokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}
	expected := Tokenizer{Path: filepath.Clean("/path/to/tokenizer"), Class: "tokenizer_class"}

	// Execute
	result := MapToTokenizerFromScriptDownloaderTokenizer(st)

	// Assert
	test.AssertEqual(t, expected, result)
}
