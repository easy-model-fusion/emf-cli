package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/sdk"
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
	if config.GetViperConfig() != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	// Remove all models
	if allFlag {
		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		} else {
			pterm.Error.Printfln("Operation failed.")
		}
		return
	}

	// Declare variables
	var selectedModels []string
	var models, err = config.GetModels()

	// Check fetched models : cannot be null or empty
	if err != nil || config.IsModelsEmpty(models) {
		pterm.Info.Printfln("No models to remove.")
		return
	}

	// No args, asks for model names
	if len(args) == 0 {
		// Get selected models from multiselect
		selectedModels = selectModelsToDelete(models)
	} else {
		// Get the selected models from the args
		selectedModels = make([]string, len(args))
		copy(selectedModels, args)
	}

	// Remove the selected models
	err = config.RemoveSelectedModelsByName(models, selectedModels)
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func selectModelsToDelete(currentModels []model.Model) []string {
	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range currentModels {
		modelNames = append(modelNames, item.Name)
	}

	checkMark := &pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	message := "Please select the model(s) to be deleted"
	modelsToDelete := utils.DisplayInteractiveMultiselect(message, modelNames, checkMark, false)
	utils.DisplaySelectedItems(modelsToDelete)
	return modelsToDelete
}

func init() {
	removeCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Remove all models")
	rootCmd.AddCommand(removeCmd)
}
