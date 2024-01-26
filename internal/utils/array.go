package utils

import "strings"

func ArrayStringContainsItem(arr []string, item string) bool {
	for _, element := range arr {
		if element == item {
			return true
		}
	}
	return false
}

func ArrayStringAsArguments(arr []string) string {
	return "[" + strings.Join(arr, "|") + "]"
}

func ArrayFromString(input string) []string {
	// Split the input string based on the space character
	result := strings.Split(input, " ")
	return result
}
