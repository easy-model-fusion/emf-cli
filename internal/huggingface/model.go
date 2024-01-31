package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/model"
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

	q := getModelUrl.Query()
	q.Add("config", "config")
	q.Add("id", id)
	getModelUrl.RawQuery = q.Encode()

	response, err := h.Client.Get(getModelUrl.String())
	if err != nil {
		return result, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("failed to fetch models. Status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}

	var models []model.Model
	if err = json.Unmarshal(body, &models); err != nil {
		return result, err
	}

	if len(models) > 1 {
		return result, fmt.Errorf("too many models returned")
	} else if len(models) == 0 {
		return result, fmt.Errorf("no model found with name = %v", id)
	}

	return models[0], nil
}

// ValidModel checks if a model exists by id
func (h HuggingFace) ValidModel(id string) (bool, error) {
	_, err := h.GetModel(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
