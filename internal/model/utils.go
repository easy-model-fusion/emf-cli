package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
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

// ConstructConfigPaths to update the model's path to elements accordingly to its configuration.
func ConstructConfigPaths(current Model) Model {
	basePath := path.Join(downloader.DirectoryPath, current.Name)
	modelPath := basePath
	if current.Module == string(TRANSFORMERS) {
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
		model.Path = utils.PathUniformize(dlModel.Path)
		model.Module = dlModel.Module
		model.Class = dlModel.Class
	}

	// Check if ScriptTokenizer is valid
	if !downloader.EmptyTokenizer(dlModel.Tokenizer) {
		tokenizer := MapToTokenizerFromDownloaderTokenizer(dlModel.Tokenizer)
		// TODO : check if tokenizer already exists
		model.Tokenizers = append(model.Tokenizers, tokenizer)
	}

	return model
}

// MapToTokenizerFromDownloaderTokenizer maps data from downloader.Tokenizer to Tokenizer.
func MapToTokenizerFromDownloaderTokenizer(dlTokenizer downloader.Tokenizer) Tokenizer {
	var modelTokenizer Tokenizer
	modelTokenizer.Path = utils.PathUniformize(dlTokenizer.Path)
	modelTokenizer.Class = dlTokenizer.Class
	return modelTokenizer
}
