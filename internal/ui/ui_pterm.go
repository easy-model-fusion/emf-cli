package ui

import "github.com/pterm/pterm"

type ptermUI struct{}

func NewPTermUI() UI {
	return &ptermUI{}
}

// AskForUsersInput asks the user for an input and returns it
func (p ptermUI) AskForUsersInput(message string) string {
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)
	result, _ := textInput.Show(message)
	pterm.Println()
	return result
}

// DisplayInteractiveMultiselect displays an interactive multiselect prompt to the user.
// It presents a message and a list of options, allowing the user to select multiple options.
// Returns the selected options.
func (p ptermUI) DisplayInteractiveMultiselect(msg string, options []string, checkMark Checkmark, optionsDefaultAll, filter bool) []string {
	// Create a new interactive multiselect printer with the options
	// Disable the filter and set the keys for confirming and selecting options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(filter).
		WithCheckmark(&pterm.Checkmark{Checked: checkMark.Checked, Unchecked: checkMark.Unchecked}).
		WithDefaultText(msg)

	if optionsDefaultAll {
		printer = printer.WithDefaultOptions(options)
	}

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	return selectedOptions
}

// DisplaySelectedItems prints the selected items in green color.
func (p ptermUI) DisplaySelectedItems(items []string) {
	// Print the selected options, highlighted in green.
	pterm.Info.Printfln("Selected options: %s", pterm.Green(items))
}

// AskForUsersConfirmation asks the user for a confirmation, returns true if the user confirms, false otherwise
func (p ptermUI) AskForUsersConfirmation(message string) bool {
	confirmation, _ := pterm.DefaultInteractiveConfirm.Show(message)
	pterm.Println()
	return confirmation
}

// StartSpinner starts a new spinner with the given message and returns a Spinner interface
func (p ptermUI) StartSpinner(message string) Spinner {
	spinner, _ := pterm.DefaultSpinner.Start(message)
	return spinner
}
