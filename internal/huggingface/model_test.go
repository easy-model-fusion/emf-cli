package huggingface

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

// TestGetModel tests GetModel
func TestGetModel(t *testing.T) {
	apiModel, err := GetModel("Xibanya/sunset_city", nil)
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertNotEqual(t, apiModel, nil, "The api call should've returned a model.")

	test.AssertEqual(t, apiModel.Name, "Xibanya/sunset_city", "Model's name should not be empty.")
}

// TestValidModel_Success tests ValidModel on valid model
func TestValidModel_Valid(t *testing.T) {
	valid, err := ValidModel("Xibanya/sunset_city")
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, valid, true, "Model's name should be valid.")
}

// TestValidModel tests ValidModel on invalid model
func TestValidModel_NotValid(t *testing.T) {
	valid, err := ValidModel("not_valid")
	test.AssertNotEqual(t, err, nil, "The api call shouldn't have passed.")
	test.AssertEqual(t, valid, false, "Model's name shouldn't be valid.")
}
