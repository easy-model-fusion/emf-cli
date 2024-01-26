package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
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
	logger := app.L().WithTime(false)

	// Load the configuration file
	err := config.Load(".")
	if err != nil {
		logger.Error("Error reading config file:" + err.Error())
	}

	var selectedModels []model.Model
	if selectFiles {
		selectedModels = selectModels(model.TEXT_TO_IMAGE)
	} else {
	}

	err = config.AddModel(selectedModels)
	if err != nil {
		logger.Error("Error while adding models into config file:" + err.Error())
	}
}

func selectModels(tag string) []model.Model {
	models, err := huggingface.GetModels(nil, tag, nil)
	if err != nil {
		app.L().Fatal("api call error")
	}

	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}

	selectedModelNames := utils.DisplayInteractiveMultiselect(modelNames)
	var selectedModels []model.Model

	for _, currentModel := range models {
		if utils.ArrayStringContainsItem(selectedModelNames, currentModel.Name) {
			selectedModels = append(selectedModels, currentModel)
		}
	}

	return selectedModels
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
	// Add --select flag to the add command
	addCmd.Flags().BoolVarP(&selectFiles, "select", "s", false, "Select files to add")
}
