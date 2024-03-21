package hfinterface

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

func TestGetModelById(t *testing.T) {
	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: "test"}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get model by id
	_, err := GetModelById("test")

	// Assertions
	test.AssertEqual(t, err.Error(), "downloading models from test library is not allowed")
}
