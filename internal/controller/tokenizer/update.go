// Package tokenizer
// This file contains the update tokenizer controller which is responsible for
// updating existing tokenizers in existing models
//

package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type UpdateTokenizerController struct{}

// TokenizerUpdateCmd TokenizerRemoveCmd runs the model update command
func (ic UpdateTokenizerController) TokenizerUpdateCmd(args []string) error {
	sdk.SendUpdateSuggestion()
	// Process remove operation with given arguments
	warningMessages, infoMessage, err := ic.processUpdateTokenizer(args)

	// Display messages to user
	for _, warningMessage := range warningMessages {
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

// processUpdateTokenizer processes tokenizers to be updated
func (ic UpdateTokenizerController) processUpdateTokenizer(args []string) (warnings []string, info string, err error) {
	// Load the configuration file
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warnings, info, err
	}

	// Get all configured models objects/names and args model
	models, err := config.GetModelsByModule(string(huggingface.TRANSFORMERS))
	if err != nil {
		return warnings, info, fmt.Errorf("error get model: %s", err.Error())
	}
	if len(models) == 0 {
		err = fmt.Errorf("no configured models found")
		return warnings, info, err
	}
	var modelToUse model.Model

	configModelsMap := models.Map()
	if len(args) == 0 {
		// Get selected models from select
		sc := SelectModelController{}
		// Get selected models from select
		modelToUse = sc.SelectTransformerModel(models)
	} else {
		// Get the selected models from the args
		selectedModelName := args[0]
		var exists bool
		modelToUse, exists = configModelsMap[selectedModelName]
		if !exists {
			err = fmt.Errorf("model is not configured")
			return warnings, info, err
		}

		// Remove model name from arguments
		args = args[1:]
	}

	var updateTokenizers model.Tokenizers
	var failedTokenizers []string

	// Extracting available tokenizers
	availableNames := modelToUse.Tokenizers.GetNames()

	// Processing arguments
	if len(args) > 0 {
		args = stringutil.SliceRemoveDuplicates(args)
		configTokenizersMap := modelToUse.Tokenizers.Map()
		// Check if selectedTokenizerNames elements exist in tokenizerNames and add them to a new list

		for _, name := range args {
			tokenizer, exists := configTokenizersMap[name]
			if !exists {
				failedTokenizers = append(failedTokenizers, name)
			} else {
				updateTokenizers = append(updateTokenizers, tokenizer)
			}
		}
	} else if len(availableNames) > 0 {
		message := "Please select the tokenizer(s) to be updated"
		tokenizerNames := app.UI().DisplayInteractiveMultiselect(message, availableNames, app.UI().BasicCheckmark(), true, true, 8)
		if len(tokenizerNames) != 0 {
			app.UI().DisplaySelectedItems(tokenizerNames)
			updateTokenizers = modelToUse.Tokenizers.FilterWithClass(tokenizerNames)
		}
	}

	// Try to update all the given models
	var updatedTokenizers model.Tokenizers
	for _, tokenizer := range updateTokenizers {

		downloaderArgs := downloadermodel.Args{
			ModelName:   modelToUse.Name,
			ModelModule: string(modelToUse.Module),
		}
		downloaderArgs.OnlyConfiguration = !modelToUse.IsDownloaded

		var success bool
		success, warnings, err = modelToUse.DownloadTokenizer(tokenizer, downloaderArgs)
		if err != nil {
			return warnings, info, err
		}
		if !success {
			failedTokenizers = append(failedTokenizers, tokenizer.Class)
		} else {
			updatedTokenizers = append(updatedTokenizers, tokenizer)
		}
	}

	// Update tokenizers' configuration
	if len(updatedTokenizers) > 0 {
		//Reset model while keeping unchanged tokenizers
		modelToUse.Tokenizers = modelToUse.Tokenizers.Difference(updatedTokenizers)
		//Adding new version of updated tokenizers
		modelToUse.Tokenizers = append(modelToUse.Tokenizers, updatedTokenizers...)

		spinner := app.UI().StartSpinner("Updating configuration file...")
		err = config.AddModels(model.Models{modelToUse})
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while updating the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}

	// Displaying the downloads that failed
	if len(failedTokenizers) > 0 {
		err = fmt.Errorf("the following tokenizer(s) couldn't be downloaded : %s", failedTokenizers)
	}
	return warnings, "Tokenizers update done", err
}
