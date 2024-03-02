package executil

import (
	"os/exec"
)

// CheckForExecutable checks if the given executable is available in the PATH, if so, it returns the path to the executable.
func CheckForExecutable(name string) (string, bool) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", false
	}
	return path, true
}
