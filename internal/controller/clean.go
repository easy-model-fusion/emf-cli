package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/pterm/pterm"
	"os"
)

const cleanDirName = "build"

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

		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		} else {
			pterm.Error.Printfln("Operation failed.")
		}

	}

	_, err := os.Stat(cleanDirName)
	if os.IsNotExist(err) {
		pterm.Success.Printfln("Operation succeeded.")
		return
	}
	err = os.RemoveAll(cleanDirName)
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}
