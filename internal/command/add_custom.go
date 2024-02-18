package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/script"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/spf13/cobra"
)

const cmdAddCustomTitle = "custom"

// addCustomCmd represents the add custom model command
var addCustomCmd = &cobra.Command{
	Use:   cmdAddCustomTitle + " <name> <diffusers|transformers>",
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

	// TODO : write command
}

func init() {

	// Bind cobra args to the downloader script args
	script.DownloaderArgsForCobra(addCustomCmd, &downloaderArgs)

	// Add the subcommands to the add command
	addCmd.AddCommand(addCustomCmd)
}