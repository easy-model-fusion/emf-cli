package cmd

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

	// Add all missing models
	err = addMissingModels(models)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Fix missing model configurations
	err = missingModelConfiguration(models)
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

/*

TODO : when can a model be downloaded ?
commands : add, update, tidy

TODO : model add : DIFFUSERS
downloaded => confirmation to overwrite the model => download the model
!downloaded => download the model

TODO : model add : TRANSFORMERS
!modelDownloaded && !tokenizerDownloaded => download the model with tokenizer
!modelDownloaded && tokenizerDownloaded => confirmation to overwrite the tokenizer => download the model with or without the --skip tokenizer
modelDownloaded && !tokenizerDownloaded => confirmation to overwrite the model => download the tokenizer with or without the --skip model
modelDownloaded && tokenizerDownloaded => confirmation to overwrite the model => confirmation to overwrite the tokenizer => download the model/tokenizer with or without the --skip model/tokenizer

TODO : model update
!modelConfigured && !modelDownloaded => confirmation to add the model => download the model with or without the --skip tokenizer
!modelConfigured && modelDownloaded => confirmation to overwrite the model => download the model with or without the --skip tokenizer
modelConfigured && !modelDownloaded => confirmation to download the model => download the model with or without the --skip tokenizer
modelConfigured && modelDownloaded => confirmation to upload the model => download the model with or without the --skip tokenizer
tokenizersConfigured => multiselect those to reinstall => download the selected tokenizers by overwriting them

TODO : tidy
!modelConfigured && modelDownloaded => get on device models => confirmation to configure them (else removed from device)
!modelConfigured && tokenizers => get on device tokenizers => !configured => confirmation to configure them (else removed from device)
modelConfigured && !modelDownloaded => confirmation to download the model => download the model
modelConfigured && tokenizers => !downloaded => confirmation to download the tokenizers => download the tokenizers
*/

// addMissingModels adds the missing models from the list of configuration file models
func addMissingModels(models []model.Model) error {
	pterm.Info.Println("Verifying if all models are downloaded...")
	// filter the models that should be added to binary
	models = model.GetModelsWithAddToBinaryFileTrue(models)

	// Search for the models that need to be downloaded
	var downloadedModels []model.Model
	var failedModels []string
	var failedTokenizersForModels []string

	// Tidying the configured but not downloaded models and tokenizers
	for _, current := range models {

		// TODO : what if there is a correct custom path that the user provided? how about the tokenizer paths?
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
			if len(current.Tokenizers) != len(missingTokenizers) {
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
		if len(missingTokenizers) == 0 {

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

	// TODO : update configuration file

	return nil
}

// missingModelConfiguration finds the downloaded models that aren't configured in the configuration file
// and then asks the user if he wants to delete them or add them to the configuration file
func missingModelConfiguration(models []model.Model) error {
	pterm.Info.Println("Verifying if all downloaded models are configured...")
	// Get the list of downloaded model names
	downloadedModelNames, err := app.GetDownloadedModelNames()
	if err != nil {
		return err
	}

	// TODO : check for every tokenizer of the model if it's downloaded

	// Get the list of configured model names
	configModelNames := model.GetNames(models)
	// Find missing models from configuration file
	missingModelNames := stringutil.SliceDifference(downloadedModelNames, configModelNames)
	if len(missingModelNames) > 0 {
		err = handleModelsWithNoConfig(missingModelNames)
		if err != nil {
			return err
		}
	} else {
		pterm.Info.Println("All downloaded models are well configured")
	}

	return nil
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

// generateModelsConfig generates models configurations
func generateModelsConfig(modelNames []string) error {
	// initialize hugging face url
	app.InitHuggingFace(huggingface.BaseUrl, "")
	// get hugging face api
	huggingFace := app.H()

	var models []model.Model
	for _, modelName := range modelNames {
		// Search for the model in hugging face
		huggingfaceModel, err := huggingFace.GetModelById(modelName)
		var currentModel model.Model
		// If not found create model configuration with only model's name
		if err != nil {
			currentModel = model.Model{Name: modelName}
			currentModel.Source = model.CUSTOM
		} else {
			// Found : Map API response to model.Model
			currentModel = model.MapToModelFromHuggingfaceModel(huggingfaceModel)
		}
		currentModel.AddToBinaryFile = true
		currentModel.IsDownloaded = true
		currentModel = model.ConstructConfigPaths(currentModel)
		models = append(models, currentModel)
	}

	// Add models to the configuration file
	err := config.AddModels(models)
	if err != nil {
		return err
	}

	return nil
}

// handleModelsWithNoConfig handles all the models with no configuration
func handleModelsWithNoConfig(missingModelNames []string) error {
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
		if yes {
			// Delete models if confirmed
			for _, modelName := range modelsToDelete {
				err := config.RemoveModelPhysically(modelName)
				if err != nil {
					return err
				}
			}
			pterm.Success.Println("Deleted models", modelsToDelete)
		} else {
			return handleModelsWithNoConfig(missingModelNames)
		}
	}

	// Configure selected models
	if len(selectedModels) > 0 {
		// Add models' configurations to config file
		err := generateModelsConfig(selectedModels)
		if err != nil {
			return err
		}
		pterm.Success.Println("Added configurations for these models", selectedModels)
	}
	return nil
}
