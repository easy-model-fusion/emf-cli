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

// Test RunModelUpdate
func TestRunModelUpdate_Success(t *testing.T) {
	// Init
	var models model.Models
	models = append(models, GetModel(1, "2021"))
	models = append(models, GetModel(2, "2022"))
	models = append(models, GetModel(3, "2022"))
	var args []string
	args = append(args, "model1")
	args = append(args, "model3")
	args = append(args, "model4")

	// Create hugging face mock
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: "."}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	RunModelUpdate(args)
	updatedModels, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected on getting all models")
	test.AssertEqual(t, len(updatedModels), 3)
	test.AssertEqual(t, updatedModels[0].Version, "2022")
	test.AssertEqual(t, updatedModels[1].Version, "2022")
	test.AssertEqual(t, updatedModels[2].Version, "2022")
}

// Test RunModelUpdate should fail
func TestRunModelUpdate_Fail(t *testing.T) {
	// Init
	var models model.Models
	models = append(models, GetModel(1, "2021"))
	models = append(models, GetModel(2, "2022"))
	models = append(models, GetModel(3, "2022"))
	var args []string
	args = append(args, "model1")
	args = append(args, "model3")
	args = append(args, "model4")

	// Create hugging face mock
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: false, UserInputResult: "."}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	RunModelUpdate(args)
	updatedModels, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected on getting all models")
	test.AssertEqual(t, len(updatedModels), 3)
	test.AssertEqual(t, updatedModels[0].Version, "2021")
	test.AssertEqual(t, updatedModels[1].Version, "2022")
	test.AssertEqual(t, updatedModels[2].Version, "2022")
}

// Tests processUpdate
func TestProcessUpdate(t *testing.T) {
	// Init
	var models model.Models
	models = append(models, GetModel(1, "2021"))
	models = append(models, GetModel(2, "2022"))
	models = append(models, GetModel(3, "2022"))
	var args []string
	args = append(args, "model1")
	args = append(args, "model3")
	args = append(args, "model4")

	// Create hugging face mock
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: "."}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	warningMessage, infoMessage, err := processUpdate(args)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warningMessage, "The following models(s) couldn't be found and were ignored : [model4]", "A warning is expected")
	test.AssertEqual(t, infoMessage, "The following model(s) are already up to date and were ignored : [model3]", "Information message expected")
}

// Tests processUpdate with no args
func TestProcessUpdate_WithNoArgs(t *testing.T) {
	// Init
	var models model.Models
	models = append(models, GetModel(1, "2021"))
	models = append(models, GetModel(2, "2022"))
	models = append(models, GetModel(3, "2022"))
	var args []string
	var expectedSelections []string
	expectedSelections = append(expectedSelections, "model1")
	expectedSelections = append(expectedSelections, "model3")

	// Create hugging face mock
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: ".", MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	warningMessage, infoMessage, err := processUpdate(args)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warningMessage, "", "No warning is expected")
	test.AssertEqual(t, infoMessage, "The following model(s) are already up to date and were ignored : [model3]", "Information message expected")
}

// Tests processUpdate with no models selected
func TestProcessUpdate_WithNoModelsSelected(t *testing.T) {
	// Init
	var models model.Models
	models = append(models, GetModel(1, "2021"))
	models = append(models, GetModel(2, "2022"))
	models = append(models, GetModel(3, "2022"))
	var args []string
	var expectedSelections []string

	// Create hugging face mock
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: ".", MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	warningMessage, infoMessage, err := processUpdate(args)

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected")
	test.AssertEqual(t, warningMessage, "", "No warning is expected")
	test.AssertEqual(t, infoMessage, "There is no models to be updated.", "Information message expected")
}

// Tests processUpdate with an error on loading configuration file
func TestProcessUpdate_WithErrorOnLoadingConfigurationFile(t *testing.T) {
	// Init
	var args []string

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: "."}
	app.SetUI(ui)

	// Process update
	warningMessage, infoMessage, err := processUpdate(args)

	// Assertions
	test.AssertNotEqual(t, err, nil, "An error is expected")
	test.AssertEqual(t, warningMessage, "", "No warning is expected")
	test.AssertEqual(t, infoMessage, "", "No information message expected")
}

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
	huggingFaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022"}}
	app.SetHuggingFace(&huggingFaceInterface)

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
	huggingFaceInterface := huggingface.MockHuggingFace{Error: fmt.Errorf("")}
	app.SetHuggingFace(&huggingFaceInterface)

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
		Name:         "model" + idStr,
		Source:       model.HUGGING_FACE,
		IsDownloaded: true,
		Version:      version,
	}
}
