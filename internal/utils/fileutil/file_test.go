package fileutil

import (
	"github.com/easy-model-fusion/emf-cli/sdk"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/cobra"
	"os"
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
	destinationFile := PathJoin(tmpDir, "destination.py")

	// Call the function to copy the embedded file
	err = CopyEmbeddedFile(sdk.EmbeddedFiles, sourceFile, destinationFile)
	if err != nil {
		t.Fatalf("CopyEmbeddedFile failed: %v", err)
	}

	// Verify that the destination file now exists
	_, err = os.Stat(destinationFile)
	if err != nil {
		t.Fatalf("DestinationDir file not created: %v", err)
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
	exists, err := IsExistingPath(PathJoin(dir, "shouldRaiseError"))
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

func TestMoveFiles(t *testing.T) {
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
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Create a temporary directory for the test
	dir2, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir2) // clean up

	// Move the file to the second directory
	err = MoveFiles(dir, dir2)
	if err != nil {
		t.Fatal(err)
	}

	// Check path existence after removal
	exists, err := IsExistingPath(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	test.AssertEqual(t, false, exists, "File should have been moved.")
}

func TestMoveFiles_RenameError(t *testing.T) {
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
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Move the file to the second directory
	err = MoveFiles(dir, PathJoin(dir, "notFound d569sdf%/**"))
	if err == nil {
		t.Fatal("Expected error")
	}
}

// TestPathUniformize_Success tests the PathUniformize to return uniformized paths.
func TestPathUniformize_Success(t *testing.T) {
	// Init
	items := []struct {
		input    string
		expected string
	}{
		{"C:/path/to/file", "C:/path/to/file"},
		{"C:/path/to/../file", "C:/path/file"},
		{"C:/path/to/dir/../file", "C:/path/to/file"},
		{"C:/path/with/double/slashes", "C:/path/with/double/slashes"},
		{"C:/path/with/dots/..", "C:/path/with"},
		{"C:/path/with/dots/../..", "C:/path"},
		{"C:/path/with/dots/.", "C:/path/with/dots"},
		{"C:/path/with/dots/./.", "C:/path/with/dots"},
		{"C:/path/with/dots/././..", "C:/path/with"},
		{"C:/path/with/dots/././../file", "C:/path/with/file"},
	}

	for _, item := range items {
		// Execute
		result := PathUniformize(item.input)

		// Assert
		test.AssertEqual(t, result, item.expected)
	}
}

// TestPathJoin_Success tests the PathJoin to return uniformized paths from joined elements.
func TestPathJoin_Success(t *testing.T) {
	// Init
	items := []struct {
		directory string
		fileName  string
		expected  string
	}{
		{"C:/path/to", "file", "C:/path/to/file"},
		{"C:/path/to/..", "file", "C:/path/file"},
		{"C:/path/to/dir/..", "file", "C:/path/to/file"},
		{"C:/path/with/double", "slashes", "C:/path/with/double/slashes"},
		{"C:/path/with/dots", "..", "C:/path/with"},
		{"C:/path/with/dots", "../..", "C:/path"},
		{"C:/path/with/dots", ".", "C:/path/with/dots"},
		{"C:/path/with/dots", "./.", "C:/path/with/dots"},
		{"C:/path/with/dots", "././..", "C:/path/with"},
		{"C:/path/with/dots", "././../file", "C:/path/with/file"},
	}

	for _, item := range items {
		// Execute
		result := PathJoin(item.directory, item.fileName)

		// Assert
		test.AssertEqual(t, result, item.expected)
	}
}
