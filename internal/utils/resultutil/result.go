package resultutil

import "github.com/easy-model-fusion/emf-cli/internal/app"

type ExecutionResult struct {
	Warnings []string
	Infos    []string
	Error    error
}

// AddWarnings adds a new warning messages
func (er *ExecutionResult) AddWarnings(warnings []string) {
	er.Warnings = append(er.Warnings, warnings...)
}

// AddInfos adds a new information messages
func (er *ExecutionResult) AddInfos(infos []string) {
	er.Infos = append(er.Infos, infos...)
}

// SetError sets an error
func (er *ExecutionResult) SetError(err error) {
	er.Error = err
}

// Display displays all the messages to the user
func (er *ExecutionResult) Display(successMessage string, errorMessage string) {
	for _, warning := range er.Warnings {
		app.UI().Warning().Printfln(warning)
	}
	for _, info := range er.Infos {
		app.UI().Warning().Printfln(info)
	}
	if er.Error == nil {
		app.UI().Success().Printfln(successMessage)
	} else {
		app.UI().Error().Printfln(errorMessage+"\n%s", er.Error.Error())
	}
}
