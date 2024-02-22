package huggingface

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

var h *HuggingFace

func init() {
	h = NewHuggingFace(BaseUrl, "")
}

// TestGetModel tests GetModel
func TestGetModel(t *testing.T) {
	apiModel, err := h.GetModel("AiPorter/DialoGPT-small-Back_to_the_future")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertNotEqual(t, apiModel, nil, "The api call should've returned a model.")

	test.AssertEqual(t, apiModel.Name, "AiPorter/DialoGPT-small-Back_to_the_future",
		"Model's name should not be empty.")
}

// TestValidModel_Success tests ValidModel on valid model
func TestValidModel_Valid(t *testing.T) {
	valid, err := h.ValidModel("AiPorter/DialoGPT-small-Back_to_the_future")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, valid, true, "Model's name should be valid.")
}

// TestValidModel_Success tests ValidModel on valid model
func TestValidModel_NotValidModule(t *testing.T) {
	valid, err := h.ValidModel("Xibanya/sunset_city")
	test.AssertNotEqual(t, err, nil, "The api call shouldn't have passed.")
	test.AssertEqual(t, valid, false, "Model's module shouldn't be valid.")
}

// TestValidModel tests ValidModel on invalid model
func TestValidModel_NotValid(t *testing.T) {
	valid, err := h.ValidModel("not_valid")
	test.AssertNotEqual(t, err, nil, "The api call shouldn't have passed.")
	test.AssertEqual(t, valid, false, "Model's name shouldn't be valid.")
}
