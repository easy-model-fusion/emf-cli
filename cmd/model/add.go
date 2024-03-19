package cmdmodel

import (
	modelcontroller "github.com/easy-model-fusion/emf-cli/internal/controller/model"
	"github.com/spf13/cobra"
)

// addCmd represents the add model by names command
var modelAddCmd = &cobra.Command{
	Use:   "add [<model name>]",
	Short: "Add model by name to your project",
	Long:  `Add model by name to your project`,
	Run:   runAdd,
}

// runAddByNames runs the add command to add models by name
func runAdd(cmd *cobra.Command, args []string) {
	modelcontroller.RunAdd(args)
}
