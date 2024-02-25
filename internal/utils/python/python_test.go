package python

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/executil"
	file2 "github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"path/filepath"
	"testing"
)

// CreateVenv creates a temporary directory and a virtual environment inside it.
// It returns the path to the temporary directory and the path to the virtual environment.
// If Python is not found, it skips the test.
func CreateVenv(t *testing.T) (string, string) {
	path, ok := CheckForPython()
	if !ok {
		t.SkipNow()
	}

	// create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	venvPath := filepath.Join(dname, "venv")
	err = CreateVirtualEnv(path, venvPath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	return dname, venvPath
}

// TestCheckForPython tests the CheckForPython function.
func TestCheckForPython(t *testing.T) {
	checkFalse := true

	if _, ok := executil.CheckForExecutable("python"); ok {
		path, ok := CheckForPython()
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
		checkFalse = false
	}
	if _, ok := executil.CheckForExecutable("python3"); ok {
		path, ok := CheckForPython()
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
		checkFalse = false
	}

	if checkFalse {
		path, ok := CheckForPython()
		test.AssertEqual(t, ok, false)
		test.AssertNotEqual(t, path, "")
	}
}

// TestCheckPythonVersion tests the CheckPythonVersion function.
func TestCheckPythonVersion(t *testing.T) {
	if _, ok := executil.CheckForExecutable("python"); ok {
		path, ok := CheckPythonVersion("python")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}
	if _, ok := executil.CheckForExecutable("python3"); ok {
		path, ok := CheckPythonVersion("python3")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}

	path, ok := CheckPythonVersion("anexecutablethatcouldnotexists-yeahhh")
	test.AssertEqual(t, ok, false)
	test.AssertEqual(t, path, "")
}

// TestCreateVirtualEnv_Success tests the CreateVirtualEnv function to successfully create a venv.
func TestCreateVirtualEnv_Success(t *testing.T) {
	path, ok := CheckForPython()
	if !ok {
		t.SkipNow()
	}

	// create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dname)

	err = CreateVirtualEnv(path, filepath.Join(dname, "venv"))
	test.AssertEqual(t, err, nil)
}

// TestFindVEnvExecutable tests the FindVEnvExecutable function with existing executable.
func TestFindVEnvExecutable_Success(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	defer os.RemoveAll(dname)

	// Execute
	pipPath, err := FindVEnvExecutable(venvPath, "pip")

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, pipPath, "")
}

// TestFindVEnvExecutable_Fail tests the FindVEnvExecutable function when it fails to find the executable.
func TestFindVEnvExecutable_Fail(t *testing.T) {
	// create temporary directory
	dname, err := os.MkdirTemp("", "emf-cli")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer os.RemoveAll(dname)

	pipPath, err := FindVEnvExecutable(filepath.Join(dname, "venv"), "pip")
	test.AssertNotEqual(t, err, nil, "Should return an error")
	test.AssertEqual(t, pipPath, "")
}

// TestExecuteScript_MissingVenv tests the ExecuteScript function when the virtual environment is missing.
func TestExecuteScript_MissingVenv(t *testing.T) {
	// Execute
	_, err, _ := ExecuteScript(".venv", "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}

// TestExecuteScript_MissingScript tests the ExecuteScript function when the script is missing.
func TestExecuteScript_MissingScript(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	defer os.RemoveAll(dname)

	// Execute
	_, err, _ := ExecuteScript(venvPath, "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}

// TestExecuteScript_EmptyResponse tests the ExecuteScript function when the script returns an empty response.
func TestExecuteScript_EmptyResponse(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	file, err := os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	defer file2.CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)
}

// TestExecuteScript_Error tests the ExecuteScript function when the script encounters an error.
func TestExecuteScript_Error(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	file, err := os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	scriptContent := []byte(`print(1/0)`)
	err = os.WriteFile(file.Name(), scriptContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file2.CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertNotEqual(t, exitCode, 0)
}

// TestExecuteScript_Success tests the ExecuteScript function when the script executes successfully.
func TestExecuteScript_Success(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	file, err := os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	scriptContent := []byte(`print('Hello, world!')`)
	err = os.WriteFile(file.Name(), scriptContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file2.CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)
}

// TestCheckAskForPython_Success tests the CheckAskForPython function when Python is installed.
func TestCheckAskForPython_Success(t *testing.T) {
	// check python
	a, ok := CheckForPython()
	if !ok {
		return
	}

	b, ok := CheckAskForPython()
	test.AssertEqual(t, ok, true, "Should return true if python is installed")
	test.AssertEqual(t, a, b, "Should return the same value as CheckForPython")
}
