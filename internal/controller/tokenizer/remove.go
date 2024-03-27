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
	"github.com/pterm/pterm"
)

type RemoveTokenizerController struct{}

// RunTokenizerRemove runs the tokenizer remove command
func (ic RemoveTokenizerController) RunTokenizerRemove(args []string) error {
	sdk.SendUpdateSuggestion()
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := ic.processRemove(args)

	// Display messages to user
	if warningMessage != "" {
		pterm.Warning.Printfln(warningMessage)
	}

	if infoMessage != "" {
		pterm.Info.Printfln(infoMessage)
		return err
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
		return nil
	} else {
		pterm.Error.Printfln("Operation failed.")
		return err
	}
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

	// Checks the presence of the model
	selectedModelName := args[0]
	configModelsMap := models.Map()
	modelToUse, exists := configModelsMap[selectedModelName]
	if !exists {
		return warning, "Model is not configured", err
	}

	// Verify model's module
	if modelToUse.Module != huggingface.TRANSFORMERS {
		return warning, info, fmt.Errorf("only transformers models have tokenizers")
	}

	configTokenizerMap := modelToUse.Tokenizers.Map()
	var tokenizersToRemove model.Tokenizers
	var invalidTokenizers []string
	var tokenizerNames []string

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
