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

func (m MockPython) CheckPythonVersion(name string) (string, bool) {
	return m.Path, m.Success
}

func (m MockPython) CheckForPython() (string, bool) {
	return m.Path, m.Success
}

func (m MockPython) CreateVirtualEnv(pythonPath, path string) error {
	return m.Error
}

func (m MockPython) FindVEnvExecutable(venvPath string, executableName string) (string, error) {
	return m.Path, m.Error
}

func (m MockPython) InstallDependencies(pipPath, path string) error {
	return m.Error
}

func (m MockPython) ExecutePip(pipPath string, args []string) error {
	return m.Error
}

func (m MockPython) ExecuteScript(_, _ string, _ []string) ([]byte, error, int) {
	return m.ScriptResult, m.Error, m.ScriptExit
}

func (m MockPython) CheckAskForPython(ui ui.UI) (string, bool) {
	return m.Path, m.Success
}
