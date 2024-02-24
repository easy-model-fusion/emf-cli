package huggingface

import (
	"fmt"
	"net/url"
)

// GetModelsByPipelineTag from hugging face api by pipeline tag
func (h HuggingFace) GetModelsByPipelineTag(tag PipelineTag, limit int) ([]Model, error) {
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
	return h.APIGet(getModelsUrl)
}

// GetModelById from hugging face api by id
func (h HuggingFace) GetModelById(id string) (Model, error) {

	getModelUrl, err := url.Parse(h.BaseUrl + modelEndpoint)
	if err != nil {
		return Model{}, err
	}

	// Prepare API call
	q := getModelUrl.Query()
	q.Add("config", "config")
	q.Add("id", id)
	getModelUrl.RawQuery = q.Encode()

	// Execute API call
	models, err := h.APIGet(getModelUrl)

	// Check response validity
	if len(models) == 0 {
		return Model{}, fmt.Errorf("no model found with name = %v", id)
	}
	if len(models) > 1 {
		return Model{}, fmt.Errorf("too many models returned")
	}

	return models[0], nil
}

// ValidModel checks if a model exists by id
func (h HuggingFace) ValidModel(id string) (bool, error) {
	_, err := h.GetModelById(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
