package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

var (
	buildDestination = "/dist"
	buildCustomName  string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long:  `Build the project.`,
	Run:   runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	if config.GetViperConfig() != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	// check if buildDestination exists
	if _, err := os.Stat(buildDestination); os.IsNotExist(err) {

		if buildDestination == "/dist" {
			err = os.Mkdir(buildDestination, os.ModeDir)
			if err != nil {
				pterm.Error.Println("Error creating dist folder")
				return
			}
		} else {
			pterm.Error.Println("Destination directory does not exist")
			return
		}

	}

}

func init() {
	buildCmd.Flags().StringVarP(&buildDestination, "out-dir", "o", "", "Destination directory where the project will be built")
	buildCmd.Flags().StringVarP(&buildCustomName, "name", "n", "", "Custom name for the executable")
	rootCmd.AddCommand(buildCmd)
}
