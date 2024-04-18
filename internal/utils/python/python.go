package python

import (
	"context"
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/executil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/pterm/pterm"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Python interface {
	CheckPythonVersion(name string) (string, bool)
	CheckForPython() (string, bool)
	CreateVirtualEnv(pythonPath, path string) error
	FindVEnvExecutable(venvPath string, executableName string) (string, error)
	InstallDependencies(pipPath, path string) error
	ExecutePip(pipPath string, args []string) error
	ExecuteScript(venvPath, filePath string, args []string, ctx context.Context) ([]byte, error, int)
	CheckAskForPython(ui ui.UI) (string, bool)
}

type python struct{}

func NewPython() Python {
	return &python{}
}

// CheckPythonVersion checks if python is found in the PATH and runs it with the --
// version flag to check if it works, and returns path to python executable and true if so.
// If python is not found, the function returns false.
func (p *python) CheckPythonVersion(name string) (string, bool) {
	path, ok := executil.CheckForExecutable(name)
	if !ok {
		return "", false
	}

	cmd := exec.Command(path, "--version")
	err := cmd.Run()
	if err == nil {
		return path, true
	}

	return "", false
}

// CheckForPython checks if python is available and works, and returns path to python executable and true if so.
func (p *python) CheckForPython() (string, bool) {
	path, ok := p.CheckPythonVersion("python")
	if ok {
		return path, true
	}
	return p.CheckPythonVersion("python3")
}

// FindVEnvExecutable searches for the requested executable within a virtual environment.
func (p *python) FindVEnvExecutable(venvPath string, executableName string) (string, error) {
	var pipPath string
	if runtime.GOOS == "windows" {
		pipPath = filepath.Join(venvPath, "Scripts", executableName+".exe")
	} else {
		pipPath = filepath.Join(venvPath, "bin", executableName)
	}

	if _, err := os.Stat(pipPath); os.IsNotExist(err) {
		return "", fmt.Errorf("'%s' executable not found in virtual environment: %s", executableName, pipPath)
	}

	return pipPath, nil
}

// InstallDependencies installs the dependencies from the given requirements.txt file
func (p *python) InstallDependencies(pipPath, path string) error {
	return p.runCommand(exec.Command(pipPath, "install", "-r", path))
}

// ExecutePip runs pip with the given arguments
func (p *python) ExecutePip(pipPath string, args []string) error {
	return p.runCommand(exec.Command(pipPath, args...))
}

// CreateVirtualEnv creates a virtual environment in the given path
func (p *python) CreateVirtualEnv(pythonPath, path string) error {
	cmd := exec.Command(pythonPath, "-m", "venv", path)
	return p.runCommand(cmd)
}

// runCommand runs the given command and returns an error if it fails
// The error message is appended with the stderr output
func (p *python) runCommand(cmd *exec.Cmd) error {
	// bind stderr to a buffer
	var errBuf strings.Builder
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		errBufStr := errBuf.String()
		if errBufStr != "" {
			return fmt.Errorf("%s: %s", err.Error(), errBufStr)
		}
		return err
	}

	return nil
}

// ExecuteScript runs the requested python file with the requested arguments
func (p *python) ExecuteScript(venvPath, filePath string, args []string, ctx context.Context) ([]byte, error, int) {

	// Find the python executable inside the venv to run the script
	pythonPath, err := p.FindVEnvExecutable(venvPath, "python")
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Error using the venv : %s", err))
		return nil, err, 1
	}

	// Checking that the script does exist
	exists, err := fileutil.IsExistingPath(filePath)
	if err != nil {
		pterm.Error.Println(fmt.Sprintf("Missing script '%s'", filePath))
		return nil, err, 1
	} else if !exists {
		err = fmt.Errorf("missing script '%s'", filePath)
		return nil, err, 1
	}

	// Create command
	var cmd = exec.CommandContext(ctx, pythonPath, append([]string{filePath}, args...)...)

	// Create pipe to capture stdout
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err, 1
	}

	// Bind stderr
	cmd.Stderr = os.Stderr

	// Start the command
	err = cmd.Start()
	if err != nil {
		return nil, err, 1
	}

	// Read output from the pipe
	output, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, err, 1
	}

	// Wait for the command to finish
	err = cmd.Wait()

	// Execution was successful but nothing returned
	if err == nil && len(output) == 0 {
		return nil, nil, 0
	}

	// Execution was successful
	if err == nil {
		return output, nil, 0
	}

	return nil, err, 1
}

// CheckAskForPython checks if python is available in the PATH
// If python is not available, a message is printed to the user and asks to specify the path to python
// Returns true if python is available and the PATH
// Returns false if python is not available
func (p *python) CheckAskForPython(ui ui.UI) (string, bool) {
	pterm.Info.Println("Checking for Python...")
	path, ok := p.CheckForPython()
	if ok {
		ui.Success().Println("Python executable found! (" + path + ")")
		return path, true
	}
	return p.askSpecificPythonPath(ui)
}

// askSpecificPythonPath asks the user if they want to specify the path to python
func (p *python) askSpecificPythonPath(ui ui.UI) (string, bool) {
	ui.Warning().Println("Python is not installed or not available in the PATH")

	if ui.AskForUsersConfirmation("Do you want to specify the path to python?") {
		result := ui.AskForUsersInput("Enter python PATH")

		if result == "" {
			ui.Error().Println("Please enter a valid path")
			return "", false
		}

		path, ok := p.CheckPythonVersion(result)
		if ok {
			ui.Success().Println("Python executable found! (" + path + ")")
			return path, true
		}

		ui.Error().Println("Could not run python with the --version flag, please check the path to python")
		return "", false
	}

	ui.Warning().Println("Please install Python 3.10 or higher and add it to the PATH")
	ui.Warning().Println("You can download Python here: https://www.python.org/downloads/")
	ui.Warning().Println("If you have already installed Python, please add it to the PATH")

	return "", false
}
