package command

import (
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

	// Users chooses a command and runs it automatically
	runCommandSelector(cmd, args, commandsList, commandsMap)
}

func init() {
	// Add the add command to the root command
	rootCmd.AddCommand(addCmd)
}
