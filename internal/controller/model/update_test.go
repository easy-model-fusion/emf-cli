package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

// Tests selectModelsToUpdate
func TestSelectModelsToUpdate(t *testing.T) {
	// Initialize model names list
	var modelNames []string
	modelNames = append(modelNames, "model1")
	modelNames = append(modelNames, "model2")
	modelNames = append(modelNames, "model3")
	var expectedSelections []string
	expectedSelections = append(expectedSelections, "model1")
	expectedSelections = append(expectedSelections, "model3")

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Select models
	selectedModels := selectModelsToUpdate(modelNames)

	// Assertions
	test.AssertEqual(t, len(selectedModels), 2, "2 models should be returned")
	test.AssertEqual(t, selectedModels[0], expectedSelections[0])
	test.AssertEqual(t, selectedModels[1], expectedSelections[1])
}

// Tests selectModelsToUpdate with empty model names list
func TestSelectModelsToUpdate_WithNoModelNames(t *testing.T) {
	// Initialize model names list
	var modelNames []string

	// Select models
	selectedModels := selectModelsToUpdate(modelNames)

	// Assertions
	test.AssertEqual(t, len(selectedModels), 0, "Empty models list should be returned")
}

// Tests getUpdatableModels
func TestGetUpdatableModels(t *testing.T) {
	// Initialize models
	var hfModelsAvailable model.Models
	hfModelsAvailable = append(hfModelsAvailable, GetModel(1, "2021"))
	hfModelsAvailable = append(hfModelsAvailable, GetModel(2, "2022"))
	hfModelsAvailable = append(hfModelsAvailable, GetModel(3, "2022"))
	var modelNames []string
	modelNames = append(modelNames, "model1")
	modelNames = append(modelNames, "model2")
	modelNames = append(modelNames, "model3")
	modelNames = append(modelNames, "model4")

	// Create hugging face mock
	huggingFace := mock.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFace)

	// get updatable models
	modelsToUpdate, notFoundModelNames, upToDateModelNames := getUpdatableModels(modelNames, hfModelsAvailable)

	// Assertions
	test.AssertEqual(t, len(modelsToUpdate), 1)
	test.AssertEqual(t, modelsToUpdate[0].Name, modelNames[0])
	test.AssertEqual(t, len(notFoundModelNames), 1)
	test.AssertEqual(t, notFoundModelNames[0], modelNames[3])
	test.AssertEqual(t, len(upToDateModelNames), 2)
	test.AssertEqual(t, upToDateModelNames[0], modelNames[1])
	test.AssertEqual(t, upToDateModelNames[1], modelNames[2])
}

// Tests getUpdatableModels with model not found in hugging face
func TestGetUpdatableModels_WithModelNotFound(t *testing.T) {
	// Initialize models
	var hfModelsAvailable model.Models
	hfModelsAvailable = append(hfModelsAvailable, GetModel(1, "2021"))
	var modelNames []string
	modelNames = append(modelNames, "model1")

	// Create hugging face mock
	huggingFace := mock.MockHuggingFace{Error: fmt.Errorf("")}
	app.SetHuggingFace(&huggingFace)

	// get updatable models
	modelsToUpdate, notFoundModelNames, upToDateModelNames := getUpdatableModels(modelNames, hfModelsAvailable)

	// Assertions
	test.AssertEqual(t, len(modelsToUpdate), 0)
	test.AssertEqual(t, len(notFoundModelNames), 1)
	test.AssertEqual(t, notFoundModelNames[0], modelNames[0])
	test.AssertEqual(t, len(upToDateModelNames), 0)
}

// TestUpdateModels_Fail tests updateModels with succeeded update
func TestUpdateModels_Success(t *testing.T) {
	// Initialize models
	var models model.Models
	models = append(models, GetModel(1, "2021"))

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Update models
	err := updateModels(models)

	// Assertions
	test.AssertEqual(t, nil, err)
}

// TestUpdateModels_Fail tests updateModels with succeeded update
func TestUpdateModels_SuccessWithConfigurationAdded(t *testing.T) {
	// Initialize models
	var models model.Models
	models = append(models, GetModel(1, "2021"))

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := config.Load(".")
	test.AssertEqual(t, nil, err, "No error expected on loading config")

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Update models
	err = updateModels(models)

	// Assertions
	test.AssertEqual(t, nil, err)
}

// TestUpdateModels_Fail tests updateModels with failed update
func TestUpdateModels_Fail(t *testing.T) {
	// Initialize models
	var models model.Models
	models = append(models, GetModel(1, "2021"))

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: false}
	app.SetUI(ui)

	// Update models
	err := updateModels(models)

	// Assertions
	test.AssertEqual(t, err.Error(), "the following models(s) couldn't be downloaded : [model1]")
}

// GetModel initiates a basic model with an id as suffix
func GetModel(id int, version string) model.Model {
	idStr := fmt.Sprint(id)
	return model.Model{
		Name:    "model" + idStr,
		Version: version,
	}
}
