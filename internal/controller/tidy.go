package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"path/filepath"
)

type TidyController struct{}

func (tc TidyController) RunTidy(yes bool, accessToken string) error {
	// get all models from config file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		app.UI().Error().Println(err.Error())
		return err
	}

	sdk.SendUpdateSuggestion()

	models, err := config.GetModels()
	if err != nil {
		app.UI().Error().Println(err.Error())
		return err
	}

	// Tidy the models configured but not physically present on the device
	app.UI().Info().Println("Verifying if all models are downloaded...")
	warningMessages, err := tc.tidyModelsConfiguredButNotDownloaded(models, accessToken)
	if len(warningMessages) > 0 {
		for _, warning := range warningMessages {
			app.UI().Warning().Println(warning)
		}
	}
	if err != nil {
		app.UI().Error().Println(err.Error())
		return err
	}

	// Tidy the models physically present on the device but not configured
	app.UI().Info().Println("Verifying if all downloaded models are configured...")
	warningMessages, err = tc.tidyModelsDownloadedButNotConfigured(models, yes, accessToken)
	if len(warningMessages) > 0 {
		for _, warning := range warningMessages {
			app.UI().Warning().Println(warning)
		}
	}
	if err != nil {
		app.UI().Error().Println(err.Error())
		return err
	}

	// Updating the models object since the configuration might have changed in between
	models, err = config.GetModels()
	if err != nil {
		app.UI().Error().Println(err.Error())
		return err
	}

	// Regenerate python code
	spinner := app.UI().StartSpinner("Generating python code...")
	err = tc.regenerateCode(models)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while generating python: %s", err))
		return err
	}
	spinner.Success()

	return nil
}

// tidyModelsConfiguredButNotDownloaded downloads any missing model and its missing tokenizers as well
func (tc TidyController) tidyModelsConfiguredButNotDownloaded(models model.Models, accessToken string) (warnings []string, err error) {
	// filter the models that should be added to binary
	models = models.FilterWithIsDownloadedTrue()

	// Search for the models that need to be downloaded
	var downloadedModels model.Models
	var failedModels []string

	// Tidying the configured but not downloaded models and also processing their tokenizers
	for _, current := range models {
		var (
			success, clean  bool
			warningMessages []string
		)
		warningMessages, success, clean, err = current.TidyConfiguredModel(accessToken)
		warnings = append(warnings, warningMessages...)
		if err != nil {
			return warnings, err
		}

		if !success {
			failedModels = append(failedModels, current.Name)
		} else if !clean {
			downloadedModels = append(downloadedModels, current)
		}

		continue
	}

	// Displaying the downloads that failed
	if len(failedModels) > 0 {
		warnings = append(warnings, fmt.Sprintf("The following models(s) couldn't be downloaded : %s", failedModels))
	}

	if len(downloadedModels) > 0 {
		// Add models to configuration file
		spinner := app.UI().StartSpinner("Writing models to configuration file...")
		err := config.AddModels(downloadedModels)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}
	return warnings, err
}

// tidyModelsDownloadedButNotConfigured configuring the downloaded models that aren't configured in the configuration file
// and then asks the user if he wants to delete them or add them to the configuration file
func (tc TidyController) tidyModelsDownloadedButNotConfigured(configModels model.Models, yes bool, accessToken string) (warnings []string, err error) {
	// Get the list of downloaded models
	downloadedModels := model.BuildModelsFromDevice(accessToken)

	// Building map for faster lookup
	mapConfigModels := configModels.Map()

	// Checking if every model is well configured
	var modelsToConfigure model.Models
	for _, current := range downloadedModels {

		// Checking if the downloaded model is already configured
		configModel, configured := mapConfigModels[current.Name]

		// Try to get model configuration
		if current.Module != "" {
			downloaderArgs := downloadermodel.Args{
				ModelName:     current.Name,
				ModelModule:   string(current.Module),
				DirectoryPath: app.DownloadDirectoryPath,
			}

			// Getting model class
			var success bool
			success, warnings, err = current.GetConfig(downloaderArgs)
			if err != nil {
				return warnings, err
			}
			if !success && current.Class == "" {
				current.Class = current.GetModuleAutoPipelineClassName()
			}
		}
		// Model not configured
		if !configured {

			// Asking for permission to configure
			configure := yes || app.UI().AskForUsersConfirmation(fmt.Sprintf("Model '%s' wasn't found in your "+
				"configuration file. Confirm to configure, otherwise it will be removed.", current.Name))

			// User chose to configure it
			if configure {
				modelsToConfigure = append(modelsToConfigure, current)
			} else {
				// User chose not to configure : removing the model
				modelPath := filepath.Join(app.DownloadDirectoryPath, current.Name)
				spinner := app.UI().StartSpinner(fmt.Sprintf("Removing model %s...", current.Name))
				err := config.RemoveItemPhysically(modelPath)
				if err != nil {
					spinner.Fail("failed to remove item")
					continue
				} else {
					spinner.Success()
				}
			}

			// Highest configuration possible : nothing more to do here
			continue
		}

		// If model is a transformer : checking tokenizers
		if current.Module == huggingface.TRANSFORMERS {

			// Building map for faster lookup
			mapConfigModelTokenizers := configModel.Tokenizers.Map()

			// Checking if every tokenizer is well configured
			var modelTokenizersToConfigure model.Tokenizers
			for _, tokenizer := range current.Tokenizers {

				// Checking if the downloaded tokenizer is already configured
				_, configured = mapConfigModelTokenizers[tokenizer.Class]

				// Tokenizer not configured
				if !configured {

					// Asking for permission to configure
					configure := yes || app.UI().AskForUsersConfirmation(fmt.Sprintf("Tokenizer '%s' for model '%s' wasn't found in your "+
						"configuration file. Confirm to configure, otherwise it will be removed.", tokenizer.Class, current.Name))

					// User chose to configure it
					if configure {
						modelTokenizersToConfigure = append(modelTokenizersToConfigure, tokenizer)
					} else {
						// User chose not to configure : removing the tokenizer
						tokenizerPath := filepath.Join(app.DownloadDirectoryPath, tokenizer.Path)
						spinner := app.UI().StartSpinner(fmt.Sprintf("Removing tokenizer %s...", tokenizer.Class))
						err := config.RemoveItemPhysically(tokenizerPath)
						if err != nil {
							spinner.Fail("failed to remove item")
							continue
						} else {
							spinner.Success()
						}
					}
				}
			}

			// If at least one tokenizer was configured
			if len(modelTokenizersToConfigure) > 0 {
				// Since model is already configured : adding missing tokenizers and reconfiguring the model
				// Note : there can't be any duplicated tokenizers in this case
				configModel.Tokenizers = append(configModel.Tokenizers, modelTokenizersToConfigure...)
				modelsToConfigure = append(modelsToConfigure, configModel)
			}
		}
	}

	if len(modelsToConfigure) > 0 {
		// Add models to configuration file
		spinner := app.UI().StartSpinner("Writing models to configuration file...")
		err := config.AddModels(modelsToConfigure)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}
	return warnings, err
}

// regenerateCode generates new default python code
func (tc TidyController) regenerateCode(models model.Models) error {
	// TODO: modify this logic when code generator is completed
	err := config.GenerateModelsPythonCode(models)
	if err != nil {
		return err
	}

	return nil
}
