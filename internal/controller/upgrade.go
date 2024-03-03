package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/pterm/pterm"
)

// RunUpgrade upgrades the sdk version of a EMF project to the latest version.
func RunUpgrade(args []string) {
	pterm.Warning.Println("All the files in the folder sdk will be replaced with the latest version of the sdk.")
	pterm.Warning.Println("Be sure to not have any custom files in the sdk folder, as they will be deleted.")
	yes := app.UI().AskForUsersConfirmation("Are you sure you want to upgrade the sdk version of this project?")
	if !yes {
		return
	}

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return
	}

	_ = sdk.Upgrade()
}
