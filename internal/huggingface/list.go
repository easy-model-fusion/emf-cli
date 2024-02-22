package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"io"
	"net/http"
	"net/url"
)

// GetModels from hugging face api
func (h HuggingFace) GetModels(tag string, limit int) ([]model.Model, error) {
	getModelsUrl, err := url.Parse(h.BaseUrl + modelEndpoint)
	if err != nil {
		return nil, err
	}

	// Prepare API call
	q := getModelsUrl.Query()
	q.Add("config", "config")
	q.Add("pipeline_tag", tag)
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	getModelsUrl.RawQuery = q.Encode()

	// Execute API call
	response, err := h.Client.Get(getModelsUrl.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models. Status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal API response
	var APIModelResponse []APIModelResponse
	if err = json.Unmarshal(body, &APIModelResponse); err != nil {
		return nil, err
	}

	// Map API responses to []model.Model
	var models []model.Model
	for _, item := range APIModelResponse {
		models = append(models, MapAPIResponseToModelObj(item))
	}

	// Filter models with handled modules
	models = getModelsByModules(models)

	return models, nil
}

// getModelsByModules filters a list of models and return only the models with handled module types
func getModelsByModules(models []model.Model) (returnedModels []model.Model) {
	modules := model.AllModules
	for _, currentModel := range models {
		if utils.SliceContainsItem(modules, currentModel.Config.Module) {
			returnedModels = append(returnedModels, currentModel)
		}
	}

	return returnedModels
}
