package utils

import (
	"os/exec"
)

// CheckForExecutable checks if the given executable is available in the PATH.
func CheckForExecutable(name string) bool {
	_, err := exec.LookPath("python")
	return err == nil
}
