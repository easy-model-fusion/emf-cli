package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/model"
	"io"
	"net/http"
	"net/url"
)

func GetModels(limit int, tag string, proxyURL *url.URL) ([]model.Model, error) {
	client := &http.Client{}
	if proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	apiURL := fmt.Sprintf("https://huggingface.co/api/models?config=config&pipeline_tag=%v&limit=%d", tag, limit)
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var models []model.Model
		if err := json.Unmarshal(body, &models); err != nil {
			return nil, err
		}

		return models, nil
	} else {
		return nil, fmt.Errorf("failed to fetch models. Status code: %d", response.StatusCode)
	}
}
