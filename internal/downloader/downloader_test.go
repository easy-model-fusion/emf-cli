package downloader

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app.SetUI(&mock.MockUI{})
	app.SetPython(&mock.MockPython{})
	os.Exit(m.Run())
}

// TestExecute_ArgsInvalid tests the Execute function with bad input arguments.
func TestExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := Args{}

	// Execute
	result, err := app.Downloader().Execute(args, app.Python())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ScriptError tests the Execute function with failing script.
func TestExecute_ScriptError(t *testing.T) {
	// Mock python script to fail
	app.Python().(*mock.MockPython).ScriptResult = []byte{}
	app.Python().(*mock.MockPython).Error = errors.New("")

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := app.Downloader().Execute(args, app.Python())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_ResponseEmpty tests the Execute function with script returning no data.
func TestExecute_ResponseEmpty(t *testing.T) {
	// Mock python script to return no data
	app.Python().(*mock.MockPython).ScriptResult = nil
	app.Python().(*mock.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := app.Downloader().Execute(args, app.Python())

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestExecute_ResponseBadFormat tests the Execute function with script returning bad data.
func TestExecute_ResponseBadFormat(t *testing.T) {
	// Mock python script to return bad data
	app.Python().(*mock.MockPython).ScriptResult = []byte("{ bad: property }")
	app.Python().(*mock.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := app.Downloader().Execute(args, app.Python())

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestExecute_Success tests the Execute function with succeeding script.
func TestExecute_Success(t *testing.T) {
	// Mock python script to succeed
	app.Python().(*mock.MockPython).ScriptResult = []byte("{}")
	app.Python().(*mock.MockPython).Error = nil

	// Init
	args := Args{ModelName: "ModelName", ModelModule: "ModelModule"}

	// Execute
	result, err := app.Downloader().Execute(args, app.Python())

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}
