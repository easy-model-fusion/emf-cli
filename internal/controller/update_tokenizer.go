package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// TokenizerUpdateCmd runs the model remove command
func TokenizerUpdateCmd(args []string) {
}

// runTokenizerUpdate runs the tokenizer update command
func runTokenizerUpdate(cmd *cobra.Command, argsModel []string, argsTokenizer []string) {
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	// Get all models from configuration file
	configModels, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Keep the downloaded models coming from huggingface (i.e. those that could potentially be updated)
	hfModels := configModels.FilterWithSourceHuggingface()
	hfModelsAvailable := hfModels.FilterWithIsDownloadedTrue()

	// Get models to update : through args or through a multiselect of models already downloaded from huggingface
	var selectedModelNames []string
	if len(argsModel) == 0 {
		// No argument provided : multiselect among the downloaded models coming from huggingface
		message := "Please select the model(s) to be updated"
		values := hfModelsAvailable.GetNames()
		checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = app.UI().DisplayInteractiveMultiselect(message, values, checkMark, false, true)
		app.UI().DisplaySelectedItems(selectedModelNames)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(argsModel)
	}
	if len(argsTokenizer) == 0 {
		// No argument provided : multiselect among the downloaded tokenizers coming from huggingface
		message := "Please select the tokenizer(s) to be updated"
		values := hfModelsAvailable.GetNames()
		checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = app.UI().DisplayInteractiveMultiselect(message, values, checkMark, false, true)
	}

	// Filter selected models to only keep those available for an update
	modelsToUpdate := filterModelsByStatusBeforeUpdate(selectedModelNames, hfModelsAvailable)
	processModelsForUpdate(configModels, modelsToUpdate)
}

// filterModelsByStatusBeforeUpdate returns the models available for an update by determining a status for each one of them
func filterModelsByStatusBeforeUpdate(modelNames []string, hfModelsAvailable model.Models) model.Models {

	// Bind the downloaded models coming from huggingface to a map for faster lookup
	// Used to check whether a model has already been downloaded
	mapHfModelsAvailable := hfModelsAvailable.Map()

	var modelsToUpdate model.Models
	var notFoundModelNames []string
	var updatedModelNames []string

	// Check which model can be updated
	for _, name := range modelNames {

		// Fetching model from huggingface
		huggingfaceModel, err := app.H().GetModelById(name)
		if err != nil {
			// Model not found : nothing more to do here, skipping to the next one
			notFoundModelNames = append(notFoundModelNames, name)
			continue
		}

		// Fetching succeeded : processing the response
		// Map API response to model.Model
		modelMapped := model.FromHuggingfaceModel(huggingfaceModel)

		// Try to find the model in the map of downloaded models coming from huggingface
		configModel, exists := mapHfModelsAvailable[name]

		if !exists {
			// Model not configured yet, offering to download it later on
			modelsToUpdate = append(modelsToUpdate, modelMapped)
		} else {
			// Get all models
			configModel.Version = modelMapped.Version
			modelsToUpdate = append(modelsToUpdate, configModel)
		}
	}

	// Indicate the models that couldn't be found
	if len(notFoundModelNames) > 0 {
		pterm.Warning.Printfln(fmt.Sprintf("The following models(s) couldn't be found "+
			"and will be ignored : %s", notFoundModelNames))
	}
	// Indicate the models that are already up-to-date
	if len(updatedModelNames) > 0 {
		pterm.Warning.Printfln(fmt.Sprintf("The following model(s) are already up to date "+
			"and will be ignored : %s", updatedModelNames))
	}

	return modelsToUpdate
}

// processModelsForUpdate
func processModelsForUpdate(configModels, modelsToUpdate model.Models) {

	// Bind config models to a map for faster lookup
	// Used to get the model's path and check if it's already configured
	mapConfigModels := configModels.Map()

	// Processing all the remaining models for an update
	var failedModels []string
	for _, current := range modelsToUpdate {

		success := current.Update(mapConfigModels)
		if !success {
			failedModels = append(failedModels, current.Name)
		} else {
			// Add models to configuration file
			spinner, _ := pterm.DefaultSpinner.Start("Updating configuration file...")
			err := config.AddModels(model.Models{current})
			if err != nil {
				spinner.Fail(fmt.Sprintf("Error while updating the configuration file: %s", err))
			} else {
				spinner.Success()
			}
		}
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		pterm.Error.Println(fmt.Sprintf("The following models(s) couldn't be downloaded : %s", failedModels))
	}
}
