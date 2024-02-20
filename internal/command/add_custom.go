package command

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/huggingface"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/script"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path"
)

const cmdAddCustomTitle = "custom"

// addCustomCmd represents the add custom model command
var addCustomCmd = &cobra.Command{
	Use:   cmdAddCustomTitle + " [flags]",
	Short: "Add a customized model to your project",
	Long:  `Add a customized model to your project by specifying properties yourself`,
	Run:   runAddCustom,
}

var downloaderArgs script.DownloaderArgs

// runAddCustom runs add command for adding a custom model
func runAddCustom(cmd *cobra.Command, args []string) {

	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// TODO: Get flags or default values

	// Searching for the currentCmd : when 'cmd' differs from 'addCustomCmd' (i.e. run through parent multiselect)
	currentCmd, found := utils.CobraFindSubCommand(cmd, cmdAddCustomTitle)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", cmdAddTitle))
		return
	}

	// Asks for the mandatory args if they are not provided
	if downloaderArgs.ModelName == "" {
		err := utils.CobraAskFlagInput(currentCmd, currentCmd.Flag(script.ModelName))
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Couldn't set the value for %s : %s", script.ModelName, err))
			return
		}
	}
	if downloaderArgs.ModelModule == "" {
		err := utils.CobraAskFlagInput(currentCmd, currentCmd.Flag(script.ModelModule))
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("Couldn't set the value for %s : %s", script.ModelModule, err))
			return
		}
	}

	// Allow the user to choose flags and specify their value
	err := utils.CobraInputAmongRemainingFlags(currentCmd)
	if err != nil {
		pterm.Error.Println(err)
		return
	}

	valid, err := validateModel(downloaderArgs.ModelName)
	if !valid {
		pterm.Warning.Println("This model is already downloaded "+
			"and should be checked manually", downloaderArgs.ModelName)
		return
	}

	// Running the script
	sdm, err := script.DownloaderExecute(downloaderArgs)
	if err != nil || sdm.IsEmpty {
		// Something went wrong or returned data is empty
		return
	}

	// Create the model for the configuration file
	modelObj := model.Model{Name: downloaderArgs.ModelName}
	modelObj.Config = model.MapToConfigFromScriptDownloaderModel(modelObj.Config, sdm)
	modelObj.ShouldBeDownloaded = true

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
	exists, err := utils.IsExistingPath(path.Join(script.DownloadModelsPath, modelName))
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

	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Bind cobra args to the downloader script args
	script.DownloaderArgsForCobra(addCustomCmd, &downloaderArgs)

	// Add the subcommands to the add command
	addCmd.AddCommand(addCustomCmd)
}
