package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
	"time"
)

func RunBuild(customName, library, destDir string, compress, includeModels, oneFile bool) {
	if config.GetViperConfig(".") != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	if library != "pyinstaller" && library != "nuitka" {
		pterm.Error.Println("Invalid library selected")
		return
	}

	// check if destDir exists
	if _, err := os.Stat(destDir); os.IsNotExist(err) {

		if destDir == "dist" {
			err = os.Mkdir(destDir, os.ModePerm)
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

	err = app.Python().ExecutePip(pipPath, []string{"install", library})
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Error installing %s", library))
		return
	}

	var libraryPath string

	switch library {
	case "pyinstaller":
		libraryPath, err = app.Python().FindVEnvExecutable(".venv", "pyinstaller")
		if err != nil {
			pterm.Error.Println("Error finding pyinstaller executable")
			return
		}
	default:
		libraryPath = pythonPath
	}

	buildArgs := createBuildArgs(customName, library, destDir, oneFile)

	pterm.Info.Println(fmt.Sprintf("Building project using %s...", library))
	pterm.Info.Println(fmt.Sprintf("Using the following arguments: %s", buildArgs))
	pterm.Info.Println(fmt.Sprintf("The project will be built to %s", destDir))

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

	if !compress {
		return
	}

	pterm.Info.Println("Compressing the output file(s) into a tarball file...")

	if includeModels {
		pterm.Info.Println("Including models in the build compressed file...")
	} else {
		pterm.Info.Println("Excluding models from the build compressed file...")
	}

	pterm.Info.Println("This may take a while... don't close the terminal")

}

// createBuildArgs creates the arguments for the build command
func createBuildArgs(customName, library, destDir string, oneFile bool) []string {
	var buildArgs []string

	if customName == "" {
		customName = viper.GetString("name")
	}

	switch library {
	case "pyinstaller":
		if oneFile {
			buildArgs = append(buildArgs, "-F")
		}

		buildArgs = append(buildArgs, fmt.Sprintf("--name=%s", customName))
		buildArgs = append(buildArgs, fmt.Sprintf("--distpath=%s", destDir))
		buildArgs = append(buildArgs, viper.GetStringSlice("build.pyinstaller.args")...)
	case "nuitka":
		buildArgs = append(buildArgs, "-m nuitka")

		if oneFile {
			buildArgs = append(buildArgs, "--onefile")
		}

		buildArgs = append(buildArgs, fmt.Sprintf("--python-flag=-o %s", customName))
		buildArgs = append(buildArgs, fmt.Sprintf("--output-dir=%s", destDir))
		buildArgs = append(buildArgs, viper.GetStringSlice("build.nuitka.args")...)
	}

	buildArgs = append(buildArgs, "main.py")

	return stringutil.SliceRemoveDuplicates(buildArgs)
}
