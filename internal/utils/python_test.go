package utils

import (
	"github.com/easy-model-fusion/client/test"
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
