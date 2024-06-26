package huggingface

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseUrl = "https://huggingface.co/api"
const modelEndpoint = "/models"

type HuggingFace interface {
	GetModelsByPipelineTag(tag PipelineTag, limit int, authorizationKey string) (Models, error)
	GetModelById(id string, authorizationKey string) (Model, error)
}

type huggingFace struct {
	BaseUrl string
	Client  *http.Client
}

// NewHuggingFace creates a new HuggingFace instance
func NewHuggingFace(baseUrl, proxyUrl string) HuggingFace {
	client := &http.Client{}
	if proxyUrl != "" {
		if pUrl, err := url.Parse(proxyUrl); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(pUrl),
			}
		}
	}

	return &huggingFace{
		BaseUrl: baseUrl,
		Client:  client,
	}
}

// Models Define a list of models to match the JSON response from the API
type Models []Model

// Model Define a model to match the JSON response from the API
type Model struct {
	Name         string      `json:"modelId"`
	PipelineTag  PipelineTag `json:"pipeline_tag"`
	LibraryName  Module      `json:"library_name"`
	LastModified string      `json:"lastModified"`
}

// apiGet performs an HTTP GET request to the specified URL.
func (h huggingFace) apiGet(getModelUrl *url.URL, authorizationKey string) ([]byte, error) {
	// Create http request
	req, err := http.NewRequest("GET", getModelUrl.String(), nil)
	// Add authorization key when needed
	if authorizationKey != "" {
		req.Header.Set("Authorization", "Bearer "+authorizationKey)
	}
	if err != nil {
		return nil, err
	}

	// Execute API call
	response, err := h.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		// Read response body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		return nil, fmt.Errorf("failed to fetch model. Status code: %s\n%s", response.Status, body)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
