package ui

type Checkmark struct {
	Checked   string
	Unchecked string
}

type UI interface {
	AskForUsersInput(message string) string
	DisplayInteractiveMultiselect(msg string, options []string, checkMark Checkmark, filter bool) []string
	DisplaySelectedItems(items []string)
	AskForUsersConfirmation(message string) bool
}
