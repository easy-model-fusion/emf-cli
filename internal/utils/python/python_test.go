package python

import (
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/executil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// CreateVenv creates a temporary directory and a virtual environment inside it.
// It returns the path to the temporary directory and the path to the virtual environment.
// If Python is not found, it skips the test.
func CreateVenv(t *testing.T) (string, string) {
	path, ok := NewPython().CheckForPython()
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
	err = NewPython().CreateVirtualEnv(path, venvPath)
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
		path, ok := NewPython().CheckForPython()
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
		checkFalse = false
	}
	if _, ok := executil.CheckForExecutable("python3"); ok {
		path, ok := NewPython().CheckForPython()
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
		checkFalse = false
	}

	if checkFalse {
		path, ok := NewPython().CheckForPython()
		test.AssertEqual(t, ok, false)
		test.AssertNotEqual(t, path, "")
	}
}

// TestCheckPythonVersion tests the CheckPythonVersion function.
func TestCheckPythonVersion(t *testing.T) {
	if _, ok := executil.CheckForExecutable("python"); ok {
		path, ok := NewPython().CheckPythonVersion("python")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}
	if _, ok := executil.CheckForExecutable("python3"); ok {
		path, ok := NewPython().CheckPythonVersion("python3")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}

	path, ok := NewPython().CheckPythonVersion("anexecutablethatcouldnotexists-yeahhh")
	test.AssertEqual(t, ok, false)
	test.AssertEqual(t, path, "")
}

// TestCreateVirtualEnv_Success tests the CreateVirtualEnv function to successfully create a venv.
func TestCreateVirtualEnv_Success(t *testing.T) {
	path, ok := NewPython().CheckForPython()
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

	err = NewPython().CreateVirtualEnv(path, filepath.Join(dname, "venv"))
	test.AssertEqual(t, err, nil)
}

// TestFindVEnvExecutable tests the FindVEnvExecutable function with existing executable.
func TestFindVEnvExecutable_Success(t *testing.T) {
	// Init
	dname, err := os.MkdirTemp("", "emf-cli")
	test.AssertEqual(t, err, nil, "Error creating temporary directory")
	venvPath := filepath.Join(dname, "venv")

	// create "virtual environment"
	if runtime.GOOS == "windows" {
		err = os.MkdirAll(filepath.Join(venvPath, "Scripts"), os.ModePerm)
		test.AssertEqual(t, err, nil, "Error creating Scripts directory")
		_, err = os.Create(filepath.Join(venvPath, "Scripts", "pip.exe"))
		test.AssertEqual(t, err, nil, "Error creating pip.exe")
	} else {
		err = os.MkdirAll(filepath.Join(venvPath, "bin"), os.ModePerm)
		test.AssertEqual(t, err, nil, "Error creating bin directory")
		_, err = os.Create(filepath.Join(venvPath, "bin", "pip"))
		test.AssertEqual(t, err, nil, "Error creating pip")
	}

	// Execute
	pipPath, err := NewPython().FindVEnvExecutable(venvPath, "pip")

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

	pipPath, err := NewPython().FindVEnvExecutable(filepath.Join(dname, "venv"), "pip")
	test.AssertNotEqual(t, err, nil, "Should return an error")
	test.AssertEqual(t, pipPath, "")
}

// TestExecuteScript_MissingVenv tests the ExecuteScript function when the virtual environment is missing.
func TestExecuteScript_MissingVenv(t *testing.T) {
	// Execute
	_, err, _ := NewPython().ExecuteScript(".venv", "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}

func TestExecuteScript(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	defer os.RemoveAll(dname)

	// ***************************
	// *** Missing script test ***
	// ***************************

	// Execute
	_, err, _ := NewPython().ExecuteScript(venvPath, "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)

	// ***************************
	// *** Empty response test ***
	// ***************************

	// Init
	file, err := os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}

	// Execute
	output, err, exitCode := NewPython().ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)

	// Cleanup
	fileutil.CloseFile(file)
	err = os.Remove(file.Name())
	if err != nil {
		t.Error(err)
	}

	// ***************************
	// *** Error response test ***
	// ***************************

	// Init
	file, err = os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	scriptContent := []byte(`print(1/0)`)
	err = os.WriteFile(file.Name(), scriptContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer fileutil.CloseFile(file)

	// Execute
	output, err, exitCode = NewPython().ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertNotEqual(t, exitCode, 0)

	// Cleanup
	fileutil.CloseFile(file)
	err = os.Remove(file.Name())
	if err != nil {
		t.Error(err)
	}

	// *****************************
	// *** Success response test ***
	// *****************************

	// Init
	file, err = os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	scriptContent = []byte(`print('Hello, world!')`)
	err = os.WriteFile(file.Name(), scriptContent, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer fileutil.CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode = NewPython().ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)

	// Cleanup
	fileutil.CloseFile(file)
	err = os.Remove(file.Name())
	if err != nil {
		t.Error(err)
	}
}

// TestCheckAskForPython_Success tests the CheckAskForPython function when Python is installed.
func TestCheckAskForPython_Success(t *testing.T) {
	// check python
	a, ok := NewPython().CheckForPython()
	if !ok {
		return
	}

	b, ok := NewPython().CheckAskForPython(ui.NewPTermUI())
	test.AssertEqual(t, ok, true, "Should return true if python is installed")
	test.AssertEqual(t, a, b, "Should return the same value as CheckForPython")
}
