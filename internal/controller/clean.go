package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"os"
)

const cleanDirName = "dist"

func RunClean(allFlagDelete bool, authorizeAllDelete bool) {

	// Send update suggestion
	sdk.SendUpdateSuggestion()

	// Only clean if config file exists (so we know it's a EMF project)
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	// Delete all models if flag --all
	if allFlagDelete {
		// Ask for confirmation
		if !authorizeAllDelete {
			yes := app.UI().AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
			if !yes {
				return
			}
		}

		spinner := app.UI().StartSpinner("Cleaning all models...")
		info, err := config.RemoveAllModels()
		if err == nil {
			spinner.Success()
			if info != "" {
				app.UI().Info().Printfln(info)
			}
		} else {
			spinner.Fail(fmt.Sprintf("Error cleaning all models: %s", err))
		}

	}

	_, err := os.Stat(cleanDirName)
	if os.IsNotExist(err) {
		app.UI().Success().Printfln("Operation succeeded.")
		return
	}

	spinner := app.UI().StartSpinner("Cleaning project files...")
	err = os.RemoveAll(cleanDirName)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error cleaning project files: %s", err))
		return
	}
	spinner.Success()
}
