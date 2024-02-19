package huggingface

import (
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/test"
	"testing"
)

func TestGetModels(t *testing.T) {
	limit := 10
	h := NewHuggingFace(BaseUrl, "")
	models, err := h.GetModels(model.TEXT_TO_IMAGE, 10)
	test.AssertEqual(t, err, nil, "The api call should've passed.")
	test.AssertEqual(t, len(models), limit, "The api call should've returned 10 models.")

	for _, apiModel := range models {
		test.AssertNotEqual(t, apiModel.Name, "", "Model's name should not be empty.")
	}
}
