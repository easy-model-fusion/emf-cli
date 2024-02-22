package cmd

import (
	"github.com/easy-model-fusion/emf-cli/cmd/model"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
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
	// Running command as palette : allowing user to choose subcommand
	err := utils.CobraRunCommandAsPalette(cmd, args, rootCommandName, []string{completionCmd.Name()})
	if err != nil {
		pterm.Error.Println("Something went wrong :", err)
	}
}

func init() {
	app.InitGit(app.Repository, "")
	// Add persistent flag for configuration file path
	rootCmd.PersistentFlags().StringVar(&config.FilePath, "path", ".", "config file path")
	rootCmd.PersistentFlags().StringVar(&app.G().AuthToken, "git-auth-token", "", "Git auth token")

	// Adding subcommands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(modelTidyCmd)
	rootCmd.AddCommand(cmdmodel.ModelCmd)
}
