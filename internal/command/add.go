package command

import (
	"fmt"
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
	Use:   "add <model name> [<other model names>...]",
	Short: "Add model(s) to your project",
	Long:  `Add model(s) to your project`,
	Args:  config.ValidModelName(),
	Run:   runAdd,
}

// displayModels indicates if the multiselect of models should be displayed or not
var displayModels bool

// runAdd runs add command
func runAdd(cmd *cobra.Command, args []string) {
	if config.GetViperConfig() != nil {
		return
	}
	var selectedModelNames []string

	var selectedModels []model.Model

	// Add models passed in args
	if len(args) > 0 {
		for _, name := range args {
			apiModel, err := huggingface.GetModel(name, nil)
			if err != nil {
				app.L().WithTime(false).Error(fmt.Sprintf("while getting model %v", name))
				return
			}
			selectedModels = append(selectedModels, *apiModel)
		}
		selectedModelNames = append(selectedModelNames, args...)
	}

	// If no models entered by user or if user entered -s/--select
	if displayModels || len(args) == 0 {
		// Get selected tags
		selectedTags := selectTags()
		if len(selectedTags) == 0 {
			app.L().WithTime(false).Warn("Please select a model type")
			runAdd(cmd, args)
		}
		// Get selected models
		selectedModels, modelsNames := selectModels(selectedTags, selectedModelNames)
		if selectedModels == nil {
			app.L().WithTime(false).Warn("No models selected")
			return
		}
		selectedModelNames = append(selectedModelNames, modelsNames...)
	}

	// User choose either to exclude or include models in binary
	selectedModels = selectExcludedModelsFromInstall(selectedModels, selectedModelNames)

	// TODO install models with addToBinary => true

	// Add models to configuration file
	err := config.AddModel(selectedModels)

	if err == nil {
		// Display the selected models
		utils.DisplaySelectedItems(selectedModelNames)
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

// selectModels displays a multiselect of models from which the user will choose to add to his project
func selectModels(tags []string, currentSelectedModels []string) ([]model.Model, []string) {
	var models []model.Model
	// Get list of models with current tags
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
	// Remove already selected models from list of models to add (in case user entered add model_name -s)
	models = config.RemoveModelsFromList(models, append(currentModelsNames, currentSelectedModels...))

	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}

	message := "Please select the model(s) to be added"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModelNames := utils.DisplayInteractiveMultiselect(message, modelNames, checkMark, true)
	if len(selectedModelNames) == 0 {
		return nil, nil
	}

	// Get models objects from models names
	var selectedModels []model.Model

	for _, currentModel := range models {
		if utils.ArrayStringContainsItem(selectedModelNames, currentModel.Name) {
			selectedModels = append(selectedModels, currentModel)
		}
	}

	return selectedModels, selectedModelNames
}

// selectTags displays a multiselect to help the user choose the model types
func selectTags() []string {
	// Build a multiselect with each tag name
	tags := model.AllTags

	message := "Please select the type of models you want to add"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedTags := utils.DisplayInteractiveMultiselect(message, tags, checkMark, true)

	return selectedTags
}

// selectExcludedModelsFromInstall returns updated models objects with excluded/included from binary
func selectExcludedModelsFromInstall(models []model.Model, modelsNames []string) []model.Model {
	// Build a multiselect with each selected model name to exclude/include in the binary
	message := "Please select the model(s) that you don't wish to install directly"
	checkMark := &pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	installsToExclude := utils.DisplayInteractiveMultiselect(message, modelsNames, checkMark, false)
	var updatedModels []model.Model
	for _, currentModel := range models {
		currentModel.AddToBinary = !utils.ArrayStringContainsItem(installsToExclude, currentModel.Name)
		updatedModels = append(updatedModels, currentModel)
	}

	return updatedModels
}

func init() {
	// Add --select flag to the add command
	addCmd.Flags().BoolVarP(&displayModels, "select", "s", false, "Select models to add")
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
