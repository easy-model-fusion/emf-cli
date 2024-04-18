package cmdmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/dmock"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"testing"
)

func TestRunModelUpdate(t *testing.T) {
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
	huggingFace := huggingface.MockHuggingFace{GetModelResult: huggingface.Model{LastModified: "2022", LibraryName: huggingface.TRANSFORMERS}}
	app.SetHuggingFace(&huggingFace)

	// Create Ui mock
	ui := mock.MockUI{UserConfirmationResult: true, UserInputResult: "."}
	app.SetUI(ui)

	// Create Downloader mock
	downloader := dmock.MockDownloader{DownloaderModel: downloadermodel.Model{Path: "test"}, DownloaderError: nil}
	app.SetDownloader(&downloader)

	// Create full test suite with a configuration file
	ts := test.TestSuite{}
	_ = ts.CreateFullTestSuite(t)
	defer ts.CleanTestSuite(t)
	err := setupConfigFile(models)
	test.AssertEqual(t, err, nil, "No error expected on setting configuration file")

	// Process update
	runModelUpdate(nil, args)
	updatedModels, err := config.GetModels()

	// Assertions
	test.AssertEqual(t, err, nil, "No error expected on getting all models")
	test.AssertEqual(t, len(updatedModels), 3)
	test.AssertEqual(t, updatedModels[0].Version, "2022")
	test.AssertEqual(t, updatedModels[1].Version, "2022")
	test.AssertEqual(t, updatedModels[2].Version, "2022")
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
