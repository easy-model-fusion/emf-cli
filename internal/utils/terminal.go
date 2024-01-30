package utils

import (
	"github.com/pterm/pterm"
)

// AskForUsersInput asks the user for an input and returns it
func AskForUsersInput(message string) string {
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)
	result, _ := textInput.Show(message)
	pterm.Println()
	return result
}

func DisplayInteractiveMultiselect(msg string, options []string, checkMark *pterm.Checkmark, filter bool) []string {
	// Create a new interactive multiselect printer with the options
	// Disable the filter and set the keys for confirming and selecting options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(filter).
		WithCheckmark(checkMark).
		WithDefaultText(msg)

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	return selectedOptions
}

func DisplaySelectedItems(items []string) {
	// Print the selected options, highlighted in green.
	pterm.Info.Printfln("Selected options: %s", pterm.Green(items))
}

// AskForUsersConfirmation asks the user for a confirmation, returns true if the user confirms, false otherwise
func AskForUsersConfirmation(message string) bool {
	confirmation, _ := pterm.DefaultInteractiveConfirm.Show(message)
	pterm.Println()
	return confirmation
}

// CheckAskForPython checks if python is available in the PATH
// If python is not available, a message is printed to the user and asks to specify the path to python
// Returns true if python is available and the PATH
// Returns false if python is not available
func CheckAskForPython() (string, bool) {
	pterm.Info.Println("Checking for Python...")
	path, ok := CheckForPython()
	if ok {
		pterm.Success.Println("Python executable found! (" + path + ")")
		return path, true
	}

	pterm.Warning.Println("Python is not installed or not available in the PATH")

	if AskForUsersConfirmation("Do you want to specify the path to python?") {
		result := AskForUsersInput("Enter python PATH")

		if result == "" {
			pterm.Error.Println("Please enter a valid path")
			return "", false
		}

		path, ok := CheckPythonVersion(result)
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
