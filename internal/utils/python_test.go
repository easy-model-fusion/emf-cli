package utils

import (
	"github.com/easy-model-fusion/client/test"
	"os"
	"path/filepath"
	"testing"
)

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

func TestCheckForPython(t *testing.T) {
	checkFalse := true

	if _, ok := CheckForExecutable("python"); ok {
		path, ok := CheckForPython()
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
		checkFalse = false
	}
	if _, ok := CheckForExecutable("python3"); ok {
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

func TestCheckPythonVersion(t *testing.T) {
	if _, ok := CheckForExecutable("python"); ok {
		path, ok := CheckPythonVersion("python")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}
	if _, ok := CheckForExecutable("python3"); ok {
		path, ok := CheckPythonVersion("python3")
		test.AssertEqual(t, ok, true)
		test.AssertNotEqual(t, path, "")
	}

	path, ok := CheckPythonVersion("anexecutablethatcouldnotexists-yeahhh")
	test.AssertEqual(t, ok, false)
	test.AssertEqual(t, path, "")
}

func TestCreateVirtualEnv(t *testing.T) {
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

func TestFindVEnvPipExecutable(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	defer os.RemoveAll(dname)

	// Execute
	pipPath, err := FindVEnvExecutable(venvPath, "pip")

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, pipPath, "")
}

func TestFindVEnvPipExecutable_Fail(t *testing.T) {
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

func TestExecuteScript_MissingVenv(t *testing.T) {
	// Execute
	_, err, _ := ExecuteScript(".venv", "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}

func TestExecuteScript_MissingScript(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	defer os.RemoveAll(dname)

	// Execute
	_, err, _ := ExecuteScript(venvPath, "script.py", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}

func TestExecuteScript_EmptyResponse(t *testing.T) {
	// Init
	dname, venvPath := CreateVenv(t)
	file, err := os.CreateTemp("", "*.y")
	if err != nil {
		t.Fatal(err)
	}
	defer CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)
}

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
	defer CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, len(output), 0)
	test.AssertNotEqual(t, exitCode, 0)
}

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
	defer CloseFile(file)
	defer os.RemoveAll(dname)

	// Execute
	output, err, exitCode := ExecuteScript(venvPath, file.Name(), []string{})

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, len(output), 0)
	test.AssertEqual(t, exitCode, 0)
}
