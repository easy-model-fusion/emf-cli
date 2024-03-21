package hfinterface

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

// GetModelsByPipelineTag from hugging face api by pipeline tag
func GetModelsByPipelineTag(tag huggingface.PipelineTag, limit int) (huggingface.Models, error) {
	// Get models from api
	models, err := app.H().GetModelsByPipelineTag(tag, limit)

	// Filter models on compatible modules
	return getModelsByModules(models), err
}

// GetModelById from hugging face api by id
func GetModelById(id string) (huggingface.Model, error) {
	// Get model from api
	model, err := app.H().GetModelById(id)
	if err != nil {
		return huggingface.Model{}, err
	}

	// Verify if the library is compatible
	modules := huggingface.AllModulesString()
	if !stringutil.SliceContainsItem(modules, string(model.LibraryName)) {
		return huggingface.Model{}, fmt.Errorf("downloading models from %s library is not allowed", model.LibraryName)
	}

	return model, err
}

// GetModelsByMultiplePipelineTags get the list of models with given types
func GetModelsByMultiplePipelineTags(tags []string) (allModelsWithTags huggingface.Models, err error) {
	// Get list of models with current tags
	for _, tag := range tags {
		huggingfaceModels, err := GetModelsByPipelineTag(huggingface.PipelineTag(tag), 0)
		if err != nil {
			return huggingface.Models{}, fmt.Errorf("error while calling api endpoint")
		}
		allModelsWithTags = append(allModelsWithTags, huggingfaceModels...)
	}

	return allModelsWithTags, err
}

// getModelsByModules filters a list of models and return only the models with handled module types
func getModelsByModules(models huggingface.Models) (returnedModels huggingface.Models) {
	modules := huggingface.AllModulesString()
	for _, currentModel := range models {
		if stringutil.SliceContainsItem(modules, string(currentModel.LibraryName)) {
			returnedModels = append(returnedModels, currentModel)
		}
	}

	return returnedModels
}
