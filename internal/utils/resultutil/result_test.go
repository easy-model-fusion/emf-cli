package resultutil

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestAddWarnings tests adding warning messages to result
func TestAddWarnings(t *testing.T) {
	// Init
	warning := "test"
	warning2 := "test2"
	var result ExecutionResult

	// add warnings
	result.AddWarnings([]string{warning, warning2})

	// Assertions
	test.AssertEqual(t, len(result.Warnings), 2)
	test.AssertEqual(t, result.Warnings[0], "test")
	test.AssertEqual(t, result.Warnings[1], "test2")
}

// TestAddInfos tests adding information messages to result
func TestAddInfos(t *testing.T) {
	// Init
	info := "test"
	info2 := "test2"
	var result ExecutionResult

	// add warnings
	result.AddInfos([]string{info, info2})

	// Assertions
	test.AssertEqual(t, len(result.Infos), 2)
	test.AssertEqual(t, result.Infos[0], "test")
	test.AssertEqual(t, result.Infos[1], "test2")
}

// TestSetError tests setting error to result
func TestSetError(t *testing.T) {
	// Init
	err := errors.New("test")
	var result ExecutionResult

	// add warnings
	test.AssertEqual(t, result.Error, nil)
	result.SetError(err)

	// Assertions
	test.AssertNotEqual(t, result.Error, nil)
	test.AssertEqual(t, result.Error.Error(), "test")
}
