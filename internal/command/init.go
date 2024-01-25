package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a EMF project",
	Long:  `Initialize a EMF project.`,
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	app.L().Info("Initialize a EMF project")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
