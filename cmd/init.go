package cmd

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/easy-model-fusion/emf-cli/sdk"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	useTorchCuda bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project name>",
	Short: "Initialize a EMF project",
	Long:  `Initialize a EMF project.`,
	Args:  fileutil.ValidFileName(1, true),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	var projectName string

	// No args, check projectName in pterm
	if len(args) == 0 {
		projectName = ptermutil.AskForUsersInput("Enter a project name")
	} else {
		projectName = args[0]
	}

	err := createProject(projectName)

	// check for errors
	if err != nil {
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

	pterm.Success.Println("Project created successfully!")
}

// createProject creates a new project with the given name
func createProject(projectName string) (err error) {
	// check if folder exists
	if _, err = os.Stat(projectName); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create folder
	spinner, _ := pterm.DefaultSpinner.Start("Creating project folder...")
	err = os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()

	// Check if user has python installed
	pythonPath, ok := python.CheckAskForPython()
	if !ok {
		os.Exit(1)
	}

	// Check the latest sdk version
	spinner, _ = pterm.DefaultSpinner.Start("Checking for latest sdk version...")
	sdkTag, err := app.G().GetLatestTag("sdk")
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error checking for latest sdk version: %s", err))
		os.Exit(1)
	}
	spinner.Success("Using latest sdk version: " + sdkTag)

	// Create project files
	spinner, _ = pterm.DefaultSpinner.Start("Creating project files...")
	if err = createProjectFiles(projectName, sdkTag); err != nil {
		spinner.Fail(err)
		return err
	}
	spinner.Success()

	// Clone SDK
	spinner, _ = pterm.DefaultSpinner.Start("Cloning sdk...")
	err = app.G().CloneSDK(sdkTag, filepath.Join(projectName, "sdk"))
	if err != nil {
		spinner.Fail("Unable to clone sdk", err)
		return err
	}
	spinner.Success()

	// Create virtual environment
	spinner, _ = pterm.DefaultSpinner.Start("Creating virtual environment...")
	err = python.CreateVirtualEnv(pythonPath, filepath.Join(projectName, ".venv"))
	if err != nil {
		spinner.Fail("Unable to create venv", err)
		return err
	}
	spinner.Success()

	// Install dependencies
	pipPath, err := python.FindVEnvExecutable(filepath.Join(projectName, ".venv"), "pip")
	if err != nil {
		return err
	}

	spinner, _ = pterm.DefaultSpinner.Start("Installing dependencies...")

	fileName := "requirements.txt"
	if useTorchCuda {
		fileName = "requirements_cuda.txt"
	}

	err = python.InstallDependencies(pipPath, filepath.Join(projectName, "sdk", fileName))
	if err != nil {
		return err
	}
	spinner.Success()

	return nil
}

// createProjectFiles creates the project files (main.py, config.yaml, .gitignore)
func createProjectFiles(projectName, sdkTag string) (err error) {
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

func init() {
	initCmd.Flags().BoolVarP(&useTorchCuda, "cuda", "c", false, "Use torch with cuda")
}
