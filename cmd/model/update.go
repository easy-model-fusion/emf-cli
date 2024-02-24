package cmdmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
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

	// Keep the models from huggingface
	hfModels := model.GetModelsWithSourceHuggingface(configModels)

	// Only keep the downloaded models : those available for an update
	hfModelsAvailable := model.GetModelsWithIsDownloadedTrue(hfModels)
	hfModelAvailableNames := model.GetNames(hfModelsAvailable)

	var selectedModelNames []string
	var notDownloadedModelNames []string

	// Get models to update
	if len(args) == 0 {
		// No argument provided : multiselect among the huggingface models that are installed
		message := "Please select the model(s) to be updated"
		checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
		selectedModelNames = ptermutil.DisplayInteractiveMultiselect(message, hfModelAvailableNames, checkMark, true)
	} else {
		// Remove all the duplicates
		args = stringutil.SliceRemoveDuplicates(args)

		// Checking if the inputted models are effectively downloaded
		for _, name := range args {
			if stringutil.SliceContainsItem(hfModelAvailableNames, name) {
				selectedModelNames = append(selectedModelNames, name)
			} else {
				notDownloadedModelNames = append(notDownloadedModelNames, name)
			}
		}
	}

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

		// TODO : get model by name from map
		var configModel model.Model

		// See if a new version is available
		if configModel.Version == modelMapped.Version {
			updatedModelNames = append(updatedModelNames, name)
		} else {
			// TODO : verify which model to add to the list
			modelsToUpdate = append(modelsToUpdate, configModel)
		}
	}

	// Indicate the models that couldn't be found
	if len(notFoundModelNames) > 0 {
		pterm.Warning.Printfln(fmt.Sprintf("The following models(s) couldn't be found "+
			"and will be ignored : %s", notFoundModelNames))
	}
	// Indicate the models that have yet to be downloaded
	if len(notDownloadedModelNames) > 0 {
		pterm.Warning.Printfln(fmt.Sprintf("The following models(s) have yet to be installed "+
			"and will be ignored : %s", notDownloadedModelNames))
	}
	// Indicate the models that are already up-to-date
	if len(updatedModelNames) > 0 {
		pterm.Warning.Printfln(fmt.Sprintf("The following model(s) are already up to date "+
			"and will be ignored : %s", updatedModelNames))
	}

	// Processing all the valid models for an update
	for _, current := range modelsToUpdate {

		// TODO: check if already downloaded
		downloaded := true
		overwrite := false

		if downloaded {
			// TODO : ask confirmation to overwrite
			pterm.Info.Println(fmt.Sprintf("A new version for '%s' is available."))
			if !overwrite {
				// User did not want to overwrite the model : skipping to the next one
				continue
			}

			if current.Module == huggingface.TRANSFORMERS {
				// TODO : multiselect for the tokenizers
				pterm.Info.Println("Select to keep and re-download (others will be removed)")
			}
		}

	}

	// TODO : if any updates : update the configuration file

}
