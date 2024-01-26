package command

import (
	"github.com/easy-model-fusion/client/internal/config"
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

		// Build a multiselect with each model name
		var modelNames []string
		for _, item := range models {
			modelNames = append(modelNames, item.Name)
		}

		// Create a new interactive multiselect printer with the options
		// Disable the filter and set the keys for confirming and selecting options
		printer := pterm.DefaultInteractiveMultiselect.
			WithOptions(modelNames).
			WithFilter(false).
			WithCheckmark(&pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")})

		// Show the interactive multiselect and get the selected options
		selectedModels, _ = printer.Show()

		// Print the selected options, highlighted in green.
		pterm.Info.Printfln("Selected options: %s", pterm.Green(selectedModels))

	} else {
		// selected models from args
		selectedModels = make([]string, len(args))
		copy(selectedModels, args)
	}

	// remove selected models
	_ = config.RemoveModels(models, selectedModels)
}

func init() {
	removeCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all models")
	rootCmd.AddCommand(removeCmd)
}
