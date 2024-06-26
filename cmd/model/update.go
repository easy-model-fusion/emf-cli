package cmdmodel

import (
	modelcontroller "github.com/easy-model-fusion/emf-cli/internal/controller/model"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the model update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update <model name> [<other model names>...]",
	Short: "Update one or more models",
	Long:  "Update one or more models",
	Run:   runModelUpdate,
}

var (
	accessToken        string
	authorizeOverwrite bool
)

func init() {
	modelUpdateCmd.Flags().BoolVarP(&authorizeOverwrite, "yes", "y", false, "Automatic yes to prompts")
	modelUpdateCmd.Flags().StringVarP(&accessToken, "access-token", "a", "", "Access token for gated models")
}

// runModelUpdate runs the model update command
func runModelUpdate(cmd *cobra.Command, args []string) {
	modelcontroller.RunModelUpdate(args, authorizeOverwrite, accessToken)
}
