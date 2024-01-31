package app

import "github.com/easy-model-fusion/client/internal/huggingface"

var huggingFace *huggingface.HuggingFace

func H() *huggingface.HuggingFace {
	return huggingFace
}

func InitHuggingFace(baseUrl, proxyUrl string) {
	huggingFace = huggingface.NewHuggingFace(baseUrl, proxyUrl)
}
