package command

import (
	"github.com/easy-model-fusion/client/internal/app"
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

	logger := app.L().WithTime(false)

	// get all commands
	var commandList []string
	for _, child := range cmd.Commands() {
		logger.Info(child.Use)
		if completionUse != child.Use {
			commandList = append(commandList, child.Use)
		}
	}

	// allow the user to choose one command
	selectedCommand, _ := pterm.DefaultInteractiveSelect.WithOptions(commandList).Show()

	// get the chosen command
	selectedChild, _, _ := cmd.Find([]string{selectedCommand})

	logger.Info(selectedChild.Use)
	if app.Name == selectedChild.Use { // avoid loops when the chosen command is the help command
		cmd.HelpFunc()(cmd, args)
		return
	} else if selectedChild != nil { // run the selected command
		selectedChild.Run(cmd, args)
	} else { // unexpected
		logger.Error("Selected command '" + selectedCommand + "' not recognized")
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
