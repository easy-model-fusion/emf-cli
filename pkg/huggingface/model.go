package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"net/url"
)

// GetModelsByPipelineTag from hugging face api by pipeline tag
func (h huggingFace) GetModelsByPipelineTag(tag PipelineTag, limit int) (Models, error) {
	getModelsUrl, err := url.Parse(h.BaseUrl + modelEndpoint)
	if err != nil {
		return nil, err
	}

	// Prepare API call
	q := getModelsUrl.Query()
	q.Add("config", "config")
	q.Add("pipeline_tag", string(tag))
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	getModelsUrl.RawQuery = q.Encode()

	// Execute API call
	response, err := h.apiGet(getModelsUrl)
	if err != nil {
		return Models{}, err
	}

	// Unmarshal API response
	var models Models
	if err = json.Unmarshal(response, &models); err != nil {
		return Models{}, err
	}

	// Execute API call
	return getModelsByModules(models), err
}

// GetModelById from hugging face api by id
func (h huggingFace) GetModelById(id string) (Model, error) {

	getModelUrl, err := url.Parse(h.BaseUrl + modelEndpoint + "/" + id)
	if err != nil {
		return Model{}, err
	}

	// Execute API call
	response, err := h.apiGet(getModelUrl)
	if err != nil {
		return Model{}, err
	}

	// Unmarshal API response
	var model Model
	if err = json.Unmarshal(response, &model); err != nil {
		return Model{}, err
	}

	// Verify if the library is compatible
	modules := AllModulesString()
	if !stringutil.SliceContainsItem(modules, string(model.LibraryName)) {
		return Model{}, fmt.Errorf("downloading models from %s library is not allowed", model.LibraryName)
	}

	return model, nil
}

// ValidModel checks if a model exists by id
func (h huggingFace) ValidModel(id string) (bool, error) {
	_, err := h.GetModelById(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetModelsByMultiplePipelineTags get the list of models with given types
func (h huggingFace) GetModelsByMultiplePipelineTags(tags []string) (allModelsWithTags Models, err error) {
	// Get list of models with current tags
	for _, tag := range tags {
		huggingfaceModels, err := h.GetModelsByPipelineTag(PipelineTag(tag), 0)
		if err != nil {
			return Models{}, fmt.Errorf("error while calling api endpoint")
		}
		allModelsWithTags = append(allModelsWithTags, huggingfaceModels...)
	}

	return allModelsWithTags, err
}

// getModelsByModules filters a list of models and return only the models with handled module types
func getModelsByModules(models []Model) (returnedModels []Model) {
	modules := AllModulesString()
	for _, currentModel := range models {
		if stringutil.SliceContainsItem(modules, string(currentModel.LibraryName)) {
			returnedModels = append(returnedModels, currentModel)
		}
	}

	return returnedModels
}
