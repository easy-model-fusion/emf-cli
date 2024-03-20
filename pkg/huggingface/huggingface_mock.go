package huggingface

type MockHuggingFace struct {
	GetModelResult   Model
	GetModelsResult  Models
	ValidModelResult bool
	Error            error
}

func (hf *MockHuggingFace) GetModelsByPipelineTag(_ PipelineTag, _ int) (Models, error) {
	return hf.GetModelsResult, hf.Error
}
func (hf *MockHuggingFace) GetModelById(id string) (Model, error) {
	hf.GetModelResult.Name = id
	return hf.GetModelResult, hf.Error
}
func (hf *MockHuggingFace) ValidModel(_ string) (bool, error) {
	return hf.ValidModelResult, hf.Error
}
func (hf *MockHuggingFace) GetModelsByMultiplePipelineTags(_ []string) (allModelsWithTags Models, err error) {
	return hf.GetModelsResult, hf.Error
}
