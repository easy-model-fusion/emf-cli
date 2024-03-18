package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

var (
	buildDestination   string
	buildCustomName    string
	buildOneFile       bool
	buildCompress      bool
	buildIncludeModels bool
	buildLibrary       string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long: `Build the project using the selected library (pyinstaller or nuitka)
			and compress the output file(s) into a tarball file.
			You can also include the models in the build compressed file.
			Note: if you want to use nuitka, you need to have a working C compiler.`,
	Run: runBuild,
}

func runBuild(cmd *cobra.Command, args []string) {
	controller.RunBuild(buildCustomName, buildLibrary, buildDestination, buildCompress, buildIncludeModels, buildOneFile)
}

func init() {
	buildCmd.Flags().StringVarP(&buildDestination, "out-dir", "o", "dist", "Destination directory where the project will be built")
	buildCmd.Flags().StringVarP(&buildCustomName, "name", "n", "", "Custom name for the executable")
	buildCmd.Flags().StringVarP(&buildLibrary, "library", "l", "pyinstaller", "Library to use for building the project (select between pyinstaller and nuitka)")
	buildCmd.Flags().BoolVarP(&buildOneFile, "one-file", "f", false, "Build the project in one file")
	buildCmd.Flags().BoolVarP(&buildCompress, "compress", "c", false, "Compress the output file(s) into a tarball file")
	buildCmd.Flags().BoolVarP(&buildIncludeModels, "include-models", "m", false, "Include models in the build compressed file")
}
