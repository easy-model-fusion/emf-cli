package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
	"os"
)

// tidyCmd represents the model tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Synchronizes the configuration file with the downloaded models",
	Long:  `Synchronizes the configuration file with the downloaded models`,
	Run:   runTidy,
}
var (
	tidyController               controller.TidyController
	authorizeAllSynchronisations bool
	accessToken                  string
)

func init() {
	tidyCmd.Flags().BoolVarP(&authorizeAllSynchronisations, "yes", "y", false, "Automatic yes to prompts")
	tidyCmd.Flags().StringVarP(&accessToken, "access-token", "a", "", "Access token for gated models")
}

// runTidy runs the model tidy command
func runTidy(cmd *cobra.Command, args []string) {
	err := tidyController.RunTidy(authorizeAllSynchronisations, accessToken)
	if err != nil {
		os.Exit(1)
	}
}
