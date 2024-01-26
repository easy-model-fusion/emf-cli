package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/easy-model-fusion/client/sdk"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project name>",
	Short: "Initialize a EMF project",
	Long:  `Initialize a EMF project.`,
	Args:  utils.ValidFileName(1, true),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	var projectName string

	// No args, check projectName in pterm
	if len(args) == 0 {
		projectName = utils.AskForUsersInput("Enter a project name")
	} else {
		projectName = args[0]
	}

	// Check if user has python installed
	path, ok := CheckForPython()
	if !ok {
		os.Exit(1)
	}

	// Check the latest sdk version
	pterm.Info.Println("Checking for latest sdk version...")
	latestSdkTag, err := utils.GetLatestTag("sdk")
	if err != nil {
		pterm.Error.Println("Error checking for latest sdk version:", err)
		os.Exit(1)
	}
	pterm.Info.Println("Using latest sdk version:", latestSdkTag)

	err = createProject(projectName, path, latestSdkTag)

	// smooth animation
	time.Sleep(1 * time.Second)

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
func createProject(projectName, pythonPath, sdkTag string) (err error) {
	// check if folder exists
	if _, err = os.Stat(projectName); err != nil && !os.IsNotExist(err) {
		return err
	}

	pterm.Info.Println("Creating project folder...")

	// Create folder
	err = os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		return err
	}

	pterm.Info.Println("Creating project files...")

	// Copy main.py & config.yaml
	err = utils.CopyEmbeddedFile(sdk.EmbeddedFiles, "main.py", filepath.Join(projectName, "main.py"))
	if err != nil {
		return err
	}

	err = utils.CopyEmbeddedFile(sdk.EmbeddedFiles, "config.yaml", filepath.Join(projectName, "config.yaml"))
	if err != nil {
		return err
	}

	// Create sdk folder
	err = os.Mkdir(filepath.Join(projectName, "sdk"), os.ModePerm)
	if err != nil {
		return err
	}

	// Clone SDK
	spinnerInfo, _ := pterm.DefaultSpinner.Start("Cloning sdk...")
	err = utils.CloneSDK(sdkTag, filepath.Join(projectName, "sdk"))
	if err != nil {
		spinnerInfo.Fail(err)
		return err
	}
	spinnerInfo.Success()

	// Create virtual environment
	spinnerInfo, _ = pterm.DefaultSpinner.Start("Creating virtual environment...")
	err = utils.CreateVirtualEnv(pythonPath, filepath.Join(projectName, ".venv"))
	if err != nil {
		spinnerInfo.Fail(err)
		return err
	}
	spinnerInfo.Success()

	// Install dependencies
	pipPath, err := utils.FindVEnvPipExecutable(filepath.Join(projectName, ".venv"))
	if err != nil {
		return err
	}

	spinnerInfo, _ = pterm.DefaultSpinner.Start("Installing dependencies...")
	err = utils.InstallDependencies(pipPath, filepath.Join(projectName, "sdk", "requirements.txt"))
	if err != nil {
		return err
	}
	spinnerInfo.Success()

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
