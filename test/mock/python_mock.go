package mock

import (
	"context"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
)

type MockPython struct {
	Path                     string
	Success                  bool
	ExecuteScriptError       error
	CreateVirtualEnvError    error
	InstallDependenciesError error
	FindVEnvExecutableError  error
	ExecutePipError          error
	ScriptResult             []byte
	ScriptExit               int
	CalledFunctions          map[string]int
}

func (m MockPython) CheckPythonVersion(_ string) (string, bool) {
	m.callFunction("CheckPythonVersion")
	return m.Path, m.Success
}

func (m MockPython) CheckForPython() (string, bool) {
	m.callFunction("CheckForPython")
	return m.Path, m.Success
}

func (m MockPython) CreateVirtualEnv(_, _ string) error {
	m.callFunction("CreateVirtualEnv")
	return m.CreateVirtualEnvError
}

func (m MockPython) FindVEnvExecutable(_ string, _ string) (string, error) {
	m.callFunction("FindVEnvExecutable")
	return m.Path, m.FindVEnvExecutableError
}

func (m MockPython) InstallDependencies(_, _ string) error {
	m.callFunction("InstallDependencies")
	return m.InstallDependenciesError
}

func (m MockPython) ExecutePip(_ string, _ []string) error {
	m.callFunction("ExecutePip")
	return m.ExecutePipError
}

func (m MockPython) ExecuteScript(_, _ string, _ []string, _ context.Context) ([]byte, error, int) {
	m.callFunction("ExecuteScript")
	return m.ScriptResult, m.ExecuteScriptError, m.ScriptExit
}

func (m MockPython) CheckAskForPython(_ ui.UI) (string, bool) {
	m.callFunction("CheckAskForPython")
	return m.Path, m.Success
}

func (m MockPython) callFunction(name string) {
	if m.CalledFunctions == nil {
		return
	}
	m.CalledFunctions[name]++
}
