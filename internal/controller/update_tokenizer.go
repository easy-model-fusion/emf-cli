package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

// TokenizerUpdateCmd TokenizerRemoveCmd runs the model remove command
func TokenizerUpdateCmd(args []string) {
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := processUpdateTokenizer(args)

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

func processUpdateTokenizer(args []string) (string, string, error) {
	// Load the configuration file
	if config.GetViperConfig(config.FilePath) != nil {
	}

	sdk.SendUpdateSuggestion()

	// Get all models from configuration file
	configModels, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
	}

	// Keep the downloaded models coming from huggingface (i.e. those that could potentially be updated)
	hfModels := configModels.FilterWithSourceHuggingface()
	hfModelsAvailable := hfModels.FilterWithIsDownloadedTrue()

	// Get all configured models objects/names
	var selectedModelNames []string
	var selectedModels string
	if len(args) < 1 {
		message := "Please select the model for which to update tokenizer"
		values := hfModelsAvailable.GetNames()
		checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = app.UI().DisplayInteractiveMultiselect(message, values, checkMark, false, true)
		app.UI().DisplaySelectedItems(selectedModelNames)

	}
	if len(selectedModelNames) == 0 && len(args) == 0 {
		pterm.Error.Println("Please select at least one Model !")
	}

	if len(selectedModelNames) > 0 {
		selectedModels = selectedModelNames[0]
	} else {
		selectedModels = args[0]
	}

	var models model.Models
	models, err = config.GetModels()
	if err != nil {
		return "", "", err
	}

	sdk.SendUpdateSuggestion()

	configModelsMap := models.Map()
	modelsToUse, exists := configModelsMap[selectedModels]
	if !exists {
		return "", "model do not exist", err
	}

	var tokenizerNames []string
	if modelsToUse.Module == huggingface.TRANSFORMERS {
		availableNames := modelsToUse.Tokenizers.GetNames()
		if len(availableNames) > 0 {
			message := "Please select the tokenizer(s) to be updated"
			checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
			tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, availableNames, checkMark, true, true)
			app.UI().DisplaySelectedItems(tokenizerNames)
		} else if len(args) > 1 {
			// Check if selectedTokenizerNames elements exist in tokenizerNames and add them to a new list
			var selectedAndAvailableTokenizerNames []string
			for _, name := range args {
				for _, availableName := range tokenizerNames {
					if name == availableName {
						selectedAndAvailableTokenizerNames = append(selectedAndAvailableTokenizerNames, name)
						break
					}
				}
			}
			tokenizerNames = selectedAndAvailableTokenizerNames
		}
	}
	// Update the selected models'
	modelsToUse.UpdateTokenizer(tokenizerNames)
	return "", "Tokenizers update done", err
}
