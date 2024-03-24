package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

var cleanDeleteAll bool
var cleanDeleteAllYes bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean project files (e.g. models, build)",
	Long:  "Clean project files (e.g. models, build)",
	Run:   runClean,
}

func init() {
	cleanCmd.Flags().BoolVarP(&cleanDeleteAll, "all", "a", false, "clean all project")
	cleanCmd.Flags().BoolVarP(&cleanDeleteAllYes, "yes", "y", false, "bypass delete all confirmation")
}

func runClean(cmd *cobra.Command, args []string) {
	controller.RunClean(cleanDeleteAll, cleanDeleteAllYes)
}
