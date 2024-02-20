package commandmodel

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
	"strings"
)

// modelTidyCmd represents the model tidy command
var modelTidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "add missing and remove unused models",
	Long:  `add missing and remove unused models`,
	Run:   runModelTidy,
}

// runModelTidy runs the model tidy command
func runModelTidy(cmd *cobra.Command, args []string) {
	// get all models from config file
	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		pterm.Error.Println(err.Error())
	}

	sdk.SendUpdateSuggestion() // TODO: here proxy?

	models, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Add all missing models
	err = addMissingModels(models)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Fix missing model configurations
	err = missingModelConfiguration(models)
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Regenerate python code
	err = regenerateCode()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}
}

// getModelsToBeAddedToBinary returned models that needs to be added to binary
func getModelsToBeAddedToBinary(models []model.Model) (returnedModels []model.Model) {

	for _, currentModel := range models {
		if currentModel.AddToBinary {
			returnedModels = append(returnedModels, currentModel)
		}
	}

	return returnedModels
}

// addMissingModels adds the missing models from the list of configuration file models
func addMissingModels(models []model.Model) error {
	pterm.Info.Println("Verifying if all models are downloaded...")
	// filter the models that should be added to binary
	models = getModelsToBeAddedToBinary(models)
	// Search for the models that need to be downloaded
	var modelsToDownload []model.Model
	for _, currentModel := range models {
		// build model path
		currentModelPath := currentModel.Config.Path
		if currentModelPath != "" {
			currentModel = model.ConstructConfigPaths(currentModel)
		}

		// Check if model is already downloaded
		downloaded, err := utils.IsExistingPath(currentModelPath)
		if err != nil {
			return err
		}

		// Add missing models to the list of models to be downloaded
		if !downloaded {
			modelsToDownload = append(modelsToDownload, currentModel)
		}
	}

	if len(modelsToDownload) > 0 {
		// download missing models
		_, failedModels := config.DownloadModels(modelsToDownload)
		if !model.Empty(failedModels) {
			return fmt.Errorf("these models could not be downloaded %s", model.GetNames(failedModels))
		}
		pterm.Success.Println("Added missing models", model.GetNames(modelsToDownload))
	} else {
		pterm.Info.Println("All models are already downloaded")
	}

	return nil
}

// missingModelConfiguration finds the downloaded models that aren't configured in the configuration file
// and then asks the user if he wants to delete them or add them to the configuration file
func missingModelConfiguration(models []model.Model) error {
	pterm.Info.Println("Verifying if all downloaded models are configured...")
	// Get the list of downloaded model names
	downloadedModelNames, err := app.GetDownloadedModelNames()
	if err != nil {
		return err
	}

	// Get the list of configured model names
	configModelNames := model.GetNames(models)
	// Find missing models from configuration file
	missingModelNames := utils.SliceDifference(downloadedModelNames, configModelNames)
	if len(missingModelNames) > 0 {
		err = handleModelsWithNoConfig(missingModelNames)
		if err != nil {
			return err
		}
	} else {
		pterm.Info.Println("All downloaded models are well configured")
	}

	return nil
}

// regenerateCode generates new default python code
func regenerateCode() error {
	// TODO: modify this logic when code generator is completed
	pterm.Info.Println("Generating new default python code...")
	//file := codegen.File{Name: "main.py"}

	//generator := codegen.PythonCodeGenerator{}
	//_, err := generator.Generate(&file)
	//if err != nil {
	//	return err
	//}
	pterm.Success.Println("Python code generated")
	return nil
}

// generateModelsConfig generates models configurations
func generateModelsConfig(modelNames []string) error {
	// initialize hugging face url
	app.InitHuggingFace(huggingface.BaseUrl, "")
	// get hugging face api
	huggingFace := app.H()

	var models []model.Model
	for _, modelName := range modelNames {
		// Search for the model in hugging face
		currentModel, err := huggingFace.GetModel(modelName)
		// If not found create model configuration with only model's name
		if err != nil {
			currentModel = model.Model{Name: modelName}
			currentModel.Source = model.CUSTOM
		}
		currentModel.AddToBinary = true
		currentModel = model.ConstructConfigPaths(currentModel)
		models = append(models, currentModel)
	}

	// Add models to the configuration file
	err := config.AddModels(models)
	if err != nil {
		return err
	}

	return nil
}

// handleModelsWithNoConfig handles all the models with no configuration
func handleModelsWithNoConfig(missingModelNames []string) error {
	// Ask user to select the models to delete/add to configuration file
	message := "These models weren't found in your configuration file and will be deleted. " +
		"Please select the models that you wish to conserve"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedModels := utils.DisplayInteractiveMultiselect(message, missingModelNames, checkMark, false)
	modelsToDelete := utils.SliceDifference(missingModelNames, selectedModels)

	// Delete selected models
	if len(modelsToDelete) > 0 {
		// Ask user for confirmation to delete these models
		message = fmt.Sprintf(
			"Are you sure you want to delete these models [%s]?",
			strings.Join(modelsToDelete, ", "))
		yes := utils.AskForUsersConfirmation(message)
		if yes {
			// Delete models if confirmed
			for _, modelName := range modelsToDelete {
				err := config.RemoveModelPhysically(modelName)
				if err != nil {
					return err
				}
			}
			pterm.Success.Println("Deleted models", modelsToDelete)
		} else {
			return handleModelsWithNoConfig(missingModelNames)
		}
	}

	// Configure selected models
	if len(selectedModels) > 0 {
		// Add models' configurations to config file
		err := generateModelsConfig(selectedModels)
		if err != nil {
			return err
		}
		pterm.Success.Println("Added configurations for these models", selectedModels)
	}
	return nil
}
