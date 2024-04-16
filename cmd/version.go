package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of emf-cli",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	app.UI().Info().Println("Client version: " + app.Version + " (" + app.BuildDate + ")")
}
