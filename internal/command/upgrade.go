package command

import (
	"github.com/spf13/cobra"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade <project name>",
	Short: "Upgrade the sdk version of a EMF project to the latest version.",
	Long:  `Upgrade the sdk version of a EMF project to the latest version.`,
	Run:   runUpgrade,
}

func runUpgrade(cmd *cobra.Command, args []string) {

}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
