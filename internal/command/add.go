package command

import (
	"github.com/easy-model-fusion/client/internal/utils"
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

	// Build objects containing all the available commands
	addSubCmd, found := utils.CobraFindSubCommand(cmd, cmdAddTitle)
	if !found {
		// technically unreachable
		return
	}
	commandsList, commandsMap := utils.CobraGetSubCommands(addSubCmd, []string{})

	// Users chooses a command and runs it automatically
	utils.CobraSelectCommandToRun(cmd, args, commandsList, commandsMap)
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
