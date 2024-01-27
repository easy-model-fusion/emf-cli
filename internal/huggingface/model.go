package huggingface

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/client/internal/model"
	"io"
	"net/http"
	"net/url"
)

// GetModel from hugging face api
func GetModel(id string, proxyURL *url.URL) (*model.Model, error) {
	client := &http.Client{}
	if proxyURL != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	apiURL := fmt.Sprintf("https://huggingface.co/api/models?config=config&id=%v", id)
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

		if len(models) > 1 {
			return nil, fmt.Errorf("too many models returned")
		} else if len(models) == 0 {
			return nil, fmt.Errorf("no model found with name = %v", id)
		}

		return &models[0], nil
	} else {
		return nil, fmt.Errorf("failed to fetch models. Status code: %d", response.StatusCode)
	}
}

func ValidModel(id string) (bool, error) {
	apiModel, err := GetModel(id, nil)
	if err != nil {
		return false, err
	}
	if apiModel == nil {
		return false, nil
	}

	return true, nil
}
