package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
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
		projectName = askForProjectName()
	} else {
		projectName = args[0]
	}

	// check if folder exists
	if _, err := os.Stat(projectName); !os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Folder %s already exists", projectName))
		return
	}

}

// askForProjectName asks the user for a project name and returns it
func askForProjectName() string {
	// Create an interactive text input with single line input mode
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

	// Show the text input and get the result
	result, _ := textInput.Show("Enter a project name")

	// Print a blank line for better readability
	pterm.Println()

	return result
}

func init() {
	rootCmd.AddCommand(initCmd)
}
