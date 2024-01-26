package utils

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

func TestArrayStringContainsItem(t *testing.T) {
	arr := []string{"apple", "banana", "orange"}

	// Test case: item is present in the array
	result := ArrayStringContainsItem(arr, "banana")
	test.AssertEqual(t, result, true, "Expected true")

	// Test case: item is not present in the array
	result = ArrayStringContainsItem(arr, "grape")
	test.AssertEqual(t, result, false, "Expected false")
}

func TestArrayStringAsArguments(t *testing.T) {
	arr := []string{"apple", "banana", "orange"}

	// Test case
	expected := "[apple|banana|orange]"
	result := ArrayStringAsArguments(arr)
	test.AssertEqual(t, result, expected, "Generated string does not match the expected format")
}

func TestArrayFromString(t *testing.T) {
	input := "apple banana orange"

	// Test case
	expected := []string{"apple", "banana", "orange"}
	result := ArrayFromString(input)

	if len(result) != len(expected) {
		test.AssertEqual(t, len(result), len(expected), "Lengths of arrays do not match")
	}

	for i := range expected {
		test.AssertEqual(t, result[i], expected[i], "Array element mismatch at index", string(rune(i)))
	}
}

func TestMapFromArrayString(t *testing.T) {
	items := []string{"apple", "banana", "orange"}

	// Test case
	expected := map[string]struct{}{"apple": {}, "banana": {}, "orange": {}}
	result := MapFromArrayString(items)

	if len(result) != len(expected) {
		test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")
	}

	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}
