package cmdmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
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

	// Keep the downloaded models coming from huggingface (i.e. those that could potentially be updated)
	hfModels := model.GetModelsWithSourceHuggingface(configModels)
	hfModelsAvailable := model.GetModelsWithIsDownloadedTrue(hfModels)
	hfModelAvailableNames := model.GetNames(hfModelsAvailable)

	// Storing the names of those wished to be updated
	var selectedModelNames []string

	// Get models to update : through args or through a multiselect of models already downloaded from huggingface
	if len(args) == 0 {
		// No argument provided : multiselect among the downloaded models coming from huggingface
		message := "Please select the model(s) to be updated"
		checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = ptermutil.DisplayInteractiveMultiselect(message, hfModelAvailableNames, []string{}, checkMark, true)
		ptermutil.DisplaySelectedItems(selectedModelNames)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(args)
	}

	// Bind the downloaded models coming from huggingface to a map for faster lookup
	// Used to check whether a model has already been downloaded
	mapHfModelsAvailable := model.ModelsToMap(hfModelsAvailable)

	var modelsToUpdate []model.Model
	var notFoundModelNames []string
	var updatedModelNames []string

	// Check which model can be updated
	for _, name := range selectedModelNames {

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

	// Bind config models to a map for faster lookup
	// Used to get the model's path and check if it's already configured
	mapConfigModels := model.ModelsToMap(configModels)

	var downloadedModels []model.Model
	var failedModels []string
	var failedTokenizersForModels []string

	// Processing all the remaining models for an update
	for _, current := range modelsToUpdate {

		// Checking if the model is already configured
		_, configured := mapConfigModels[current.Name]

		// TODO : what if there is a correct custom path that the user provided? how about the tokenizer paths?
		// Check if model is physically present on the device
		current = model.ConstructConfigPaths(current)
		downloaded, err := model.ModelDownloadedOnDevice(current)
		if err != nil {
			failedModels = append(failedModels, current.Name)
			continue
		}

		// Process internal state of the model
		install := false
		if !configured && !downloaded {
			install = ptermutil.AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be added. "+
				"Would you like to add it?", current.Name))
		} else if configured && !downloaded {
			install = ptermutil.AskForUsersConfirmation(fmt.Sprintf("Model '%s' has yet to be downloaded. "+
				"Would you like to download it?", current.Name))
		} else if !configured && downloaded {
			install = ptermutil.AskForUsersConfirmation(fmt.Sprintf("Model '%s' already exists. "+
				"Would you like to overwrite it?", current.Name))
		} else {
			// Model already configured and downloaded : a new version is available
			install = ptermutil.AskForUsersConfirmation(fmt.Sprintf("New version of '%s' is available. "+
				"Would you like to overwrite its old version?", current.Name))
		}

		// Model will not be downloaded or overwritten, nothing more to do here
		if !install {
			continue
		}

		// Downloader script to skip the tokenizers download process if none selected
		var skip string

		// If transformers : select the tokenizers to update using a multiselect
		var tokenizerNames []string
		if current.Module == huggingface.TRANSFORMERS {

			// Get tokenizer names for the model
			availableNames := model.GetTokenizerNames(current)

			// Allow to select only if at least one tokenizer is available
			if len(availableNames) > 0 {

				// Prepare the tokenizers multiselect
				message := "Please select the tokenizer(s) to be updated"
				checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
				tokenizerNames = ptermutil.DisplayInteractiveMultiselect(message, availableNames, availableNames, checkMark, true)
				ptermutil.DisplaySelectedItems(tokenizerNames)

				// No tokenizer is selected : skipping so that it doesn't overwrite the default one
				if len(tokenizerNames) > 0 {
					skip = downloader.SkipValueTokenizer
				}
			}
		}

		// TODO : options model => Waiting for issue 74 to be completed : [Client] Model options to config
		// Prepare the script arguments
		downloaderArgs := downloader.Args{
			ModelName:    current.Name,
			ModelModule:  string(current.Module),
			ModelClass:   current.Class,
			ModelOptions: []string{},
			Skip:         skip,
		}

		// Running the script
		dlModel, err := downloader.Execute(downloaderArgs)

		// Something went wrong or no data has been returned
		if err != nil || dlModel.IsEmpty {
			failedModels = append(failedModels, current.Name)
			continue
		}

		// Update the model for the configuration file
		current = model.MapToModelFromDownloaderModel(current, dlModel)
		current.AddToBinaryFile = true
		current.IsDownloaded = true

		// Bind the model tokenizers to a map for faster lookup
		mapModelTokenizers := model.TokenizersToMap(current)

		var failedTokenizers []string
		for _, tokenizerName := range tokenizerNames {
			tokenizer := mapModelTokenizers[tokenizerName]

			// TODO : options tokenizer => Waiting for issue 74 to be completed : [Client] Model options to config
			// Building downloader args for the tokenizer
			downloaderArgs.Skip = downloader.SkipValueModel
			downloaderArgs.TokenizerClass = tokenizer.Class
			downloaderArgs.TokenizerOptions = []string{}

			// Running the script for the tokenizer only
			dlModelTokenizer, err := downloader.Execute(downloaderArgs)

			// Something went wrong or no data has been returned
			if err != nil || dlModelTokenizer.IsEmpty {
				failedTokenizers = append(failedTokenizers, tokenizer.Class)
				continue
			}

			// Update the model with the tokenizer for the configuration file
			current = model.MapToModelFromDownloaderModel(current, dlModelTokenizer)
		}

		if len(failedTokenizers) > 0 {
			failedTokenizersForModels = append(failedModels, fmt.Sprintf("These tokenizers could not be downloaded for '%s': %s", current.Name, failedTokenizers))
		}

		downloadedModels = append(downloadedModels, current)
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		pterm.Error.Println(fmt.Sprintf("These models could not be downloaded : %s", failedModels))
	}
	for _, failedTokenizers := range failedTokenizersForModels {
		pterm.Error.Println(failedTokenizers)
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
