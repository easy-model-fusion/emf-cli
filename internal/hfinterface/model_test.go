package hfinterface

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// Tests GetModelById with valid module
func TestGetModelById_WithValidModule(t *testing.T) {
	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LibraryName: huggingface.DIFFUSERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get model by id
	fetchedModel, err := GetModelById("test")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, fetchedModel.Name, "test")
}

// Tests GetModelById with invalid module
func TestGetModelById_WithInvalidModule(t *testing.T) {
	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: "test"}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get model by id
	_, err := GetModelById("test")

	// Assertions
	test.AssertEqual(t, err.Error(), "downloading models from test library is not allowed")
}

// Tests GetModelsByPipelineTag
func TestGetModelsByPipelineTag(t *testing.T) {
	// Init
	var hfModels huggingface.Models
	hfModels = append(hfModels, huggingface.Model{Name: "model1", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model2", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model3", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model4", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model5", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model6", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model7", LibraryName: huggingface.DIFFUSERS})

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelsResult: hfModels}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get models by pipeline tag
	fetchedModels, err := GetModelsByPipelineTag(huggingface.TextToImage, 0)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(fetchedModels), 5)
}

// Tests GetModelsByMultiplePipelineTags
func TestGetModelsByMultiplePipelineTags(t *testing.T) {
	// Init
	var hfModels huggingface.Models
	hfModels = append(hfModels, huggingface.Model{Name: "model1", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model2", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model3", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model4", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model5", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model6", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model7", LibraryName: huggingface.DIFFUSERS})

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelsResult: hfModels}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get models by multiple pipeline tags
	fetchedModels, err := GetModelsByMultiplePipelineTags([]string{"tag1", "tag2"})

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(fetchedModels), 10)
}

// Tests getModelsByModules
func Test_getModelsByModules(t *testing.T) {
	// Init
	var hfModels huggingface.Models
	hfModels = append(hfModels, huggingface.Model{Name: "model1", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model2", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model3", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model4", LibraryName: huggingface.DIFFUSERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model5", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model6", LibraryName: "invalid"})
	hfModels = append(hfModels, huggingface.Model{Name: "model7", LibraryName: huggingface.DIFFUSERS})

	// Get models by modules
	filteredModels := getModelsByModules(hfModels)

	// Assertions
	test.AssertEqual(t, len(filteredModels), 5)
	test.AssertEqual(t, filteredModels[0].Name, "model1")
	test.AssertEqual(t, filteredModels[1].Name, "model2")
	test.AssertEqual(t, filteredModels[2].Name, "model3")
	test.AssertEqual(t, filteredModels[3].Name, "model4")
	test.AssertEqual(t, filteredModels[4].Name, "model7")
}
