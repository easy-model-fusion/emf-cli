package cmd

import (
	"github.com/easy-model-fusion/emf-cli/cmd/model"
	"github.com/easy-model-fusion/emf-cli/cmd/tokenizer"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
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

func init() {
	app.InitGit(app.Repository, "")
	// Add persistent flag for configuration file path
	rootCmd.PersistentFlags().StringVar(&config.FilePath, "config-path", ".", "config file path")
	rootCmd.PersistentFlags().StringVar(app.G().GetAuthToken(), "git-auth-token", "", "Git auth token")

	// Adding subcommands
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(tidyCmd)
	rootCmd.AddCommand(cmdmodel.ModelCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(cmdtokenizer.TokenizerCmd)
}

func runRoot(cmd *cobra.Command, args []string) {
	// Running command as palette : allowing user to choose subcommand
	err := cobrautil.RunCommandAsPalette(cmd, args, rootCommandName, []string{completionCmd.Name()})
	if err != nil {
		app.UI().Error().Println("Something went wrong :", err)
	}
}
