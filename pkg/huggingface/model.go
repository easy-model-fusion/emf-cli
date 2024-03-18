package huggingface

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetModelsByPipelineTag from hugging face api by pipeline tag
func (h huggingFace) GetModelsByPipelineTag(tag PipelineTag, limit int) ([]Model, error) {
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
		return []Model{}, err
	}

	// Unmarshal API response
	var models []Model
	if err = json.Unmarshal(response, &models); err != nil {
		return []Model{}, err
	}

	// Execute API call
	return models, err
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
