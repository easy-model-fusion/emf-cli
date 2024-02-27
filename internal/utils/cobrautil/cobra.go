package cobrautil

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strconv"
)

// FindSubCommand searches for a sub-command within a Cobra command.
func FindSubCommand(cmd *cobra.Command, cmdSearchName string) (*cobra.Command, bool) {

	// Running from the searched command
	if cmd.Name() == cmdSearchName {
		return cmd, true
	}

	for _, subCmd := range cmd.Commands() {
		if subCmd.Name() == cmdSearchName {
			return subCmd, true
		}
	}

	return nil, false
}

// GetSubCommands retrieves sub-commands and hides the ones specified.
func GetSubCommands(cmd *cobra.Command, cmdsToHide []string) ([]string, map[string]func(*cobra.Command, []string)) {

	// Variables for the commands data
	var commandsList []string
	var commandsMap = make(map[string]func(*cobra.Command, []string)) // key: command.Use; value: command.Run

	// Iterating over all sub-commands
	for _, child := range cmd.Commands() {
		if !stringutil.SliceContainsItem(cmdsToHide, child.Use) {
			commandsList = append(commandsList, child.Use)
			commandsMap[child.Use] = child.Run
		}
	}
	return commandsList, commandsMap
}

// MultiselectSubcommands presents an interactive selection of available sub-commands and executes the chosen one.
func MultiselectSubcommands(cmd *cobra.Command, args []string, commandsList []string, commandsMap map[string]func(*cobra.Command, []string)) {
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandsList).Show()

	if runCommand, exists := commandsMap[selectedCommand]; exists {
		runCommand(cmd, args)
	} else {
		pterm.Error.Println(fmt.Sprintf("Selected command '%s' not recognized", selectedCommand))
	}
}

// GetNonProvidedLocalFlags retrieves flags that have not been provided for a Cobra command.
func GetNonProvidedLocalFlags(cmd *cobra.Command) []*pflag.Flag {
	var nonProvided []*pflag.Flag

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "help" && !flag.Changed {
			nonProvided = append(nonProvided, flag)
		}
	})

	return nonProvided
}

// MultiselectRemainingFlags presents an interactive multiselect for remaining flags and returns selected ones.
func MultiselectRemainingFlags(cmd *cobra.Command) (map[string]*pflag.Flag, []string) {

	// Get all flags that were not already provided
	remainingFlags := GetNonProvidedLocalFlags(cmd)

	// Process the flag's usage property into slice and map
	var remainingFlagsUsages []string
	remainingFlagsMap := make(map[string]*pflag.Flag)
	for _, flag := range remainingFlags {
		if flag.Name != "help" {
			remainingFlagsMap[flag.Usage] = flag
			remainingFlagsUsages = append(remainingFlagsUsages, flag.Usage)
		}
	}

	// User multi-selects the flags he wishes to use
	message := "Select any property you wish to set"
	checkMark := &pterm.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
	selectedFlags := ptermutil.DisplayInteractiveMultiselect(message, remainingFlagsUsages, []string{}, checkMark, false)
	ptermutil.DisplaySelectedItems(selectedFlags)

	return remainingFlagsMap, selectedFlags
}

// AskFlagInput prompts the user for input for a specific flag of a Cobra command.
func AskFlagInput(cmd *cobra.Command, flag *pflag.Flag) error {

	// Prepare value
	var inputValue string

	// Iterate until the user provides a valid value
	for inputValue == "" {

		// Ask for different types of input
		switch flag.Value.Type() {
		case "bool":
			inputValue = strconv.FormatBool(ptermutil.AskForUsersConfirmation(flag.Usage))
		default:
			inputValue = ptermutil.AskForUsersInput(flag.Usage)
		}
	}

	// Set the flag's value
	return cmd.Flags().Set(flag.Name, inputValue)
}

// AllowInputAmongRemainingFlags presents remaining flags for input selection.
func AllowInputAmongRemainingFlags(cmd *cobra.Command) error {

	// User chooses among the remaining flags
	remainingFlagsMap, selectedFlags := MultiselectRemainingFlags(cmd)

	// User inputs data for the chosen flags
	for _, flag := range selectedFlags {
		err := AskFlagInput(cmd, remainingFlagsMap[flag])
		if err != nil {
			return fmt.Errorf("couldn't set the value for %s : %s", remainingFlagsMap[flag].Name, err)
		}
	}

	return nil
}

// RunCommandAsPalette allows the user to run a subcommand
func RunCommandAsPalette(cmd *cobra.Command, args []string, cmdSearchName string, cmdsToHide []string) error {

	// Searching for the commandName : when 'cmd' differs from 'cmdSearchName'
	currentCmd, found := FindSubCommand(cmd, cmdSearchName)
	if !found {
		return fmt.Errorf("the '%s' command was not found", cmdSearchName)
	}

	// Retrieve all the subcommands of the current command
	commandsList, commandsMap := GetSubCommands(currentCmd, cmdsToHide)

	// Users chooses a command and runs it automatically
	MultiselectSubcommands(currentCmd, args, commandsList, commandsMap)

	return nil
}
