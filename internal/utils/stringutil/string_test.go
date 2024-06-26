package stringutil

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestSplit_Success tests the Split function to split correctly on space characters.
func TestSplit_Success(t *testing.T) {
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

// TestSplitPath_Success tests the SplitPath function to return the correct elements composing the path.
func TestSplitPath_Success(t *testing.T) {

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

// TestParseOptions_KeyValueClassic tests parsing key-value pairs in classic format.
// It initializes an input string with a single key-value pair in classic format.
// It then executes the ParseOptions function with the input string and asserts that the result contains the specified option.
func TestParseOptions_KeyValueClassic(t *testing.T) {
	// Init
	option1 := "option=value"
	input := option1

	// Execute
	result := ParseOptions(input)

	// Assert
	test.AssertEqual(t, SliceContainsItem(result, option1), true, "should work with "+option1)
}

// TestParseOptions_KeyValueStrings tests parsing key-value pairs with values containing spaces.
// It initializes an input string with multiple key-value pairs, some with values containing spaces and enclosed in double or single quotes.
// It then executes the ParseOptions function with the input string and asserts that the result contains each specified option.
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

// TestParseOptions_ValueClassic tests parsing options without explicit values.
// It initializes an input string with multiple options specified without explicit values.
// It then executes the ParseOptions function with the input string and asserts that the result contains each specified option.
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

// TestParseOptions_ValueStrings tests parsing options with values containing spaces.
// It initializes an input string with multiple options specified with values containing spaces, some enclosed in double or single quotes.
// It then executes the ParseOptions function with the input string and asserts that the result contains each specified option.
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

func TestOptionsMapToSlice(t *testing.T) {
	// Init
	optionsMap := map[string]string{
		"key1": "value",
		"key2": "True",
		"key3": "module.value",
		"key4": "'text'",
	}

	// Execute
	optionsSlice := OptionsMapToSlice(optionsMap)

	// Assert
	test.AssertEqual(t, len(optionsSlice), len(optionsMap))
	test.AssertEqual(t, SliceContainsItem(optionsSlice, "key1=value"), true)
	test.AssertEqual(t, SliceContainsItem(optionsSlice, "key2=True"), true)
	test.AssertEqual(t, SliceContainsItem(optionsSlice, "key3=module.value"), true)
	test.AssertEqual(t, SliceContainsItem(optionsSlice, "key4='text'"), true)
}

func TestPathRemoveSpecialCharacter(t *testing.T) {
	// Init
	testPath := "models\\FredZhang7\\anime-anything-promptgen-v2\\model"

	// Execute
	updatedPath := PathRemoveSpecialCharacter(testPath)
	expectedPath := "models/FredZhang7/anime-anything-promptgen-v2/model"

	test.AssertEqual(t, updatedPath, expectedPath)
}
