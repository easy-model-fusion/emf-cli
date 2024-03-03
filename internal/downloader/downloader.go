package downloader

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/pterm/pterm"
)

const ScriptPath = "sdk/downloader.py"

// Model represents a model returned by the downloader script.
type Model struct {
	Path      string    `json:"path"`
	Module    string    `json:"module"`
	Class     string    `json:"class"`
	Tokenizer Tokenizer `json:"tokenizer"`
	IsEmpty   bool
}

// Tokenizer represents a tokenizer returned by the downloader script.
type Tokenizer struct {
	Path  string `json:"path"`
	Class string `json:"class"`
}

// Tags for the arguments
const TagPrefix = "--"
const ModelName = "model-name"
const ModelModule = "model-module"
const ModelClass = "model-class"
const ModelOptions = "model-options"
const TokenizerClass = "tokenizer-class"
const TokenizerOptions = "tokenizer-options"
const Overwrite = "overwrite"
const Skip = "skip"
const EmfClient = "emf-client"

// Args represents the arguments for the script
type Args struct {
	ModelName        string
	ModelModule      string
	ModelClass       string
	ModelOptions     []string
	TokenizerClass   string
	TokenizerOptions []string
	Skip             string
}

// Execute runs the downloader script and handles the result
func Execute(downloaderArgs Args) (Model, error) {

	// Check arguments validity
	err := ArgsValidate(downloaderArgs)
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Arguments provided are invalid : %s", err))
		return Model{}, err
	}

	// Building args for the python script
	args := ArgsProcessForPython(downloaderArgs)

	// Run the script to download the model
	spinner := app.UI().StartSpinner(fmt.Sprintf("Downloading model '%s'...", downloaderArgs.ModelName))
	scriptModel, err, exitCode := python.ExecuteScript(".venv", ScriptPath, args)

	// An error occurred while running the script
	if err != nil {
		spinner.Fail(err)
		switch exitCode {
		case 2:
			pterm.Info.Println("Use command 'add custom' to customize the model to download.")
		}
		return Model{}, err
	}

	// No data was returned by the script
	if scriptModel == nil {
		spinner.Fail(fmt.Sprintf("The script didn't return any data for '%s'", downloaderArgs.ModelName))
		return Model{IsEmpty: true}, nil
	}

	// Unmarshall JSON response
	var model Model
	err = json.Unmarshal(scriptModel, &model)
	if err != nil {
		return Model{}, err
	}

	// Download was successful
	spinner.Success(fmt.Sprintf("Successfully downloaded model '%s'", downloaderArgs.ModelName))

	return model, nil
}
