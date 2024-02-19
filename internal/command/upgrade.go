package command

import (
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the sdk version of a EMF project to the latest version.",
	Long:  `Upgrade the sdk version of a EMF project to the latest version.`,
	Run:   runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) {
	pterm.Warning.Println("All the files in the folder sdk will be replaced with the latest version of the sdk.")
	pterm.Warning.Println("Be sure to not have any custom files in the sdk folder, as they will be deleted.")
	yes := utils.AskForUsersConfirmation("Are you sure you want to upgrade the sdk version of this project?")
	if !yes {
		return
	}

	err := config.GetViperConfig(config.FilePath)
	if err != nil {
		return
	}

	_ = sdk.Upgrade()
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
