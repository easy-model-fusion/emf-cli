package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
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
	logger := app.L().WithTime(false)

	// Load the configuration file
	err := config.Load(".")
	if err != nil {
		logger.Error("Error reading config file:" + err.Error())
	}

	var selectedModels []model.Model
	if selectFiles || len(args) == 0 {
		selectedTags := selectTags()
		if len(selectedTags) == 0 {
			app.L().WithTime(false).Warn("Please select a model type")
			runAdd(cmd, args)
		}
		selectedModels = selectModels(selectedTags)
		if selectedModels == nil {
			app.L().WithTime(false).Warn("No models selected")
			return
		}
	} else {
	}

	// TODO install models with addToBinary => true

	err = config.AddModel(selectedModels)
	if err != nil {
		logger.Error("Error while adding models into config file:" + err.Error())
	}
}

func selectModels(tags []string) []model.Model {
	// Get list of models with current tag
	var models []model.Model
	for _, tag := range tags {
		apiModels, err := huggingface.GetModels(nil, tag, nil)
		if err != nil {
			app.L().Fatal("error while calling api endpoint")
		}
		models = append(models, apiModels...)
	}

	// Get existent models from configuration file
	currentModelsNames, err := config.GetModelsNames()
	if err != nil {
		app.L().Fatal("error while getting current models")
	}

	// Remove existent models from list of models to add
	models = config.RemoveModelsFromList(models, currentModelsNames)

	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}

	message := "Please select the model(s) to be added"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModelNames := utils.DisplayInteractiveMultiselect(message, modelNames, checkMark, true)
	if len(selectedModelNames) == 0 {
		return nil
	}
	message = "Please select the model(s) that you don't wish to install directly"
	checkMark = &pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	installsToExclude := utils.DisplayInteractiveMultiselect(message, selectedModelNames, checkMark, false)

	utils.DisplaySelectedItems(selectedModelNames)
	var selectedModels []model.Model

	for _, currentModel := range models {
		if utils.ArrayStringContainsItem(selectedModelNames, currentModel.Name) {
			currentModel.AddToBinary = !utils.ArrayStringContainsItem(installsToExclude, currentModel.Name)
			selectedModels = append(selectedModels, currentModel)
		}
	}

	return selectedModels
}

func selectTags() []string {
	// Build a multiselect with each model name
	tags := model.AllTags

	message := "Please select the type of models you want to add"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedTags := utils.DisplayInteractiveMultiselect(message, tags, checkMark, true)

	return selectedTags
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
	// Add --select flag to the add command
	addCmd.Flags().BoolVarP(&selectFiles, "select", "s", false, "Select files to add")
}
