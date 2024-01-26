package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// addCmd represents the add model(s) command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add model(s) to your project",
	Long:  `Add model(s) to your project`,
	Run:   runAdd,
}

var selectFiles bool

func runAdd(cmd *cobra.Command, args []string) {
	var selectedModels []string
	if selectFiles {
		selectedModels = selectModels()
	} else {
	}

	config.Load(".")
	config.AddModel(selectedModels)
}

func selectModels() []string {
	limit := 10
	models, err := huggingface.GetModels(&limit, model.TEXT_TO_IMAGE, nil)
	if err != nil {
		app.L().Fatal("api call error")
	}

	var modelsNames []string
	for i := 0; i < len(models); i++ {
		modelsNames = append(modelsNames, models[i].Name)
	}

	selectedModels, _ := pterm.DefaultInteractiveMultiselect.WithOptions(modelsNames).Show()

	// Print the selected models, highlighted in green.
	pterm.Info.Printfln("Selected models: %s", pterm.Green(selectedModels))

	return selectedModels
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
	// Add --select flag to the add command
	addCmd.Flags().BoolVarP(&selectFiles, "select", "s", false, "Select files to add")
}
