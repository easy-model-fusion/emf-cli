package controller

import (
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

		info, err := config.RemoveAllModels()
		if info != "" {
			app.UI().Info().Printfln(info)
		}

		if err == nil {
			app.UI().Success().Printfln("Operation succeeded.")
		} else {
			app.UI().Error().Printfln("Operation failed.")
		}

	}

	_, err := os.Stat(cleanDirName)
	if os.IsNotExist(err) {
		app.UI().Success().Printfln("Operation succeeded.")
		return
	}
	err = os.RemoveAll(cleanDirName)
	if err == nil {
		app.UI().Success().Printfln("Operation succeeded.")
	} else {
		app.UI().Error().Printfln("Operation failed.")
	}
}
