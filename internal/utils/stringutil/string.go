package stringutil

import (
	"path/filepath"
	"regexp"
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

// ParseOptions parses a string containing options in various formats
// Returns a slice of strings where each string represents an option.
func ParseOptions(input string) []string {
	var result []string

	// Regular expression to match key-value pairs
	// Also catches value as single or double quotes strings, with or without spaces
	re := regexp.MustCompile(`(\S+)=((?:"[^"]+")|(?:'[^']+')|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		result = append(result, match[0])
		// Remove processed option from the input string
		input = strings.Replace(input, match[0], "", 1)
	}

	// Find all inputs as single or double quotes strings, with or without spaces
	reQuotes := regexp.MustCompile(`(?:"[^"]+")|(?:'[^']+')`)
	matchesQuotes := reQuotes.FindAllStringSubmatch(input, -1)
	for _, match := range matchesQuotes {
		result = append(result, match[0])
		// Remove processed option from the input string
		input = strings.Replace(input, match[0], "", 1)
	}

	// Split the input by spaces and add any parts that are not key-value pairs
	parts := strings.Fields(input)
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			result = append(result, part)
		}
	}

	return result
}
