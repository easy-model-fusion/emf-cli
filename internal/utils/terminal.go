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

// AskForUsersConfirmation asks the user for a confirmation, returns true if the user confirms, false otherwise
func AskForUsersConfirmation(message string) bool {
	confirmation, _ := pterm.DefaultInteractiveConfirm.Show(message)
	pterm.Println()
	return confirmation
}
