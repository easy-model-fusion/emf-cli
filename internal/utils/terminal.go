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
