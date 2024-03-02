package cmd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strings"
)

// tidyCmd represents the model tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "add missing and remove unused models",
	Long:  `add missing and remove unused models`,
	Run:   runTidy,
}

// runTidy runs the model tidy command
func runTidy(cmd *cobra.Command, args []string) {
	// get all models from config file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		pterm.Error.Println(err.Error())
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	models, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Tidy the models configured but not physically present on the device
	err = tidyModelsConfiguredButNotDownloaded(models)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Tidy the models physically present on the device but not configured
	tidyModelsDownloadedButNotConfigured(models)

	// Updating the models object since the configuration might have changed in between
	models, err = config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Regenerate python code
	err = regenerateCode(models)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}
}

// tidyModelsConfiguredButNotDownloaded downloads any missing model and its missing tokenizers as well
func tidyModelsConfiguredButNotDownloaded(models []model.Model) error {
	pterm.Info.Println("Verifying if all models are downloaded...")
	// filter the models that should be added to binary
	models = model.GetModelsWithAddToBinaryFileTrue(models)

	// Search for the models that need to be downloaded
	var downloadedModels []model.Model
	var failedModels []string
	var failedTokenizersForModels []string

	// Tidying the configured but not downloaded models and tokenizers
	for _, current := range models {

		// TODO : what if there is a correct custom path that the user provided?
		// Check if model is physically present on the device
		current = model.ConstructConfigPaths(current)
		downloaded, err := model.ModelDownloadedOnDevice(current)
		if err != nil {
			failedModels = append(failedModels, current.Name)
			continue
		}

		// Get all the configured but not downloaded tokenizers
		missingTokenizers := model.TokenizersNotDownloadedOnDevice(current)

		// Model is clean, nothing more to do here
		if downloaded && len(missingTokenizers) == 0 {
			downloadedModels = append(downloadedModels, current)
			continue
		}

		// TODO : options model => Waiting for issue 74 to be completed : [Client] Model options to config
		// Prepare the script arguments
		downloaderArgs := downloader.Args{
			ModelName:    current.Name,
			ModelModule:  string(current.Module),
			ModelClass:   current.Class,
			ModelOptions: []string{},
		}

		// Model has yet to be downloaded
		if !downloaded {

			// If at least one tokenizer is already installed : skipping the default tokenizer
			if len(current.Tokenizers) >= 0 && len(current.Tokenizers) > len(missingTokenizers) {
				downloaderArgs.Skip = downloader.SkipValueTokenizer
			}

			// TODO : write a DownloadModel function

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
		}

		// Some tokenizers are missing
		if len(missingTokenizers) != 0 {

			// Downloading the missing tokenizers
			var failedTokenizers []string
			for _, tokenizer := range missingTokenizers {

				// TODO : write a DownloadTokenizer function

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

			// The process failed for at least one tokenizer
			if len(failedTokenizers) > 0 {
				failedTokenizersForModels = append(failedModels, fmt.Sprintf("The following tokenizer(s) couldn't be downloaded for '%s': %s", current.Name, failedTokenizers))
			}
		}

		downloadedModels = append(downloadedModels, current)
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		pterm.Error.Println(fmt.Sprintf("The following models(s) couldn't be downloaded : %s", failedModels))
	}
	for _, failedTokenizers := range failedTokenizersForModels {
		pterm.Error.Println(failedTokenizers)
	}

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err := config.AddModels(downloadedModels)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

	return nil
}

// tidyModelsDownloadedButNotConfigured configuring the downloaded models that aren't configured in the configuration file
// and then asks the user if he wants to delete them or add them to the configuration file
func tidyModelsDownloadedButNotConfigured(models []model.Model) {
	pterm.Info.Println("Verifying if all downloaded models are configured...")

	// Get the list of configured model names
	configModelNames := model.GetNames(models)

	// TODO : is there a way to get back the options used for downloading a model/tokenizer?
	// Get the list of downloaded models
	downloadedModels := model.BuildModelsFromDevice()
	downloadedModelsNames := model.GetNames(downloadedModels)

	// Find missing models from configuration file
	missingModelNames := stringutil.SliceDifference(downloadedModelsNames, configModelNames)

	// Everything is fine, nothing more to do here
	if len(missingModelNames) == 0 {
		pterm.Info.Println("All downloaded models are well configured")
		return
	}

	// Asking the user to choose which model to remove or to configure
	modelNamesToConfigure, modelNamesToRemove := handleMissingModels(missingModelNames)

	// Removing models the user chose not to keep
	for _, name := range modelNamesToRemove {
		err := config.RemoveModelPhysically(name)
		if err != nil {
			continue
		}
	}

	// Retrieving the selected models to be configured
	modelsToConfigure := model.GetModelsByNames(downloadedModels, modelNamesToConfigure)

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err := config.AddModels(modelsToConfigure)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
}

// handleMissingModels handles all the models with no configuration
func handleMissingModels(missingModelNames []string) ([]string, []string) {
	// Ask user to select the models to delete/add to configuration file
	message := "These models weren't found in your configuration file and will be deleted. " +
		"Please select the models that you wish to conserve"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModels := ptermutil.DisplayInteractiveMultiselect(message, missingModelNames, []string{}, checkMark, false)
	modelsToDelete := stringutil.SliceDifference(missingModelNames, selectedModels)

	// Delete selected models
	if len(modelsToDelete) > 0 {
		// Ask user for confirmation to delete these models
		message = fmt.Sprintf(
			"Are you sure you want to delete these models [%s]?",
			strings.Join(modelsToDelete, ", "))
		yes := ptermutil.AskForUsersConfirmation(message)
		if !yes {
			// Re-asking the user since he changed his mind
			return handleMissingModels(missingModelNames)
		}
	}

	return selectedModels, modelsToDelete
}

// regenerateCode generates new default python code
func regenerateCode(models []model.Model) error {
	// TODO: modify this logic when code generator is completed
	pterm.Info.Println("Generating new default python code...")

	err := config.GenerateModelsPythonCode(models)
	if err != nil {
		return err
	}

	pterm.Success.Println("Python code generated")
	return nil
}
