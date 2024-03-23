package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

// tidyCmd represents the model tidy command
var tidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "synchronizes the configuration file with the downloaded models",
	Long:  `synchronizes the configuration file with the downloaded models`,
	Run:   runTidy,
}

var yes bool

func init() {
	tidyCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic yes to prompts")
}

// runTidy runs the model tidy command
func runTidy(cmd *cobra.Command, args []string) {
	controller.RunTidy(yes)
}
