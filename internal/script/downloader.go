package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const DownloaderName = "downloader.py"

// DownloaderModel represents a model obtained from the download downloader.
type DownloaderModel struct {
	Path      string              `json:"path"`
	Module    string              `json:"module"`
	Class     string              `json:"class"`
	Tokenizer DownloaderTokenizer `json:"tokenizer"`
}

// DownloaderTokenizer represents a tokenizer obtained the download downloader.
type DownloaderTokenizer struct {
	Path  string `json:"path"`
	Class string `json:"class"`
}

// IsDownloaderScriptModelEmpty checks if a DownloaderScriptModel is empty.
func IsDownloaderScriptModelEmpty(dsm DownloaderModel) bool {
	return dsm.Path == "" && dsm.Module == "" && dsm.Class == ""
}

// IsDownloaderScriptTokenizer checks if a DownloaderScriptTokenizer is empty.
func IsDownloaderScriptTokenizer(dst DownloaderTokenizer) bool {
	return dst.Path == "" && dst.Class == ""
}

// Download runs the download script for a specific model
func Download(pythonPath, downloadPath, modelName, moduleName, className string, overwrite bool) (DownloaderModel, error, int) {

	// Create command
	// TODO : pass arguments
	var cmd *exec.Cmd
	if overwrite {
		cmd = exec.Command(pythonPath, DownloaderName, "--emf-client", downloadPath, modelName, moduleName, "--model-class", className, "--skip", "model", "--overwrite")
	} else {
		cmd = exec.Command(pythonPath, DownloaderName, "--emf-client", downloadPath, modelName, moduleName, "--model-class", className, "--skip", "model")
	}

	// Bind stderr to a buffer
	var errBuf strings.Builder
	cmd.Stderr = &errBuf

	// Run command
	output, err := cmd.Output()

	// Prepare return values
	exitCode := 0
	var result DownloaderModel

	// Download was successful
	if err == nil {
		err = json.Unmarshal(output, &result)
		if err != nil {
			return result, err, exitCode
		}
		return result, nil, exitCode
	}

	// If there was an error running the command, check if it's a command execution error
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}

	// Log the errors back
	errBufStr := errBuf.String()
	if errBufStr != "" {
		return result, fmt.Errorf("%s", errBufStr), exitCode
	}

	return result, err, exitCode
}
