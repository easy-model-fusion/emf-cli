package cmd

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of emf-cli",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	app.L().WithTime(false).Info("Client version: " + app.Version + " (" + app.BuildDate + ")")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
