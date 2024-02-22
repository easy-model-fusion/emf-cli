package downloader

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"testing"
)

// TestDownloaderExecute_ArgsInvalid tests the Execute function with bad input.
func TestDownloaderExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := Args{}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestDownloaderExecute_ScriptError tests the Execute function with failing script.
func TestDownloaderExecute_ScriptError(t *testing.T) {
	// TODO : mock utils.ExecuteScript to return ([]byte{}, errors.New(""), 0)

	// Init
	args := Args{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestDownloaderExecute_ResponseEmpty tests the Execute function with script returning no data.
func TestDownloaderExecute_ResponseEmpty(t *testing.T) {
	// TODO : mock utils.ExecuteScript to return (nil, nil, 0)
	t.Skip()

	// Init
	args := Args{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestDownloaderExecute_Success tests the Execute function with succeeding script.
func TestDownloaderExecute_Success(t *testing.T) {

	/*responseBadString := "{ \"bad\": \"property\" }"
	bytes := []byte(responseBadString)*/

	// TODO : mock utils.ExecuteScript to return (bytes, nil, 0)
	t.Skip()

	// Init
	args := Args{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := Execute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}
