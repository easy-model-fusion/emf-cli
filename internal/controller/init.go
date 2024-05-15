// Package controller
// This file contains the init controller which is responsible for initializing a new project.
// Whenever a user wants to create a new project, he will run the init command that uses the init controller.
// The init controller will create a new project folder with the given name.
// It will check if the user has python installed and clone the sdk into the project.
// After that, it will create a virtual environment, install the dependencies, and create the project files.
//
// The final project folder should contain the following files:
// new-project/
// ├── .gitignore
// ├── config.yaml
// ├── main.py
// ├── requirements.txt
// ├── sdk/
// ├── models/
// └── .venv/
package controller

import (
	"errors"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/sdk"
	"github.com/spf13/viper"
	"os"
)

type InitController struct{}

var initDependenciesPath = fileutil.PathJoin("sdk", "requirements.txt")

// Run runs the init command
func (ic InitController) Run(args []string, useTorchCuda bool, customTag string) error {
	var projectName string

	// No args, check projectName in ui
	if len(args) == 0 {
		projectName = app.UI().AskForUsersInput("Enter a project name")
	} else {
		projectName = args[0]
	}

	err := ic.createProject(projectName, useTorchCuda, customTag)

	// check for errors
	if err == nil {
		app.UI().Success().Println("Project created successfully!")
		return nil
	}

	if !os.IsExist(err) {
		removeErr := os.RemoveAll(projectName)
		if removeErr != nil {
			app.UI().Warning().Println(fmt.Sprintf("Error deleting folder '%s': %s", projectName, removeErr))
			return removeErr
		}
	}

	app.UI().Error().Println(fmt.Sprintf("Error creating project '%s': %s", projectName, err))
	return err
}

// createProject creates a new project with the given name
func (ic InitController) createProject(projectName string, useTorchCuda bool, customTag string) (err error) {
	// Create project folder
	if err = ic.createProjectFolder(projectName); err != nil {
		return err
	}

	// Check if user has python installed
	pythonPath, ok := app.Python().CheckAskForPython(app.UI())
	if !ok {
		return errors.New("python not found")
	}

	// Clone sdk
	if err = ic.cloneSDK(projectName, customTag); err != nil {
		return err
	}

	// Create virtual environment
	spinner := app.UI().StartSpinner("Creating virtual environment")
	err = app.Python().CreateVirtualEnv(pythonPath, fileutil.PathJoin(projectName, ".venv"))
	if err != nil {
		spinner.Fail("Unable to create venv: ", err)
		return err
	}
	spinner.Success()

	// Install dependencies
	if err = ic.installDependencies(projectName, useTorchCuda); err != nil {
		return err
	}

	return nil
}

// createProjectFolder creates the project folder
func (ic InitController) createProjectFolder(projectName string) (err error) {
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
func (ic InitController) createProjectFiles(projectName, sdkTag string) (err error) {
	spinner := app.UI().StartSpinner("Creating project files")
	defer func() {
		if err != nil {
			spinner.Fail(err)
		} else {
			spinner.Success()
		}
	}()

	// Copy main.py, config.yaml & .gitignore
	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "main.py", fileutil.PathJoin(projectName, "main.py"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "config.yaml", fileutil.PathJoin(projectName, "config.yaml"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, ".gitignore", fileutil.PathJoin(projectName, ".gitignore"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "README.md", fileutil.PathJoin(projectName, "README.md"))
	if err != nil {
		return err
	}

	err = fileutil.CopyEmbeddedFile(sdk.EmbeddedFiles, "requirements.txt", fileutil.PathJoin(projectName, "requirements.txt"))
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
	err = os.Mkdir(fileutil.PathJoin(projectName, "sdk"), os.ModePerm)
	if err != nil {
		return err
	}

	// Create models folder
	err = os.Mkdir(fileutil.PathJoin(projectName, "models"), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// installDependencies installs the dependencies for the project
func (ic InitController) installDependencies(projectName string, useTorchCuda bool) (err error) {
	// Install dependencies
	pipPath, err := app.Python().FindVEnvExecutable(fileutil.PathJoin(projectName, ".venv"), "pip")
	if err != nil {
		return err
	}

	spinner := app.UI().StartSpinner("Installing dependencies")
	err = app.Python().InstallDependencies(pipPath, fileutil.PathJoin(projectName, initDependenciesPath))
	if err != nil {
		spinner.Fail("Unable to install dependencies: ", err)
		return err
	}
	spinner.Success()

	if useTorchCuda {
		spinner = app.UI().StartSpinner("Installing torch cuda")
		err = app.Python().ExecutePip(pipPath, []string{"uninstall", "-y", "torch"})
		if err != nil {
			spinner.Fail("Unable to uninstall torch: ", err)
			return err
		}

		err = app.Python().ExecutePip(pipPath, []string{"install", "torch", "-f", app.TorchCudaURL})
		if err != nil {
			spinner.Fail("Unable to install torch cuda: ", err)
			return err
		}
		spinner.Success()
	}

	return nil
}

// cloneSDK clones the sdk into the project
func (ic InitController) cloneSDK(projectName, tag string) (err error) {
	// Check the latest sdk version
	if tag != "" {
		app.UI().Info().Println("Using custom sdk version: " + tag)
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
	if err = ic.createProjectFiles(projectName, tag); err != nil {
		return err
	}

	// Clone SDK
	spinner := app.UI().StartSpinner("Cloning SDK")
	err = app.G().CloneSDK(tag, fileutil.PathJoin(projectName, "sdk"))
	if err != nil {
		spinner.Fail("Unable to clone sdk: ", err)
		return err
	}
	spinner.Success()

	spinner = app.UI().StartSpinner("Reorganizing SDK files")

	// Move files from sdk/sdk to sdk/
	err = fileutil.MoveFiles(fileutil.PathJoin(projectName, "sdk", "sdk"), fileutil.PathJoin(projectName, "sdk"))
	if err != nil {
		spinner.Fail("Unable to move SDK files: ", err)
		return err
	}

	// remove sdk/sdk folder
	err = os.RemoveAll(fileutil.PathJoin(projectName, "sdk", "sdk"))
	if err != nil {
		spinner.Fail("Unable to remove sdk/sdk folder: ", err)
		return err
	}

	// remove .github/ folder
	err = os.RemoveAll(fileutil.PathJoin(projectName, "sdk", ".github"))
	if err != nil {
		spinner.Fail("Unable to remove .github folder: ", err)
		return err
	}
	spinner.Success()

	return nil
}
