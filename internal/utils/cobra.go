package utils

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func CobraFindSubCommand(cmd *cobra.Command, cmdToSearch string) (*cobra.Command, bool) {
	for _, child := range cmd.Commands() {
		title := Split(child.Use)[0]
		if title == cmdToSearch {
			return child, true
		}
	}
	return nil, false
}

// CobraGetSubCommands retrieves sub-commands and hides the ones specified.
func CobraGetSubCommands(cmd *cobra.Command, cmdsToHide []string) ([]string, map[string]func(*cobra.Command, []string)) {

	// Variables for the commands data
	var commandsList []string
	var commandsMap = make(map[string]func(*cobra.Command, []string)) // key: command.Use; value: command.Run

	// Iterating over all sub-commands
	for _, child := range cmd.Commands() {
		if !SliceContainsItem(cmdsToHide, child.Use) {
			commandsList = append(commandsList, child.Use)
			commandsMap[child.Use] = child.Run
		}
	}
	return commandsList, commandsMap
}

func CobraSelectCommandToRun(cmd *cobra.Command, args []string, commandsList []string, commandsMap map[string]func(*cobra.Command, []string)) {
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandsList).Show()

	if runCommand, exists := commandsMap[selectedCommand]; exists {
		runCommand(cmd, args)
	} else {
		pterm.Error.Println(fmt.Sprintf("Selected command '%s' not recognized", selectedCommand))
	}
}
