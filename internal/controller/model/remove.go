package modelcontroller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
)

// RunModelRemove runs the model remove command
func RunModelRemove(args []string, modelRemoveAllFlag bool) {
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := processRemove(args, modelRemoveAllFlag)

	// Display messages to user
	if warningMessage != "" {
		app.UI().Warning().Printfln(warningMessage)
	}

	if infoMessage != "" {
		app.UI().Info().Printfln(infoMessage)
	} else if err == nil {
		app.UI().Success().Printfln("Operation succeeded.")
	} else {
		app.UI().Error().Printfln("Operation failed.")
	}
}

// processRemove processes the remove model operation
func processRemove(args []string, modelRemoveAllFlag bool) (string, string, error) {
	// Load the configuration file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return "", "", err
	}

	// Get all configured models objects/names
	var selectedModels []string
	var models model.Models
	models, err = config.GetModels()
	if err != nil {
		return "", "", err
	}
	modelNames := models.GetNames()

	sdk.SendUpdateSuggestion()

	// No args, asks for model names
	if len(args) == 0 {
		// Get selected models from multiselect
		selectedModels = selectModelsToDelete(modelNames, modelRemoveAllFlag)
	} else {
		// Get the selected models from the args
		selectedModels = stringutil.SliceRemoveDuplicates(args)
	}

	// Remove the selected models
	return removeModels(models, selectedModels)
}

// selectModelsToDelete displays an interactive multiselect so the user can choose the models to remove
func selectModelsToDelete(modelNames []string, selectAllModels bool) []string {
	// Displays the multiselect only if the user has previously configured some models but hasn't selected all of them
	if !selectAllModels && len(modelNames) > 0 {
		checkMark := ui.Checkmark{Checked: app.UI().Red("x"), Unchecked: app.UI().Blue("-")}
		message := "Please select the model(s) to be deleted"
		modelNames = app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false, 8)
		app.UI().DisplaySelectedItems(modelNames)
	}
	return modelNames
}

// removeModels processes the selected models and removes them
func removeModels(models model.Models, selectedModels []string) (warning string, info string, err error) {
	if models.Empty() || len(selectedModels) == 0 {
		info = "There is no models to be removed."
	} else {
		warning, info, err = config.RemoveModelsByNames(models, selectedModels)
	}

	return warning, info, err
}
