package utils

import (
	"github.com/easy-model-fusion/client/sdk"
	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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

func TestCopyEmbeddedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	sourceFile := "main.py"
	destinationFile := filepath.Join(tmpDir, "destination.py")

	// Call the function to copy the embedded file
	err = CopyEmbeddedFile(sdk.EmbeddedFiles, sourceFile, destinationFile)
	if err != nil {
		t.Fatalf("CopyEmbeddedFile failed: %v", err)
	}

	// Verify that the destination file now exists
	_, err = os.Stat(destinationFile)
	if err != nil {
		t.Fatalf("Destination file not created: %v", err)
	}

	// Read the content of the destination file
	_, err = os.ReadFile(destinationFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

}
