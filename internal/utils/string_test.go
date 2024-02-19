package utils

import (
	"github.com/easy-model-fusion/client/test"
	"path/filepath"
	"testing"
)

// TestSplit tests the Split function.
func TestSplit(t *testing.T) {
	// Sample string array
	input := "apple banana orange"

	// Test case
	expected := []string{"apple", "banana", "orange"}
	result := Split(input)

	// Check if lengths match
	test.AssertEqual(t, len(result), len(expected), "Lengths of arrays do not match")

	// Check each element
	for i := range expected {
		test.AssertEqual(t, result[i], expected[i], "Array element mismatch at index", string(rune(i)))
	}
}

// TestSplitPath tests the SplitPath function.
func TestSplitPath(t *testing.T) {

	// Init
	expected := []string{"input", "to", "test"}
	input := expected[0] + "/" + expected[1] + "/" + expected[2]

	// Test
	result := SplitPath(input)

	// Assert
	test.AssertEqual(t, len(expected), len(result), "Lengths do not match")
	for i := 0; i < len(expected); i++ {
		test.AssertEqual(t, expected[i], result[i], "Values do not match")
	}
}

// TestPathUniformize tests the PathUniformize to return uniformized paths.
func TestPathUniformize(t *testing.T) {
	// Init
	items := []struct {
		input    string
		expected string
	}{
		{"C:\\path\\to\\file", "C:/path/to/file"},
		{"C:\\path\\to\\..\\file", "C:/path/file"},
		{"C:\\path\\to\\dir\\..\\file", "C:/path/to/file"},
		{"C:\\path\\with\\double\\\\slashes", "C:/path/with/double/slashes"},
		{"C:\\path\\with\\dots\\..", "C:/path/with"},
		{"C:\\path\\with\\dots\\..\\..", "C:/path"},
		{"C:\\path\\with\\dots\\.", "C:/path/with/dots"},
		{"C:\\path\\with\\dots\\.\\.", "C:/path/with/dots"},
		{"C:\\path\\with\\dots\\.\\.\\..", "C:/path/with"},
		{"C:\\path\\with\\dots\\.\\.\\..\\file", "C:/path/with/file"},
	}

	for _, item := range items {
		// Execute
		result := PathUniformize(item.input)

		// Assert
		test.AssertEqual(t, result, filepath.Clean(item.expected))
	}
}

func TestParseOptions_KeyValueClassic(t *testing.T) {
	// Init
	option1 := "option=value"
	input := option1

	// Execute
	result := ParseOptions(input)

	// Assert
	test.AssertEqual(t, SliceContainsItem(result, option1), true, "should work with "+option1)
}

func TestParseOptions_KeyValueStrings(t *testing.T) {
	// Init
	option1 := "option1=\"value 1\""
	option2 := "option2=\"value 2 with spaces\""
	option3 := "option3='value3'"
	option4 := "option4='value 4 with spaces'"
	input := option1 + " " + option2 + " " + option3 + " " + option4

	// Execute
	result := ParseOptions(input)

	// Assert
	test.AssertEqual(t, SliceContainsItem(result, option1), true, "should work with "+option1)
	test.AssertEqual(t, SliceContainsItem(result, option2), true, "should work with "+option2)
	test.AssertEqual(t, SliceContainsItem(result, option3), true, "should work with "+option3)
	test.AssertEqual(t, SliceContainsItem(result, option4), true, "should work with "+option4)
}

func TestParseOptions_ValueClassic(t *testing.T) {
	// Init
	option1 := "option5"
	option2 := "--option6"
	input := option1 + " " + option2

	// Execute
	result := ParseOptions(input)

	// Assert
	test.AssertEqual(t, SliceContainsItem(result, option1), true, "should work with "+option1)
	test.AssertEqual(t, SliceContainsItem(result, option2), true, "should work with "+option2)
}

func TestParseOptions_ValueStrings(t *testing.T) {
	// Init
	option1 := "\"value 1\""
	option2 := "\"value 2 with spaces\""
	option3 := "'value3'"
	option4 := "'value 4 with spaces'"
	input := option1 + " " + option2 + " " + option3 + " " + option4

	// Execute
	result := ParseOptions(input)

	// Assert
	test.AssertEqual(t, SliceContainsItem(result, option1), true, "should work with "+option1)
	test.AssertEqual(t, SliceContainsItem(result, option2), true, "should work with "+option2)
	test.AssertEqual(t, SliceContainsItem(result, option3), true, "should work with "+option3)
	test.AssertEqual(t, SliceContainsItem(result, option4), true, "should work with "+option4)
}
