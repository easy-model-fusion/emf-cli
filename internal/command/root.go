package command

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	commandmodel "github.com/easy-model-fusion/emf-cli/internal/command/model"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
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

const rootCommandName string = app.Name

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   rootCommandName,
	Short: "emf-cli is a command line tool to manage a EMF project easily",
	Long:  `emf-cli is a command line tool to manage a EMF project easily.`,
	Run:   runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {

	// Build objects containing all the available commands
	commandsList, commandsMap := utils.CobraGetSubCommands(cmd, []string{completionCmd.Use})

	// Users chooses a command and runs it automatically
	utils.CobraSelectCommandToRun(cmd, args, commandsList, commandsMap)
}

func init() {
	// Add persistent flag for configuration file path
	rootCmd.PersistentFlags().StringVar(&config.FilePath, "path", ".", "config file path")

	// Adding subcommands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(commandmodel.ModelCmd)
}
