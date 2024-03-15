package test

import (
	"github.com/easy-model-fusion/emf-cli/internal/ui"
)

type MockPython struct {
	Path         string
	Success      bool
	Error        error
	ScriptResult []byte
	ScriptExit   int
}

func (m MockPython) CheckPythonVersion(_ string) (string, bool) {
	return m.Path, m.Success
}

func (m MockPython) CheckForPython() (string, bool) {
	return m.Path, m.Success
}

func (m MockPython) CreateVirtualEnv(_, _ string) error {
	return m.Error
}

func (m MockPython) FindVEnvExecutable(_ string, _ string) (string, error) {
	return m.Path, m.Error
}

func (m MockPython) InstallDependencies(_, _ string) error {
	return m.Error
}

func (m MockPython) ExecutePip(_ string, _ []string) error {
	return m.Error
}

func (m MockPython) ExecuteScript(_, _ string, _ []string) ([]byte, error, int) {
	return m.ScriptResult, m.Error, m.ScriptExit
}

func (m MockPython) CheckAskForPython(_ ui.UI) (string, bool) {
	return m.Path, m.Success
}
