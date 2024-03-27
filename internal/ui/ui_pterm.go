package ui

import "github.com/pterm/pterm"

type ptermUI struct {
	infoPrinter       ptermPrinter
	successPrinter    ptermPrinter
	errorPrinter      ptermPrinter
	warningPrinter    ptermPrinter
	defaultBoxPrinter ptermDefaultBoxPrinter
}

type ptermPrinter struct {
	printer pterm.PrefixPrinter
}

type ptermDefaultBoxPrinter struct{}

// NewPTermUI creates a new ptermUI instance
func NewPTermUI() UI {
	return &ptermUI{
		warningPrinter:    newPTermPrinter(pterm.Warning),
		errorPrinter:      newPTermPrinter(pterm.Error),
		successPrinter:    newPTermPrinter(pterm.Success),
		infoPrinter:       newPTermPrinter(pterm.Info),
		defaultBoxPrinter: ptermDefaultBoxPrinter{},
	}
}

// newPTermPrinter creates a new ptermPrinter with the given prefix printer
func newPTermPrinter(printer pterm.PrefixPrinter) ptermPrinter {
	return ptermPrinter{printer: printer}
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
func (p ptermUI) DisplayInteractiveMultiselect(msg string, options []string, checkMark Checkmark, optionsDefaultAll, filter bool, maxHeight int) []string {
	// Create a new interactive multiselect printer with the options
	// Disable the filter and set the keys for confirming and selecting options
	printer := pterm.DefaultInteractiveMultiselect.
		WithOptions(options).
		WithFilter(filter).
		WithCheckmark(&pterm.Checkmark{Checked: checkMark.Checked, Unchecked: checkMark.Unchecked}).
		WithDefaultText(msg)

	if maxHeight > 0 {
		printer.MaxHeight = maxHeight
	} else {
		printer.MaxHeight = 5
	}

	if optionsDefaultAll {
		printer = printer.WithDefaultOptions(options)
	}

	// Show the interactive multiselect and get the selected options
	selectedOptions, _ := printer.Show()

	return selectedOptions
}

// DisplayInteractiveSelect displays an interactive select (only one selectable option)
func (p ptermUI) DisplayInteractiveSelect(msg string, options []string, filter bool, maxHeight int) string {
	interactiveSelect := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithDefaultText(msg).
		WithFilter(filter)

	if maxHeight > 0 {
		interactiveSelect.MaxHeight = maxHeight
	} else {
		interactiveSelect.MaxHeight = 5
	}

	selectedOption, _ := interactiveSelect.Show()
	return selectedOption
}

// DisplaySelectedItems prints the selected items in green color.
func (p ptermUI) DisplaySelectedItems(items []string) {
	// Print the selected options, highlighted in green.
	p.Info().Printfln("Selected options: %s", pterm.Green(items))
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

// Info returns a Printer interface for printing info messages
func (p ptermUI) Info() Printer {
	return &p.infoPrinter
}

// Success returns a Printer interface for printing success messages
func (p ptermUI) Success() Printer {
	return &p.successPrinter
}

// Error returns a Printer interface for printing error messages
func (p ptermUI) Error() Printer {
	return &p.errorPrinter
}

// Warning returns a Printer interface for printing warning messages
func (p ptermUI) Warning() Printer {
	return &p.warningPrinter
}

// DefaultBox returns a Printer interface for printing messages in a default box
func (p ptermUI) DefaultBox() Printer {
	return &p.defaultBoxPrinter
}

// Printfln prints the given arguments with a newline
func (p ptermPrinter) Printfln(format string, a ...interface{}) {
	p.printer.Printfln(format, a...)
}

// Printf prints the given arguments with a newline
func (p ptermPrinter) Printf(format string, a ...interface{}) {
	p.printer.Printf(format, a...)
}

// Println prints the given arguments with a newline
func (p ptermPrinter) Println(a ...interface{}) {
	p.printer.Println(a...)
}

// Print prints the given arguments
func (p ptermPrinter) Print(a ...interface{}) {
	p.printer.Print(a...)
}

// Println prints the given arguments into a default box with a newline
func (p ptermDefaultBoxPrinter) Println(a ...interface{}) {
	pterm.DefaultBox.Println(a...)
}

// Printf prints the given arguments into a default box
func (p ptermDefaultBoxPrinter) Printf(format string, a ...interface{}) {
	pterm.DefaultBox.Printf(format, a...)
}

// Printfln prints the given arguments into a default box with a newline
func (p ptermDefaultBoxPrinter) Printfln(format string, a ...interface{}) {
	pterm.DefaultBox.Printfln(format, a...)
}

// Print prints the given arguments into a default box
func (p ptermDefaultBoxPrinter) Print(a ...interface{}) {
	pterm.DefaultBox.Print(a...)
}
