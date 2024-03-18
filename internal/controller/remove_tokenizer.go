package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/pterm/pterm"
)

// TokenizerRemoveCmd runs the model remove command
func TokenizerRemoveCmd(args []string) {
	
}

func selectModel(currentModels []model.Model) []string {
	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range currentModels {
		modelNames = append(modelNames, item.Name)
	}

	checkMark := ui.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	message := "Please select a model"
	modelsToChoice := app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false)
	app.UI().DisplaySelectedItems(modelsToChoice)
	return modelsToChoice
}
