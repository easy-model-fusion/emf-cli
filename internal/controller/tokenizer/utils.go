package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

// selectModel displays a selector of models from which the user will choose to add to his project
func selectModel(models model.Models, configModelsMap map[string]model.Model) (model.Model, string, error) {
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
	// If no compatible models are found, return an error
	if len(compatibleModels) == 0 {
		return model.Model{}, "no compatible models found",
			fmt.Errorf("only transformers models have tokenizers")
	}
	message := "Please select the model for which to modify tokenizers "
	selectedModelName := app.UI().DisplayInteractiveSelect(message, compatibleModels, true, 8)
	// Get newly selected model
	selectedModels := models.FilterWithNames([]string{selectedModelName})
	if len(selectedModels) == 0 {
		return model.Model{}, "please select a model",
			fmt.Errorf("please select a model")
	}
	// Return newly selected model along with selected model name and no error
	return selectedModels[0], selectedModelName, nil
}
