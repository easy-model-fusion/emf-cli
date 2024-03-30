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

// RunTokenizerAdd runs the tokenizer add command
func RunTokenizerAdd(args []string, customArgs downloadermodel.Args, yes bool) {
	sdk.SendUpdateSuggestion()

	// Process add operation with given arguments
	warningMessage, infoMessage, err := processAddTokenizer(args, customArgs, yes)

	// Display messages to user
	if warningMessage != "" {
		pterm.Warning.Printfln(warningMessage)
	}

	if infoMessage != "" {
		pterm.Info.Printfln(infoMessage)
	} else if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

// processAddTokenizer processes tokenizers to be added
func processAddTokenizer(
	args []string,
	customArgs downloadermodel.Args,
	yes bool,
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
		return warning, info, fmt.Errorf("error get model: %s", err.Error())
	}

	// Checks the presence of the model
	selectedModel := args[0]
	configModelsMap := models.Map()
	modelToUse, exists := configModelsMap[selectedModel]
	if !exists {
		return warning, "Model is not configured", err
	}

	// Verify model's module
	if modelToUse.Module != huggingface.TRANSFORMERS {
		return warning, info, fmt.Errorf("only transformers models have tokzenizers")
	}
	// No tokenizer name in args
	if len(args) < 2 {
		return warning, info, fmt.Errorf("enter a tokenizer in argument")
	}

	// Setting tokenizer name from args
	//tokenizerName := args[1]

	// Extracting available tokenizers

	//var downloadTokenizers []string
	//tokenizerToDl
	//TODO Add options
	var addedTokenizer = model.Tokenizer{
		Path:  modelToUse.Path,
		Class: customArgs.TokenizerClass,
	}
	println("Class is ", customArgs.TokenizerClass)
	modelToUse.Tokenizers = append(modelToUse.Tokenizers, addedTokenizer)

	println("addind tok to model ", model.Tokenizers{})
	// Try to update all the given models
	var downloadedTokenizers model.Tokenizers
	var failedTokenizers []string

	for tokenizer := range modelToUse.Tokenizers {

		downloaderArgs := downloadermodel.Args{
			ModelName:   modelToUse.Name,
			ModelModule: string(modelToUse.Module),
		}
		println("calling download")
		success := modelToUse.DownloadTokenizer(addedTokenizer, downloaderArgs)
		if !success {
			println("fail dl ")
			failedTokenizers = append(failedTokenizers, string(rune(tokenizer)))
		} else {
			println("success download")
			downloadedTokenizers = append(downloadedTokenizers, addedTokenizer)
		}
	}
	println("here")

	// Update tokenizers' configuration
	if len(downloadedTokenizers) > 0 {
		//Reset model while keeping unchanged tokenizers
		modelToUse.Tokenizers = modelToUse.Tokenizers.Difference(downloadedTokenizers)
		//Adding new version of downloaded tokenizers
		modelToUse.Tokenizers = append(modelToUse.Tokenizers, downloadedTokenizers...)

		spinner, _ := pterm.DefaultSpinner.Start("Updating configuration file...")
		err := config.AddModels(model.Models{modelToUse})
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while updating the configuration file: %s", err))
		} else {
			spinner.Success()
		}
	}

	// Displaying the downloads that failed
	if len(failedTokenizers) > 0 {
		err = fmt.Errorf("the following tokenizers(s) couldn't be downloaded : %s", failedTokenizers)
	}
	return warning, "Tokenizers update done", err
}
