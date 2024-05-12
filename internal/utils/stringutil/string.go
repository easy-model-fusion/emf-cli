package stringutil

import (
	"fmt"
	"path/filepath"
	"regexp"
	"runtime"
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

// PathRemoveSpecialCharacter adds a slash before a special character
func PathRemoveSpecialCharacter(path string) string {
	//List of special characters that need to be escaped
	specialChars := map[string]string{
		"\\": "/",
	}

	// Replace special characters with their escaped form
	for oldChar, newChar := range specialChars {
		if strings.Contains(path, oldChar) {
			path = strings.ReplaceAll(path, oldChar, newChar)
		}
	}

	return path
}

// PathUniformize returns uniformized path regarding the device OS.
func PathUniformize(path string) string {
	// Replace backslashes with forward slashes
	// Resolve dots and double slashes

	// Handling platform-specific behavior since filepath.Clean behaves differently for each
	switch runtime.GOOS {
	case "windows":
		path = filepath.Clean(path)
		path = PathRemoveSpecialCharacter(path)
	case "linux", "darwin":
		path = PathRemoveSpecialCharacter(path)
		path = filepath.Clean(path)
	}

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

// OptionsMapToSlice transforms a map to a slice of strings under the key=value format.
func OptionsMapToSlice(optionsMap map[string]string) []string {
	var optionsSlice []string
	for key, value := range optionsMap {
		option := fmt.Sprintf("%s=%s", key, value)
		optionsSlice = append(optionsSlice, option)
	}
	return optionsSlice
}
