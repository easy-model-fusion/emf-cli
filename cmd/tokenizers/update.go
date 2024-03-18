package tokenizers

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the model update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update <model name> [<tokenizers>...]",
	Short: "Update one or more tokenizers",
	Long:  "Update one or more tokenizers",
	Run:   runTokenizerUpdate,
}

// runTokenizerUpdate runs the model remove command
func runTokenizerUpdate(cmd *cobra.Command, args []string) {
	controller.TokenizerUpdateCmd(args)
}
