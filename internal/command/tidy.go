package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"path/filepath"
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
	models, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

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
			pterm.Error.Println(err.Error())
			return
		}

		// Add missing models to the list of models to be downloaded
		if !downloaded {
			modelsToDownload = append(modelsToDownload, currentModel)
		}
	}

	// download missing models
	err, _ = config.DownloadModels(modelsToDownload)
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

func init() {
	// Add the tidy command to the root command
	rootCmd.AddCommand(tidyCmd)
}
