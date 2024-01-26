package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var shells = []string{"bash", "zsh", "fish", "powershell"}

var arguments = utils.ArrayStringAsArguments(shells)

// completionCmd represents the init command
var completionCmd = &cobra.Command{
	Use:   "ccompletion " + arguments,
	Short: "Generate shell autocompletion scripts",
	Long:  `Generate shell autocompletion scripts`,
	Run:   runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) {
	logger := app.L().WithTime(false)

	var selectedShell string

	// No args, asking for a shell input
	if len(args) == 0 {
		selectedShell = askForShell()
	} else {
		selectedShell = args[0]
	}

	// Checks whether the input shell is handled
	if utils.ArrayStringContainsItem(shells, selectedShell) {
		logger.Info(selectedShell)
	} else {
		logger.Error(fmt.Sprintf("Shell %s not recognized. Expected "+arguments, selectedShell))
		return
	}
}

// askForShell asks the user for a shell name and returns it
func askForShell() string {
	// Create an interactive text input with single line input mode
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

	// Show the text input and get the result
	result, _ := textInput.Show("Enter a shell name " + arguments)

	// Print a blank line for better readability
	pterm.Println()

	return result
}

func init() {
	// disables the defaults completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
}
