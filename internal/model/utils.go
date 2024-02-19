package model

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
)

// Empty checks if the models slice is empty.
func Empty(models []Model) bool {
	// No models currently downloaded
	if len(models) == 0 {
		pterm.Info.Println("Models list is empty.")
		return true
	}
	return false
}

// Contains checks if a models slice contains the requested model
func Contains(models []Model, model Model) bool {
	for _, item := range models {
		if model == item {
			return true
		}
	}
	return false
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
	namesMap := utils.MapFromArrayString(namesSlice)

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
