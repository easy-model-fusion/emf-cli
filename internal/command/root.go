package command

import (
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "emf-client",
	Short: "EMF client is a command line tool to manage a EMF project easily",
	Long:  `EMF client is a command line tool to manage a EMF project easily.`,
	Run:   runBase,
}

func runBase(cmd *cobra.Command, args []string) {

	// get all commands
	var commandList []string
	for _, child := range cmd.Commands() {
		commandList = append(commandList, child.Use)
	}

	// allow the user to choose one command
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandList).Show()

	// run the selected command
	selectedChild, _, _ := cmd.Find([]string{selectedCommand})
	if selectedChild != nil {
		selectedChild.Run(cmd, args)
	} else {
		app.L().WithTime(false).Error("Selected command " + selectedCommand + " is not recognized")
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
