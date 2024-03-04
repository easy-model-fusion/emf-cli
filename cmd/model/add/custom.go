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
	modelObj = model.ConstructConfigPaths(modelObj)

	// Validate the model : if model is already downloaded
	downloaded, err := model.ModelDownloadedOnDevice(modelObj)
	if err != nil {
		pterm.Error.Println(err)
		return
	} else if downloaded {
		message := fmt.Sprintf("Model '%s' is already downloaded. Do you wish to overwrite it?", modelObj.Name)
		overwrite := app.UI().AskForUsersConfirmation(message)
		if !overwrite {
			pterm.Warning.Println("This model is already downloaded and should be checked manually", modelObj.Name)
			return
		}
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

	// TODO : check if tokenizer already exists => Waiting for issue #63 : [Client] Validate models for download

	// Downloading model
	modelObj, success := model.Download(modelObj, addCustomDownloaderArgs)
	if !success {
		return
	}

	// Add models to configuration file
	spinner := app.UI().StartSpinner("Writing model to configuration file...")
	err = config.AddModels([]model.Model{modelObj})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

}
