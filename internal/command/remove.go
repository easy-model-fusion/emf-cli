package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/spf13/cobra"
)

var allFlag bool

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <model name> [<other model names>...]",
	Short: "Remove one or more models",
	Run:   runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {

	// If allFlag is true, remove all models
	if allFlag {
		err := config.RemoveAllModels()
		if err != nil {
			return
		}
		return
	}

	var modelsString string
	var modelsSlice []string

	// No args, asks for model names
	if len(args) == 0 {
		// TODO : multiselect of all downloaded models
		modelsString = utils.AskForUsersInput("Indicate the models to remove")
		modelsSlice = utils.ArrayFromString(modelsString)
	} else {
		modelsSlice = make([]string, len(args))
		copy(modelsSlice, args)
	}

	err := config.RemoveModel(modelsSlice)
	if err != nil {
		return
	}
}

func init() {
	removeCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all models")
	rootCmd.AddCommand(removeCmd)
}
