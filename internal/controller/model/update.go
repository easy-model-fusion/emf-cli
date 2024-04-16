package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/hfinterface"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/resultutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
)

// RunModelUpdate runs the model update command
func RunModelUpdate(args []string, yes bool, accessToken string) {
	// Process update operation with given arguments
	result := processUpdate(args, yes, accessToken)

	// Display messages to user
	result.Display("Operation succeeded.", "Operation failed.")
}

// processUpdate processes the update model operation
func processUpdate(args []string, yes bool, accessToken string) (result resultutil.ExecutionResult) {
	// Load the configuration file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		result.SetError(err)
		return result
	}

	// Request an update suggestion of the client when needed
	sdk.SendUpdateSuggestion()

	// Get all models from configuration file
	configModels, err := config.GetModels()
	if err != nil {
		result.SetError(err)
		return result
	}

	// Keep the downloaded models coming from huggingface (i.e. those that could potentially be updated)
	hfModels := configModels.FilterWithSourceHuggingface()
	hfModelsAvailable := hfModels.FilterWithIsDownloadedTrue()

	// Get models to update : through args or through a multiselect of models already downloaded from huggingface
	var selectedModelNames []string
	if len(args) == 0 {
		// No argument provided : multiselect among the downloaded models coming from huggingface
		modelNames := hfModelsAvailable.GetNames()
		selectedModelNames = selectModelsToUpdate(modelNames)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(args)
	}

	// Verify if the user selected some models to update
	if len(selectedModelNames) > 0 {
		var warning string
		// Filter selected models to only keep those available for an update
		modelsToUpdate, notFoundModelNames, upToDateModelNames := getUpdatableModels(selectedModelNames, hfModelsAvailable, accessToken)

		// Indicate the models that couldn't be found
		if len(notFoundModelNames) > 0 {
			warning = fmt.Sprintf("The following models(s) couldn't be found "+
				"and were ignored : %s", notFoundModelNames)
			result.AddWarnings([]string{warning})
		}
		// Indicate the models that are already up-to-date
		if len(upToDateModelNames) > 0 {
			result.AddInfos([]string{fmt.Sprintf("The following model(s) are already up to date "+
				"and were ignored : %s", upToDateModelNames)})
		}

		// Processing filtered models for an update
		var warningMessages []string
		warningMessages, err = updateModels(modelsToUpdate, yes, accessToken)
		result.AddWarnings(warningMessages)
	} else {
		result.AddInfos([]string{"There is no models to be updated."})
	}

	return result
}

// getUpdatableModels returns the models available for an update
func getUpdatableModels(modelNames []string, hfModelsAvailable model.Models, accessToken string) (
	modelsToUpdate model.Models, notFoundModelNames, upToDateModelNames []string) {

	// Bind the downloaded models coming from huggingface to a map for faster lookup
	// Used to check whether a model has already been downloaded
	mapHfModelsAvailable := hfModelsAvailable.Map()

	// Check which model can be updated
	for _, name := range modelNames {
		// Try to find the model in the map of downloaded models coming from huggingface
		configModel, exists := mapHfModelsAvailable[name]

		if !exists {
			// Model not configured
			notFoundModelNames = append(notFoundModelNames, name)
			continue
		}

		// Fetching model from huggingface
		token := configModel.AccessToken
		if accessToken != "" {
			token = accessToken
		}
		huggingfaceModel, err := hfinterface.GetModelById(name, token)
		if err != nil {
			// Model not found : nothing more to do here, skipping to the next one
			notFoundModelNames = append(notFoundModelNames, name)
			continue
		}

		// Fetching succeeded : processing the response
		// Map API response to model.Model
		modelMapped := model.FromHuggingfaceModel(huggingfaceModel)

		if configModel.Version != modelMapped.Version {
			// Model already configured but not up-to-date
			configModel.Version = modelMapped.Version
			modelsToUpdate = append(modelsToUpdate, configModel)
		} else {
			// Model already up-to-date, nothing more to do here
			upToDateModelNames = append(upToDateModelNames, name)
		}
	}

	return modelsToUpdate, notFoundModelNames, upToDateModelNames
}

// updateModels updates the given models
func updateModels(modelsToUpdate model.Models, yes bool, accessToken string) (warnings []string, err error) {
	// Try to update all the given models
	var failedModels []string
	var updatedModels model.Models
	for _, current := range modelsToUpdate {
		var success bool
		warnings, success, err = current.Update(yes, accessToken)
		if err != nil {
			return warnings, err
		}

		if !success {
			failedModels = append(failedModels, current.Name)
		} else {
			updatedModels = append(updatedModels, current)
		}
	}

	// Update models' configuration
	if len(updatedModels) > 0 {
		spinner := app.UI().StartSpinner("Updating configuration file...")
		err := config.AddModels(updatedModels)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while updating the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		err = fmt.Errorf("the following models(s) couldn't be downloaded : %s", failedModels)
	}

	return warnings, err
}

// selectModelsToUpdate displays an interactive multiselect so the user can choose the models to update
func selectModelsToUpdate(modelNames []string) (selectedModelNames []string) {
	if len(modelNames) > 0 {
		message := "Please select the model(s) to be updated"
		selectedModelNames = app.UI().DisplayInteractiveMultiselect(message, modelNames, app.UI().BasicCheckmark(), false, true, 8)
		app.UI().DisplaySelectedItems(selectedModelNames)
	}
	return selectedModelNames
}
