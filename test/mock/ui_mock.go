package mock

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
)

type MockUI struct {
	UserInputResult        string
	MultiselectResult      []string
	SelectResult           string
	UserConfirmationResult bool
}

type mockPrinter struct {
	printerType string
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

func (m MockUI) DisplayInteractiveSelect(_ string, _ []string, _ bool) string {
	return m.SelectResult
}

func (m MockUI) DisplaySelectedItems(_ []string) {
}

func (m MockUI) AskForUsersConfirmation(_ string) bool {
	return m.UserConfirmationResult
}

func (m MockUI) Info() ui.Printer {
	return &mockPrinter{
		printerType: "info",
	}
}

func (m MockUI) Success() ui.Printer {
	return &mockPrinter{
		printerType: "success",
	}
}

func (m MockUI) Error() ui.Printer {
	return &mockPrinter{
		printerType: "error",
	}
}

func (m MockUI) Warning() ui.Printer {
	return &mockPrinter{
		printerType: "warning",
	}
}

func (m MockUI) DefaultBox() ui.Printer {
	return &mockPrinter{
		printerType: "default-box",
	}
}

func (m mockPrinter) Printfln(format string, a ...interface{}) {
	fmt.Printf("[%s] %s\n", m.printerType, fmt.Sprintf(format, a...))
}

func (m mockPrinter) Printf(format string, a ...interface{}) {
	fmt.Printf("[%s] %s", m.printerType, fmt.Sprintf(format, a...))
}

func (m mockPrinter) Println(a ...interface{}) {
	fmt.Printf("[%s] %s\n", m.printerType, fmt.Sprint(a...))
}

func (m mockPrinter) Print(a ...interface{}) {
	fmt.Printf("[%s] %s", m.printerType, fmt.Sprint(a...))
}
