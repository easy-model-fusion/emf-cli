package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   app.Name,
	Short: "emf-cli is a command line tool to manage a EMF project easily",
	Long:  `emf-cli is a command line tool to manage a EMF project easily.`,
	Run:   runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {

	logger := app.L().WithTime(false)

	// Variables for the commands data
	var commandsList []string
	var commandsMap = make(map[string]func(*cobra.Command, []string)) // key: command.Use; value: command.Run

	// get all the commands data
	for _, child := range cmd.Commands() {
		// Hiding the completion command inside the root command
		if completionUse != child.Use {
			commandsList = append(commandsList, child.Use)
			commandsMap[child.Use] = child.Run
		}
	}

	// allow the user to choose one command
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandsList).Show()

	// Check if the selected command exists and runs it
	if runCommand, exists := commandsMap[selectedCommand]; exists {
		runCommand(cmd, args)
	} else {
		logger.Error("Selected command '" + selectedCommand + "' not recognized")
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
