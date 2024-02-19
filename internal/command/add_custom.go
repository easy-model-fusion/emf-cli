package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/script"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
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

	if config.GetViperConfig() != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// TODO: Get flags or default values
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Asks for the mandatory args if they are not provided
	utils.CobraAskFlagInput(cmd, cmd.Flag(script.ModelName))
	utils.CobraAskFlagInput(cmd, cmd.Flag(script.ModelModule))

	// Allow the user to choose flags and specify their value
	utils.CobraInputAmongRemainingFlags(cmd)

	// TODO : errors
	// TODO : options : split and encapsulate
	// TODO : unit tests

	// TODO : validate model to download

	// Running the script
	sdm, err := script.DownloaderExecute(downloaderArgs)
	if err != nil || sdm.IsEmpty {
		// Something went wrong or returned data is empty
		return
	}

	// Create the model for the configuration file
	modelObj := model.Model{Name: downloaderArgs.ModelName}
	modelObj.Config = model.MapToConfigFromScriptDownloaderModel(modelObj.Config, sdm)
	modelObj.AddToBinary = true

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing model to configuration file...")
	err = config.AddModel([]model.Model{modelObj})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the model to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

}

func init() {

	// Bind cobra args to the downloader script args
	script.DownloaderArgsForCobra(addCustomCmd, &downloaderArgs)

	// Add the subcommands to the add command
	addCmd.AddCommand(addCustomCmd)
}
