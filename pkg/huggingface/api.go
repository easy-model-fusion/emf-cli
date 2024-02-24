package huggingface

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseUrl = "https://huggingface.co/api"
const modelEndpoint = "/models"

type HuggingFace struct {
	BaseUrl string
	Client  *http.Client
}

// NewHuggingFace creates a new HuggingFace instance
func NewHuggingFace(baseUrl, proxyUrl string) *HuggingFace {
	client := &http.Client{}
	if pUrl, err := url.Parse(proxyUrl); err != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(pUrl),
		}
	}

	return &HuggingFace{
		BaseUrl: baseUrl,
		Client:  client,
	}
}

// Model Define a model to match the JSON response from the API
type Model struct {
	Name         string      `json:"modelId"`
	PipelineTag  PipelineTag `json:"pipeline_tag"`
	LibraryName  Module      `json:"library_name"`
	LastModified string      `json:"lastModified"`
}

// APIGet performs an HTTP GET request to the specified URL.
func (h HuggingFace) APIGet(getModelUrl *url.URL) ([]byte, error) {
	// Execute API call
	var response, err = h.Client.Get(getModelUrl.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch model. Status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}