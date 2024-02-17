package model

import (
	"github.com/easy-model-fusion/client/internal/script"
	"github.com/easy-model-fusion/client/internal/utils"
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

// GetNames retrieves the names from the models.
func GetNames(models []Model) []string {
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}
	return modelNames
}

// GetModelsByNames retrieves the models by their names given an input slice.
func GetModelsByNames(models []Model, namesSlice []string) []Model {
	// Create a map for faster lookup
	namesMap := utils.SliceToMap(namesSlice)

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

// MapToConfigFromScriptDownloadModel maps data from script.DownloaderModel to Config.
func MapToConfigFromScriptDownloadModel(config Config, dsm script.DownloaderModel) Config {

	// Check if ScriptModel is valid
	if !script.IsDownloaderScriptModelEmpty(dsm) {
		config.Path = utils.PathUniformize(dsm.Path)
		config.Module = dsm.Module
		config.Class = dsm.Class
	}

	// Check if ScriptTokenizer is valid
	if !script.IsDownloaderScriptTokenizer(dsm.Tokenizer) {
		tokenizer := MapToTokenizerFromScriptDownloaderTokenizer(dsm.Tokenizer)
		config.Tokenizers = append(config.Tokenizers, tokenizer)
	}

	return config
}

// MapToTokenizerFromScriptDownloaderTokenizer maps data from script.DownloaderTokenizer to Tokenizer.
func MapToTokenizerFromScriptDownloaderTokenizer(dst script.DownloaderTokenizer) Tokenizer {
	var modelTokenizer Tokenizer
	modelTokenizer.Path = utils.PathUniformize(dst.Path)
	modelTokenizer.Class = dst.Class
	return modelTokenizer
}
