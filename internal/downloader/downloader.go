package downloader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
)

type Downloader interface {
	Execute(downloaderArgs downloadermodel.Args, python python.Python, ctx context.Context) (downloadermodel.Model, error)
}

type scriptDownloader struct{}

// NewScriptDownloader initialize a new script downloader
func NewScriptDownloader() Downloader {
	return &scriptDownloader{}
}

// Execute runs the downloader script and handles the result
func (downloader *scriptDownloader) Execute(downloaderArgs downloadermodel.Args, python python.Python, ctx context.Context) (downloadermodel.Model, error) {

	// Check arguments validity
	err := downloaderArgs.Validate()
	if err != nil {
		return downloadermodel.Model{}, fmt.Errorf("arguments provided are invalid : %s", err)
	}

	// Building args for the python script
	args := downloaderArgs.ToPython()

	// Run the script to download the model
	scriptModel, err, _ := python.ExecuteScript(".venv", downloadermodel.ScriptPath, args, ctx)

	// An error occurred while running the script
	if err != nil {
		return downloadermodel.Model{}, err
	}

	// No data was returned by the script
	if scriptModel == nil {
		return downloadermodel.Model{IsEmpty: true}, fmt.Errorf("the script didn't return any data")
	}

	// Unmarshall JSON response
	var model downloadermodel.Model
	err = json.Unmarshal(scriptModel, &model)
	if err != nil {
		return downloadermodel.Model{}, fmt.Errorf("failed to process the script return data : %s", err)
	}

	// Download was successful
	return model, nil
}
