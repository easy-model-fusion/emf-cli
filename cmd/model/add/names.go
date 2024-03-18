package cmdmodeladd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// addByNamesCmd represents the add model by names command
var addByNamesCmd = &cobra.Command{
	Use:   "names <model name> [<other model names>...]",
	Short: "Add model(s) by name to your project",
	Long:  `Add model(s) by name to your project`,
	Run:   runAddByNames,
}

// displayModels indicates if the multiselect of models should be displayed or not
var displayModels bool

func init() {
	// Add --select flag to the model add command
	addByNamesCmd.Flags().BoolVarP(&displayModels, "select", "s", false, "Select models to add")
}

// runAddByNames runs the add command to add models by name
func runAddByNames(cmd *cobra.Command, args []string) {
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// Get all existing models
	existingModels, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	var selectedModelNames []string
	var selectedModels model.Models

	// Add models passed in args
	if len(args) > 0 {

		// Remove all the duplicates
		args = stringutil.SliceRemoveDuplicates(args)

		var notFoundModelNames []string
		var existingModelNames []string

		// Fetching the requested models
		for _, name := range args {
			// Verify if model already exists in the project
			exist := existingModels.ContainsByName(name)
			if exist {
				// Model already exists : skipping to the next one
				existingModelNames = append(existingModelNames, name)
				continue
			}

			// Verify if the model is a valid hugging face model
			huggingfaceModel, err := app.H().GetModelById(name)
			if err != nil {
				// Model not found : skipping to the next one
				pterm.Warning.Printfln("Model %s not valid : "+err.Error(), name)
				notFoundModelNames = append(notFoundModelNames, name)
				continue
			}
			// Map API response to model.Model
			modelMapped := model.FromHuggingfaceModel(huggingfaceModel)
			// Saving the model data in the variables
			selectedModels = append(selectedModels, modelMapped)
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
		if err != nil {
			pterm.Error.Println(err.Error())
			return
		}
	}

	// Verify if selected models is not empty
	if selectedModels == nil {
		pterm.Warning.Println("No models selected")
		return
	}

	// Get all selected models names
	selectedModelNames = selectedModels.GetNames()

	// User choose the models he wishes to install now
	selectedModels = selectModelsToInstall(selectedModels, selectedModelNames)
	app.UI().DisplaySelectedItems(selectedModelNames)

	// Download the models
	var models model.Models
	var failedModels model.Models
	for _, currentModel := range selectedModels {

		// Validate model for download
		if !config.Validate(currentModel) {
			failedModels = append(failedModels, currentModel)
			continue
		}

		// Prepare the script arguments
		downloaderArgs := downloadermodel.Args{
			ModelName:     currentModel.Name,
			ModelModule:   string(currentModel.Module),
			ModelClass:    currentModel.Class,
			DirectoryPath: app.DownloadDirectoryPath,
		}

		var success bool
		if currentModel.AddToBinaryFile {
			// Downloading model
			success = currentModel.Download(downloaderArgs)
		} else {
			// Getting model configuration
			success = currentModel.GetConfig(downloaderArgs)
		}

		if !success {
			// Reset in case the download fails
			currentModel.AddToBinaryFile = false
			failedModels = append(failedModels, currentModel)
		} else {
			models = append(models, currentModel)
		}
	}

	// Indicate models that failed to download
	if !failedModels.Empty() {
		pterm.Error.Println("These models couldn't be downloaded", failedModels.GetNames())
		return
	}
	// No models were downloaded : stopping there
	if models.Empty() {
		pterm.Info.Println("There isn't any model to add to the configuration file.")
		return
	}

	// Add models to configuration file
	spinner := app.UI().StartSpinner("Writing models to configuration file...")
	err = config.AddModels(models)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}

	// Attempt to generate code
	spinner = app.UI().StartSpinner("Generating python code...")
	err = config.GenerateExistingModelsPythonCode()
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while generating python code for added models: %s", err))
	} else {
		spinner.Success()
	}
}

// selectModels displays a multiselect of models from which the user will choose to add to his project
func selectModels(tags []string, currentSelectedModels model.Models, existingModels model.Models) (model.Models, error) {
	spinner := app.UI().StartSpinner("Listing all models with selected tags...")
	var allModelsWithTags model.Models
	// Get list of models with current tags
	for _, tag := range tags {
		huggingfaceModels, err := app.H().GetModelsByPipelineTag(huggingface.PipelineTag(tag), 0)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error while fetching the models from hugging face api: %s", err))
			return nil, fmt.Errorf("error while calling api endpoint")
		}
		// Map API responses to model.Models
		var mappedModels model.Models
		for _, huggingfaceModel := range huggingfaceModels {
			mappedModel := model.FromHuggingfaceModel(huggingfaceModel)
			mappedModels = append(mappedModels, mappedModel)
		}
		allModelsWithTags = append(allModelsWithTags, mappedModels...)
	}
	spinner.Success()

	// Excluding models entered in args + configuration file models
	modelsToExclude := append(currentSelectedModels, existingModels...)
	availableModels := allModelsWithTags.Difference(modelsToExclude)

	// Build a multiselect with each model name
	availableModelNames := availableModels.GetNames()
	message := "Please select the model(s) to be added"
	checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModelNames := app.UI().DisplayInteractiveMultiselect(message, availableModelNames, checkMark, false, true)

	// No new model was selected : returning the input state
	if len(selectedModelNames) == 0 {
		return currentSelectedModels, nil
	}

	// Get newly selected models
	selectedModels := availableModels.FilterWithNames(selectedModelNames)

	// returns newly selected models + models entered in args
	return append(currentSelectedModels, selectedModels...), nil
}

// selectTags displays a multiselect to help the user choose the model types
func selectTags() []string {
	// Build a multiselect with each tag name
	message := "Please select the type of models you want to add"
	checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedTags := app.UI().DisplayInteractiveMultiselect(message, huggingface.AllTagsString(), checkMark, false, true)

	return selectedTags
}

// selectModelsToInstall returns updated models objects with excluded/included from binary
func selectModelsToInstall(models model.Models, modelNames []string) model.Models {
	// Build a multiselect with each selected model name to exclude/include in the binary
	message := "Please select the model(s) to install later"
	checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Blue("-")}
	installsToExclude := app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false)
	var updatedModels model.Models
	for _, currentModel := range models {
		currentModel.AddToBinaryFile = !stringutil.SliceContainsItem(installsToExclude, currentModel.Name)
		updatedModels = append(updatedModels, currentModel)
	}

	return updatedModels
}
