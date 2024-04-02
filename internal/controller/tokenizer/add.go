package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

type AddController struct{}

// Run the tokenizer add command
func (ic AddController) Run(args []string,
	customArgs downloadermodel.Args) error {
	sdk.SendUpdateSuggestion()

	// Process add operation with given arguments
	warningMessage, infoMessage, err := ic.processAddTokenizer(args, customArgs)

	// Display messages to user
	if warningMessage != "" {
		pterm.Warning.Printfln(warningMessage)
	}

	if infoMessage != "" {
		pterm.Info.Printfln(infoMessage)
		return err
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
		return err
	} else {
		pterm.Error.Printfln("Operation failed.")
		return err
	}
}

// processAddTokenizer processes tokenizers to be added
func (ic AddController) processAddTokenizer(
	args []string,
	customArgs downloadermodel.Args,
) (warning, info string, err error) {
	// Load the configuration file
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warning, info, err
	}

	// No model name in args
	if len(args) < 1 {
		return warning, info, fmt.Errorf("enter a model in argument")
	}

	// Get all configured models objects/names and args model
	models, err := config.GetModels()
	if err != nil {
		return warning, info, fmt.Errorf("error getting model: %s", err.Error())
	}

	// Checks the presence of the model
	selectedModel := args[0]
	configModelsMap := models.Map()
	modelToUse, exists := configModelsMap[selectedModel]
	if !exists {
		return warning, "Model is not configured", err
	}
	// No tokenizer name in args
	if len(args) < 2 {
		return warning, info, fmt.Errorf("enter a tokenizer in argument")
	}

	// Setting tokenizer name from args
	tokenizerName := args[1]

	tokenizerFound := modelToUse.Tokenizers.ContainsByClass(tokenizerName)

	if tokenizerFound {
		err = fmt.Errorf("the following tokenizer is already downloaded :%s",
			tokenizerName)
		return warning, "Tokenizer add failed, already downloaded", err
	}

	// Verify model's module
	if modelToUse.Module != huggingface.TRANSFORMERS {
		return warning, info, fmt.Errorf("only transformers models have tokenizers")
	}

	var addedTokenizer = model.Tokenizer{
		Path:  modelToUse.Path,
		Class: tokenizerName,
	}
	modelToUse.Tokenizers = append(modelToUse.Tokenizers, addedTokenizer)

	customArgs.ModelName = modelToUse.Name
	customArgs.ModelModule = string(modelToUse.Module)

	success := modelToUse.DownloadTokenizer(addedTokenizer, customArgs)
	if !success {
		err = fmt.Errorf("the following tokenizer"+
			" couldn't be downloaded : %s", tokenizerName)
	} else {

		spinner, _ := pterm.DefaultSpinner.Start("Updating configuration file...")
		err := config.AddModels(model.Models{modelToUse})
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while updating the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}
	return warning, "Tokenizers add done", err
}
