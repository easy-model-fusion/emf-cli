package tokenizer

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

type AddController struct{}

// Run the tokenizer add command
func (ic AddController) Run(args []string,
	customArgs downloadermodel.Args) error {
	sdk.SendUpdateSuggestion()

	// Process add operation with given arguments
	warningMessages, infoMessage, err := ic.processAddTokenizer(args, customArgs)

	// Display messages to user
	for _, warningMessage := range warningMessages {
		app.UI().Warning().Printfln(warningMessage)
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
) (warnings []string, info string, err error) {
	// Load the configuration file
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warnings, info, err
	}

	// No model name in args
	if len(args) < 2 {
		return warnings, info, fmt.Errorf("please provide a model and tokenizer name to add")
	}
	// Get all configured models objects/names and args model
	models, err := config.GetModelsByModule(string(huggingface.TRANSFORMERS))
	if err != nil {
		return warnings, info, fmt.Errorf("error getting model: %s", err.Error())
	}
	if len(models) == 0 {
		err = fmt.Errorf("no models to choose from")
		return warnings, info, err
	}

	var modelToUse model.Model
	configModelsMap := models.Map()

	// Checks the presence of the model
	selectedModel := args[0]
	modelToUse, exists := configModelsMap[selectedModel]
	if !exists {
		return warnings, "Model is not configured", err
	}

	// Remove model name from arguments
	args = args[1:]

	var selectedTokenizersToUse []string

	// Setting tokenizer name from args
	selectedTokenizersToUse = stringutil.SliceRemoveDuplicates(args)

	var tokenizerName string
	for _, tokenizerName = range selectedTokenizersToUse {
		tokenizerFound := modelToUse.Tokenizers.ContainsByClass(tokenizerName)
		if tokenizerFound {
			err = fmt.Errorf("the following tokenizer is already downloaded :%s",
				tokenizerName)
			return warnings, info, err
		}
		addedTokenizer := model.Tokenizer{
			Class: tokenizerName,
		}
		customArgs.ModelName = modelToUse.Name
		customArgs.DirectoryPath, err = modelToUse.GetModelDirectory()
		if err != nil {
			return warnings, info, err
		}

		var success bool
		success, warnings, err = modelToUse.DownloadTokenizer(addedTokenizer, customArgs)
		if err != nil {
			return warnings, info, err
		}
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
			modelToUse.Tokenizers = append(modelToUse.Tokenizers, addedTokenizer)
		}

	}
	return warnings, "Tokenizers add done", err
}
