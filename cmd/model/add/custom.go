package cmdmodeladd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const addCustomCommandName = "custom"

// addCustomCmd represents the add custom model command
var addCustomCmd = &cobra.Command{
	Use:   addCustomCommandName + " [flags]",
	Short: "Add a customized model to your project",
	Long:  `Add a customized model to your project by specifying properties yourself`,
	Run:   runAddCustom,
}

var addCustomDownloaderArgs downloader.Args

func init() {
	// Bind cobra args to the downloader script args
	downloader.ArgsGetForCobra(addCustomCmd, &addCustomDownloaderArgs)
}

// runAddCustom runs add command to add a custom model
func runAddCustom(cmd *cobra.Command, args []string) {

	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// Searching for the currentCmd : when 'cmd' differs from 'addCustomCmd' (i.e. run through parent multiselect)
	currentCmd, found := cobrautil.FindSubCommand(cmd, addCustomCommandName)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", addCustomCommandName))
		return
	}

	// Model name is mandatory
	if addCustomDownloaderArgs.ModelName == "" {
		err := cobrautil.AskFlagInput(currentCmd, currentCmd.Flag(downloader.ModelName))
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Couldn't set the value for %s : %s", downloader.ModelName, err))
			return
		}
	}

	// Get model from huggingface : verify its existence and get mandatory data
	huggingfaceModel, err := app.H().GetModelById(addCustomDownloaderArgs.ModelName)
	if err != nil {
		// Model not found
		pterm.Warning.Printfln("Model %s not valid : "+err.Error(), addCustomDownloaderArgs.ModelName)
		return
	}
	// Map API response to model.Model
	modelObj := model.MapToModelFromHuggingfaceModel(huggingfaceModel)

	// Validate model for download
	modelObj.AddToBinaryFile = true
	if !config.Validate(modelObj) {
		return
	}

	// Allow the user to choose flags and set their values
	err = cobrautil.AllowInputAmongRemainingFlags(currentCmd)
	if err != nil {
		pterm.Error.Println(err)
		return
	}

	// Module not provided : get it from the API
	if addCustomDownloaderArgs.ModelModule == "" {
		addCustomDownloaderArgs.ModelModule = string(modelObj.Module)
	}

	// Class not provided : get it from the API
	if addCustomDownloaderArgs.ModelClass == "" {
		addCustomDownloaderArgs.ModelClass = modelObj.Class
	}

	// Downloading model
	var success bool
	modelObj, success = model.Download(modelObj, addCustomDownloaderArgs)
	if !success {
		modelObj.AddToBinaryFile = false
	}

	// Add models to configuration file
	spinner := app.UI().StartSpinner("Writing model to configuration file...")
	err = config.AddModels([]model.Model{modelObj})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

	// Attempt to generate code
	spinner = app.UI().StartSpinner("Generating python code...")
	err = config.GenerateExistingModelsPythonCode()
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while generating python code for added model: %s", err))
	} else {
		spinner.Success()
	}
}
