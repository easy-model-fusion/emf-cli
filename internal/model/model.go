package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type Models []Model
type Model struct {
	Name            string
	Path            string
	Module          huggingface.Module
	Class           string
	Options         map[string]string
	Tokenizers      Tokenizers
	PipelineTag     huggingface.PipelineTag
	Source          string
	AddToBinaryFile bool
	IsDownloaded    bool
	Version         string
}

// Sources
const (
	HUGGING_FACE = "hugging_face"
	CUSTOM       = "custom"
)

// Empty checks if the models slice is empty.
func (m Models) Empty() bool {
	// No models currently downloaded
	return len(m) == 0
}

// ContainsByName checks if a models slice contains the requested model name
func (m Models) ContainsByName(name string) bool {
	for _, currentModel := range m {
		if currentModel.Name == name {
			return true
		}
	}
	return false
}

// Difference returns the models in that are not present in `slice`
func (m Models) Difference(slice Models) Models {
	var difference Models
	for _, item := range m {
		if !slice.ContainsByName(item.Name) {
			difference = append(difference, item)
		}
	}
	return difference
}

// Union returns the models present in `slice` as well
func (m Models) Union(slice Models) Models {
	var union Models
	for _, item := range m {
		if slice.ContainsByName(item.Name) {
			union = append(union, item)
		}
	}
	return union
}

// ToMap creates a map from models for faster lookup.
func (m Models) ToMap() map[string]Model {
	modelsMap := make(map[string]Model)
	for _, current := range m {
		modelsMap[current.Name] = current
	}
	return modelsMap
}

// GetNames retrieves the names from the models.
func (m Models) GetNames() []string {
	var modelNames []string
	for _, item := range m {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames
}

// GetByNames retrieves the models by their names given an input slice.
func (m Models) GetByNames(namesSlice []string) Models {
	// Create a map for faster lookup
	namesMap := stringutil.SliceToMap(namesSlice)

	// Slice of all the models that were found
	var namesModels Models

	// Find the requested models
	for _, existingModel := range m {
		// Check if this model exists and adds it to the result
		if _, exists := namesMap[existingModel.Name]; exists {
			namesModels = append(namesModels, existingModel)
		}
	}

	return namesModels
}

// FilterWithSourceHuggingface return a sub-slice of models sourcing from huggingface.
func (m Models) FilterWithSourceHuggingface() Models {
	var huggingfaceModels Models
	for _, current := range m {
		if current.Source == HUGGING_FACE {
			huggingfaceModels = append(huggingfaceModels, current)
		}
	}
	return huggingfaceModels
}

// FilterWithIsDownloadedTrue return a sub-slice of models with IsDownloaded to true.
func (m Models) FilterWithIsDownloadedTrue() Models {
	var downloadedModels Models
	for _, current := range m {
		if current.IsDownloaded {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// FilterWithAddToBinaryFileTrue return a sub-slice of models with AddToBinaryFile to true.
func (m Models) FilterWithAddToBinaryFileTrue() Models {
	var downloadedModels Models
	for _, current := range m {
		if current.AddToBinaryFile {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// DownloadedOnDevice returns true if the model is physically present on the device.
func (m *Model) DownloadedOnDevice() (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(m.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(m.Path)
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
