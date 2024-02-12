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
	"path/filepath"
)

// addCmd represents the add model(s) command
var addCmd = &cobra.Command{
	Use:   "add <model name> [<other model names>...]",
	Short: "Add model(s) to your project",
	Long:  `Add model(s) to your project`,
	Args:  config.ValidModelName(), // TODO: Do this validation in the run function, bc proxy could not be initialized
	Run:   runAdd,
}

// displayModels indicates if the multiselect of models should be displayed or not
var displayModels bool

// runAdd runs add command
func runAdd(cmd *cobra.Command, args []string) {
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
		for _, name := range args {
			apiModel, err := app.H().GetModel(name)
			if err != nil {
				pterm.Error.Println("while getting model " + name)
				return
			}
			selectedModels = append(selectedModels, apiModel)
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
		var modelNames []string
		selectedModels, modelNames = selectModels(selectedTags, selectedModelNames)
		if selectedModels == nil {
			app.L().WithTime(false).Warn("No models selected")
			return
		}
		selectedModelNames = append(selectedModelNames, modelNames...)
	}

	// User choose either to exclude or include models in binary
	selectedModels = selectExcludedModelsFromInstall(selectedModels, selectedModelNames)

	// Display the selected models
	utils.DisplaySelectedItems(selectedModelNames)

	// Download the models
	err, selectedModels := downloadModels(selectedModels)
	if err != nil {
		return
	}

	// Add models to configuration file
	spinner, _ := pterm.DefaultSpinner.Start("Writing models to configuration file...")
	err = config.AddModel(selectedModels)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the models to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
}

func downloadModels(models []model.Model) (error, []model.Model) {

	// Find the python executable inside the venv to run the scripts
	pythonPath, err := utils.FindVEnvExecutable(".venv", "python")
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Error using the venv : %s", err))
		return err, nil
	}

	// Iterate over every model for instant download
	for i := range models {

		// Get mandatory model data for the download script
		modelName := models[i].Name
		moduleName := models[i].Config.ModuleName
		className := models[i].Config.ClassName
		overwrite := false

		// TODO : get Config.ModuleName & Config.ClassName
		moduleName = "diffusers"
		// className = "StableDiffusionXLPipeline"
		// moduleName = "transformers"
		// className = "AutoModelForCausalLM"

		// Local path where the model will be downloaded
		downloadPath := app.ModelsDownloadPath
		modelPath := filepath.Join(downloadPath, modelName)

		// Check if the model_path already exists
		if exists, err := utils.IsExistingPath(modelPath); err != nil {
			// Skipping model : an error occurred
			continue
		} else if exists {
			// Model path already exists : ask the user if he would like to overwrite it
			overwrite, _ = pterm.DefaultInteractiveConfirm.Show(fmt.Sprintf("Model already exists at '%s'. Do you want to overwrite it?", modelPath))

			// User does not want to overwrite : skipping to the next model
			if !overwrite {
				continue
			}
		}

		// Run the script to download the model
		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Downloading model '%s'...", modelName))
		err, exitCode := utils.DownloadModel(pythonPath, downloadPath, modelName, moduleName, className, overwrite)
		if err != nil {
			spinner.Fail(err)
			switch exitCode {
			case 2:
				// TODO : Update the log message once the command is implemented
				pterm.Info.Println("Run the 'add --single' command to manually add the model.")
			}
			continue
		}
		spinner.Success(fmt.Sprintf("Successfully downloaded model '%s'", modelName))

		// Update the directory path to the model
		models[i].DirectoryPath = downloadPath
	}

	return nil, models
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
	currentModelNames, err := config.GetModelNames()
	if err != nil {
		app.L().Fatal("error while getting current models")
	}

	// Remove existent models from list of models to add
	// Remove already selected models from list of models to add (in case user entered add model_name -s)
	existingModels := config.GetModelsByNames(models, append(currentModelNames, currentSelectedModels...))
	models = config.Difference(models, existingModels)

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
func selectExcludedModelsFromInstall(models []model.Model, modelNames []string) []model.Model {
	// Build a multiselect with each selected model name to exclude/include in the binary
	message := "Please select the model(s) that you don't wish to install directly"
	checkMark := &pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	installsToExclude := utils.DisplayInteractiveMultiselect(message, modelNames, checkMark, false)
	var updatedModels []model.Model
	for _, currentModel := range models {
		currentModel.AddToBinary = !utils.ArrayStringContainsItem(installsToExclude, currentModel.Name)
		updatedModels = append(updatedModels, currentModel)
	}

	return updatedModels
}

func init() {
	app.InitHuggingFace(huggingface.BaseUrl, "")
	// Add --select flag to the add command
	addCmd.Flags().BoolVarP(&displayModels, "select", "s", false, "Select models to add")
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
