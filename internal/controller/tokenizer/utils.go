package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type SelectModelController struct {
}

// SelectTransformerModel displays a selector of models from which the user will choose to add to his project
func (ic SelectModelController) SelectTransformerModel(models model.Models, configModelsMap map[string]model.Model) model.Model {
	// Build a selector with each model name
	availableModelNames := models.GetNames()

	// List of models that accept tokenizers
	var compatibleModels []string
	// Check for valid tokenizers
	for _, modelName := range availableModelNames {
		module := configModelsMap[modelName]
		if module.Module == huggingface.TRANSFORMERS {
			compatibleModels = append(compatibleModels, modelName)
		}
	}
	message := "Please select the model for which to add tokenizers "
	selectedModelName := app.UI().DisplayInteractiveSelect(message, compatibleModels, true, 8)
	// Get newly selected model
	selectedModels := models.FilterWithNames([]string{selectedModelName})
	// Return newly selected model along with selected model name and no error
	return selectedModels[0]
}
