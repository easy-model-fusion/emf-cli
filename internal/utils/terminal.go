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

func DisplayInteractiveMultiselect(options []string) []string {
	// Create a new interactive multiselect printer with the options
	// Disable the filter and set the keys for confirming and selecting options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(false).
		WithCheckmark(&pterm.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")})

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	// Print the selected options, highlighted in green.
	pterm.Info.Printfln("Selected options: %s", pterm.Green(selectedOptions))

	return selectedOptions
}
