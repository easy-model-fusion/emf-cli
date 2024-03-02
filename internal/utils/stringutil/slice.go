package stringutil

import (
	"strings"
)

// SliceContainsItem checks if an item is present in a slice of strings.
// True if the item is found, otherwise false.
func SliceContainsItem(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

// SliceToArgsFormat converts a slice of strings into a string with elements separated by '|'.
func SliceToArgsFormat(arr []string) string {
	return "[" + strings.Join(arr, "|") + "]"
}

// SliceToMap creates a map from a slice of strings for faster lookup.
func SliceToMap(slice []string) map[string]struct{} {
	stringMap := make(map[string]struct{})
	for _, item := range slice {
		stringMap[item] = struct{}{}
	}
	return stringMap
}

// SliceRemoveDuplicates returns a slice with all duplicates removed.
func SliceRemoveDuplicates(slice []string) []string {

	// Prepare variables
	itemsMap := make(map[string]bool)
	var result []string

	// Looking for duplicates
	for _, item := range slice {
		// Item not contained yet
		if !itemsMap[item] {
			// Adding it to the result and indicating it as seen inside the map
			itemsMap[item] = true
			result = append(result, item)
		}
	}
	return result
}

// SliceDifference returns the elements in `slice1` that are not present in `slice2`.
func SliceDifference(slice1, slice2 []string) []string {
	var difference []string
	for _, item := range slice1 {
		if !SliceContainsItem(slice2, item) {
			difference = append(difference, item)
		}
	}
	return difference
}

// SliceRemoveValue returns the slice with the specified value removed.
func SliceRemoveValue(slice []string, value string) []string {
	for i, item := range slice {
		if item == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
