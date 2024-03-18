package test

import "github.com/easy-model-fusion/emf-cli/internal/ui"

type MockUI struct {
	UserInputResult        string
	MultiselectResult      []string
	UserConfirmationResult bool
}

type MockSpinner struct{}

func (m MockSpinner) Success(message ...interface{}) {
}

func (m MockSpinner) Warning(message ...interface{}) {
}

func (m MockSpinner) Fail(message ...interface{}) {
}

func (m MockUI) StartSpinner(message string) ui.Spinner {
	return &MockSpinner{}
}

func (m MockUI) AskForUsersInput(message string) string {
	return m.UserInputResult
}

func (m MockUI) DisplayInteractiveMultiselect(msg string, options []string, checkMark ui.Checkmark, optionsDefaultAll, filter bool) []string {
	return m.MultiselectResult
}

func (m MockUI) DisplaySelectedItems(items []string) {
}

func (m MockUI) AskForUsersConfirmation(message string) bool {
	return m.UserConfirmationResult
}
