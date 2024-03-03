package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
)

const cleanDirName = "build"

func RunClean(allFlagDelete bool, authorizeAllDelete bool) {
	// Delete all models if flag --all
	if allFlagDelete {
		// Ask for confirmation
		if !authorizeAllDelete {
			yes := app.UI().AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
			if !yes {
				return
			}
		}
		if config.GetViperConfig(config.FilePath) != nil {
			return
		}
		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		} else {
			pterm.Error.Printfln("Operation failed.")
		}

	}

	// Get the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}
	buildDir := filepath.Join(currentDir, cleanDirName)

	_, err = os.Stat(buildDir)
	if os.IsNotExist(err) {
		pterm.Success.Printfln("Operation succeeded.")
		return
	}
	err = os.RemoveAll(buildDir)
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}
