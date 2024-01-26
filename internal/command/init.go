package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/spf13/cobra"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project name>",
	Short: "Initialize a EMF project",
	Long:  `Initialize a EMF project.`,
	Args:  utils.ValidFileName(1, true),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	logger := app.L().WithTime(false)

	var projectName string

	// No args, check projectName in pterm
	if len(args) == 0 {
		projectName = utils.AskForUsersInput("Enter a project name")
	} else {
		projectName = args[0]
	}

	// check if folder exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Folder %s already exists", projectName))
		return
	}

}

func init() {
	rootCmd.AddCommand(initCmd)
}
