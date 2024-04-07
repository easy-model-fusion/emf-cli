package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestGetModelsByPipelineTag_Success tests the GetModelsByPipelineTag method of the HuggingFace type.
// It initializes a HuggingFace instance and calls GetModelsByPipelineTag to retrieve models by pipeline tag.
// It asserts that the API call is successful, the expected number of models are returned, and that each model has a non-empty name.
func TestGetModelsByPipelineTag_Success(t *testing.T) {
	limit := 10
	h := NewHuggingFace(BaseUrl, "")
	models, err := h.GetModelsByPipelineTag(TextToImage, 10, "")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, len(models), limit, "The api call should've returned 10 models.")

	for _, apiModel := range models {
		test.AssertNotEqual(t, apiModel.Name, "", "Model's name should not be empty.")
	}
}

// TestGetModelsByPipelineTag_Failure tests the GetModelsByPipelineTag method of the HuggingFace type.
// It initializes a HuggingFace instance and calls GetModelsByPipelineTag to retrieve models by pipeline tag.
// It asserts that the API call fails and that no models are returned.
func TestGetModelsByPipelineTag_Failure(t *testing.T) {
	h := NewHuggingFace("% xw*cbadurl", "")
	models, err := h.GetModelsByPipelineTag(TextToImage, 10, "")
	test.AssertNotEqual(t, err, nil, "The api call should've failed.")
	test.AssertEqual(t, len(models), 0, "The api call should've returned 0 models.")
}

// TestGetModelById tests the GetModelById method of the HuggingFace type.
// It initializes a HuggingFace instance and calls GetModelById to retrieve a model by its ID.
// It asserts that the API call is successful, a model is returned, and that the model has a non-empty name matching the specified ID.
func TestGetModelById(t *testing.T) {
	h := NewHuggingFace(BaseUrl, "")
	apiModel, err := h.GetModelById("Xibanya/sunset_city", "")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertNotEqual(t, apiModel, nil, "The api call should've returned a model.")

	test.AssertEqual(t, apiModel.Name, "Xibanya/sunset_city", "Model's name should not be empty.")

	// set with bad url
	h = NewHuggingFace("% xw*cbadurl", "")
	apiModel, err = h.GetModelById("Xibanya/sunset_city", "")
	test.AssertNotEqual(t, err, nil, "The api call should've failed.")
	test.AssertEqual(t, apiModel, Model{}, "The api call should've returned an empty model.")
}
