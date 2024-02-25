package stringutil

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestSliceContainsItem_True tests the SliceContainsItem function to return true when item contained.
func TestSliceContainsItem_True(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case: item is present in the array
	result := SliceContainsItem(arr, "banana")
	test.AssertEqual(t, result, true, "Expected true")
}

// TestSliceContainsItem_False tests the SliceContainsItem function to return false when item not contained.
func TestSliceContainsItem_False(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case: item is not present in the array
	result := SliceContainsItem(arr, "grape")
	test.AssertEqual(t, result, false, "Expected false")
}

// TestSliceToArgsFormat_Success tests the SliceToArgsFormat function to return the slice formatted as arguments.
func TestSliceToArgsFormat_Success(t *testing.T) {
	// Sample string array
	arr := []string{"apple", "banana", "orange"}

	// Test case
	expected := "[apple|banana|orange]"
	result := SliceToArgsFormat(arr)
	test.AssertEqual(t, result, expected, "Generated string does not match the expected format")
}

// TestSliceToMap_Success tests the SliceToMap function to return a map from a slice of strings.
func TestSliceToMap_Success(t *testing.T) {
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

// TestSliceRemoveDuplicates_WithDuplicates tests the SliceRemoveDuplicates function to succeed with duplicates.
func TestSliceRemoveDuplicates_WithDuplicates(t *testing.T) {

	// Init
	expected := []string{"input", "to", "test"}
	input := append(expected, expected...)

	// Test
	result := SliceRemoveDuplicates(input)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths do not match")
}

// TestSliceRemoveDuplicates_WithoutDuplicates tests the SliceRemoveDuplicates function to succeed with duplicates.
func TestSliceRemoveDuplicates_WithoutDuplicates(t *testing.T) {

	// Init
	expected := []string{"input", "to", "test"}
	input := expected

	// Test
	result := SliceRemoveDuplicates(input)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths do not match")
}

// TestSliceDifference_WithDifference tests the SliceDifference function to succeed with differences.
func TestSliceDifference_WithDifference(t *testing.T) {
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

// TestSliceDifference_WithFirstEmpty tests the SliceDifference function to succeed with empty first slice.
func TestSliceDifference_WithFirstEmpty(t *testing.T) {
	// Init
	elements := []string{"item0", "item1", "item2", "item3", "item4"}
	var expected []string

	// Execute
	difference := SliceDifference([]string{}, elements)

	// Assert
	test.AssertEqual(t, len(expected), len(difference), "Lengths should be equal.")
}

// TestSliceDifference_WithSecondEmpty tests the SliceDifference function to succeed with empty second slice.
func TestSliceDifference_WithSecondEmpty(t *testing.T) {
	// Init
	elements := []string{"item0", "item1", "item2", "item3", "item4"}

	// Execute
	difference := SliceDifference(elements, []string{})

	// Assert
	test.AssertEqual(t, len(elements), len(difference), "Lengths should be equal.")
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
