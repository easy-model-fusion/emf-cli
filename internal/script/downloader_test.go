package script

import (
	"github.com/easy-model-fusion/client/test"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"testing"
)

// TestIsDownloaderScriptModelEmpty_True tests the IsDownloaderScriptModelEmpty to return true.
func TestIsDownloaderScriptModelEmpty_True(t *testing.T) {
	// Init
	sm := DownloaderModel{}

	// Execute
	result := IsDownloaderScriptModelEmpty(sm)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestIsDownloaderScriptModelEmpty_False tests the IsDownloaderScriptModelEmpty to return true.
func TestIsDownloaderScriptModelEmpty_False(t *testing.T) {
	// Init
	sm := DownloaderModel{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
	}

	// Execute
	result := IsDownloaderScriptModelEmpty(sm)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestIsScriptTokenizerEmpty_True tests the IsDownloaderScriptTokenizer to return true.
func TestIsScriptTokenizerEmpty_True(t *testing.T) {
	// Init
	st := DownloaderTokenizer{}

	// Execute
	result := IsDownloaderScriptTokenizer(st)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestDownloaderArgsForCobra tests the DownloaderArgsForCobra.
func TestDownloaderArgsForCobra(t *testing.T) {
	// Init
	cmd := &cobra.Command{}
	args := &DownloaderArgs{}

	// Execute
	DownloaderArgsForCobra(cmd, args)

	// Assert
	test.AssertNotEqual(t, cmd.Flags().Lookup(ModelName), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(ModelModule), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(ModelClass), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(ModelOptions), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(TokenizerClass), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(TokenizerOptions), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(Overwrite), nil)
	test.AssertNotEqual(t, cmd.Flags().Lookup(Skip), nil)

	test.AssertEqual(t, args.ModelName, "")
	test.AssertEqual(t, args.ModelModule, "")
	test.AssertEqual(t, args.ModelClass, "")
	test.AssertEqual(t, len(args.ModelOptions), 0)
	test.AssertEqual(t, args.TokenizerClass, "")
	test.AssertEqual(t, len(args.TokenizerOptions), 0)
	test.AssertEqual(t, args.Overwrite, false)
	test.AssertEqual(t, args.Skip, "")
}

// TestDownloaderArgsForPython tests the DownloaderArgsForPython.
func TestDownloaderArgsForPython(t *testing.T) {
	// Init
	args := DownloaderArgs{
		DownloadPath:     "/path/to/download",
		ModelName:        "model",
		ModelModule:      "module",
		ModelClass:       "class",
		ModelOptions:     []string{"opt1=val1", "opt2=val2"},
		TokenizerClass:   "tokenizer",
		TokenizerOptions: []string{"tok_opt1=val1"},
		Skip:             "model",
		Overwrite:        true,
	}
	expected := []string{
		TagPrefix + EmfClient, "/path/to/download", "model", "module",
		TagPrefix + ModelClass, "class",
		TagPrefix + ModelOptions, "opt1=val1", "opt2=val2",
		TagPrefix + TokenizerClass, "tokenizer",
		TagPrefix + TokenizerOptions, "tok_opt1=val1",
		TagPrefix + Skip, "model",
		TagPrefix + Overwrite,
	}

	// Execute
	result := DownloaderArgsForPython(args)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}

// TestIsScriptTokenizerEmpty_False tests the IsDownloaderScriptTokenizer to return true.
func TestIsScriptTokenizerEmpty_False(t *testing.T) {
	// Init
	st := DownloaderTokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}

	// Execute
	result := IsDownloaderScriptTokenizer(st)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestDownloaderArgsValidate_MissingName tests the DownloaderArgsValidate function to return an error.
func TestDownloaderArgsValidate_MissingName(t *testing.T) {
	// Init
	args := DownloaderArgs{ModelModule: "present"}

	// Execute
	result := DownloaderArgsValidate(args)

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestDownloaderArgsValidate_MissingModule tests the DownloaderArgsValidate function to return an error.
func TestDownloaderArgsValidate_MissingModule(t *testing.T) {
	// Init
	args := DownloaderArgs{ModelName: "present"}

	// Execute
	result := DownloaderArgsValidate(args)

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestDownloaderArgsValidate_Success tests the DownloaderArgsValidate function to succeed.
func TestDownloaderArgsValidate_Success(t *testing.T) {
	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result := DownloaderArgsValidate(args)

	// Assert
	test.AssertEqual(t, result, nil)
}

// TestDownloaderArgsProcess_ValidateFail tests the DownloaderArgsProcess function to fail.
func TestDownloaderArgsProcess_ValidateFail(t *testing.T) {
	// Init
	args := DownloaderArgs{}

	// Execute
	_, err := DownloaderArgsProcess(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
}

// TestDownloaderArgsProcess_OverwriteFalse tests the DownloaderArgsProcess function to succeed.
func TestDownloaderArgsProcess_OverwriteFalse(t *testing.T) {
	// TODO : mock utils.AskForUsersConfirmation to return false
	t.Skip()

	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}
	path := filepath.Join(DownloadModelsPath, args.ModelName)
	err := os.MkdirAll(path, 0750)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(path)

	// Execute
	result, err := DownloaderArgsProcess(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result, nil)
}

// TestDownloaderArgsProcess_Success tests the DownloaderArgsValidate function to succeed.
func TestDownloaderArgsProcess_Success(t *testing.T) {
	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := DownloaderArgsProcess(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertNotEqual(t, result, nil)
}

// TestDownloaderExecute_ArgsInvalid tests the DownloaderExecute function with bad input.
func TestDownloaderExecute_ArgsInvalid(t *testing.T) {
	// Init
	args := DownloaderArgs{}

	// Execute
	result, err := DownloaderExecute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestDownloaderExecute_OverwriteFalse tests the DownloaderExecute function with already downloaded model.
func TestDownloaderExecute_OverwriteFalse(t *testing.T) {
	// TODO : mock DownloaderArgsProcess to return (nil, nil)
	t.Skip()

	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := DownloaderExecute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestDownloaderExecute_ScriptError tests the DownloaderExecute function with failing script.
func TestDownloaderExecute_ScriptError(t *testing.T) {
	// TODO : mock utils.ExecuteScript to return ([]byte{}, errors.New(""), 0)

	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := DownloaderExecute(args)

	// Assert
	test.AssertNotEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}

// TestDownloaderExecute_ResponseEmpty tests the DownloaderExecute function with script returning no data.
func TestDownloaderExecute_ResponseEmpty(t *testing.T) {
	// TODO : mock utils.ExecuteScript to return (nil, nil, 0)
	t.Skip()

	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := DownloaderExecute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, true)
}

// TestDownloaderExecute_Success tests the DownloaderExecute function with succeeding script.
func TestDownloaderExecute_Success(t *testing.T) {

	/*responseBadString := "{ \"bad\": \"property\" }"
	bytes := []byte(responseBadString)*/

	// TODO : mock utils.ExecuteScript to return (bytes, nil, 0)
	t.Skip()

	// Init
	args := DownloaderArgs{ModelName: "present", ModelModule: "present"}

	// Execute
	result, err := DownloaderExecute(args)

	// Assert
	test.AssertEqual(t, err, nil)
	test.AssertEqual(t, result.IsEmpty, false)
}
