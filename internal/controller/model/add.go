package modelcontroller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/hfinteface"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
)

// RunAdd runs the add command to add models by name
func RunAdd(args []string, customArgs downloadermodel.Args, yes bool) {
	selectedModel, err := getRequestedModel(args)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}
	if selectedModel.Name == "" {
		pterm.Warning.Println("Please select a model type")
		RunAdd(args, customArgs, yes)
		return
	}

	warningMessage, err := processAdd(selectedModel, customArgs, yes)
	if warningMessage != "" {
		pterm.Warning.Println(warningMessage)
	}
	if err != nil {
		pterm.Error.Println(err.Error())
	}
}

// getRequestedModel returns the model to be added
func getRequestedModel(args []string) (model.Model, error) {
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return model.Model{}, err
	}

	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	// Get all existing models
	existingModels, err := config.GetModels()
	if err != nil {
		return model.Model{}, err
	}

	var selectedModel model.Model

	// Add models passed in args
	if len(args) == 1 {
		name := args[0]
		// Verify if model already exists in the project
		exist := existingModels.ContainsByName(name)
		if exist {
			// Model already exists
			return model.Model{}, fmt.Errorf("the following model already exist and will be ignored : %s", name)
		}

		// Verify if the model is a valid hugging face model
		huggingfaceModel, err := hfinteface.GetModelById(name)
		if err != nil {
			// Model not found
			return model.Model{}, fmt.Errorf("Model %s not valid : "+err.Error(), name)
		}

		// Map API response to model.Model
		selectedModel = model.FromHuggingfaceModel(huggingfaceModel)
	} else if len(args) > 1 {
		return model.Model{}, fmt.Errorf("you can enter only one model at a time")
	} else {
		// If no models entered by user or if user entered -s/--select
		// Get selected tags
		var selectedTags []string
		selectedTags = selectTags()
		if len(selectedTags) == 0 {
			return model.Model{}, nil
		}
		// Get selected models
		spinner := app.UI().StartSpinner("Listing all models with selected tags...")
		availableModels, err := getModelsList(selectedTags, existingModels)
		if err != nil {
			spinner.Fail(err.Error())
			return model.Model{}, err
		}
		spinner.Success()
		selectedModel = selectModel(availableModels)
	}
	return selectedModel, nil
}

func processAdd(selectedModel model.Model, customArgs downloadermodel.Args, yes bool) (warning string, err error) {
	// User choose if he wishes to install the model directly
	message := fmt.Sprintf("Do you wish to directly download %s?", selectedModel.Name)
	selectedModel.AddToBinaryFile = yes || app.UI().AskForUsersConfirmation(message)

	// Validate model for download
	warningMessage, valid, err := config.Validate(selectedModel, yes)
	if !valid {
		return warningMessage, err
	}

	// Try to download model
	updatedModel, err := downloadModel(selectedModel, customArgs)
	if err != nil {
		return warning, err
	}

	// Add models to configuration file
	spinner := app.UI().StartSpinner("Adding model to configuration file...")
	err = config.AddModels(model.Models{updatedModel})
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

// downloadModel tries to download the selected model
func downloadModel(selectedModel model.Model, downloaderArgs downloadermodel.Args) (model.Model, error) {
	// Prepare the script arguments
	downloaderArgs.ModelName = selectedModel.Name
	if downloaderArgs.ModelClass == "" {
		downloaderArgs.ModelClass = selectedModel.Class
	}
	if downloaderArgs.ModelModule == "" {
		downloaderArgs.ModelModule = string(selectedModel.Module)
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
		return model.Model{}, fmt.Errorf("this model %s couldn't be downloaded", selectedModel.Name)
	}

	return selectedModel, nil
}

// getModelsList get list of models to display
func getModelsList(tags []string, existingModels model.Models) (model.Models, error) {
	allModelsWithTags, err := hfinteface.GetModelsByMultiplePipelineTags(tags)
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

// selectTags displays a multiselect to help the user choose the model types
func selectTags() []string {
	// Build a multiselect with each tag name
	message := "Please select the type of models you want to add"
	checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedTags := app.UI().DisplayInteractiveMultiselect(message, huggingface.AllTagsString(), checkMark, false, true)

	return selectedTags
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
