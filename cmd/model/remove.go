package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/model"
	"github.com/spf13/cobra"
)

var modelRemoveAllFlag bool

// modelRemoveCmd represents the model remove command
var modelRemoveCmd = &cobra.Command{
	Use:   "remove <model name> [<other model names>...]",
	Short: "Remove one or more models",
	Long:  "Remove one or more models",
	Run:   runModelRemove,
}

func init() {
	// Adding the command's flags
	modelRemoveCmd.Flags().BoolVarP(&modelRemoveAllFlag, "all", "a", false, "Remove all models")
}

// runModelRemove runs the model remove command
func runModelRemove(cmd *cobra.Command, args []string) {
	modelcontroller.RunModelRemove(args, modelRemoveAllFlag)
}
