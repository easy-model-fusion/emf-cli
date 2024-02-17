package utils

import (
	"github.com/easy-model-fusion/client/test"
	"testing"
)

// TestSliceContainsItem tests the SliceContainsItem function.
func TestSliceContainsItem(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case: item is present in the array
	result := SliceContainsItem(arr, "banana")
	test.AssertEqual(t, result, true, "Expected true")

	// Test case: item is not present in the array
	result = SliceContainsItem(arr, "grape")
	test.AssertEqual(t, result, false, "Expected false")
}

// TestSliceToArgsFormat tests the SliceToArgsFormat function.
func TestSliceToArgsFormat(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case
	expected := "[apple|banana|orange]"
	result := SliceToArgsFormat(arr)
	test.AssertEqual(t, result, expected, "Generated string does not match the expected format")
}

// TestSliceToMap tests the SliceToMap function.
func TestSliceToMap(t *testing.T) {
	// Sample string array
	items := []string{"apple", "banana", "orange"}

	// Test case
	expected := map[string]struct{}{"apple": {}, "banana": {}, "orange": {}}
	result := SliceToMap(items)

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of maps do not match")

	// Check each key
	for key := range expected {
		_, exists := result[key]
		test.AssertEqual(t, exists, true, "Key not found in the result map:", key)
	}
}

// TestSliceRemoveDuplicates tests the SliceRemoveDuplicates function.
func TestSliceRemoveDuplicates(t *testing.T) {

	// Init
	expected := []string{"input", "to", "test"}
	input := append(expected, expected...)

	// Test
	result := SliceRemoveDuplicates(input)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths do not match")
}

// TestSliceDifference tests the SliceDifference function to return the correct difference.
func TestSliceDifference(t *testing.T) {
	// Init
	elements := []string{"item0", "item1", "item2", "item3", "item4"}
	index := 2
	sub := elements[:index]
	expected := elements[index:]

	// Execute
	difference := SliceDifference(elements, sub)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

// TestSliceRemoveValue_Present tests the SliceRemoveValue function with a value present in the slice.
func TestSliceRemoveValue_Present(t *testing.T) {
	// Init
	value := "value"
	slice := []string{value}

	// Execute
	result := SliceRemoveValue(slice, value)

	// Assert
	test.AssertEqual(t, len(result), len(slice)-1)
}

// TestSliceRemoveValue_NotPresent tests the SliceRemoveValue function with a value not present in the slice.
func TestSliceRemoveValue_NotPresent(t *testing.T) {
	// Init
	nonExistentValue := "nonExistentValue"
	value := "value"
	slice := []string{value}

	// Execute
	result := SliceRemoveValue(slice, nonExistentValue)

	// Assert
	test.AssertEqual(t, len(result), len(slice))
}
