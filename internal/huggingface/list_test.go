package huggingface

import (
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/test"
	"testing"
)

func TestGetModels(t *testing.T) {
	limit := 10
	models, err := GetModels(limit, model.TEXT_TO_IMAGE, nil)
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, len(models), limit, "The api call should've passed.")

	for _, apiModel := range models {
		test.AssertNotEqual(t, apiModel.Name, "", "Model's name should not be empty.")
	}
}
