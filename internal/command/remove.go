package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var allFlag bool

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <model name> [<other model names>...]",
	Short: "Remove one or more models",
	Long:  "Remove one or more models",
	Run:   runRemove,
}

func runRemove(cmd *cobra.Command, args []string) {
	logger := app.L().WithTime(false)

	// Load the configuration file
	test := config.Load(".")
	if test != nil {
		logger.Error("Error reading config file:" + test.Error())
	}

	// remove all models
	if allFlag {
		_ = config.RemoveAllModels()
		return
	}

	// Declare variables
	var selectedModels []string
	var models, err = config.GetModels()

	// Check fetched models : cannot be null or empty
	if err != nil || config.IsModelsEmpty(models) {
		return
	}

	// No args, asks for model names
	if len(args) == 0 {
		// Get selected models from multiselect
		selectedModels = selectModelsToDelete(models)
	} else {
		// selected models from args
		selectedModels = make([]string, len(args))
		copy(selectedModels, args)
	}

	// remove selected models
	_ = config.RemoveModels(models, selectedModels)
}

func selectModelsToDelete(currentModels []model.Model) []string {
	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range currentModels {
		modelNames = append(modelNames, item.Name)
	}

	checkMark := &pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	return utils.DisplayInteractiveMultiselect(modelNames, checkMark, false)
}

func init() {
	removeCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all models")
	rootCmd.AddCommand(removeCmd)
}
