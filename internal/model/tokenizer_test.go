package model

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"testing"
)

// GetTokenizer initiates a basic tokenizer with an id as suffix
func GetTokenizer(suffix int) Tokenizer {
	idStr := fmt.Sprint(suffix)
	return Tokenizer{
		Class: "tokenizer" + idStr,
		Path:  "path" + idStr,
	}
}

func GetTokenizers(length int) Tokenizers {
	var tokenizers Tokenizers
	for i := 1; i <= length; i++ {
		tokenizers = append(tokenizers, GetTokenizer(i-1))
	}
	return tokenizers
}

// TestTokenizersToMap_Success tests the TokenizersToMap function to return a map from a slice of tokenizers.
func TestTokenizersToMap_Success(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Tokenizers = GetTokenizers(3)
	expected := map[string]Tokenizer{
		model.Tokenizers[0].Class: model.Tokenizers[0],
		model.Tokenizers[1].Class: model.Tokenizers[1],
		model.Tokenizers[2].Class: model.Tokenizers[2],
	}

	// Execute
	result := model.Tokenizers.ToMap()

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestGetTokenizerNames_Success tests the GetNames function to return the correct names.
func TestGetTokenizerNames_Success(t *testing.T) {
	// Init
	input := GetModel(0)
	input.Tokenizers = Tokenizers{{Class: "tokenizer1"}, {Class: "tokenizer2"}, {Class: "tokenizer3"}}
	expected := []string{
		input.Tokenizers[0].Class,
		input.Tokenizers[1].Class,
		input.Tokenizers[2].Class,
	}

	// Execute
	names := input.Tokenizers.GetNames()

	// Assert
	test.AssertEqual(t, len(expected), len(names), "Lengths should be equal.")
}

// TestTokenizerDownloadedOnDevice_FalseMissing tests the TokenizerDownloadedOnDevice function to return false upon missing.
func TestTokenizerDownloadedOnDevice_FalseMissing(t *testing.T) {
	// Init
	tokenizer := Tokenizer{Path: ""}

	// Execute
	exists, err := tokenizer.DownloadedOnDevice()

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
	exists, err := tokenizer.DownloadedOnDevice()

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
	exists, err := tokenizer.DownloadedOnDevice()

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, exists, true)
}

// TestTokenizersNotDownloadedOnDevice_NotTransformers tests the TokenizersNotDownloadedOnDevice function while not using a transformer model.
func TestTokenizersNotDownloadedOnDevice_NotTransformers(t *testing.T) {
	// Init
	model := GetModel(0)
	model.Module = huggingface.DIFFUSERS

	// Execute
	var expected Tokenizers
	result := model.GetTokenizersNotDownloadedOnDevice()

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestTokenizersNotDownloadedOnDevice_Missing tests the TokenizersNotDownloadedOnDevice function with no missing tokenizer.
func TestTokenizersNotDownloadedOnDevice_Missing(t *testing.T) {
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
	model := GetModel(0)
	model.Module = huggingface.TRANSFORMERS
	model.Tokenizers = Tokenizers{{Path: tokenizerDirectory}}

	// Execute
	var expected Tokenizers
	result := model.GetTokenizersNotDownloadedOnDevice()

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}
