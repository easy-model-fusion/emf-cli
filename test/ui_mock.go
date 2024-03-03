package test

import "github.com/easy-model-fusion/emf-cli/internal/ui"

type MockUI struct {
	UserInputResult        string
	MultiselectResult      []string
	UserConfirmationResult bool
}

func NewMockUI() ui.UI {
	return &MockUI{}
}

func (m MockUI) AskForUsersInput(message string) string {
	return m.UserInputResult
}

func (m MockUI) DisplayInteractiveMultiselect(msg string, options []string, checkMark ui.Checkmark, filter bool) []string {
	return m.MultiselectResult
}

func (m MockUI) DisplaySelectedItems(items []string) {
}

func (m MockUI) AskForUsersConfirmation(message string) bool {
	return m.UserConfirmationResult
}
