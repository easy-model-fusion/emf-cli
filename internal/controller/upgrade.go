package controller

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
)

// RunUpgrade upgrades the sdk version of a EMF project to the latest version.
func RunUpgrade(yes bool) {
	app.UI().Warning().Println("All the files in the folder sdk will be replaced with the latest version of the sdk.")
	app.UI().Warning().Println("Be sure to not have any custom files in the sdk folder, as they will be deleted.")

	if !yes {
		yes = app.UI().AskForUsersConfirmation("Are you sure you want to upgrade the sdk version of this project?")
		if !yes {
			return
		}
	}

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return
	}

	err = sdk.Upgrade()
	if err != nil {
		app.UI().Error().Println("Error upgrading sdk:", err)
		return
	}

	models, err := config.GetModels()
	if err != nil {
		app.UI().Error().Println("Error regenerating code:", err)
		return
	}

	err = regenerateCode(models)
	if err != nil {
		app.UI().Error().Println("Error regenerating code:", err)
		return
	}

}
