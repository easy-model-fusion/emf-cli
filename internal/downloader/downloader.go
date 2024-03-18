package downloader

import (
	"encoding/json"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/pterm/pterm"
)

type scriptDownloader struct{}

// NewScriptDownloader initialize a new script downloader
func NewScriptDownloader() Downloader {
	return &scriptDownloader{}
}

// Constants related to the downloader script python arguments.
const (
	ScriptPath         = "sdk/downloader.py"
	TagPrefix          = "--"
	ModelName          = "model-name"
	ModelModule        = "model-module"
	ModelClass         = "model-class"
	ModelOptions       = "model-options"
	TokenizerClass     = "tokenizer-class"
	TokenizerOptions   = "tokenizer-options"
	Overwrite          = "overwrite"
	Skip               = "skip"
	SkipValueModel     = "model"
	SkipValueTokenizer = "tokenizer"
	EmfClient          = "emf-client"
	OnlyConfiguration  = "only-configuration"
)

// Args represents the arguments for the script.
type Args struct {
	ModelName         string
	ModelModule       string
	ModelClass        string
	ModelOptions      []string
	TokenizerClass    string
	TokenizerOptions  []string
	Skip              string
	OnlyConfiguration bool
}

// Model represents a model returned by the downloader script.
type Model struct {
	Path      string            `json:"path"`
	Module    string            `json:"module"`
	Class     string            `json:"class"`
	Options   map[string]string `json:"options"`
	Tokenizer Tokenizer         `json:"tokenizer"`
	IsEmpty   bool
}

// Tokenizer represents a tokenizer returned by the downloader script.
type Tokenizer struct {
	Path    string            `json:"path"`
	Class   string            `json:"class"`
	Options map[string]string `json:"options"`
}

// Execute runs the downloader script and handles the result
func (downloader *scriptDownloader) Execute(downloaderArgs Args) (Model, error) {

	// Check arguments validity
	err := downloaderArgs.Validate()
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Arguments provided are invalid : %s", err))
		return Model{}, err
	}

	// Building args for the python script
	args := downloaderArgs.ToPython()

	// Preparing spinner message
	var downloaderItemMessage string
	switch downloaderArgs.Skip {
	case SkipValueModel:
		downloaderItemMessage = fmt.Sprintf("tokenizer '%s' for model '%s'...", downloaderArgs.TokenizerClass, downloaderArgs.ModelName)
	default:
		downloaderItemMessage = fmt.Sprintf("model '%s'...", downloaderArgs.ModelName)
	}

	// Run the script to download the model
	customMessage := "Downloading "
	if downloaderArgs.OnlyConfiguration {
		customMessage = "Getting configuration for "
	}
	spinner := app.UI().StartSpinner(customMessage + downloaderItemMessage)
	scriptModel, err, _ := app.Python().ExecuteScript(".venv", ScriptPath, args)

	// An error occurred while running the script
	if err != nil {
		spinner.Fail(err)
		return Model{}, err
	}

	// No data was returned by the script
	if scriptModel == nil {
		spinner.Fail("The script didn't return any data when processing " + downloaderItemMessage)
		return Model{IsEmpty: true}, nil
	}

	// Unmarshall JSON response
	var model Model
	err = json.Unmarshal(scriptModel, &model)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to process the script return data : %s", err))
		return Model{}, err
	}

	// Download was successful
	if downloaderArgs.OnlyConfiguration {
		customMessage = "got configuration for "
	} else {
		customMessage = "downloaded "
	}
	spinner.Success("Successfully " + customMessage + downloaderItemMessage)

	return model, nil
}
