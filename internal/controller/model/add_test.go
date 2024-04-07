package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/dotenv"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

// Tests selectModel
func TestSelectModel(t *testing.T) {
	// Initialize the controller
	ac := AddController{}
	// Initialize models list
	var models model.Models
	models = append(models, model.Model{Name: "model1"})
	models = append(models, model.Model{Name: "model2"})
	models = append(models, model.Model{Name: "model3"})

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2"}
	app.SetUI(ui)

	// Select models
	selectedModel := ac.selectModel(models)

	// Assertions
	test.AssertEqual(t, models[1].Name, selectedModel.Name)
}

// Tests selectTags
func TestSelectTags(t *testing.T) {
	// Initialize the controller
	ac := AddController{}
	// Initialize expected selections
	var expectedSelections []string
	expectedSelections = append(expectedSelections, string(huggingface.TextToImage))

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Select models
	selectedTags := ac.selectTags()

	// Assertions
	test.AssertEqual(t, len(selectedTags), 1, "1 tag should be returned")
	test.AssertEqual(t, selectedTags[0], expectedSelections[0])
}

// Tests getModelsList
func TestGetModelsList(t *testing.T) {
	// Initialize the controller
	ac := AddController{}
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
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelsResult: hfModels}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get expectedModels list
	models, err := ac.getModelsList(tags, existingModels, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), len(expectedModels))
	for i, currentModel := range models {
		test.AssertEqual(t, currentModel.Name, expectedModels[i].Name)
	}
}

// Tests getModelsList throws error on failed api call
func TestGetModelsList_Fail(t *testing.T) {
	// Initialize the controller
	ac := AddController{}
	// Initialize tags list
	tags := []string{"tag1"}
	// Initialize existing expectedModels list
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{Error: fmt.Errorf("test")}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get expectedModels list
	models, err := ac.getModelsList(tags, existingModels, "")

	// Assertions
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, len(models), 0)
}

// Tests downloadModel with addToBinary = true
func TestDownloadModel(t *testing.T) {
	// Init
	ac := AddController{}
	selectedModel := model.Model{Name: "model1", AddToBinaryFile: true}
	var downloaderArgs downloadermodel.Args
	var returnedModel downloadermodel.Model

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: returnedModel}
	app.SetDownloader(&downloader)

	// Download model
	downloadedModel, err := ac.downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, downloadedModel.Name, selectedModel.Name)
}

// Tests downloadModel with addToBinary = false
func TestDownloadModel_OnlyConfiguration(t *testing.T) {
	// Init
	ac := AddController{}
	selectedModel := model.Model{Name: "model1"}
	var downloaderArgs downloadermodel.Args
	var returnedModel downloadermodel.Model

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: returnedModel}
	app.SetDownloader(&downloader)

	// Get model's config
	downloadedModel, err := ac.downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, downloadedModel.Name, selectedModel.Name)
}

// Tests downloadModel failure
func TestDownloadModel_Fail(t *testing.T) {
	// Init
	ac := AddController{}
	selectedModel := model.Model{Name: "model1"}
	var downloaderArgs downloadermodel.Args

	// Create downloader mock
	downloader := mock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Download model
	_, err := ac.downloadModel(selectedModel, downloaderArgs)

	// Assertions
	test.AssertNotEqual(t, err, nil)
}

// Tests getRequestedModel with valid model passed in arguments
func TestGetRequestedModel_WithValidArg(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model2"}
	expectedModel := model.Model{Name: "model2"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get requested model
	requestedModel, err := ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, requestedModel.Name, expectedModel.Name)
}

// Tests getRequestedModel with valid model passed in arguments and single file enabled
func TestGetRequestedModel_WithSingleFile(t *testing.T) {
	// Init
	ac := AddController{
		SingleFile: true,
	}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model2"}
	expectedModel := model.Model{Name: "model2"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Get requested model
	requestedModel, err := ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, requestedModel.Name, expectedModel.Name)
	test.AssertEqual(t, requestedModel.Source, model.CUSTOM)
	test.AssertEqual(t, requestedModel.AddToBinaryFile, true)
	test.AssertEqual(t, requestedModel.IsDownloaded, true)
}

// Tests getRequestedModel with existing model requested
func TestGetRequestedModel_WithInvalidArg(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model1"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get requested model
	_, err = ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err.Error(), "the following model already exist and will be ignored : model1")
}

// Tests getRequestedModel with model not found
func TestGetRequestedModel_WithModelNotFound(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model2"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{Error: fmt.Errorf("test")}
	app.SetHuggingFace(&huggingfaceInterface)

	// Get requested model
	_, err = ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err.Error(), "Model model2 not valid : test")
}

// Tests getRequestedModel with no arguments
func TestGetRequestedModel_WithNoArgs(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	var args []string
	selectedTags := []string{"tag1"}
	// Initialize api expectedModels list
	var hfModels huggingface.Models
	hfModels = append(hfModels, huggingface.Model{Name: "model1", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model3", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model4", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model5", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model6", LibraryName: huggingface.TRANSFORMERS})
	hfModels = append(hfModels, huggingface.Model{Name: "model7", LibraryName: huggingface.TRANSFORMERS})

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelsResult: hfModels}
	app.SetHuggingFace(&huggingfaceInterface)

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2", MultiselectResult: selectedTags}
	app.SetUI(ui)

	// Get requested model
	requestedModel, err := ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, requestedModel.Name, "model2")
}

