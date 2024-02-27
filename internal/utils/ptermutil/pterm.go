package ptermutil

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

// DisplayInteractiveMultiselect displays an interactive multiselect prompt to the user.
// It presents a message and a list of options, allowing the user to select multiple options.
// Returns the selected options.
func DisplayInteractiveMultiselect(msg string, options, defaultOptions []string, checkMark *pterm.Checkmark, filter bool) []string {
	// Create a new interactive multiselect printer with the options
	// Disable the filter and set the keys for confirming and selecting options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(filter).
		WithCheckmark(checkMark).
		WithDefaultText(msg).
		WithDefaultOptions(defaultOptions)

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	return selectedOptions
}

// DisplaySelectedItems prints the selected items in green color.
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
