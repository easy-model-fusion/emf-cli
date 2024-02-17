package utils

import (
	"path/filepath"
	"strings"
)

// Split splits a string into a slice of strings based on space characters.
func Split(input string) []string {
	return strings.Split(input, " ")
}

// SplitPath splits a filepath into its individual elements.
func SplitPath(path string) []string {
	var elements []string
	for {
		dir, file := filepath.Split(path)
		if len(dir) > 0 {
			elements = append([]string{file}, elements...)
			path = filepath.Clean(dir)
		} else {
			if len(file) > 0 {
				elements = append([]string{file}, elements...)
			}
			break
		}
	}
	return elements
}

// PathUniformize returns uniformized path regarding the device OS.
func PathUniformize(path string) string {
	// Replace backslashes with forward slashes
	path = strings.ReplaceAll(path, "\\", "/")

	// Resolve dots and double slashes
	path = filepath.Clean(path)

	return path
}
