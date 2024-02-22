package cmdmodeladd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path"
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

// runAddCustom runs add command to add a custom model
func runAddCustom(cmd *cobra.Command, args []string) {

	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// TODO: Get flags or default values

	// Searching for the currentCmd : when 'cmd' differs from 'addCustomCmd' (i.e. run through parent multiselect)
	currentCmd, found := utils.CobraFindSubCommand(cmd, addCustomCommandName)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", addCustomCommandName))
		return
	}

	// Asks for the mandatory args if they are not provided
	if addCustomDownloaderArgs.ModelName == "" {
		err := utils.CobraAskFlagInput(currentCmd, currentCmd.Flag(downloader.ModelName))
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Couldn't set the value for %s : %s", downloader.ModelName, err))
			return
		}
	}
	if addCustomDownloaderArgs.ModelModule == "" {
		err := utils.CobraAskFlagInput(currentCmd, currentCmd.Flag(downloader.ModelModule))
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Couldn't set the value for %s : %s", downloader.ModelModule, err))
			return
		}
	}

	// Allow the user to choose flags and specify their value
	err := utils.CobraInputAmongRemainingFlags(currentCmd)
	if err != nil {
		pterm.Error.Println(err)
		return
	}

	valid, err := validateModel(addCustomDownloaderArgs.ModelName)
	if err != nil {
		pterm.Error.Println(err)
		return
	}
	if !valid {
		pterm.Warning.Println("This model is already downloaded "+
			"and should be checked manually", addCustomDownloaderArgs.ModelName)
		return
	}

	// Running the script
	dlModel, err := downloader.Execute(addCustomDownloaderArgs)
	if err != nil || dlModel.IsEmpty {
		// Something went wrong or returned data is empty
		return
	}

	// Create the model for the configuration file
	modelObj := model.Model{Name: addCustomDownloaderArgs.ModelName}
	modelObj = model.MapToModelFromDownloaderModel(modelObj, dlModel)
	modelObj.AddToBinaryFile = true
	modelObj.IsDownloaded = true

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing model to configuration file...")
	err = config.AddModels([]model.Model{modelObj})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

}

func validateModel(modelName string) (bool, error) {
	exists, err := utils.IsExistingPath(path.Join(downloader.DirectoryPath, modelName))
	if err != nil {
		return false, err
	}
	if exists {
		message := fmt.Sprintf("This model %s is already downloaded do you wish to overwrite it?", modelName)
		valid := utils.AskForUsersConfirmation(message)
		return valid, nil
	}
	return true, nil
}

func init() {
	// Bind cobra args to the downloader script args
	downloader.ArgsGetForCobra(addCustomCmd, &addCustomDownloaderArgs)
}
