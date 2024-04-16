package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

// TestTokenizerAddCmd_WrongModule tests the Add command with wrong module
func TestTokenizerAddCmd_WrongModule(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.DIFFUSERS,
	})

	// Create Downloader mock
	downloader := mock.MockDownloader{
		DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Initialize selected models list
	args := []string{"model1", "tokenizer1"}

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	expectedMessage := "no models to choose from"
	test.AssertEqual(t, err.Error(), expectedMessage, "error")

}

// TestTokenizerAddCmd_NoArgs tests the Add command with no args
func TestTokenizerAddCmd_NoArgs(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Create Downloader mock
	downloader := mock.MockDownloader{
		DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)
	// Initialize selected models list
	var args []string

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	test.AssertEqual(t, err.Error(), "please provide a model and tokenizer name to add", "error")

}

// TestTokenizerAddCmd_NoTokenizerArg tests the Add
// command with no tokenizer in args
func TestTokenizerAddCmd_NoTokenizerArg(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Initialize selected models list
	args := []string{"model1"}
	// Create ui mock
	ui := mock.MockUI{SelectResult: "model1"}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{
		DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)
	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	test.AssertEqual(t, err.Error(), "please provide a model and tokenizer name to add", "error")
}

// TestTokenizerAddCmd_TokenizerDl tests the Add
// command with the tokenizer already downloaded
func TestTokenizerAddCmd_TokenizerDl(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
		Tokenizers: model.Tokenizers{
			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
		},
	})

	// Initialize selected models list
	args := []string{"model1", "tokenizer1"}

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	expectedErrMsg := "the following tokenizer is already downloaded :tokenizer1"
	test.AssertEqual(t, err.Error(), expectedErrMsg, "Unexpected error message")

}

// TestTokenizerAddCmd_DownloadTokenizerSuccess
// tests the Add command dl tokenizer Success
func TestTokenizerAddCmd_DownloadTokenizerSuccess(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Initialize selected models list
	args := []string{"model1", "tokenizer1"}

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	test.AssertEqual(t, err, nil)
}

// TestTokenizerAddCmd_DownloadTokenizerFail
// tests the Add command dl tokenizer Fail
func TestTokenizerAddCmd_DownloadTokenizerFail(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Initialize selected models list
	args := []string{"model1", "tokenizer1"}

	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderError: fmt.Errorf("")}
	app.SetDownloader(&downloader)

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	// Assertions
	expectedErrMsg := "the following tokenizer couldn't be downloaded : tokenizer1"
	test.AssertEqual(t, err.Error(), expectedErrMsg, "Unexpected error message")
}

// TestTokenizerAddCmd_NoModels tests the Add command with no models to choose from
func TestTokenizerAddCmd_NoModels(t *testing.T) {
	var models model.Models
	// Create ui mock
	ui := mock.MockUI{SelectResult: "model1"}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := mock.MockDownloader{
		DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)
	// Initialize selected models list
	args := []string{"model1", "tokenizer1"}

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
	ic := AddController{}

	var customArgs downloadermodel.Args
	// Process update
	err = ic.Run(args, customArgs)
	test.AssertEqual(t, err.Error(), "no models to choose from", "error")

}
