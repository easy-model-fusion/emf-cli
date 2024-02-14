package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// CheckPythonVersion checks if python is found in the PATH and runs it with the --
// version flag to check if it works, and returns path to python executable and true if so.
// If python is not found, the function returns false.
func CheckPythonVersion(name string) (string, bool) {
	path, ok := CheckForExecutable(name)
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
func CheckForPython() (string, bool) {
	path, ok := CheckPythonVersion("python")
	if ok {
		return path, true
	}
	return CheckPythonVersion("python3")
}

// CreateVirtualEnv creates a virtual environment in the given path
func CreateVirtualEnv(pythonPath, path string) error {
	cmd := exec.Command(pythonPath, "-m", "venv", path)
	return cmd.Run()
}

// FindVEnvExecutable searches for the requested executable within a virtual environment.
func FindVEnvExecutable(venvPath string, executableName string) (string, error) {
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
func InstallDependencies(pipPath, path string) error {
	cmd := exec.Command(pipPath, "install", "-r", path)

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

// DownloadModel runs the download script for a specific model
func DownloadModel(pythonPath, downloadPath, modelName, moduleName, className string, overwrite bool) (error, int) {

	// Create command
	var cmd *exec.Cmd
	if overwrite {
		cmd = exec.Command(pythonPath, "download.py", downloadPath, modelName, moduleName, className, "--overwrite")
	} else {
		cmd = exec.Command(pythonPath, "download.py", downloadPath, modelName, moduleName, className)
	}

	// bind stderr to a buffer
	var errBuf strings.Builder
	cmd.Stderr = &errBuf

	// cmd exit code
	exitCode := 0

	err := cmd.Run()
	if err == nil {
		return nil, exitCode
	}

	// If there was an error running the command, check if it's a command execution error
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}

	// Log the errors back
	errBufStr := errBuf.String()
	if errBufStr != "" {
		return fmt.Errorf("%s", errBufStr), exitCode
	}

	return err, exitCode
}
