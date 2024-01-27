package utils

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

// TestArrayStringContainsItem tests the ArrayStringContainsItem function.
func TestArrayStringContainsItem(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case: item is present in the array
	result := ArrayStringContainsItem(arr, "banana")
	test.AssertEqual(t, result, true, "Expected true")

	// Test case: item is not present in the array
	result = ArrayStringContainsItem(arr, "grape")
	test.AssertEqual(t, result, false, "Expected false")
}

// TestArrayStringAsArguments tests the ArrayStringAsArguments function.
func TestArrayStringAsArguments(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case
	expected := "[apple|banana|orange]"
	result := ArrayStringAsArguments(arr)
	test.AssertEqual(t, result, expected, "Generated string does not match the expected format")
}

// TestArrayFromString tests the ArrayFromString function.
func TestArrayFromString(t *testing.T) {
	// Sample string array
	input := "apple banana orange"

	// Test case
	expected := []string{"apple", "banana", "orange"}
	result := ArrayFromString(input)

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of arrays do not match")

	// Check each element
	for i := range expected {
		test.AssertEqual(t, result[i], expected[i], "Array element mismatch at index", string(rune(i)))
	}
}

// TestMapFromArrayString tests the MapFromArrayString function.
func TestMapFromArrayString(t *testing.T) {
	// Sample string array
	items := []string{"apple", "banana", "orange"}

	// Test case
	expected := map[string]struct{}{"apple": {}, "banana": {}, "orange": {}}
	result := MapFromArrayString(items)

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}
