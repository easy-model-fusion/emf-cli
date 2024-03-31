package tokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

// TestRunTokenizerAdd_Success tests the RunTokenizerAdd function
func TestRunTokenizerAdd_ModelNotFound(t *testing.T) {
	// Setup test data
	args := []string{"model13", "tokenizer1"}

	var models model.Models

	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected "+
		"while adding models to configuration file")

	downloaderArgs := downloadermodel.Args{}
	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	ic := AddTokenizerController{}
	if err := ic.RunTokenizerAdd(args, downloaderArgs); err != nil {
		expectedErrMsg := "Model is not configured"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
}

// TestRunTokenizerAdd_Success tests the RunTokenizerAdd function
func TestRunTokenizerAdd_NoArgs(t *testing.T) {
	// Setup test data
	var args []string

	var models model.Models

	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected "+
		"while adding models to configuration file")

	downloaderArgs := downloadermodel.Args{}
	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	ic := AddTokenizerController{}
	if err := ic.RunTokenizerAdd(args, downloaderArgs); err != nil {
		expectedErrMsg := "enter a model in argument"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
}

// TestRunTokenizerAdd_Success tests the RunTokenizerAdd function
func TestRunTokenizerAdd_NoTokenizerArgs(t *testing.T) {
	// Setup test data
	args := []string{"model1"}

	var models model.Models

	models = append(models, model.Model{
		Name:   "model1",
		Module: huggingface.TRANSFORMERS,
	})

	// Create temporary configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected "+
		"while adding models to configuration file")

	downloaderArgs := downloadermodel.Args{}
	//Create downloader mock
	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
	app.SetDownloader(&downloader)

	// Process add
	ic := AddTokenizerController{}
	if err := ic.RunTokenizerAdd(args, downloaderArgs); err != nil {
		expectedErrMsg := "enter a tokenizer in argument"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
}

// TestRemoveTokenizer_WithModuleNotTransformers tests the RunTokenizerRemove function with no transformers module
func TestAddTokenizer_WithModuleNotTransformers(t *testing.T) {
	var models model.Models
	models = append(models, model.Model{
		Name: "model1",
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

	ic := AddTokenizerController{}
	// Process add
	downloaderArgs := downloadermodel.Args{}
	if err := ic.RunTokenizerAdd(args, downloaderArgs); err != nil {
		expectedErrMsg := "only transformers models have tokenizers"
		if err.Error() != expectedErrMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
		}
	}
	test.AssertEqual(t, err, nil, "No error expected while processing remove")
	newModels, err := config.GetModels()
	test.AssertEqual(t, err, nil, "No error expected on getting models")

	//Assertions
	test.AssertEqual(t, len(newModels[0].Tokenizers), 1, "Only one model should be left.")
}

//
//// TestRemoveTokenizer_Success tests the RunTokenizerRemove function
//func TestAddTokenizer_Success(t *testing.T) {
//	var models model.Models
//	models = append(models, model.Model{
//		Name:   "model1",
//		Module: huggingface.TRANSFORMERS,
//		Tokenizers: model.Tokenizers{
//			{Path: "path1", Class: "tokenizer1", Options: map[string]string{"option1": "value1"}},
//		},
//	})
//	// Initialize selected models list
//	var args []string
//	args = append(args, "model1")
//	args = append(args, "tokenizer1")
//
//	// Create temporary configuration file
//	ts := test.TestSuite{}
//	_ = ts.CreateFullTestSuite(t)
//	defer ts.CleanTestSuite(t)
//	err := setupConfigFile(models)
//	test.AssertEqual(t, err, nil, "No error expected while adding models to configuration file")
//
//	//Create downloader mock
//	downloader := mock.MockDownloader{DownloaderModel: downloadermodel.Model{Module: "diffusers", Class: "test"}}
//	app.SetDownloader(&downloader)
//
//	ic := AddTokenizerController{}
//	// Process add
//	downloaderArgs := downloadermodel.Args{}
//
//	if err := ic.RunTokenizerAdd(args, downloaderArgs); err != nil {
//		test.AssertEqual(t, err, nil, "Error on update")
//	}
//	test.AssertEqual(t, err, nil, "No error expected while processing remove")
//	newModels, err := config.GetModels()
//	test.AssertEqual(t, err, nil, "No error expected on getting models")
//
//	//Assertions
//	test.AssertEqual(t, len(newModels[0].Tokenizers), 0, "Only one model should be left.")
//}
