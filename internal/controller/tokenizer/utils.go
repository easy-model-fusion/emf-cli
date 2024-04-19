package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

type SelectModelController struct {
}

// SelectTransformerModel displays a selector of models from which the user will choose to add to his project
func (ic SelectModelController) SelectTransformerModel(models model.Models) model.Model {
	// Build a selector with each model name
	availableModelNames := models.GetNames()

	message := "Please select the model for which to add tokenizers "
	selectedModelName := app.UI().DisplayInteractiveSelect(message, availableModelNames, true, 8)
	// Get newly selected model
	selectedModels := models.FilterWithNames([]string{selectedModelName})
	// Return newly selected model along with selected model name and no error
	return selectedModels[0]
}
