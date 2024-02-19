package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/model"
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

	q := getModelsUrl.Query()
	q.Add("config", "config")
	q.Add("pipeline_tag", tag)
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	getModelsUrl.RawQuery = q.Encode()

	response, err := h.Client.Get(getModelsUrl.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models. Status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var models []model.Model
	if err = json.Unmarshal(body, &models); err != nil {
		return nil, err
	}

	return models, nil
}
