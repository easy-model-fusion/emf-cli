package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"path"
)

// Empty checks if the models slice is empty.
func Empty(models []Model) bool {
	// No models currently downloaded
	return len(models) == 0
}

// ContainsByName checks if a models slice contains the requested model by name
func ContainsByName(models []Model, name string) bool {
	for _, currentModel := range models {
		if currentModel.Name == name {
			return true
		}
	}
	return false
}

// Difference returns the models in `parentSlice` that are not present in `subSlice`
func Difference(parentSlice, subSlice []Model) []Model {
	var difference []Model
	for _, item := range parentSlice {
		if !ContainsByName(subSlice, item.Name) {
			difference = append(difference, item)
		}
	}
	return difference
}

// Union returns the models present in both `slice1` and `slice2`
func Union(slice1, slice2 []Model) []Model {
	var union []Model
	for _, item := range slice1 {
		if ContainsByName(slice2, item.Name) {
			union = append(union, item)
		}
	}
	return union
}

// ModelsToMap creates a map from a slice of models for faster lookup.
func ModelsToMap(models []Model) map[string]Model {
	modelsMap := make(map[string]Model)
	for _, current := range models {
		modelsMap[current.Name] = current
	}
	return modelsMap
}

// TokenizersToMap creates a map from a slice of tokenizers for faster lookup.
func TokenizersToMap(model Model) map[string]Tokenizer {
	tokenizersMap := make(map[string]Tokenizer)
	for _, current := range model.Tokenizers {
		tokenizersMap[current.Class] = current
	}
	return tokenizersMap
}

// GetNames retrieves the names from the models.
func GetNames(models []Model) []string {
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames
}

// GetTokenizerNames retrieves the tokenizer names from a model.
func GetTokenizerNames(model Model) []string {
	var names []string
	for _, current := range model.Tokenizers {
		names = append(names, current.Class)
	}
	return names
}

// GetModelsByNames retrieves the models by their names given an input slice.
func GetModelsByNames(models []Model, namesSlice []string) []Model {
	// Create a map for faster lookup
	namesMap := stringutil.SliceToMap(namesSlice)

	// Slice of all the models that were found
	var namesModels []Model

	// Find the requested models
	for _, existingModel := range models {
		// Check if this model exists and adds it to the result
		if _, exists := namesMap[existingModel.Name]; exists {
			namesModels = append(namesModels, existingModel)
		}
	}

	return namesModels
}

// GetModelsWithSourceHuggingface return a sub-slice of models sourcing from huggingface.
func GetModelsWithSourceHuggingface(models []Model) []Model {
	var huggingfaceModels []Model
	for _, current := range models {
		if current.Source == HUGGING_FACE {
			huggingfaceModels = append(huggingfaceModels, current)
		}
	}
	return huggingfaceModels
}

// GetModelsWithIsDownloadedTrue return a sub-slice of models with isdownloaded to true.
func GetModelsWithIsDownloadedTrue(models []Model) []Model {
	var downloadedModels []Model
	for _, current := range models {
		if current.IsDownloaded {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// GetModelsWithAddToBinaryFileTrue return a sub-slice of models with AddToBinaryFile to true.
func GetModelsWithAddToBinaryFileTrue(models []Model) []Model {
	var downloadedModels []Model
	for _, current := range models {
		if current.AddToBinaryFile {
			downloadedModels = append(downloadedModels, current)
		}
	}
	return downloadedModels
}

// ConstructConfigPaths to update the model's path to elements accordingly to its configuration.
func ConstructConfigPaths(current Model) Model {
	basePath := path.Join(downloader.DirectoryPath, current.Name)
	modelPath := basePath
	if current.Module == huggingface.TRANSFORMERS {
		modelPath = path.Join(modelPath, "model")
		for i, tokenizer := range current.Tokenizers {
			current.Tokenizers[i].Path = path.Join(basePath, tokenizer.Class)
		}
	}
	current.Path = modelPath

	return current
}

// MapToModelFromDownloaderModel maps data from downloader.Model to Model.
func MapToModelFromDownloaderModel(model Model, dlModel downloader.Model) Model {

	// Check if ScriptModel is valid
	if !downloader.EmptyModel(dlModel) {
		model.Path = stringutil.PathUniformize(dlModel.Path)
		model.Module = huggingface.Module(dlModel.Module)
		model.Class = dlModel.Class
	}

	// Check if ScriptTokenizer is valid
	if !downloader.EmptyTokenizer(dlModel.Tokenizer) {
		tokenizer := MapToTokenizerFromDownloaderTokenizer(dlModel.Tokenizer)

		// Check if tokenizer already configured and replace it
		var replaced bool
		for i := range model.Tokenizers {
			if model.Tokenizers[i].Class == tokenizer.Class {
				model.Tokenizers[i] = tokenizer
				replaced = true
			}
		}

		// Tokenizer was already found and replaced : nothing to append
		if replaced {
			return model
		}

		// Tokenizer not found : adding it to the list
		model.Tokenizers = append(model.Tokenizers, tokenizer)
	}

	return model
}

// MapToTokenizerFromDownloaderTokenizer maps data from downloader.Tokenizer to Tokenizer.
func MapToTokenizerFromDownloaderTokenizer(dlTokenizer downloader.Tokenizer) Tokenizer {
	var modelTokenizer Tokenizer
	modelTokenizer.Path = stringutil.PathUniformize(dlTokenizer.Path)
	modelTokenizer.Class = dlTokenizer.Class
	return modelTokenizer
}

// MapToModelFromHuggingfaceModel map the Huggingface API model to a model
func MapToModelFromHuggingfaceModel(huggingfaceModel huggingface.Model) Model {
	var model Model
	model.Name = huggingfaceModel.Name
	model.PipelineTag = huggingfaceModel.PipelineTag
	model.Module = huggingfaceModel.LibraryName
	model.Source = HUGGING_FACE
	model.Version = huggingfaceModel.LastModified
	return model
}

// ModelDownloadedOnDevice returns true if the model is physically present on the device.
func ModelDownloadedOnDevice(model Model) (bool, error) {

	// Check if model is already downloaded
	downloaded, err := fileutil.IsExistingPath(model.Path)
	if err != nil {
		// An error occurred
		return false, err
	} else if !downloaded {
		// Model is not downloaded on the device
		return false, nil
	}

	// Check if the model directory is empty
	empty, err := fileutil.IsDirectoryEmpty(model.Path)
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

func TokenizersNotDownloadedOnDevice(model Model) []Tokenizer {
	var notDownloadedTokenizers []Tokenizer
	for _, tokenizer := range model.Tokenizers {

		// Check if tokenizer is already downloaded
		downloaded, err := fileutil.IsExistingPath(tokenizer.Path)
		if err != nil {
			// An error occurred
			continue
		} else if downloaded {
			// Tokenizer is downloaded on the device
			continue
		}

		// Check if the tokenizer directory is empty
		empty, err := fileutil.IsDirectoryEmpty(model.Path)
		if err != nil {
			// An error occurred
			continue
		} else if !empty {
			// Tokenizer is downloaded on the device
			continue
		}

		// Tokenizer is not downloaded on the device
		notDownloadedTokenizers = append(notDownloadedTokenizers, tokenizer)
	}

	return notDownloadedTokenizers
}
