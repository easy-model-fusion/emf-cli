package downloader

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

var pythonObject python.Python
var downloaderObject Downloader

func TestMain(m *testing.M) {
	pythonObject = &mock.MockPython{}
	downloaderObject = &scriptDownloader{}
	os.Exit(m.Run())
}

// TestExecute_ArgsInvalid tests the Execute function with bad input arguments.
func TestExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := downloadermodel.Args{}

	// Execute
	result, err := downloaderObject.Execute(args, pythonObject)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ScriptError tests the Execute function with failing script.
func TestExecute_ScriptError(t *testing.T) {
	// Mock python script to fail
	pythonObject.(*mock.MockPython).ScriptResult = []byte{}
	pythonObject.(*mock.MockPython).Error = errors.New("")

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderObject.Execute(args, pythonObject)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ResponseEmpty tests the Execute function with script returning no data.
func TestExecute_ResponseEmpty(t *testing.T) {
	// Mock python script to return no data
	pythonObject.(*mock.MockPython).ScriptResult = nil
	pythonObject.(*mock.MockPython).Error = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderObject.Execute(args, pythonObject)

	// Assert
	test.AssertEqual(t, err.Error(), "the script didn't return any data")
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestExecute_ResponseBadFormat tests the Execute function with script returning bad data.
func TestExecute_ResponseBadFormat(t *testing.T) {
	// Mock python script to return bad data
	pythonObject.(*mock.MockPython).ScriptResult = []byte("{ bad: property }")
	pythonObject.(*mock.MockPython).Error = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderObject.Execute(args, pythonObject)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_Success tests the Execute function with succeeding script.
func TestExecute_Success(t *testing.T) {
	// Mock python script to succeed
	pythonObject.(*mock.MockPython).ScriptResult = []byte("{}")
	pythonObject.(*mock.MockPython).Error = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderObject.Execute(args, pythonObject)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestNewScriptDownloader tests NewScriptDownloader
func TestNewScriptDownloader(t *testing.T) {
	test.AssertEqual(t, &scriptDownloader{}, NewScriptDownloader())
}
