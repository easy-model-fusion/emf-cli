package test

import (
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type MockHuggingFace struct {
	GetModelResult   huggingface.Model
	GetModelsResult  huggingface.Models
	ValidModelResult bool
}

func (hf *MockHuggingFace) GetModelsByPipelineTag(_ huggingface.PipelineTag, _ int) (huggingface.Models, error) {
	return hf.GetModelsResult, nil
}
func (hf *MockHuggingFace) GetModelById(_ string) (huggingface.Model, error) {
	return hf.GetModelResult, nil
}
func (hf *MockHuggingFace) ValidModel(_ string) (bool, error) {
	return hf.ValidModelResult, nil
}
