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
		projectName = askForProjectName()
	} else {
		projectName = args[0]
	}

	createProjectSpinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Creating project '%s'...", projectName))

	// smooth animation
	time.Sleep(500 * time.Millisecond)

	err := createProject(projectName, func(msg string) {
		createProjectSpinner.UpdateText(msg)
	})

	// smooth animation
	time.Sleep(1 * time.Second)

	// check for errors
	if err != nil {
		if !os.IsExist(err) {
			removeErr := os.RemoveAll(projectName)
			if removeErr != nil {
				createProjectSpinner.Fail(fmt.Sprintf("Error deleting folder '%s': %s", projectName, removeErr))
				os.Exit(1)
			}
		}
		createProjectSpinner.Fail(fmt.Sprintf("Error creating project '%s': %s", projectName, err))
		os.Exit(1)
	}

	createProjectSpinner.Success("Project created successfully")
}

// createProject creates a new project with the given name
func createProject(projectName string, logger func(msg string)) (err error) {
	// check if folder exists
	if _, err = os.Stat(projectName); err != nil && !os.IsNotExist(err) {
		return err
	}

	logger("Creating project folder...")

	// Create folder
	err = os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		return err
	}

	logger("Creating project files...")

	// Copy main.py & config.yaml
	err = utils.CopyEmbeddedFile(sdk.EmbeddedFiles, "main.py", filepath.Join(projectName, "main.py"))
	if err != nil {
		return err
	}

	err = utils.CopyEmbeddedFile(sdk.EmbeddedFiles, "config.yaml", filepath.Join(projectName, "config.yaml"))
	if err != nil {
		return err
	}

	return nil
}

// askForProjectName asks the user for a project name and returns it
func askForProjectName() string {
	// Create an interactive text input with single line input mode
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

	// Show the text input and get the result
	result, _ := textInput.Show("Enter a project name")

	// Print a blank line for better readability
	pterm.Println()

	return result
}

func init() {
	rootCmd.AddCommand(initCmd)
}
