package controller

import (
	"context"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
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
	pythonPath, err := bc.InstallDependencies(library)
	if err != nil {
		return err
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

	// Build the project
	err = bc.Build(customName, library, destDir, libraryPath, oneFile)
	if err != nil {
		return err
	}

	// Compress the output file(s) into a tarball file
	if compress {
		return bc.Compress(includeModels)
	}

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

// InstallDependencies installs the dependencies for the project
// returns the path to the python executable
func (bc BuildController) InstallDependencies(library string) (string, error) {
	pythonPath, err := app.Python().FindVEnvExecutable(".venv", "python")
	if err != nil {
		return "", fmt.Errorf("error finding python executable: %s", err.Error())
	}

	pipPath, err := app.Python().FindVEnvExecutable(".venv", "pip")
	if err != nil {
		return "", fmt.Errorf("error finding pip executable: %s", err.Error())
	}

	err = app.Python().ExecutePip(pipPath, []string{"install", library})
	if err != nil {
		return "", fmt.Errorf("error installing %s: %s", library, err.Error())
	}

	return pythonPath, nil
}

// Build builds the project
func (bc BuildController) Build(customName, library, destDir, libraryPath string, oneFile bool) (err error) {
	buildArgs := bc.createBuildArgs(customName, library, destDir, oneFile)

	app.UI().Info().Println(fmt.Sprintf("Building project using %s...", library))
	app.UI().Info().Println(fmt.Sprintf("Using the following arguments: %s", buildArgs))
	app.UI().Info().Println(fmt.Sprintf("The project will be built to %s", destDir))

	// Setup signal catching
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	command := exec.CommandContext(ctx, libraryPath, buildArgs...)

	var errBuf strings.Builder
	command.Stderr = &errBuf

	spinner := app.UI().StartSpinner("Building project...")
	now := time.Now()

	// Running the build command in a goroutine (to handle cancellation, since the build can take a long time)
	go func(err error) {
		err = command.Run()
		// Sending signal to the main goroutine that the script has finished
		done <- syscall.SIGQUIT
	}(err)

	switch code := <-done; {
	case code == syscall.SIGQUIT:
		// Do nothing
	case code == syscall.SIGINT:
		fallthrough
	case code == syscall.SIGTERM:
		cancel()
		spinner.Fail("Build cancelled manually after " + time.Since(now).String())
		return err
	}

	// make sure that the context is cancelled, even if the build has finished
	cancel()

	spinner.Success(fmt.Sprintf("Project built successfully in %s", time.Since(now).String()))
	return nil
}

// Compress compresses the output file(s) into a tarball file
func (bc BuildController) Compress(includeModels bool) error {
	app.UI().Info().Println("Compressing the output file(s) into a tarball file...")

	if includeModels {
		app.UI().Info().Println("Including models in the build compressed file...")
	} else {
		app.UI().Info().Println("Excluding models from the build compressed file...")
	}

	app.UI().Info().Println("This may take a while... don't close the terminal")
	return nil
}
