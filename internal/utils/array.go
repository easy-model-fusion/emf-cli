package utils

import "strings"

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
