package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/dmock"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.Init("", "")
	os.Exit(m.Run())
}

// Sets the configuration file with the given models
func setupConfigFile(models model.Models) error {
	config.FilePath = "."
	// Load configuration file
	err := config.GetViperConfig(".")
	if err != nil {
		return err
	}
	// Write models to the config file
	viper.Set("models", models)
	return config.WriteViperConfig()
}

// Tests tidyModelsConfiguredButNotDownloaded
func TestTidyModelsConfiguredButNotDownloaded_Success(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:         "model1",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model2",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: false,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model5",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := config.GetViperConfig(".")
	test.AssertEqual(t, err, nil, "No error expected on loading configuration file")

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsConfiguredButNotDownloaded(existingModels, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
}

// Tests tidyModelsConfiguredButNotDownloaded with no configuration file loaded
func TestTidyModelsConfiguredButNotDownloaded_SuccessWithNoConfFile(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:         "model1",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model2",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: false,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model5",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsConfiguredButNotDownloaded(existingModels, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
}

// Tests tidyModelsConfiguredButNotDownloaded with download failure
func TestTidyModelsConfiguredButNotDownloaded_Fail(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:         "model1",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model2",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: false,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model5",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := config.GetViperConfig(".")
	test.AssertEqual(t, err, nil, "No error expected on loading configuration file")

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsConfiguredButNotDownloaded(existingModels, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 1)
	test.AssertEqual(t, warnings[0], "The following models(s) couldn't be downloaded : [model5]")

}

// Tests tidyModelsConfiguredButNotDownloaded with tokenizer download failure
func TestTidyModelsConfiguredButNotDownloaded_WithTokenizerFailure(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:         "model1/name",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: true,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model2/name",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		IsDownloaded: false,
	})
	existingModels = append(existingModels, model.Model{
		Name:   "model4/name",
		Path:   "./models/model4/name/model",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			model.Tokenizer{
				Class: "tokenizer",
				Path:  "models/model4/name/tokenizer",
			},
			model.Tokenizer{
				Class: "tokenizer2",
				Path:  "invalid/Path",
			},
		},
		IsDownloaded: true,
	})

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsConfiguredButNotDownloaded(existingModels, "")

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 1)
	test.AssertEqual(t, warnings[0], "The following tokenizer(s) couldn't be downloaded for 'model4/name': [tokenizer2]")
}

// Tests tidyModelsDownloadedButNotConfigured
func TestTidyModelsDownloadedButNotConfigured(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:   "model1/name",
		Module: huggingface.DIFFUSERS,
		Class:  "test",
	})
	existingModels = append(existingModels, model.Model{
		Name:   "model2/name",
		Module: huggingface.DIFFUSERS,
		Class:  "test",
	})
	existingModels = append(existingModels, model.Model{
		Name:   "model4/name",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			model.Tokenizer{
				Class: "tokenizer",
			},
			model.Tokenizer{
				Class: "tokenizer2",
			},
		},
	})

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsDownloadedButNotConfigured(existingModels, true, "")
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 4)
	test.AssertEqual(t, models[3].Name, "model3/name")
}

// Tests tidyModelsDownloadedButNotConfigured with no user confirmation
func TestTidyModelsDownloadedButNotConfigured_WithNoConfirmation(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:   "model1/name",
		Module: huggingface.DIFFUSERS,
		Class:  "test",
	})
	existingModels = append(existingModels, model.Model{
		Name:   "model2/name",
		Module: huggingface.DIFFUSERS,
		Class:  "test",
	})
	existingModels = append(existingModels, model.Model{
		Name:   "model4/name",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			model.Tokenizer{
				Class: "tokenizer",
			},
			model.Tokenizer{
				Class: "tokenizer2",
			},
		},
	})

	// Create huggingface mock
	huggingfaceInterface := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingfaceInterface)

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Create Downloader mock
	ui := mock.MockUI{UserConfirmationResult: false}
	app.SetUI(&ui)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Download missing models
	var tidyController TidyController
	warnings, err := tidyController.tidyModelsDownloadedButNotConfigured(existingModels, false, "")
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(warnings), 0)
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 3)
}

// Tests RunTidy
func TestRunTidy(t *testing.T) {
	// Init
	var existingModels model.Models
	existingModels = append(existingModels, model.Model{
		Name:         "model1/name",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		PipelineTag:  huggingface.TextToImage,
		IsDownloaded: true,
	})
	existingModels = append(existingModels, model.Model{
		Name:         "model2/name",
		Module:       huggingface.DIFFUSERS,
		Class:        "test",
		PipelineTag:  huggingface.TextToImage,
		IsDownloaded: false,
	})
	existingModels = append(existingModels, model.Model{
		Name:        "model4/name",
		Path:        "./models/model4/model",
		Module:      huggingface.TRANSFORMERS,
		PipelineTag: huggingface.TextToImage,
		Tokenizers: model.Tokenizers{
			model.Tokenizer{
				Class: "tokenizer",
				Path:  "models/model4/tokenizer",
			},
			model.Tokenizer{
				Class: "tokenizer2",
				Path:  "invalid/Path",
			},
		},
		IsDownloaded: true,
	})

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateModelsFolderFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(existingModels)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Download missing models
	var tidyController TidyController
	_ = tidyController.RunTidy(true, "")
	models, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(models), 4)
	test.AssertEqual(t, models[3].Name, "model3/name")
}

// Tests RunTidy with no configuration file
func TestRunTidy_WithNoConfigurationFile(t *testing.T) {
	// Download missing models
	var tidyController TidyController
	err := tidyController.RunTidy(true, "")
	test.AssertNotEqual(t, err, nil, "An error expected on synchronizing models")
}
