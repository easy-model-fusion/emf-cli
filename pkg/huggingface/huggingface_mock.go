package huggingface

type MockHuggingFace struct {
	GetModelResult  Model
	GetModelsResult Models
	Error           error
}

func (hf *MockHuggingFace) GetModelsByPipelineTag(_ PipelineTag, _ int, _ string) (Models, error) {
	return hf.GetModelsResult, hf.Error
}
func (hf *MockHuggingFace) GetModelById(id string, _ string) (Model, error) {
	hf.GetModelResult.Name = id
	return hf.GetModelResult, hf.Error
}
