package add

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const modelAddCommandName string = "add"

// ModelAddCmd represents the add model command
var ModelAddCmd = &cobra.Command{
	Use:   modelAddCommandName,
	Short: "Palette that contains add model based commands",
	Long:  "Palette that contains add model based commands",
	Run:   runModelAdd,
}

// runModelAdd runs model add command
func runModelAdd(cmd *cobra.Command, args []string) {

	// Searching for the currentCmd : when 'cmd' differs from 'addCmd' (i.e. run through parent multiselect)
	currentCmd, found := utils.CobraFindSubCommand(cmd, modelAddCommandName)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", modelAddCommandName))
		return
	}

	// Retrieve all the subcommands of the current command
	commandsList, commandsMap := utils.CobraGetSubCommands(currentCmd, []string{})

	// Users chooses a command and runs it automatically
	utils.CobraSelectCommandToRun(currentCmd, args, commandsList, commandsMap)
}

func init() {
	// Adding subcommands
	ModelAddCmd.AddCommand(addCustomCmd)
	ModelAddCmd.AddCommand(addByNamesCmd)
}
