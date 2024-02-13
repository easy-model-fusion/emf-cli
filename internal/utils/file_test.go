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

func TestCloseFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "emf-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Close the file
	CloseFile(tmpFile)

	// Verify that the file is closed
	err = tmpFile.Close()
	if err == nil {
		t.Fatal("File should be closed")
	}
}

// TestIsExistingPath_True tests the IsExistingPath function with an existing path.
func TestIsExistingPath_True(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Check path existence
	exists, err := IsExistingPath(dir)
	if err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, true, exists, "Path should be found as existing.")
}

// TestIsExistingPath_False tests the IsExistingPath function with a non-existing path.
func TestIsExistingPath_False(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Check path existence
	exists, err := IsExistingPath(filepath.Join(dir, "shouldRaiseError"))
	if err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, false, exists, "Path should be found as not existing.")
}

// TestIsDirectoryEmpty_True tests the IsDirectoryEmpty function with an empty directory.
func TestIsDirectoryEmpty_True(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Check directory emptiness
	exists, err := IsDirectoryEmpty(dir)
	if err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, true, exists, "Path should be found as empty.")
}

// TestIsDirectoryEmpty_False tests the IsDirectoryEmpty function with a non-empty directory.
func TestIsDirectoryEmpty_False(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Create a temporary file in dir for the test
	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Check directory emptiness
	exists, err := IsDirectoryEmpty(dir)
	if err != nil {
		t.Fatal(err)
	}

	test.AssertEqual(t, false, exists, "Path should be found as not empty.")
}

// TestDeleteDirectoryIfEmpty_Empty tests the DeleteDirectoryIfEmpty function with an empty directory.
func TestDeleteDirectoryIfEmpty_Empty(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Check directory emptiness
	err = DeleteDirectoryIfEmpty(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Check path existence after removal
	exists, err := IsExistingPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "Directory should have been removed.")
}

// TestDeleteDirectoryIfEmpty_NonEmpty tests the DeleteDirectoryIfEmpty function with a non-empty directory.
func TestDeleteDirectoryIfEmpty_NonEmpty(t *testing.T) {
	// Create a temporary directory for the test
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	// Create a temporary file in dir for the test
	file, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Check directory emptiness
	err = DeleteDirectoryIfEmpty(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Check path existence after removal
	exists, err := IsExistingPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, true, exists, "Directory should not have been removed.")
}
