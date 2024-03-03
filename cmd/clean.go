package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

var allFlagDelete bool
var authorizeAllDelete bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean project",
	Long:  "Clean project",
	Run:   runClean,
}

func init() {
	cleanCmd.Flags().BoolVarP(&allFlagDelete, "all", "a", false, "clean all project")
	cleanCmd.Flags().BoolVarP(&authorizeAllDelete, "yes", "y", false, "authorize all deletions")
}

func runClean(cmd *cobra.Command, args []string) {
	controller.RunClean(allFlagDelete, authorizeAllDelete)
}
