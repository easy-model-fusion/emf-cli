package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

// RunAdd runs the add command to add models by name
func RunAdd(args []string) {

}

func processAdd(args []string) (warning string, err error) {
	err = config.GetViperConfig(config.FilePath)
	if err != nil {
		return warning, err
	}

	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// Get all existing models
	existingModels, err := config.GetModels()
	if err != nil {
		return warning, err
	}

	var selectedModel model.Model

	// Add models passed in args
	if len(args) == 1 {
		name := args[0]
		// Verify if model already exists in the project
		exist := existingModels.ContainsByName(name)
		if exist {
			// Model already exists
			err = fmt.Errorf("the following model already exist and will be ignored : %s", name)
			return warning, err
		}

		// Verify if the model is a valid hugging face model
		huggingfaceModel, err := app.H().GetModelById(name)
		if err != nil {
			// Model not found : skipping to the next one
			err = fmt.Errorf("Model %s not valid : "+err.Error(), name)
			return warning, err
		}

		// Map API response to model.Model
		selectedModel = model.FromHuggingfaceModel(huggingfaceModel)
	} else if len(args) > 1 {
		err = fmt.Errorf("you can enter only one model at a time")
		return warning, err
	} else {
		// If no models entered by user or if user entered -s/--select
		// Get selected tags
		selectedTags := selectTags()
		if len(selectedTags) == 0 {
			warning = "Please select a model type"
			RunAdd(args)
		}
		// Get selected models
		availableModels, err := getModelsList(selectedTags, existingModels)
		if err != nil {
			return warning, err
		}
		selectedModel = selectModel(availableModels)
	}

	// User choose if he wishes to install the model directly
	message := fmt.Sprintf("Do you wish to install %s directly?", selectedModel.Name)
	selectedModel.AddToBinaryFile = app.UI().AskForUsersConfirmation(message)

	// Validate model for download
	warningMessage, valid, err := config.Validate(selectedModel)
	if !valid {
		return warningMessage, err
	}

	// Try to download model
	err = downloadModel(selectedModel)
	if err != nil {
		return warning, err
	}

	// Add models to configuration file
	spinner := app.UI().StartSpinner("Adding model to configuration file...")
	err = config.AddModels(model.Models{selectedModel})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while adding the model to the configuration file: %s", err))
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

	return warning, err
}

// selectModel displays a selector of models from which the user will choose to add to his project
func selectModel(models model.Models) model.Model {
	// Build a selector with each model name
	availableModelNames := models.GetNames()
	message := "Please select the model(s) to be added"
	selectedModelName := app.UI().DisplayInteractiveSelect(message, availableModelNames, true)

	// Get newly selected models
	selectedModels := models.FilterWithNames([]string{selectedModelName})

	// returns newly selected models + models entered in args
	return selectedModels[0]
}

// selectTags displays a multiselect to help the user choose the model types
func selectTags() []string {
	// Build a multiselect with each tag name
	message := "Please select the type of models you want to add"
	checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedTags := app.UI().DisplayInteractiveMultiselect(message, huggingface.AllTagsString(), checkMark, false, true)

	return selectedTags
}

// getModelsList get list of models to display
func getModelsList(tags []string, existingModels model.Models) (model.Models, error) {
	allModelsWithTags, err := app.H().GetModelsByMultiplePipelineTags(tags)
	// Map API responses to model.Models
	var mappedModels model.Models
	for _, huggingfaceModel := range allModelsWithTags {
		mappedModel := model.FromHuggingfaceModel(huggingfaceModel)
		mappedModels = append(mappedModels, mappedModel)
	}
	if err != nil {
		return model.Models{}, fmt.Errorf("error while calling api endpoint")
	}

	return mappedModels.Difference(existingModels), nil
}

func downloadModel(selectedModel model.Model) error {
	// Prepare the script arguments
	downloaderArgs := downloadermodel.Args{
		ModelName:     selectedModel.Name,
		ModelModule:   string(selectedModel.Module),
		ModelClass:    selectedModel.Class,
		DirectoryPath: app.DownloadDirectoryPath,
	}

	var success bool
	if selectedModel.AddToBinaryFile {
		// Downloading model
		success = selectedModel.Download(downloaderArgs)
	} else {
		// Getting model configuration
		success = selectedModel.GetConfig(downloaderArgs)
	}

	if !success {
		return fmt.Errorf("this model %s couldn't be downloaded", selectedModel.Name)
	}

	return nil
}
