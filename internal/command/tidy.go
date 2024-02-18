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
	"path/filepath"
	"strings"
)

// addCmd represents the add model(s) command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "add missing and remove unused models",
	Long:  `add missing and remove unused models`,
	Run:   runTidy,
}

// runAdd runs add command
func runTidy(cmd *cobra.Command, args []string) {
	// get all models from config file
	err := config.GetViperConfig()
	if err != nil {
		pterm.Error.Println(err.Error())
	}
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
func getModelsToBeAddedToBinary(models []model.Model) []model.Model {
	var returnedModels []model.Model

	for _, currentModel := range models {
		if currentModel.AddToBinary {
			returnedModels = append(returnedModels, currentModel)
		}
	}

	return returnedModels
}

// addMissingModels adds the missing models from the list of configuration file models
func addMissingModels(models []model.Model) error {
	pterm.Info.Println("verifying if all models are downloaded...")
	// filter the models that should be added to binary
	models = getModelsToBeAddedToBinary(models)
	// Search for the models that need to be downloaded
	var modelsToDownload []model.Model
	for _, currentModel := range models {
		// Check if download path is stored
		if currentModel.DirectoryPath == "" {
			currentModel.DirectoryPath = filepath.Join(app.ModelsDownloadPath, currentModel.Name)
		}

		// Check if model is already downloaded
		downloaded, err := utils.IsExistingPath(currentModel.DirectoryPath)
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
		//err, _ := config.DownloadModels(modelsToDownload)
		//if err != nil {
		//	return err
		//}
		pterm.Success.Println("added missing models", model.GetNames(modelsToDownload))
	} else {
		pterm.Info.Println("all models are already downloaded")
	}

	return nil
}

// missingModelConfiguration finds the downloaded models that aren't configured in the configuration file
// and then asks the user if he wants to delete them or add them to the configuration file
func missingModelConfiguration(models []model.Model) error {
	pterm.Info.Println("verifying if all downloaded models are configured...")
	// Get the list of downloaded model names
	downloadedModelNames, err := app.GetDownloadedModelNames()
	if err != nil {
		return err
	}

	// Get the list of configured model names
	configModelNames := model.GetNames(models)
	// Find missing models from configuration file
	missingModelNames := utils.StringDifference(downloadedModelNames, configModelNames)
	if len(missingModelNames) > 0 {
		// Ask user for confirmation to delete these models
		message := fmt.Sprintf(
			"These models %s weren't found in your configuration file. Do you wish to delete these models?",
			strings.Join(missingModelNames, ", "))
		yes := utils.AskForUsersConfirmation(message)
		// Delete models if confirmed
		if yes {
			for _, modelName := range missingModelNames {
				_ = config.RemoveModelPhysically(modelName)
			}
			pterm.Success.Println("deleted models %s", missingModelNames)
		} else { // Add models' configurations to config file
			err = generateModelsConfig(missingModelNames)
			if err != nil {
				return err
			}
			pterm.Success.Println("added the configurations for these models", missingModelNames)
		}
	} else {
		pterm.Info.Println("all downloaded models are well configured")
	}

	return nil
}

// regenerateCode generates new default python code
func regenerateCode() error {
	// TODO: modify this logic when code generator is completed
	pterm.Info.Println("generating new default python code...")
	//file := codegen.File{Name: "main.py"}

	//generator := codegen.PythonCodeGenerator{}
	//_, err := generator.Generate(&file)
	//if err != nil {
	//	return err
	//}
	pterm.Success.Println("python code generated")
	return nil
}

func init() {
	// Add the tidy command to the root command
	rootCmd.AddCommand(tidyCmd)
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
		}
		currentModel.AddToBinary = true
		currentModel.DirectoryPath = app.ModelsDownloadPath
		models = append(models, currentModel)
	}

	// Add models to the configuration file
	err := config.AddModels(models)
	if err != nil {
		return err
	}

	return nil
}
