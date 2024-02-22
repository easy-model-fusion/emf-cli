package cmdmodel

import (
	cmdmodeladd "github.com/easy-model-fusion/emf-cli/cmd/model/add"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const modelCommandName string = "model"

// ModelCmd represents the model command
var ModelCmd = &cobra.Command{
	Use:   modelCommandName,
	Short: "Palette that contains model based commands",
	Long:  "Palette that contains model based commands",
	Run:   runModel,
}

func init() {
	// Preparing to use the hugging face API
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Adding the subcommands
	ModelCmd.AddCommand(modelRemoveCmd)
	ModelCmd.AddCommand(cmdmodeladd.ModelAddCmd)
}

// runModel runs model command
func runModel(cmd *cobra.Command, args []string) {

	// Running command as palette : allowing user to choose subcommand
	err := cobrautil.RunCommandAsPalette(cmd, args, modelCommandName, []string{})
	if err != nil {
		pterm.Error.Println("Something went wrong :", err)
	}
}
