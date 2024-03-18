package downloader

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.SetUI(&test.MockUI{})
	app.SetPython(&test.MockPython{})
	os.Exit(m.Run())
}

// TestExecute_ArgsInvalid tests the Execute function with bad input arguments.
func TestExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := Args{}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ScriptError tests the Execute function with failing script.
func TestExecute_ScriptError(t *testing.T) {
	// Mock python script to fail
	app.Python().(*test.MockPython).ScriptResult = []byte{}
	app.Python().(*test.MockPython).Error = errors.New("")

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ResponseEmpty tests the Execute function with script returning no data.
func TestExecute_ResponseEmpty(t *testing.T) {
	// Mock python script to return no data
	app.Python().(*test.MockPython).ScriptResult = nil
	app.Python().(*test.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestExecute_ResponseBadFormat tests the Execute function with script returning bad data.
func TestExecute_ResponseBadFormat(t *testing.T) {
	// Mock python script to return bad data
	app.Python().(*test.MockPython).ScriptResult = []byte("{ bad: property }")
	app.Python().(*test.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_Success tests the Execute function with succeeding script.
func TestExecute_Success(t *testing.T) {
	// Mock python script to succeed
	app.Python().(*test.MockPython).ScriptResult = []byte("{}")
	app.Python().(*test.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}
