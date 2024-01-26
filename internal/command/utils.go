package command

import (
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
)

// CheckForPython checks if python is available in the PATH
// If python is not available, a message is printed to the user and asks to specify the path to python
// Returns true if python is available and the PATH
// Returns false if python is not available
func CheckForPython() (string, bool) {
	pterm.Info.Println("Checking for Python...")
	path, ok := utils.CheckForPython()
	if ok {
		pterm.Success.Println("Python executable found! (" + path + ")")
		return path, true
	}

	pterm.Warning.Println("Python is not installed or not available in the PATH")

	if AskConfirmation("Do you want to specify the path to python?") {
		result := AskInput("Enter python PATH")

		if result == "" {
			pterm.Error.Println("Please enter a valid path")
			return "", false
		}

		path, ok := utils.CheckPythonVersion(result)
		if ok {
			pterm.Success.Println("Python executable found! (" + path + ")")
			return path, true
		}

		pterm.Error.Println("Could not run python with the --version flag, please check the path to python")
		return "", false
	}

	pterm.Warning.Println("Please install Python 3.10 or higher and add it to the PATH")
	pterm.Warning.Println("You can download Python here: https://www.python.org/downloads/")
	pterm.Warning.Println("If you have already installed Python, please add it to the PATH")

	return "", false
}

// AskInput asks the user for an input
func AskInput(message string) string {
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)
	result, _ := textInput.Show(message)
	pterm.Println()
	return result
}

// AskConfirmation asks the user for a confirmation, returns true if the user confirms, false otherwise
func AskConfirmation(message string) bool {
	confirmation, _ := pterm.DefaultInteractiveConfirm.Show(message)
	pterm.Println()
	return confirmation
}
