package mock

import "github.com/easy-model-fusion/emf-cli/internal/ui"

type MockUI struct {
	UserInputResult        string
	MultiselectResult      []string
	UserConfirmationResult bool
}

type MockSpinner struct{}

func (m MockSpinner) Success(_ ...interface{}) {
}

func (m MockSpinner) Warning(_ ...interface{}) {
}

func (m MockSpinner) Fail(_ ...interface{}) {
}

func (m MockUI) StartSpinner(_ string) ui.Spinner {
	return &MockSpinner{}
}

func (m MockUI) AskForUsersInput(_ string) string {
	return m.UserInputResult
}

func (m MockUI) DisplayInteractiveMultiselect(_ string, _ []string, _ ui.Checkmark, _, _ bool) []string {
	return m.MultiselectResult
}

func (m MockUI) DisplaySelectedItems(_ []string) {
}

func (m MockUI) AskForUsersConfirmation(_ string) bool {
	return m.UserConfirmationResult
}
