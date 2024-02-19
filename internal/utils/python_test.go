package utils

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"path/filepath"
	"testing"
)

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
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	pipPath, err := FindVEnvExecutable(filepath.Join(dname, "venv"), "pip")
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
