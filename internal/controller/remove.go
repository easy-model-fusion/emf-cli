package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/pterm/pterm"
)

// RunModelRemove runs the model remove command
func RunModelRemove(args []string, modelRemoveAllFlag bool) {
	// Process remove operation with given arguments
	warningMessage, infoMessage, err := processRemove(args, modelRemoveAllFlag)

	// Display messages to user
	if warningMessage != "" {
		pterm.Warning.Printfln(warningMessage)
	}

	if infoMessage != "" {
		pterm.Info.Printfln(infoMessage)
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
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
	modelNames := models.GetNames()
	if err != nil {
		return "", "", err
	}

	sdk.SendUpdateSuggestion()

	// No args, asks for model names
	if len(args) == 0 {
		// Get selected models from multiselect
		selectedModels = selectModelsToDelete(modelNames, modelRemoveAllFlag)
	} else {
		// Get the selected models from the args
		selectedModels = args
	}

	// Remove the selected models
	return removeModels(models, selectedModels)
}

// selectModelsToDelete displays an interactive multiselect so the user can choose the models to remove
func selectModelsToDelete(modelNames []string, selectAllModels bool) []string {
	// Displays the multiselect only if the user has previously configured some models but hasn't selected all of them
	if !selectAllModels && len(modelNames) > 0 {
		checkMark := ui.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
		message := "Please select the model(s) to be deleted"
		modelNames = app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false)
		app.UI().DisplaySelectedItems(modelNames)
	}
	return modelNames
}

// removeModels processes the selected models and removes them
func removeModels(models model.Models, selectedModels []string) (warning string, info string, err error) {
	if models.Empty() || len(selectedModels) == 0 {
		info = "There is no selected models to be removed."
	} else {
		warning, info, err = config.RemoveModelsByNames(models, selectedModels)
	}

	return warning, info, err
}
