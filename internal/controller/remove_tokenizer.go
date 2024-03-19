package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"os"
)

// TokenizerRemoveCmd runs the model remove command
func TokenizerRemoveCmd(args []string) {
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := processRemoveTokenizer(args)

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

func processRemoveTokenizer(args []string) (string, string, error) {
	// Load the configuration file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return "", "", err
	}

	// No model name in args
	if len(args) < 1 {
		return "", "Enter a model in argument", err
	}

	// Get all configured models objects/names and args model
	selectedModel := args[0]
	var models model.Models
	models, err = config.GetModels()
	if err != nil {
		return "", "Error get model", err
	}

	sdk.SendUpdateSuggestion()

	// checks the presence of the model
	configModelsMap := models.Map()
	modelsToUse, exists := configModelsMap[selectedModel]
	if !exists {
		return "", "model do not exist", err
	}

	// load tokenizer for the chosen model
	var tokenizerNames []string
	if modelsToUse.Module == huggingface.TRANSFORMERS {
		availableNames := modelsToUse.Tokenizers.GetNames()
		// No tokenizer, asks for tokenizers names
		if len(availableNames) > 0 && len(args) == 1 {
			message := "Please select the tokenizer(s) to be deleted"
			checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
			tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, availableNames, checkMark, false, true)
			app.UI().DisplaySelectedItems(tokenizerNames)
		} else if len(args) > 1 {
			// Check if selectedTokenizerNames elements exist in tokenizerNames and add them to a new list
			var selectedAndAvailableTokenizerNames []string

			for i, name := range args {
				if i == 0 {
					continue
				}
				for _, availableName := range availableNames {
					if name == availableName {
						selectedAndAvailableTokenizerNames = append(selectedAndAvailableTokenizerNames, name)
						break
					}
				}
			}
			tokenizerNames = selectedAndAvailableTokenizerNames
		}
	}

	// if one or more tokenizers selected
	if len(tokenizerNames) > 0 {
		classesMap := make(map[string]bool)
		for _, class := range tokenizerNames {
			classesMap[class] = true
		}

		// delete tokenizer file and remove tokenizer to config file
		for index, tokenizer := range modelsToUse.Tokenizers {
			if classesMap[tokenizer.Class] {
				err := os.RemoveAll(tokenizer.Path)
				if err != nil {
					pterm.Error.Println(err.Error())
					return "", "Error remove tokenizer dir", err
				}
				modelsToUse.Tokenizers = append(modelsToUse.Tokenizers[:index], modelsToUse.Tokenizers[index+1:]...)
			}
		}
		spinner := app.UI().StartSpinner("Writing model to configuration file...")
		// update config file
		err = config.AddModels(model.Models{modelsToUse})
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}
	return "", "", err
}
