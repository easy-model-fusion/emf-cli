package huggingface

import (
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

// APIModelResponse Define a struct to match the JSON response
type APIModelResponse struct {
	Name        string `json:"modelId"`
	PipelineTag string `json:"pipeline_tag"`
	LibraryName string `json:"library_name"`
}
