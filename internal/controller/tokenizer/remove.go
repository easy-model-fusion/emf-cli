package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
)

// TokenizerRemoveCmd runs the model remove command
func TokenizerRemoveCmd(args []string) {
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := processRemove(args)

	// Display messages to user
	if warningMessage != "" {
		pterm.Warning.Printfln(warningMessage)
	}

	if infoMessage != "" {
		pterm.Info.Printfln(infoMessage)
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

// processRemove processes the remove tokenizer operation
func processRemove(args []string) (warning, info string, err error) {
	// Load the configuration file
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warning, info, err
	}

	// No model name in args
	if len(args) < 1 {
		return warning, info, fmt.Errorf("enter a model in argument")
	}

	// Get all configured models objects/names and args model
	selectedModel := args[0]
	var models model.Models
	models, err = config.GetModels()
	if err != nil {
		return warning, info, fmt.Errorf("error get model: %s", err.Error())
	}

	sdk.SendUpdateSuggestion()

	// checks the presence of the model
	configModelsMap := models.Map()
	modelsToUse, exists := configModelsMap[selectedModel]
	if !exists {
		return warning, "Model is not configured", err
	}
	if modelsToUse.Module != huggingface.TRANSFORMERS {
		return warning, info, fmt.Errorf("only transformers models have tokzenizers")
	}

	configTokenizerMap := modelsToUse.Tokenizers.Map()
	var tokenizersToRemove model.Tokenizers
	var invalidTokenizers []string
	var tokenizerNames []string
	// load tokenizer for the chosen model

	// No tokenizer, asks for tokenizers names
	if len(args) == 1 {
		availableNames := modelsToUse.Tokenizers.GetNames()
		tokenizerNames = selectTokenizersToDelete(availableNames)

	} else if len(args) > 1 {
		// Check if selectedTokenizerNames elements exist in tokenizerNames and add them to a new list
		tokenizerNames = stringutil.SliceRemoveDuplicates(args)
	}
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
		info = fmt.Sprintf("no selected tokenizers to remove")
		return warning, info, err
	}

	// if one or more tokenizers selected
	var removedTokenizers model.Tokenizers

	// delete tokenizer file and remove tokenizer to config file
	for _, tokenizer := range tokenizersToRemove {
		err := os.RemoveAll(tokenizer.Path)
		if err != nil {
			pterm.Error.Println(err.Error())
			return "", "Error remove tokenizer dir", err
		}
		// Successfully removed tokenizer
		removedTokenizers = append(removedTokenizers, tokenizer)

	}
	spinner := app.UI().StartSpinner("Writing model to configuration file...")
	// update config file
	err = config.AddModels(model.Models{modelsToUse})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

	return warning, info, err
}

// selectTokenizersToDelete displays an interactive multiselect so the user can choose the tokenizers to remove
func selectTokenizersToDelete(tokenizerNames []string) []string {
	// Displays the multiselect only if the user has previously configured some models but hasn't selected all of them
	if len(tokenizerNames) > 0 {
		message := "Please select the tokenizer(s) to be deleted"
		checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, tokenizerNames, checkMark, false, true)
		app.UI().DisplaySelectedItems(tokenizerNames)
	}
	return tokenizerNames
}
