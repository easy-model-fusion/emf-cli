package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	buildDestination   = "/dist"
	buildCustomName    string
	buildOneFile       bool
	buildCompress      bool
	buildIncludeModels bool
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the project",
	Long:  `Build the project.`,
	Run:   runBuild,
}

func init() {
	buildCmd.Flags().StringVarP(&buildDestination, "out-dir", "o", "", "Destination directory where the project will be built")
	buildCmd.Flags().StringVarP(&buildCustomName, "name", "n", "", "Custom name for the executable")
	buildCmd.Flags().BoolVarP(&buildOneFile, "one-file", "f", false, "Build the project in one file")
	buildCmd.Flags().BoolVarP(&buildCompress, "compress", "c", false, "Compress the output file(s) into a tarball file")
	buildCmd.Flags().BoolVarP(&buildIncludeModels, "include-models", "m", false, "Include models in the build compressed file")
}

func runBuild(cmd *cobra.Command, args []string) {
	if config.GetViperConfig(".") != nil {
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

	// Install dependencies
	pythonPath, err := python.FindVEnvExecutable(".venv", "python")
	if err != nil {
		pterm.Error.Println("Error finding python executable")
		return
	}

	var pyinstaller *exec.Cmd

	if buildOneFile {
		pyinstaller = exec.Command(pythonPath, "-m", "pyinstaller", "--onefile", "main.py")
	} else {
		pyinstaller = exec.Command(pythonPath, "-m", "pyinstaller", "main.py")
	}

	err = pyinstaller.Run()
	if err != nil {
		pterm.Error.Println("Error building project")
		return
	}

	// if buildCompress {
	// }

	pterm.Success.Println("Project built successfully!")
}