// Tests getRequestedModel with no arguments and hugging face models fetch error
func TestGetRequestedModel_WithNoArgsWithFailedModelsFetch(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	var args []string
	selectedTags := []string{"tag1"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{Error: fmt.Errorf("")}
	app.SetHuggingFace(&huggingfaceInterface)

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2", MultiselectResult: selectedTags}
	app.SetUI(ui)

	// Get requested model
	_, err = ac.getRequestedModel(args, "")

	// Assertions
	test.AssertNotEqual(t, err, nil)
}

// Tests getRequestedModel with no arguments and no tags selected by the user
func TestGetRequestedModel_WithNoArgsAndNoTags(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	var args []string
	var selectedTags []string

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model2", MultiselectResult: selectedTags}
	app.SetUI(ui)

	// Get requested model
	requestedModel, err := ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, requestedModel.Name, "")
}

// Tests getRequestedModel with more than 1 argument
func TestGetRequestedModel_WithTooManyArgs(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model2", "model4"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Get requested model
	_, err = ac.getRequestedModel(args, "")

	// Assertions
	test.AssertEqual(t, err.Error(), "you can enter only one model at a time")
}

// Tests getRequestedModel with invalid configuration path
func TestGetRequestedModel_WithInvalidConfigPath(t *testing.T) {
	// Initialize the controller
	ac := AddController{}
	// Create mock UI
	ui := mock.MockUI{UserInputResult: "invalid"}
	app.SetUI(ui)

	// Get requested model
	_, err := ac.getRequestedModel([]string{}, "")

	// Assertions
	test.AssertNotEqual(t, err, nil)
}

// Tests process add for single file
func TestProcessAdd_SingleFile(t *testing.T) {
	// Init
	ac := AddController{
		SingleFile: true,
	}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.HUGGING_FACE, Class: "test"})
	existingModels = append(existingModels, model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.HUGGING_FACE, Class: "test"})
	downloaderArgs := downloadermodel.Args{
		ModelModule: string(huggingface.DIFFUSERS),
		ModelClass:  "Test",
		ModelName:   "model2",
	}
	selectedModel := model.Model{Name: "model2", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.CUSTOM, Class: "test", Path: "test.safetensors"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// create a sample file
	err = os.WriteFile("test.safetensors", []byte("test"), 0644)
	test.AssertEqual(t, err, nil)

	// Process add
	warning, err := ac.processAdd(selectedModel, downloaderArgs)
	test.AssertEqual(t, err, nil)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, warning, "")
	test.AssertEqual(t, len(models), 3)
	test.AssertEqual(t, models[2].Name, "model2")
}

// Tests process add
func TestProcessAdd_HuggingFace(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.HUGGING_FACE, Class: "test"})
	existingModels = append(existingModels, model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.HUGGING_FACE, Class: "test"})
	downloaderArgs := downloadermodel.Args{}
	selectedModel := model.Model{Name: "model2", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Source: model.HUGGING_FACE, Class: "test"}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	warning, err := ac.processAdd(selectedModel, downloaderArgs)
	test.AssertEqual(t, err, nil)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, warning, "")
	test.AssertEqual(t, len(models), 3)
	test.AssertEqual(t, models[2].Name, "model2")
}

// Tests process add with access token
func TestProcessAdd_WithAccessToken(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	existingModels = append(existingModels, model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	downloaderArgs := downloadermodel.Args{AccessToken: "testToken"}
	selectedModel := model.Model{Name: "model2", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	warning, err := ac.processAdd(selectedModel, downloaderArgs)
	test.AssertEqual(t, err, nil)
	token, err := dotenv.GetEnvValue("ACCESS_TOKEN_MODEL2")
	test.AssertEqual(t, err, nil)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, warning, "")
	test.AssertEqual(t, len(models), 3)
	test.AssertEqual(t, models[2].Name, "model2")
	test.AssertEqual(t, token, "testToken")
}

// Tests process add with invalid model
func TestProcessAdd_WithInvalidModel(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	existingModels = append(existingModels, model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	downloaderArgs := downloadermodel.Args{}
	selectedModel := model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	warning, err := ac.processAdd(selectedModel, downloaderArgs)
	test.AssertEqual(t, err, nil)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, warning, "Model 'model3' is already configured")
	test.AssertEqual(t, len(models), 2)
}

// Tests process add with failed download
func TestProcessAdd_WithFailedDownload(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	existingModels = append(existingModels, model.Model{Name: "model3", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE})
	downloaderArgs := downloadermodel.Args{}
	selectedModel := model.Model{Name: "model2", PipelineTag: huggingface.TextToImage, Module: huggingface.DIFFUSERS, Class: "test", Source: model.HUGGING_FACE}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Process add
	warning, err := ac.processAdd(selectedModel, downloaderArgs)
	test.AssertNotEqual(t, err, nil)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, warning, "")
	test.AssertEqual(t, len(models), 2)
}

// Tests Run
func TestAddController_Run(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model2"}
	downloaderArgs := downloadermodel.Args{}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Run add method
	_ = ac.Run(args, downloaderArgs)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 3)
	test.AssertEqual(t, models[2].Name, "model2")
}

// Tests RunAdd with invalid model
func TestAddController_Run_WithInvalidModel(t *testing.T) {
	// Init
	ac := AddController{}
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{Name: "model1"})
	existingModels = append(existingModels, model.Model{Name: "model3"})
	args := []string{"model3"}
	downloaderArgs := downloadermodel.Args{}

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{Name: "model2", LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Run add method
	err = ac.Run(args, downloaderArgs)
	test.AssertNotEqual(t, err, nil)

	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 2)
}
