package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/appselec"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.Init("", "")
	app.InitGit("https://github.com/SchawnnDev", "")
	os.Exit(m.Run())
}

// TestTokenizerUpdateCmd_ValidArgs tests the update command with valid args
func TestTokenizerUpdateCmd_ValidArgs(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	// Initialize selected models list
	var args []string
	args = append(args, "model1")
	args = append(args, "tokenizer1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := UpdateTokenizerController{}
	// Process update
	err = ic.TokenizerUpdateCmd(args)
	expectedErrMsg := "the following tokenizer(s) couldn't be downloaded : [tokenizer1]"
	test.AssertEqual(t, err.Error(), expectedErrMsg, "Unexpected error message")

	_, err = config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")
}

// TestTokenizerUpdateCmd_NoModuleTransformersUpdate tests the update command
// with non-existing tokenizer path
func TestTokenizerUpdateCmd_NoModuleTransformersUpdate(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Source: "CUSTOM",
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})
	// Initialize selected models list
	var args []string
	args = append(args, "model1")
	args = append(args, "tokenizer2")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := UpdateTokenizerController{}
	// Process update

	err = ic.TokenizerUpdateCmd(args)
	expectedErrMsg := "only transformers models have tokenizers"
	test.AssertEqual(t, err.Error(), expectedErrMsg, "Unexpected error message")

	_, err = config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")
}

// TestTokenizerUpdateCmd_WrongModelNameUpdate tests the update command
// with the wrong model name
func TestTokenizerUpdateCmd_WrongModelNameUpdate(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})
	// Initialize selected models list
	var args []string
	args = append(args, "modelX")
	args = append(args, "tokenizer1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	ic := UpdateTokenizerController{}
	// Process update
	err = ic.TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "Operation failed, no model found")
}

// TestTokenizerUpdateCmd_NoArgs tests the update command
// with no args
func TestTokenizerUpdateCmd_NoArgs(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})
	// Initialize selected models list
	var args []string

	//Create model selector mock
	mockModel := model.Model{Name: "model1", Module: huggingface.TRANSFORMERS, Tokenizers: model.Tokenizers{
		{Path: "path1",
			Class:   "tokenizer1",
			Options: map[string]string{"option1": "value1"}},
	}}
	// Create Ui mock
	selector := mock.MockModelSelector{SelectorModel: mockModel,
		SelectorError:   nil,
		SelectorWarning: "",
	}
	appselec.SetSelector(&selector)

	// Create ui mock
	ui := mock.MockUI{SelectResult: "tokenizer1"}
	app.SetUI(ui)

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := UpdateTokenizerController{}
	// Process update
	err = ic.TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "Operation succeeded.")

}

// TestTokenizerUpdateCmd_NoTokenizerInArgs tests the update command
// with no tokenizers in args
func TestTokenizerUpdateCmd_NoTokenizerInArgs(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	// Create ui mock
	ui := mock.MockUI{SelectResult: "model1"}
	app.SetUI(ui)

	// Initialize selected models list
	var args []string
	args = append(args, "model1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	ic := UpdateTokenizerController{}
	// Process update
	err = ic.TokenizerUpdateCmd(args)
	_, err = config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

}

// TestTokenizerUpdateCmd_NoTokenizerInArgsDownload tests the update command
// with no tokenizers in args and downloading
func TestTokenizerUpdateCmd_NoTokenizerInArgsDownload(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	var expectedSelections []string
	expectedSelections = append(expectedSelections, "tokenizer1")

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Initialize selected models list
	var args []string
	args = append(args, "model1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")

	ic := UpdateTokenizerController{}
	// Process update
	err = ic.TokenizerUpdateCmd(args)
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
}

// TestTokenizerUpdateCmd_UpdateError tests the error return of the update command
func TestTokenizerUpdateCmd_UpdateError(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	var expectedSelections []string
	expectedSelections = append(expectedSelections, "tokenizer1")

	// Create ui mock
	ui := mock.MockUI{MultiselectResult: expectedSelections}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{
		DownloaderModel: downloadermodel.Model{Path: "test"},
		DownloaderError: os.ErrClosed,
	}
	app.SetDownloader(&downloader)

	// Initialize selected models list
	var args []string
	args = append(args, "model1")

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := UpdateTokenizerController{}
	// Process update

	err = ic.TokenizerUpdateCmd(args)
	expectedErrMsg := "the following tokenizer(s) couldn't be downloaded : [tokenizer1]"
	test.AssertEqual(t, err.Error(), expectedErrMsg, "Unexpected error message")
}
