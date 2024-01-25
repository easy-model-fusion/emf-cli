package utils

import (
	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/cobra"
	"testing"
)

func TestIsFileNameValid(t *testing.T) {
	// Test with valid file name
	test.AssertEqual(t, IsFileNameValid("test"), true, "Expected true")

	// Test with invalid file name
	test.AssertEqual(t, IsFileNameValid("test/test"), false, "Expected false")
}

func TestValidFileName(t *testing.T) {
	// First arg, not optional
	validator := ValidFileName(1, false)

	// Test with no args
	err := validator(&cobra.Command{}, []string{})
	test.AssertNotEqual(t, err, nil, "Expected error")

	// Test with one arg
	err = validator(&cobra.Command{}, []string{"test"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Test with invalid arg
	err = validator(&cobra.Command{}, []string{"test/test"})
	test.AssertNotEqual(t, err, nil, "Expected error")

	// First arg, optional
	validator = ValidFileName(1, true)

	// Test with no args
	err = validator(&cobra.Command{}, []string{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Test with one arg
	err = validator(&cobra.Command{}, []string{"test"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
