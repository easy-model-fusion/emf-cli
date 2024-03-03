package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
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
	controller.RunUpgrade(args)
}
