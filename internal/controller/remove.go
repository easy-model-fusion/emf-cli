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
	// Remove the selected models
	infoMessage, err := processRemove(args, modelRemoveAllFlag)

	if infoMessage != "" {
		pterm.Success.Printfln(infoMessage)
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func processRemove(args []string, modelRemoveAllFlag bool) (string, error) {
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return "", err
	}

	var selectedModels []string
	var models model.Models
	models, err = config.GetModels()
	modelNames := models.GetNames()
	if err != nil {
		return "", err
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

func selectModelsToDelete(modelNames []string, selectAllModels bool) []string {
	if !selectAllModels && len(modelNames) > 0 {
		checkMark := ui.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
		message := "Please select the model(s) to be deleted"
		modelNames = app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false)
		app.UI().DisplaySelectedItems(modelNames)
	}
	return modelNames
}

func removeModels(models model.Models, selectedModels []string) (message string, err error) {
	if models.Empty() {
		message = "There is no models to be removed."
	}

	err = config.RemoveModelsByNames(models, selectedModels)

	return message, err
}
