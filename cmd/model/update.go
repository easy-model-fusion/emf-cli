package cmdmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
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

	// Keep the models coming from huggingface
	hfModels := model.GetModelsWithSourceHuggingface(configModels)

	// Keep the downloaded models coming from huggingface (i.e. those that could potentially be updated)
	hfModelsAvailable := model.GetModelsWithIsDownloadedTrue(hfModels)
	hfModelAvailableNames := model.GetNames(hfModelsAvailable)

	// Storing the names of those wished to be updated
	var selectedModelNames []string

	// Get models to update
	if len(args) == 0 {
		// No argument provided : multiselect among the downloaded models coming from huggingface
		message := "Please select the model(s) to be updated"
		checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = app.UI().DisplayInteractiveMultiselect(message, hfModelAvailableNames, checkMark, true)
		app.UI().DisplaySelectedItems(selectedModelNames)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(args)
	}

	// Bind the downloaded models coming from huggingface to a map for faster lookup
	mapHfModelsAvailable := model.ModelsToMap(hfModelsAvailable)

	var modelsToUpdate []model.Model
	var notFoundModelNames []string
	var updatedModelNames []string

	// Check which model can be updated
	for _, name := range selectedModelNames {

		// Fetching model from huggingface
		huggingfaceModel, err := app.H().GetModelById(name)
		if err != nil {
			// Model not found : skipping to the next one
			notFoundModelNames = append(notFoundModelNames, name)
			continue
		}

		// Map API response to model.Model
		modelMapped := model.MapToModelFromHuggingfaceModel(huggingfaceModel)

		// Try to find the model in the map of downloaded models coming from huggingface
		configModel, exists := mapHfModelsAvailable[name]

		if !exists {
			// Model not configured yet
			modelsToUpdate = append(modelsToUpdate, modelMapped)
		} else if configModel.Version != modelMapped.Version {
			// Model already configured but not up-to-date
			modelsToUpdate = append(modelsToUpdate, configModel)
		} else {
			// Model already up-to-date
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

	// Bind config models to a map for faster lookup
	mapConfigModels := model.ModelsToMap(configModels)

	var downloadedModels []model.Model

	// Processing all the remaining models for an update
	for _, current := range modelsToUpdate {

		// Checking if the model is already configured
		_, configured := mapConfigModels[current.Name]

		// Checking if the model is already physically downloaded
		current = model.ConstructConfigPaths(current)
		downloaded, err := fileutil.IsExistingPath(current.Path)
		if err != nil {
			continue
		}

		install := false

		// Process internal state of the model
		if !configured && !downloaded {
			install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be added. "+
				"Would you like to add it?", current.Name))
		} else if configured && !downloaded {
			install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be downloaded. "+
				"Would you like to download it?", current.Name))
		} else if !configured && downloaded {
			install = app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' already exists. "+
				"Would you like to overwrite it?", current.Name))
		} else {
			// Model already configured and downloaded : a new version is available
			install = app.UI().AskForUsersConfirmation(fmt.Sprintf("New version of '%s' is available. "+
				"Would you like to overwrite its old version?", current.Name))
		}

		// Model will not be downloaded
		if !install {
			continue
		}

		// If transformers : select the tokenizers to update through a multiselect
		var tokenizerNames []string
		if current.Module == huggingface.TRANSFORMERS {
			tokenizerNames = model.GetTokenizerNames(current)
			message := "Please select the tokenizer(s) to be updated"
			checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
			tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, tokenizerNames, checkMark, true)
			app.UI().DisplaySelectedItems(tokenizerNames)
		}

		// TODO : options model and tokenizer => Waiting for issue 74 to be completed : [Client] Model options to config
		// Prepare the script arguments
		downloaderArgs := downloader.Args{
			ModelName:   current.Name,
			ModelModule: string(current.Module),
			ModelClass:  current.Class,
		}

		// Running the script
		dlModel, err := downloader.Execute(downloaderArgs)

		// Something went wrong or no data has been returned
		if err != nil || dlModel.IsEmpty {
			continue
		}

		// Update the model for the configuration file
		current = model.MapToModelFromDownloaderModel(current, dlModel)
		current.AddToBinaryFile = true
		current.IsDownloaded = true

		downloadedModels = append(downloadedModels, current)

		// TODO : download tokenizers => Waiting for issue 55 to be completed : [Client] Edit downloader execute
	}

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err = config.AddModels(downloadedModels)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

}
