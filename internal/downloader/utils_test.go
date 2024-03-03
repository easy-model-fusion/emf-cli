package downloader

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/cobra"
	"testing"
)

// TestEmptyModel_True tests the EmptyModel to return true.
func TestEmptyModel_True(t *testing.T) {
	// Init
	sm := Model{}

	// Execute
	result := EmptyModel(sm)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestEmptyModel_False tests the EmptyModel to return true.
func TestEmptyModel_False(t *testing.T) {
	// Init
	sm := Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
	}

	// Execute
	result := EmptyModel(sm)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestEmptyTokenizer_True tests the EmptyTokenizer to return true.
func TestEmptyTokenizer_True(t *testing.T) {
	// Init
	st := Tokenizer{}

	// Execute
	result := EmptyTokenizer(st)

	// Assert
	test.AssertEqual(t, true, result)
}

// TestEmptyTokenizer_False tests the EmptyTokenizer to return true.
func TestEmptyTokenizer_False(t *testing.T) {
	// Init
	st := Tokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}

	// Execute
	result := EmptyTokenizer(st)

	// Assert
	test.AssertEqual(t, false, result)
}

// TestArgsValidate_MissingName tests the ArgsValidate function to return an error.
func TestArgsValidate_MissingName(t *testing.T) {
	// Init
	args := Args{ModelModule: "present"}

	// Execute
	result := ArgsValidate(args)

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestArgsValidate_MissingModule tests the ArgsValidate function to return an error.
func TestArgsValidate_MissingModule(t *testing.T) {
	// Init
	args := Args{ModelName: "present"}

	// Execute
	result := ArgsValidate(args)

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestArgsValidate_Success tests the ArgsValidate function to succeed.
func TestArgsValidate_Success(t *testing.T) {
	// Init
	args := Args{ModelName: "present", ModelModule: "present"}

	// Execute
	result := ArgsValidate(args)

	// Assert
	test.AssertEqual(t, result, nil)
}

// TestArgsGetForCobra tests the ArgsGetForCobra.
func TestArgsGetForCobra(t *testing.T) {
	// Init
	cmd := &cobra.Command{}
	args := &Args{}

	// Execute
	ArgsGetForCobra(cmd, args)

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
	test.AssertEqual(t, args.Skip, "")
}

// TestArgsProcessForPython tests the ArgsProcessForPython.
func TestArgsProcessForPython(t *testing.T) {
	// Init
	args := Args{
		ModelName:        "model",
		ModelModule:      "module",
		ModelClass:       "class",
		ModelOptions:     []string{"opt1=val1", "opt2=val2"},
		TokenizerClass:   "tokenizer",
		TokenizerOptions: []string{"tok_opt1=val1"},
		Skip:             "model",
	}
	expected := []string{
		TagPrefix + EmfClient, TagPrefix + Overwrite,
		"/path/to/download", "model", "module",
		TagPrefix + ModelClass, "class",
		TagPrefix + ModelOptions, "opt1=val1", "opt2=val2",
		TagPrefix + TokenizerClass, "tokenizer",
		TagPrefix + TokenizerOptions, "tok_opt1=val1",
		TagPrefix + Skip, "model",
	}

	// Execute
	result := ArgsProcessForPython(args)

	// Assert
	test.AssertEqual(t, len(result), len(expected))
}
