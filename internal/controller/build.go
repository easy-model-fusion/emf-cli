package controller

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
	"time"
)

type BuildController struct{}

func (bc BuildController) RunBuild(customName, library, destDir string, compress, includeModels, oneFile bool) error {
	if err := config.GetViperConfig("."); err != nil {
		return err
	}

	sdk.SendUpdateSuggestion()

	if library != "pyinstaller" && library != "nuitka" {
		return fmt.Errorf("invalid library selected")
	}

	// check if destDir exists
	if _, err := os.Stat(destDir); os.IsNotExist(err) {

		if destDir == "dist" {
			err = os.Mkdir(destDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("error creating dist folder: %s", err.Error())
			}
		} else {
			return fmt.Errorf("destination directory does not exist")
		}
	}

	// Install dependencies
	pythonPath, err := app.Python().FindVEnvExecutable(".venv", "python")
	if err != nil {
		return fmt.Errorf("error finding python executable: %s", err.Error())
	}

	pipPath, err := app.Python().FindVEnvExecutable(".venv", "pip")
	if err != nil {
		return fmt.Errorf("error finding pip executable: %s", err.Error())
	}

	err = app.Python().ExecutePip(pipPath, []string{"install", library})
	if err != nil {
		return fmt.Errorf("error installing %s: %s", library, err.Error())
	}

	var libraryPath string

	switch library {
	case "pyinstaller":
		libraryPath, err = app.Python().FindVEnvExecutable(".venv", "pyinstaller")
		if err != nil {
			return fmt.Errorf("error finding pyinstaller executable: %s", err.Error())
		}
	default:
		libraryPath = pythonPath
	}

	buildArgs := bc.createBuildArgs(customName, library, destDir, oneFile)

	app.UI().Info().Println(fmt.Sprintf("Building project using %s...", library))
	app.UI().Info().Println(fmt.Sprintf("Using the following arguments: %s", buildArgs))
	app.UI().Info().Println(fmt.Sprintf("The project will be built to %s", destDir))

	command := exec.Command(libraryPath, buildArgs...)

	var errBuf strings.Builder
	command.Stderr = &errBuf

	now := time.Now()
	spinner := app.UI().StartSpinner("Building project...")
	err = command.Run()
	if err != nil {
		spinner.Fail("Error building project")
		app.UI().Error().Println(errBuf.String())
		return fmt.Errorf("error building project: %s", err.Error())
	}
	spinner.Success(fmt.Sprintf("Project built successfully in %s", time.Since(now).String()))

	if !compress {
		return nil
	}

	app.UI().Info().Println("Compressing the output file(s) into a tarball file...")

	if includeModels {
		app.UI().Info().Println("Including models in the build compressed file...")
	} else {
		app.UI().Info().Println("Excluding models from the build compressed file...")
	}

	app.UI().Info().Println("This may take a while... don't close the terminal")

	return nil
}

// createBuildArgs creates the arguments for the build command
func (bc BuildController) createBuildArgs(customName, library, destDir string, oneFile bool) []string {
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
