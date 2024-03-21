package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

// Tests selectModel
func TestSelectModel(t *testing.T) {
	// Initialize models list
	var models model.Models
	models = append(models, model.Model{Name: "model1"})
	models = append(models, model.Model{Name: "model2"})
	models = append(models, model.Model{Name: "model3"})

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2"}
	app.SetUI(ui)

	// Select models
	selectedModel := selectModel(models)

	// Assertions
	test.AssertEqual(t, models[1].Name, selectedModel.Name)
}

// Tests selectTags
func TestSelectTags(t *testing.T) {
	// Initialize expected selections
	var expectedSelections []string
	expectedSelections = append(expectedSelections, string(huggingface.TextToImage))

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Select models
	selectedTags := selectTags()

	// Assertions
	test.AssertEqual(t, len(selectedTags), 1, "1 tag should be returned")
	test.AssertEqual(t, selectedTags[0], expectedSelections[0])
}

// Tests getModelsList
func TestGetModelsList(t *testing.T) {
	// Initialize tags list
	tags := []string{"tag1"}
	// Initialize existing expectedModels list
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	// Initialize expected expectedModels list
	var expectedModels model.Models
	expectedModels = append(expectedModels, model.Model{Name: "model2"})
	expectedModels = append(expectedModels, model.Model{Name: "model4"})
	expectedModels = append(expectedModels, model.Model{Name: "model5"})
	expectedModels = append(expectedModels, model.Model{Name: "model6"})
	expectedModels = append(expectedModels, model.Model{Name: "model7"})
	// Initialize api expectedModels list
	var hfModels huggingface.Models
	hfModels = append(hfModels, huggingface.Model{Name: "model1", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model3", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model4", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model5", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model6", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model7", LibraryName: huggingface.TRANSFORMERS})

	// Create huggingface mock
	huggingfaceInterface := mock.MockHuggingFace{GetModelsResult: hfModels}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get expectedModels list
	models, err := getModelsList(tags, existingModels)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), len(expectedModels))
	for i, currentModel := range models {
		test.AssertEqual(t, currentModel.Name, expectedModels[i].Name)
	}
}

// Tests getModelsList throws error on failed api call
func TestGetModelsList_Fail(t *testing.T) {
	// Initialize tags list
	tags := []string{"tag1"}
	// Initialize existing expectedModels list
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})

	// Create huggingface mock
	huggingfaceInterface := mock.MockHuggingFace{Error: fmt.Errorf("test")}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get expectedModels list
	models, err := getModelsList(tags, existingModels)

	// Assertions
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, len(models), 0)
}

// Tests downloadModel with addToBinary = true
func TestDownloadModel(t *testing.T) {
	// Init
	selectedModel := model.Model{Name: "model1", AddToBinaryFile: true}
	var downloaderArgs downloadermodel.Args
	var returnedModel downloadermodel.Model

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: returnedModel}
	app.SetDownloader(&downloader)

	// Download model
	downloadedModel, err := downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, downloadedModel.Name, selectedModel.Name)
}

// Tests downloadModel with addToBinary = false
func TestDownloadModel_OnlyConfiguration(t *testing.T) {
	// Init
	selectedModel := model.Model{Name: "model1"}
	var downloaderArgs downloadermodel.Args
	var returnedModel downloadermodel.Model

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: returnedModel}
	app.SetDownloader(&downloader)

	// Get model's config
	downloadedModel, err := downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, downloadedModel.Name, selectedModel.Name)
}

// Tests downloadModel failure
func TestDownloadModel_Fail(t *testing.T) {
	// Init
	selectedModel := model.Model{Name: "model1"}
	var downloaderArgs downloadermodel.Args

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Download model
	_, err := downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertNotEqual(t, err, nil)
}
