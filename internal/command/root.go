package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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

	// Build objects containing all the available commands
	commandsList, commandsMap = getAllCommands(cmd, commandsList, commandsMap)
	commandsList, commandsMap = hideCommands(commandsList, commandsMap, []string{completionCmd.Use, addCmd.Use})

	// Users chooses a command and runs it automatically
	runCommandSelector(cmd, args, commandsList, commandsMap)
}

// getAllCommands retrieves all commands and their subcommands recursively.
func getAllCommands(cmd *cobra.Command, commandsList []string, commandsMap map[string]func(*cobra.Command, []string)) ([]string, map[string]func(*cobra.Command, []string)) {
	for _, child := range cmd.Commands() {
		commandsList = append(commandsList, child.Use)
		commandsMap[child.Use] = child.Run
		commandsList, commandsMap = getAllCommands(child, commandsList, commandsMap)
	}
	return commandsList, commandsMap
}

// hideCommands hides specified commands from the given list of commands.
func hideCommands(commandsList []string, commandsMap map[string]func(*cobra.Command, []string), commands []string) ([]string, map[string]func(*cobra.Command, []string)) {
	for _, command := range commands {
		commandsList = utils.SliceRemoveValue(commandsList, command)
		delete(commandsMap, command)
	}
	return commandsList, commandsMap
}

func runCommandSelector(cmd *cobra.Command, args []string, commandsList []string, commandsMap map[string]func(*cobra.Command, []string)) {
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandsList).Show()

	if runCommand, exists := commandsMap[selectedCommand]; exists {
		runCommand(cmd, args)
	} else {
		pterm.Error.Println(fmt.Sprintf("Selected command '%s' not recognized", selectedCommand))
	}
}
