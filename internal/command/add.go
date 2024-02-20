package command

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const cmdAddTitle string = "add"

// addCmd represents the add model(s) command
var addCmd = &cobra.Command{
	Use:   cmdAddTitle,
	Short: "Add model(s) to your project",
	Long:  `Add model(s) to your project`,
	Run:   runAdd,
}

// runAdd runs add command
func runAdd(cmd *cobra.Command, args []string) {

	// Searching for the currentCmd : when 'cmd' differs from 'addCmd' (i.e. run through parent multiselect)
	currentCmd, found := utils.CobraFindSubCommand(cmd, cmdAddTitle)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", cmdAddTitle))
		return
	}

	// Retrieve all the subcommands of the current command
	commandsList, commandsMap := utils.CobraGetSubCommands(currentCmd, []string{})

	// Users chooses a command and runs it automatically
	utils.CobraSelectCommandToRun(currentCmd, args, commandsList, commandsMap)
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
