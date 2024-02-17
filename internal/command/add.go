package command

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const cmdAddTitle string = "add"

// addCmd represents the add model(s) command
var addCmd = &cobra.Command{
	Use:   cmdAddTitle + " (" + cmdAddCustomTitle + " | " + cmdAddNamesTitle + ")",
	Short: "Add model(s) to your project",
	Long:  `Add model(s) to your project`,
	Run:   runAdd,
}

// runAdd runs add command
func runAdd(cmd *cobra.Command, args []string) {

	commandsList := []string{addCustomCmd.Use, addByNamesCmd.Use}
	commandsMap := map[string]func(*cobra.Command, []string){
		addCustomCmd.Use:  addCustomCmd.Run,
		addByNamesCmd.Use: addByNamesCmd.Run,
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

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
