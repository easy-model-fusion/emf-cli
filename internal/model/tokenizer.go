package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type Tokenizers []Tokenizer
type Tokenizer struct {
	Path    string
	Class   string
	Options map[string]string
}

// ToMap creates a map from tokenizers for faster lookup.
func (t Tokenizers) ToMap() map[string]Tokenizer {
	tokenizersMap := make(map[string]Tokenizer)
	for _, current := range t {
		tokenizersMap[current.Class] = current
	}
	return tokenizersMap
}

// GetNames retrieves the names from the tokenizers.
func (t Tokenizers) GetNames() []string {
	var tokenizerNames []string
	for _, current := range t {
		tokenizerNames = append(tokenizerNames, current.Class)
	}
	return tokenizerNames
}

// DownloadedOnDevice returns true if the tokenizer is physically present on the device.
func (t *Tokenizer) DownloadedOnDevice() (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(t.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(t.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if empty {
		// Model is not downloaded on the device
		return false, nil
	}

	// Model is downloaded on the device
	return true, nil
}

// GetTokenizersNotDownloadedOnDevice returns the list of tokenizers that should but are not physically present on the device.
func (m *Model) GetTokenizersNotDownloadedOnDevice() Tokenizers {

	// Model can't have any tokenizer
	if m.Module != huggingface.TRANSFORMERS {
		return Tokenizers{}
	}

	// Processing the configured tokenizers
	var notDownloadedTokenizers Tokenizers
	for _, tokenizer := range m.Tokenizers {

		// Check if tokenizer is already downloaded
		downloaded, err := tokenizer.DownloadedOnDevice()
		if err != nil {
			// An error occurred
			continue
		} else if !downloaded {
			// Tokenizer is not downloaded on the device
			notDownloadedTokenizers = append(notDownloadedTokenizers, tokenizer)
			continue
		}
	}

	return notDownloadedTokenizers
}
