package cmd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
	"time"
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
	if config.GetViperConfig(".") != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	if buildLibrary != "pyinstaller" && buildLibrary != "nuitka" {
		pterm.Error.Println("Invalid library selected")
		return
	}

	// check if buildDestination exists
	if _, err := os.Stat(buildDestination); os.IsNotExist(err) {

		if buildDestination == "dist" {
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
	pythonPath, err := app.Python().FindVEnvExecutable(".venv", "python")
	if err != nil {
		pterm.Error.Println("Error finding python executable")
		return
	}

	pipPath, err := app.Python().FindVEnvExecutable(".venv", "pip")
	if err != nil {
		pterm.Error.Println("Error finding pip executable")
		return
	}

	err = app.Python().ExecutePip(pipPath, []string{"install", buildLibrary})
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Error installing %s", buildLibrary))
		return
	}

	var libraryPath string

	switch buildLibrary {
	case "pyinstaller":
		libraryPath, err = app.Python().FindVEnvExecutable(".venv", "pyinstaller")
		if err != nil {
			pterm.Error.Println("Error finding pyinstaller executable")
			return
		}
	default:
		libraryPath = pythonPath
	}

	buildArgs := createBuildArgs()

	pterm.Info.Println(fmt.Sprintf("Building project using %s...", buildLibrary))
	pterm.Info.Println(fmt.Sprintf("Using the following arguments: %s", buildArgs))
	pterm.Info.Println(fmt.Sprintf("The project will be built to %s", buildDestination))

	command := exec.Command(libraryPath, buildArgs...)

	var errBuf strings.Builder
	command.Stderr = &errBuf

	now := time.Now()
	spinner := app.UI().StartSpinner("Building project...")
	err = command.Run()
	if err != nil {
		spinner.Fail("Error building project")
		pterm.Error.Println(errBuf.String())
		return
	}
	spinner.Success(fmt.Sprintf("Project built successfully in %s", time.Since(now).String()))

	if !buildCompress {
		return
	}

	pterm.Info.Println("Compressing the output file(s) into a tarball file...")

	if buildIncludeModels {
		pterm.Info.Println("Including models in the build compressed file...")
	} else {
		pterm.Info.Println("Excluding models from the build compressed file...")
	}

	pterm.Info.Println("This may take a while... don't close the terminal")

}

// createBuildArgs creates the arguments for the build command
func createBuildArgs() []string {
	var buildArgs []string

	if buildCustomName == "" {
		buildCustomName = viper.GetString("name")
	}

	switch buildLibrary {
	case "pyinstaller":
		if buildOneFile {
			buildArgs = append(buildArgs, "-F")
		}

		buildArgs = append(buildArgs, fmt.Sprintf("-n %s", buildCustomName))
		buildArgs = append(buildArgs, fmt.Sprintf("--distpath=%s", buildDestination))
		buildArgs = append(buildArgs, viper.GetStringSlice("build.pyinstaller.args")...)
	case "nuitka":
		buildArgs = append(buildArgs, "-m nuitka")

		if buildOneFile {
			buildArgs = append(buildArgs, "--onefile")
		}

		buildArgs = append(buildArgs, fmt.Sprintf("--python-flag=-o %s", buildCustomName))
		buildArgs = append(buildArgs, fmt.Sprintf("--output-dir=%s", buildDestination))
		buildArgs = append(buildArgs, viper.GetStringSlice("build.nuitka.args")...)
	}

	buildArgs = append(buildArgs, "main.py")

	return stringutil.SliceRemoveDuplicates(buildArgs)
}

func init() {
	buildCmd.Flags().StringVarP(&buildDestination, "out-dir", "o", "dist", "Destination directory where the project will be built")
	buildCmd.Flags().StringVarP(&buildCustomName, "name", "n", "", "Custom name for the executable")
	buildCmd.Flags().StringVarP(&buildLibrary, "library", "l", "pyinstaller", "Library to use for building the project (select between pyinstaller and nuitka)")
	buildCmd.Flags().BoolVarP(&buildOneFile, "one-file", "f", false, "Build the project in one file")
	buildCmd.Flags().BoolVarP(&buildCompress, "compress", "c", false, "Compress the output file(s) into a tarball file")
	buildCmd.Flags().BoolVarP(&buildIncludeModels, "include-models", "m", false, "Include models in the build compressed file")
}
