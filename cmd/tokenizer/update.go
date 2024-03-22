package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the model update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update <model name> [<other model names>...]",
	Short: "Update one or more models",
	Long:  "Update one or more models",
	Args:  cobra.MinimumNArgs(1),
	Run:   runModelUpdate,
}

// runModelUpdate runs the model update command
func runModelUpdate(cmd *cobra.Command, args []string) {
	sdk.SendUpdateSuggestion()
	modelcontroller.RunModelUpdate(args)
}
