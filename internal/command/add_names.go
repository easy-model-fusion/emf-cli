package command

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/huggingface"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
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
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// TODO: Get flags or default values

	// Get all existing models
	existingModels, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	var selectedModelNames []string
	var selectedModels []model.Model

	// Add models passed in args
	if len(args) > 0 {

		// Remove all the duplicates
		args = utils.SliceRemoveDuplicates(args)

		var notFoundModelNames []string
		var existingModelNames []string

		// Fetching the requested models
		for _, name := range args {
			// Verify if model already exists in the project
			exist := model.ContainsByName(existingModels, name)
			if exist {
				// Model already exists : skipping to the next one
				existingModelNames = append(existingModelNames, name)
				continue
			}

			// Verify if the model is a valid hugging face model
			apiModel, err := app.H().GetModel(name)
			if err != nil {
				// Model not found : skipping to the next one
				notFoundModelNames = append(notFoundModelNames, name)
				continue
			}

			// Adding valid models
			selectedModels = append(selectedModels, apiModel)
		}

		// Indicate the models that couldn't be found
		if len(notFoundModelNames) > 0 {
			pterm.Warning.Printfln(fmt.Sprintf("The following models(s) couldn't be found "+
				"and will be ignored : %s", notFoundModelNames))
		}
		// Indicate the models that already exists
		if len(existingModelNames) > 0 {
			pterm.Warning.Printfln(fmt.Sprintf("The following model(s) already exist(s) "+
				"and will be ignored : %s", existingModelNames))
		}
	}

	// If no models entered by user or if user entered -s/--select
	if displayModels || len(args) == 0 {
		// Get selected tags
		selectedTags := selectTags()
		if len(selectedTags) == 0 {
			pterm.Warning.Println("Please select a model type")
			runAddByNames(cmd, args)
		}
		// Get selected models
		selectedModels, err = selectModels(selectedTags, selectedModels, existingModels)
	}

	// Verify if selected models is not empty
	if selectedModels == nil {
		pterm.Warning.Println("No models selected")
		return
	}

	// Get all selected models names
	selectedModelNames = model.GetNames(selectedModels)

	// User choose the models he wishes to install now
	selectedModels = selectModelsToInstall(selectedModels, selectedModelNames)
	utils.DisplaySelectedItems(selectedModelNames)

	// Download the models
	models, failedModels := config.DownloadModels(selectedModels)

	// No models were downloaded : stopping there
	if !model.Empty(failedModels) {
		pterm.Error.Println("These models couldn't be downloaded", model.GetNames(failedModels))
		return
	}
	if model.Empty(models) {
		pterm.Info.Println("There isn't any model to add to the configuration file.")
		return
	}

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err = config.AddModels(models)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
}

// selectModels displays a multiselect of models from which the user will choose to add to his project
func selectModels(tags []string, currentSelectedModels []model.Model, existingModels []model.Model) ([]model.Model, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Listing all models with selected tags...")
	var allModelsWithTags []model.Model
	// Get list of models with current tags
	for _, tag := range tags {
		apiModels, err := app.H().GetModels(tag, 0)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while fetching the models from hugging face api: %s", err))
			return nil, fmt.Errorf("error while calling api endpoint")
		}
		allModelsWithTags = append(allModelsWithTags, apiModels...)
	}
	spinner.Success()

	// Excluding models entered in args + configuration file models
	modelsToExclude := append(currentSelectedModels, existingModels...)
	availableModels := model.Difference(allModelsWithTags, modelsToExclude)

	// Build a multiselect with each model name
	availableModelNames := model.GetNames(availableModels)
	message := "Please select the model(s) to be added"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModelNames := utils.DisplayInteractiveMultiselect(message, availableModelNames, checkMark, true)

	// No new model was selected : returning the input state
	if len(selectedModelNames) == 0 {
		return currentSelectedModels, nil
	}

	// Get newly selected models
	selectedModels := model.GetModelsByNames(availableModels, selectedModelNames)

	// returns newly selected models + models entered in args
	return append(currentSelectedModels, selectedModels...), nil
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
		currentModel.IsDownloaded = utils.SliceContainsItem(installsToExclude, currentModel.Name)
		updatedModels = append(updatedModels, currentModel)
	}

	return updatedModels
}

func init() {
	app.InitHuggingFace(huggingface.BaseUrl, "")
	// Add --select flag to the add default command
	addByNamesCmd.Flags().BoolVarP(&displayModels, "select", "s", false, "Select models to add")
	// Group the add by names subcommand to the add command
	addCmd.AddCommand(addByNamesCmd)
}
