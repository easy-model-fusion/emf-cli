// Package tokenizer
// This file contains the remove tokenizer controller which is responsible for removing
// existing tokenizers in existing models
package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type RemoveTokenizerController struct{}

// RunTokenizerRemove runs the tokenizer remove command
func (ic RemoveTokenizerController) RunTokenizerRemove(args []string) error {
	sdk.SendUpdateSuggestion()
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := ic.processRemove(args)

	// Display messages to user
	if warningMessage != "" {
		app.UI().Warning().Printfln(warningMessage)
	}

	if infoMessage != "" {
		app.UI().Info().Printfln(infoMessage)
		return err
	} else if err == nil {
		app.UI().Success().Printfln("Operation succeeded.")
		return nil
	} else {
		app.UI().Error().Printfln("Operation failed.")
		return err
	}
}

// selectModel displays a selector of models from which the user will choose to add to his project
func selectModel(models model.Models, configModelsMap map[string]model.Model) model.Model {
	// Build a selector with each model name
	availableModelNames := models.GetNames()
	// List of models that accept tokenizers
	var compatibleModels []string
	// Check for valid tokenizers
	for _, modelName := range availableModelNames {
		module := configModelsMap[modelName]
		if module.Module == huggingface.TRANSFORMERS {
			compatibleModels = append(compatibleModels, modelName)
		}
	}
	message := "Please select the model for which to remove tokenizer"
	selectedModelName := app.UI().DisplayInteractiveSelect(message, compatibleModels, true, 8)

	// Get newly selected model
	selectedModels := models.FilterWithNames([]string{selectedModelName})

	// returns newly selected models + models entered in args
	return selectedModels[0]
}

// processRemove processes the remove tokenizer operation
func (ic RemoveTokenizerController) processRemove(args []string) (warning, info string, err error) {
	// Load the configuration file
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warning, info, err
	}

	// Get all configured models objects/names and args model
	var models model.Models
	models, err = config.GetModels()
	if err != nil {
		return warning, info, fmt.Errorf("error get model: %s", err.Error())
	}

	var tokenizersToRemove model.Tokenizers
	var invalidTokenizers []string
	var tokenizerNames []string
	var modelToUse model.Model
	configModelsMap := models.Map()
	// No args, asks for model names
	if len(args) == 0 {
		// Get selected models from select
		modelToUse = selectModel(models, configModelsMap)
		// No tokenizer, asks for tokenizers names
		availableNames := modelToUse.Tokenizers.GetNames()
		tokenizerNames = selectTokenizersToDelete(availableNames)
	} else {
		// Get the selected models from the args
		selectedModelName := args[0]
		var exists bool
		modelToUse, exists = configModelsMap[selectedModelName]
		if !exists {
			return warning, "Model is not configured", err
		}
		// Verify model's module
		if modelToUse.Module != huggingface.TRANSFORMERS {
			return warning, info, fmt.Errorf("only transformers models have tokenizers")
		}
		// Remove model name from arguments
		args = args[1:]
		if len(args) == 0 {
			// No tokenizer, asks for tokenizers names
			availableNames := modelToUse.Tokenizers.GetNames()
			tokenizerNames = selectTokenizersToDelete(availableNames)
		} else if len(args) > 0 {
			// Check for duplicates
			tokenizerNames = stringutil.SliceRemoveDuplicates(args)
		}
	}

	configTokenizerMap := modelToUse.Tokenizers.Map()

	// Check for valid tokenizers
	for _, name := range tokenizerNames {
		tokenizer, exists := configTokenizerMap[name]
		if !exists {
			invalidTokenizers = append(invalidTokenizers, name)
		} else {
			tokenizersToRemove = append(tokenizersToRemove, tokenizer)
		}
	}

	if len(invalidTokenizers) > 0 {
		warning = fmt.Sprintf("those tokenizers are invalid and will be ignored: %s", invalidTokenizers)
	}

	if len(tokenizersToRemove) == 0 {
		return warning, "no selected tokenizers to remove", err
	}

	// Delete tokenizer file and remove tokenizer to config file
	failedTokenizersRemove, err := config.RemoveTokenizersByName(modelToUse, tokenizersToRemove)
	if err != nil {
		return warning, info, err
	}
	if len(failedTokenizersRemove) > 0 {
		info = fmt.Sprintf("failed to remove these tokenizers: %s", failedTokenizersRemove)
	}

	return warning, info, err
}

// selectTokenizersToDelete displays an interactive multiselect so the user can choose the tokenizers to remove
func selectTokenizersToDelete(tokenizerNames []string) []string {
	// Displays the multiselect only if the user has previously configured some tokenizers
	if len(tokenizerNames) > 0 {
		message := "Please select the tokenizer(s) to be deleted"
		tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, tokenizerNames, app.UI().BasicCheckmark(), false, true, 8)
		app.UI().DisplaySelectedItems(tokenizerNames)
	}
	return tokenizerNames
}
