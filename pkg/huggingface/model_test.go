package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

var h *HuggingFace

func init() {
	h = NewHuggingFace(BaseUrl, "")
}

// TestGetModelsByPipelineTag_Success tests the GetModelsByPipelineTag method of the HuggingFace type.
// It initializes a HuggingFace instance and calls GetModelsByPipelineTag to retrieve models by pipeline tag.
// It asserts that the API call is successful, the expected number of models are returned, and that each model has a non-empty name.
func TestGetModelsByPipelineTag_Success(t *testing.T) {
	limit := 10
	h := NewHuggingFace(BaseUrl, "")
	models, err := h.GetModelsByPipelineTag(TextToImage, 10)
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, len(models), limit, "The api call should've returned 10 models.")

	for _, apiModel := range models {
		test.AssertNotEqual(t, apiModel.Name, "", "Model's name should not be empty.")
	}
}

// TestGetModelById tests the GetModelById method of the HuggingFace type.
// It initializes a HuggingFace instance and calls GetModelById to retrieve a model by its ID.
// It asserts that the API call is successful, a model is returned, and that the model has a non-empty name matching the specified ID.
func TestGetModelById(t *testing.T) {
	apiModel, err := h.GetModelById("Xibanya/sunset_city")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertNotEqual(t, apiModel, nil, "The api call should've returned a model.")

	test.AssertEqual(t, apiModel.Name, "Xibanya/sunset_city", "Model's name should not be empty.")
}

// TestValidModel_Valid tests the ValidModel method of the HuggingFace type with a valid model ID.
// It initializes a HuggingFace instance and calls ValidModel to check if a model with the specified ID exists.
// It asserts that the API call is successful and that the model ID is valid.
func TestValidModel_Valid(t *testing.T) {
	valid, err := h.ValidModel("Xibanya/sunset_city")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, valid, true, "Model's name should be valid.")
}

// TestValidModel_NotValid tests the ValidModel method of the HuggingFace type with an invalid model ID.
// It initializes a HuggingFace instance and calls ValidModel to check if a model with the specified ID exists.
// It asserts that the API call fails and that the model ID is not valid.
func TestValidModel_NotValid(t *testing.T) {
	valid, err := h.ValidModel("not_valid")
	test.AssertNotEqual(t, err, nil, "The api call shouldn't have passed.")
	test.AssertEqual(t, valid, false, "Model's name shouldn't be valid.")
}
