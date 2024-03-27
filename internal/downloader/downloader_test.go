package downloader

import (
	"context"
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/python"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

var pythonInterface python.Python
var downloaderInterface Downloader

func TestMain(m *testing.M) {
	pythonInterface = &mock.MockPython{}
	downloaderInterface = &scriptDownloader{}
	os.Exit(m.Run())
}

// TestExecute_ArgsInvalid tests the Execute function with bad input arguments.
func TestExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := downloadermodel.Args{}

	// Execute
	result, err := downloaderInterface.Execute(args, pythonInterface, context.Background())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ScriptError tests the Execute function with failing script.
func TestExecute_ScriptError(t *testing.T) {
	// Mock python script to fail
	pythonInterface.(*mock.MockPython).ScriptResult = []byte{}
	pythonInterface.(*mock.MockPython).ExecuteScriptError = errors.New("")

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderInterface.Execute(args, pythonInterface, context.Background())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ResponseEmpty tests the Execute function with script returning no data.
func TestExecute_ResponseEmpty(t *testing.T) {
	// Mock python script to return no data
	pythonInterface.(*mock.MockPython).ScriptResult = nil
	pythonInterface.(*mock.MockPython).ExecuteScriptError = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderInterface.Execute(args, pythonInterface, context.Background())

	// Assert
	test.AssertEqual(t, err.Error(), "the script didn't return any data")
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestExecute_ResponseBadFormat tests the Execute function with script returning bad data.
func TestExecute_ResponseBadFormat(t *testing.T) {
	// Mock python script to return bad data
	pythonInterface.(*mock.MockPython).ScriptResult = []byte("{ bad: property }")
	pythonInterface.(*mock.MockPython).ExecuteScriptError = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderInterface.Execute(args, pythonInterface, context.Background())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_Success tests the Execute function with succeeding script.
func TestExecute_Success(t *testing.T) {
	// Mock python script to succeed
	pythonInterface.(*mock.MockPython).ScriptResult = []byte("{}")
	pythonInterface.(*mock.MockPython).ExecuteScriptError = nil

	// Init
	args := downloadermodel.Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := downloaderInterface.Execute(args, pythonInterface, context.Background())

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestNewScriptDownloader tests NewScriptDownloader
func TestNewScriptDownloader(t *testing.T) {
	test.AssertEqual(t, &scriptDownloader{}, NewScriptDownloader())
}
