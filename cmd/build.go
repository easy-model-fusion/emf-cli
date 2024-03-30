package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
	"os"
)

var (
	buildDestination   string
	buildCustomName    string
	buildOneFile       bool
	buildModelsSymlink bool
	buildLibrary       string
	buildController    = controller.BuildController{}
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the project using the selected library (pyinstaller or nuitka).
			Note: if you want to use nuitka, you need to have a working C compiler.`,
	Run: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	err := buildController.Run(buildCustomName, buildLibrary, buildDestination, buildOneFile, buildModelsSymlink)
	if err != nil {
		app.UI().Error().Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	buildCmd.Flags().StringVarP(&buildDestination, "out-dir", "o", "dist", "Destination directory where the project will be built")
	buildCmd.Flags().StringVarP(&buildCustomName, "name", "n", "", "Custom name for the executable")
	buildCmd.Flags().StringVarP(&buildLibrary, "library", "l", "pyinstaller", "Library to use for building the project (select between pyinstaller and nuitka)")
	buildCmd.Flags().BoolVarP(&buildOneFile, "one-file", "f", false, "Build the project in one file")
	buildCmd.Flags().BoolVarP(&buildModelsSymlink, "models-symlink", "s", false, "Symlink the models directory to the build directory")
}
