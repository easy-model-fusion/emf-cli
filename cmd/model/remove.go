package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/pterm/pterm"
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
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	// Remove all models
	if modelRemoveAllFlag {
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
	if err != nil || model.Empty(models) {
		pterm.Info.Printfln("There is no models to be removed.")
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
	err = config.RemoveModelsByNames(models, selectedModels)
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

	checkMark := ui.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	message := "Please select the model(s) to be deleted"
	modelsToDelete := app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false)
	app.UI().DisplaySelectedItems(modelsToDelete)
	return modelsToDelete
}
