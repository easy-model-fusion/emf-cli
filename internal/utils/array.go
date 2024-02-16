package utils

import (
	"path/filepath"
	"strings"
)

// ArrayStringContainsItem checks if an item is present in a slice of strings.
func ArrayStringContainsItem(arr []string, item string) bool {
	for _, element := range arr {
		if element == item {
			return true
		}
	}
	return false
}

// ArrayStringAsArguments converts a slice of strings into a string with elements separated by '|'.
func ArrayStringAsArguments(arr []string) string {
	return "[" + strings.Join(arr, "|") + "]"
}

// ArrayFromString splits a string into a slice of strings based on space characters.
func ArrayFromString(input string) []string {
	// Split the input string based on the space character
	result := strings.Split(input, " ")
	return result
}

// MapFromArrayString creates a map from a slice of strings for faster lookup.
func MapFromArrayString(items []string) map[string]struct{} {
	stringMap := make(map[string]struct{})
	for _, item := range items {
		stringMap[item] = struct{}{}
	}
	return stringMap
}

// ArrayFromPath splits a filepath into its individual elements
func ArrayFromPath(path string) []string {
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

// StringRemoveDuplicates returns a slice in which every element only appears once
func StringRemoveDuplicates(items []string) []string {

	// Prepare variables
	itemsMap := make(map[string]bool)
	var result []string

	// Looking for duplicates
	for _, item := range items {
		// Item not contained yet
		if !itemsMap[item] {
			// Adding it to the result and indicating it as seen inside the map
			itemsMap[item] = true
			result = append(result, item)
		}
	}
	return result
}

// StringDifference returns the elements in `parentSlice` that are not present in `subSlice`
func StringDifference(parentSlice, subSlice []string) []string {
	var difference []string
	for _, item := range parentSlice {
		if !ArrayStringContainsItem(subSlice, item) {
			difference = append(difference, item)
		}
	}
	return difference
}

func UniformizePath(path string) string {
	// Replace backslashes with forward slashes
	path = strings.ReplaceAll(path, "\\", "/")

	// Resolve dots and double slashes
	path = filepath.Clean(path)

	return path
}
