// Package controller
// This file contains the install controller which is responsible for installing an already existing project.
// Basically combining a slightly different init and the tidy controller features.
//
// Whenever a user wants to install an existing project, he will clone an existing project.
// The existing project file structure should look like this:
// cloned-project/
// ├── .gitignore
// ├── config.yaml
// ├── main.py
// ├── requirements.txt
// └── any user project related files
// However, if the user wants to run the project, there are some files missing.
// Here the install controller comes into play:
// - It will create a virtual environment
// - Clone the configured sdk version
// - Install the dependencies (requirements.txt & sdk/requirements.txt)
// - Install torch with cuda if the user wants to
// - Download the missing models
// - Generate the needed python code
//
// In the end, the user should have a fully working project that should look like this:
// cloned-project/
// ├── .gitignore
// ├── config.yaml
// ├── main.py
// ├── requirements.txt
// ├── sdk/
// ├── models/
// ├── .venv/
// └── any user project related files
package controller

import (
	"errors"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

type InstallController struct{}

// Run runs the install command
func (ic InstallController) Run(args []string, useTorchCuda bool, accessToken string) error {
	start := time.Now()

	// Only clean if config file exists (so we know it's a EMF project)
	if err := config.GetViperConfig(config.FilePath); err != nil {
		return err
	}

	// Check if user has python installed
	pythonPath, ok := app.Python().CheckAskForPython(app.UI())
	if !ok {
		return errors.New("python not found")
	}

	// Create missing directories
	if err := ic.createMissingDirectories(); err != nil {
		return err
	}

	// Clone SDK & move files
	if err := ic.cloneSDK(); err != nil {
		return err
	}

	// Create virtual environment & install dependencies
	if err := ic.installDependencies(pythonPath, useTorchCuda); err != nil {
		return err
	}

	// handle errors in run tidy (new structure)
	if err := tidyController.RunTidy(false, accessToken); err != nil {
		return err
	}

	app.UI().Success().Printfln("Project installed successfully in %v", time.Since(start))

	return nil
}

// createMissingDirectories creates the missing directories (sdk, models)
func (ic InstallController) createMissingDirectories() (err error) {
	spinner := app.UI().StartSpinner("Creating missing directories")
	defer func() {
		if err != nil {
			spinner.Fail(err)
		} else {
			spinner.Success()
		}
	}()

	// Create sdk folder
	err = os.Mkdir("sdk", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Create models folder
	err = os.Mkdir("models", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

// installDependencies installs the dependencies for the project
func (ic InstallController) installDependencies(pythonPath string, useTorchCuda bool) (err error) {
	// check if a venv is already installed
	app.UI().Info().Println("Checking if a venv is already installed")

	_, err = app.Python().FindVEnvExecutable(".venv", "python")
	if err != nil {
		app.UI().Info().Println("No venv found, creating a new one")

		// Create virtual environment
		spinner := app.UI().StartSpinner("Creating virtual environment")

		err = app.Python().CreateVirtualEnv(pythonPath, ".venv")
		if err != nil {
			spinner.Fail("Unable to create venv: ", err)
			return err
		}
		spinner.Success()

	} else {
		app.UI().Info().Println("Venv found, installing dependencies")
	}

	// Install dependencies
	spinner := app.UI().StartSpinner("Installing sdk dependencies")
	pipPath, err := app.Python().FindVEnvExecutable(".venv", "pip")
	if err != nil {
		spinner.Fail("Unable to find pip: ", err)
		return err
	}

	// First install the sdk dependencies
	err = app.Python().InstallDependencies(pipPath, "sdk/requirements.txt")
	if err != nil {
		spinner.Fail("Unable to install sdk dependencies: ", err)
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

	// Install the project dependencies
	spinner = app.UI().StartSpinner("Installing project dependencies")
	err = app.Python().InstallDependencies(pipPath, "requirements.txt")
	if err != nil {
		spinner.Warning("Unable to install project dependencies: ", err)
	} else {
		spinner.Success()
	}

	return nil
}

// cloneSDK clones the sdk into the project
func (ic InstallController) cloneSDK() (err error) {
	// Get sdk tag
	tag := viper.GetString("sdk-tag")

	// Clone SDK
	retry := false
clone:
	spinner := app.UI().StartSpinner("Cloning SDK")
	err = app.G().CloneSDK(tag, "sdk")
	if err != nil {
		spinner.Fail("Unable to clone sdk: ", err)

		if !retry && app.UI().AskForUsersConfirmation("Do you want to remove the sdk folder and try again?") {

			// Remove sdk folder
			err = os.RemoveAll("sdk")
			if err != nil {
				return fmt.Errorf("unable to remove sdk folder: %w", err)
			}

			// Retry cloning
			retry = true
			goto clone
		}

		if retry {
			app.UI().Warning().Println("Cloning the SDK failed twice, please retry the whole process.")
		}

		return err
	}
	spinner.Success()

	spinner = app.UI().StartSpinner("Reorganizing SDK files")

	// Move files from sdk/sdk to sdk/
	err = fileutil.MoveFiles(filepath.Join("sdk", "sdk"), "sdk")
	if err != nil {
		spinner.Fail("Unable to move SDK files: ", err)
		return err
	}

	// remove sdk/sdk folder
	err = os.RemoveAll(filepath.Join("sdk", "sdk"))
	if err != nil {
		spinner.Fail("Unable to remove sdk/sdk folder: ", err)
		return err
	}

	// remove .github/ folder
	err = os.RemoveAll(filepath.Join("sdk", ".github"))
	if err != nil {
		spinner.Fail("Unable to remove .github folder: ", err)
		return err
	}
	spinner.Success()

	return nil
}
