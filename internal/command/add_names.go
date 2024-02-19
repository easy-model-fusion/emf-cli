package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/huggingface"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const cmdAddNamesTitle = "names"

// addByNamesCmd represents the add model by names command
var addByNamesCmd = &cobra.Command{
	Use:   cmdAddNamesTitle + " <model name> [<other model names>...]",
	Short: "Add model(s) by name to your project",
	Long:  `Add model(s) by name to your project`,
	Run:   runAddByNames,
}

// displayModels indicates if the multiselect of models should be displayed or not
var displayModels bool

// runAddByNames runs the add command for adding model by names
func runAddByNames(cmd *cobra.Command, args []string) {
	if config.GetViperConfig() != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// TODO: Get flags or default values
	app.InitHuggingFace(huggingface.BaseUrl, "")

	var selectedModelNames []string
	var selectedModels []model.Model

	// Add models passed in args
	if len(args) > 0 {
		// Fetching the requested models
		for _, name := range args {
			apiModel, err := app.H().GetModel(name)
			if err != nil {
				// Model not recognized : skipping to the next one
				continue
			}
			// Saving the model data in the variables
			selectedModels = append(selectedModels, apiModel)
			selectedModelNames = append(selectedModelNames, name)
		}

		// Remove all the duplicates
		selectedModelNames = utils.SliceRemoveDuplicates(selectedModelNames)

		// Indicate the models that couldn't be found
		notFound := utils.SliceDifference(args, selectedModelNames)
		if len(notFound) != 0 {
			pterm.Warning.Printfln(fmt.Sprintf("The following models couldn't be found and will be ignored : %s", notFound))
		}
	}

	// If no models entered by user or if user entered -s/--select
	if displayModels || len(args) == 0 {
		// Get selected tags
		selectedTags := selectTags()
		if len(selectedTags) == 0 {
			app.L().WithTime(false).Warn("Please select a model type")
			runAddByNames(cmd, args)
		}
		// Get selected models
		var modelNames []string
		selectedModels, modelNames = selectModels(selectedTags, selectedModelNames)
		if selectedModels == nil {
			app.L().WithTime(false).Warn("No models selected")
			return
		}
		selectedModelNames = append(selectedModelNames, modelNames...)
	}

	// Process the models to only keep the valid ones
	selectedModels, err := processSelectedModels(selectedModels)
	if err != nil {
		return
	}

	// Check if any model is still valid
	if len(selectedModels) == 0 {
		pterm.Info.Println("None of the requested models can be added.")
		return
	}

	// Update the selected model names
	selectedModelNames = model.GetNames(selectedModels)

	// User choose the models he wishes to install now
	selectedModels = selectModelsToInstall(selectedModels, selectedModelNames)
	utils.DisplaySelectedItems(selectedModelNames)

	// Download the models
	var models []model.Model
	for _, item := range selectedModels {
		result, ok := config.DownloadModel(item)
		if !ok {
			continue
		}
		models = append(models, result)
	}

	// No models were downloaded : stopping there
	if model.Empty(models) {
		pterm.Info.Println("There isn't any model to add to the configuration file.")
		return
	}

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err = config.AddModel(models)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
}

// processSelectedModels process the selected models and only keep the valid ones
func processSelectedModels(selectedModels []model.Model) ([]model.Model, error) {
	// Get the models from the configuration file
	configModels, err := config.GetModels()
	if err != nil {
		return nil, err
	}

	// Filter the requested models that have already been added
	alreadyAdded := model.Union(configModels, selectedModels)
	if len(alreadyAdded) != 0 {
		pterm.Warning.Println(fmt.Sprintf("The following models have already been added and will be ignored : %s", model.GetNames(alreadyAdded)))
	}

	// Filter the ones that haven't been added yet
	toBeAdded := model.Difference(selectedModels, alreadyAdded)

	return toBeAdded, nil
}

// selectModels displays a multiselect of models from which the user will choose to add to his project
func selectModels(tags []string, currentSelectedModels []string) ([]model.Model, []string) {
	var models []model.Model
	// Get list of models with current tags
	for _, tag := range tags {
		apiModels, err := app.H().GetModels(tag, 0)
		if err != nil {
			app.L().Fatal("error while calling api endpoint")
		}
		models = append(models, apiModels...)
	}

	// Get existent models from configuration file
	configModels, err := config.GetModels()
	if err != nil {
		app.L().Fatal("error while getting current models")
	}
	configModelNames := model.GetNames(configModels)

	// Remove existent models from list of models to add
	// Remove already selected models from list of models to add (in case user entered add model_name -s)
	existingModels := model.GetModelsByNames(models, append(configModelNames, currentSelectedModels...))
	models = model.Difference(models, existingModels)

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
		if utils.SliceContainsItem(selectedModelNames, currentModel.Name) {
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

// selectModelsToInstall returns updated models objects with excluded/included from binary
func selectModelsToInstall(models []model.Model, modelNames []string) []model.Model {
	// Build a multiselect with each selected model name to exclude/include in the binary
	message := "Please select the model(s) to install now"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Blue("-")}
	installsToExclude := utils.DisplayInteractiveMultiselect(message, modelNames, checkMark, false)
	var updatedModels []model.Model
	for _, currentModel := range models {
		currentModel.AddToBinary = utils.SliceContainsItem(installsToExclude, currentModel.Name)
		updatedModels = append(updatedModels, currentModel)
	}

	return updatedModels
}

func init() {
	// Add --select flag to the add default command
	addByNamesCmd.Flags().BoolVarP(&displayModels, "select", "s", false, "Select models to add")
	// Group the add by names subcommand to the add command
	addCmd.AddCommand(addByNamesCmd)
}
