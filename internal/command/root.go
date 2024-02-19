package command

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
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
	} else { // technically unreachable
		pterm.Error.Println(fmt.Sprintf("Selected command '%s' not recognized", selectedCommand))
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

func init() {
	// Add persistent flag for configuration file path
	rootCmd.PersistentFlags().StringVar(&config.FilePath, "path", ".", "config file path")
}
