package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"io"
	"net/http"
	"net/url"
)

// GetModel from hugging face api by id
func (h HuggingFace) GetModel(id string) (model.Model, error) {
	var result model.Model

	getModelUrl, err := url.Parse(h.BaseUrl + modelEndpoint)
	if err != nil {
		return result, err
	}

	// Prepare API call
	q := getModelUrl.Query()
	q.Add("config", "config")
	q.Add("id", id)
	getModelUrl.RawQuery = q.Encode()

	// Execute API call
	response, err := h.Client.Get(getModelUrl.String())
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to fetch model. Status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	// Unmarshal API response
	var APIModelResponse []APIModelResponse
	if err = json.Unmarshal(body, &APIModelResponse); err != nil {
		return result, err
	}

	// Check response validity
	if len(APIModelResponse) > 1 {
		return result, fmt.Errorf("too many models returned")
	} else if len(APIModelResponse) == 0 {
		return result, fmt.Errorf("no model found with name = %v", id)
	}

	// Map API response to model.Model
	return MapAPIResponseToModelObj(APIModelResponse[0]), nil
}

// ValidModel checks if a model exists by id
func (h HuggingFace) ValidModel(id string) (bool, error) {
	_, err := h.GetModel(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

// MapAPIResponseToModelObj map a model API response to a model object
func MapAPIResponseToModelObj(response APIModelResponse) model.Model {
	var modelObj model.Model
	modelObj.Name = response.Name
	modelObj.PipelineTag = model.PipelineTag(response.PipelineTag)
	modelObj.Config.Module = response.LibraryName
	modelObj.Source = model.HUGGING_FACE
	return modelObj
}
