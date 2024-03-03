package cmdmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the model update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update <model name> [<other model names>...]",
	Short: "Update one or more models",
	Long:  "Update one or more models",
	Run:   runModelUpdate,
}

// runModelUpdate runs the model update command
func runModelUpdate(cmd *cobra.Command, args []string) {
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
	hfModels := model.GetModelsWithSourceHuggingface(configModels)
	hfModelsAvailable := model.GetModelsWithIsDownloadedTrue(hfModels)

	// Get models to update : through args or through a multiselect of models already downloaded from huggingface
	var selectedModelNames []string
	if len(args) == 0 {
		// No argument provided : multiselect among the downloaded models coming from huggingface
		message := "Please select the model(s) to be updated"
		values := model.GetNames(hfModelsAvailable)
		checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = ptermutil.DisplayInteractiveMultiselect(message, values, []string{}, checkMark, true)
		ptermutil.DisplaySelectedItems(selectedModelNames)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(args)
	}

	// Filter selected models to only keep those available for an update
	modelsToUpdate := filterModelsByStatusBeforeUpdate(selectedModelNames, hfModelsAvailable)

	// Processing filtered models for an update
	downloadedModels := processModelsForUpdate(configModels, modelsToUpdate)

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err = config.AddModels(downloadedModels)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

}

// filterModelsByStatusBeforeUpdate returns the models available for an update by determining a status for each one of them
func filterModelsByStatusBeforeUpdate(modelNames []string, hfModelsAvailable []model.Model) []model.Model {

	// Bind the downloaded models coming from huggingface to a map for faster lookup
	// Used to check whether a model has already been downloaded
	mapHfModelsAvailable := model.ModelsToMap(hfModelsAvailable)

	var modelsToUpdate []model.Model
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
		modelMapped := model.MapToModelFromHuggingfaceModel(huggingfaceModel)

		// Try to find the model in the map of downloaded models coming from huggingface
		configModel, exists := mapHfModelsAvailable[name]

		if !exists {
			// Model not configured yet, offering to download it later on
			modelsToUpdate = append(modelsToUpdate, modelMapped)
		} else if configModel.Version != modelMapped.Version {
			// Model already configured but not up-to-date
			configModel.Version = modelMapped.Version
			modelsToUpdate = append(modelsToUpdate, configModel)
		} else {
			// Model already up-to-date, nothing more to do here
			updatedModelNames = append(updatedModelNames, name)
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
func processModelsForUpdate(configModels, modelsToUpdate []model.Model) []model.Model {

	// Bind config models to a map for faster lookup
	// Used to get the model's path and check if it's already configured
	mapConfigModels := model.ModelsToMap(configModels)

	var downloadedModels []model.Model
	var failedModels []string

	// Processing all the remaining models for an update
	for _, current := range modelsToUpdate {

		success := model.Update(current, mapConfigModels)
		if success {
			downloadedModels = append(downloadedModels, current)
		} else {
			failedModels = append(failedModels, current.Name)
		}
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		pterm.Error.Println(fmt.Sprintf("The following models(s) couldn't be downloaded : %s", failedModels))
	}

	return downloadedModels
}
