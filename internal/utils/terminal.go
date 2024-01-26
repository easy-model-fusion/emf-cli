package utils

import (
	"github.com/pterm/pterm"
)

// AskForUsersInput asks the user for an input and returns it
func AskForUsersInput(inputDirective string) string {
	// Create an interactive text input with single line input mode
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

	// Show the text input and get the result
	result, _ := textInput.Show(inputDirective)

	// Print a blank line for better readability
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
