package mock

import (
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
)

type MockHuggingFace struct {
	GetModelResult   huggingface.Model
	GetModelsResult  huggingface.Models
	ValidModelResult bool
	Error            error
}

func (hf *MockHuggingFace) GetModelsByPipelineTag(_ huggingface.PipelineTag, _ int) (huggingface.Models, error) {
	return hf.GetModelsResult, hf.Error
}
func (hf *MockHuggingFace) GetModelById(_ string) (huggingface.Model, error) {
	return hf.GetModelResult, hf.Error
}
func (hf *MockHuggingFace) ValidModel(_ string) (bool, error) {
	return hf.ValidModelResult, hf.Error
}
func (hf *MockHuggingFace) GetModelsByMultiplePipelineTags(_ []string) (allModelsWithTags huggingface.Models, err error) {
	return hf.GetModelsResult, hf.Error
}
