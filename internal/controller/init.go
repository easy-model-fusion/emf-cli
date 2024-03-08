package controller

import (
	"errors"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/easy-model-fusion/emf-cli/sdk"
	"github.com/pterm/pterm"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var initDependenciesPath = filepath.Join("sdk", "requirements.txt")

// RunInit runs the init command
func RunInit(args []string, useTorchCuda bool, customTag string) {
	var projectName string

	// No args, check projectName in ui
	if len(args) == 0 {
		projectName = app.UI().AskForUsersInput("Enter a project name")
	} else {
		projectName = args[0]
	}

	err := createProject(projectName, useTorchCuda, customTag)

	// check for errors
	if err == nil {
		pterm.Success.Println("Project created successfully!")
		return
	}

	if !os.IsExist(err) {
		removeErr := os.RemoveAll(projectName)
		if removeErr != nil {
			pterm.Warning.Println(fmt.Sprintf("Error deleting folder '%s': %s", projectName, removeErr))
			os.Exit(1)
		}
	}

	pterm.Error.Println(fmt.Sprintf("Error creating project '%s': %s", projectName, err))
	os.Exit(1)
}

// createProject creates a new project with the given name
func createProject(projectName string, useTorchCuda bool, customTag string) (err error) {
	// Create project folder
	if err = createProjectFolder(projectName); err != nil {
		return err
	}

	// Check if user has python installed
	pythonPath, ok := python.CheckAskForPython(app.UI())
	if !ok {
		os.Exit(1)
	}

	// Clone sdk
	if err = cloneSDK(projectName, customTag); err != nil {
		return err
	}

	// Create virtual environment
	spinner := app.UI().StartSpinner("Creating virtual environment")
	err = python.CreateVirtualEnv(pythonPath, filepath.Join(projectName, ".venv"))
	if err != nil {
		spinner.Fail("Unable to create venv: ", err)
		return err
	}
	spinner.Success()

	// Install dependencies
	if err = installDependencies(projectName, useTorchCuda); err != nil {
		return err
	}

	return nil
}

// createProjectFolder creates the project folder
func createProjectFolder(projectName string) (err error) {
	// check if folder exists
	if _, err = os.Stat(projectName); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create folder
	spinner := app.UI().StartSpinner("Creating project folder")
	err = os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()
	return nil
}

// createProjectFiles creates the project files (main.py, config.yaml, .gitignore)
func createProjectFiles(projectName, sdkTag string) (err error) {
	spinner := app.UI().StartSpinner("Creating project files")
	defer func() {
		if err != nil {
			spinner.Fail(err)
		} else {
			spinner.Success()
		}
	}()

	// Copy main.py, config.yaml & .gitignore
	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "main.py", filepath.Join(projectName, "main.py"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "config.yaml", filepath.Join(projectName, "config.yaml"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, ".gitignore", filepath.Join(projectName, ".gitignore"))
	if err != nil {
		return err
	}

	err = config.GetViperConfig(projectName)
	if err != nil {
		return err
	}

	// Write project name and sdk tag to config
	viper.Set("name", projectName)
	viper.Set("sdk-tag", sdkTag)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	// Create sdk folder
	err = os.Mkdir(filepath.Join(projectName, "sdk"), os.ModePerm)
	if err != nil {
		return err
	}

	// Create models folder
	err = os.Mkdir(filepath.Join(projectName, "models"), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// installDependencies installs the dependencies for the project
func installDependencies(projectName string, useTorchCuda bool) (err error) {
	// Install dependencies
	pipPath, err := python.FindVEnvExecutable(filepath.Join(projectName, ".venv"), "pip")
	if err != nil {
		return err
	}

	spinner := app.UI().StartSpinner("Installing dependencies")
	err = python.InstallDependencies(pipPath, filepath.Join(projectName, initDependenciesPath))
	if err != nil {
		spinner.Fail("Unable to install dependencies: ", err)
		return err
	}
	spinner.Success()

	spinner = app.UI().StartSpinner("Installing torch")
	if useTorchCuda { // TODO: refactor this
		err = python.ExecutePip(pipPath, []string{"install", "torch", "-f", "https://download.pytorch.org/whl/torch_stable.html"})
		if err != nil {
			spinner.Fail("Unable to install torch cuda: ", err)
			return err
		}
	}
	spinner.Success()

	return nil
}

// cloneSDK clones the sdk into the project
func cloneSDK(projectName, tag string) (err error) {
	// Check the latest sdk version
	if tag != "" {
		pterm.Info.Println("Using custom sdk version: " + tag)
	} else {
		spinner := app.UI().StartSpinner("Checking for latest sdk version")
		tag, err = app.G().GetLatestTag("sdk")
		if err != nil {
			spinner.Fail(fmt.Sprintf("Error checking for latest sdk version: %s", err))
			return errors.New("error checking for latest sdk version")
		}
		spinner.Success("Using latest sdk version: " + tag)
	}

	// Create project files
	if err = createProjectFiles(projectName, tag); err != nil {
		return err
	}

	// Clone SDK
	spinner := app.UI().StartSpinner("Cloning SDK")
	err = app.G().CloneSDK(tag, filepath.Join(projectName, "sdk"))
	if err != nil {
		spinner.Fail("Unable to clone sdk: ", err)
		return err
	}
	spinner.Success()

	spinner = app.UI().StartSpinner("Reorganizing SDK files")

	// Move files from sdk/sdk to sdk/
	err = fileutil.MoveFiles(filepath.Join(projectName, "sdk", "sdk"), filepath.Join(projectName, "sdk"))
	if err != nil {
		spinner.Fail("Unable to move SDK files: ", err)
		return err
	}

	// remove sdk/sdk folder
	err = os.RemoveAll(filepath.Join(projectName, "sdk", "sdk"))
	if err != nil {
		spinner.Fail("Unable to remove sdk/sdk folder: ", err)
		return err
	}

	// remove .github/ folder
	err = os.RemoveAll(filepath.Join(projectName, "sdk", ".github"))
	if err != nil {
		spinner.Fail("Unable to remove .github folder: ", err)
		return err
	}
	spinner.Success()

	return nil
}
