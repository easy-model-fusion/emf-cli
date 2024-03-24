package downloadermodel

import (
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/spf13/cobra"
	"testing"
)

// TestEmpty_Model_True tests the EmptyModel to return true.
func TestEmpty_Model_True(t *testing.T) {
	// Init
	model := Model{}

	// Execute
	result := model.Empty()

	// Assert
	test.AssertEqual(t, true, result)
}

// TestEmpty_Model_False tests the EmptyModel to return true.
func TestEmpty_Model_False(t *testing.T) {
	// Init
	model := Model{
		Path:   "/path/to/model",
		Module: "module_name",
		Class:  "class_name",
	}

	// Execute
	result := model.Empty()

	// Assert
	test.AssertEqual(t, false, result)
}

// TestEmpty_Tokenizer_True tests the EmptyTokenizer to return true.
func TestEmpty_Tokenizer_True(t *testing.T) {
	// Init
	tokenizer := Tokenizer{}

	// Execute
	result := tokenizer.Empty()

	// Assert
	test.AssertEqual(t, true, result)
}

// TestEmpty_Tokenizer_False tests the EmptyTokenizer to return true.
func TestEmpty_Tokenizer_False(t *testing.T) {
	// Init
	tokenizer := Tokenizer{
		Path:  "/path/to/tokenizer",
		Class: "tokenizer_class",
	}

	// Execute
	result := tokenizer.Empty()

	// Assert
	test.AssertEqual(t, false, result)
}

// TestValidate_MissingName tests the ArgsValidate function to return an error.
func TestValidate_MissingName(t *testing.T) {
	// Init
	args := Args{ModelModule: "present"}

	// Execute
	result := args.Validate()

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestValidate_MissingModule tests the ArgsValidate function to return an error.
func TestValidate_MissingModule(t *testing.T) {
	// Init
	args := Args{ModelName: "present"}

	// Execute
	result := args.Validate()

	// Assert
	test.AssertNotEqual(t, result, nil)
}

// TestValidate_Success tests the ArgsValidate function to succeed.
func TestValidate_Success(t *testing.T) {
	// Init
	args := Args{ModelName: "present", ModelModule: "present"}

	// Execute
	result := args.Validate()

	// Assert
	test.AssertEqual(t, result, nil)
}

// TestToCobra tests the ArgsGetForCobra.
func TestToCobra(t *testing.T) {
	// Init
	cmd := &cobra.Command{}
	args := &Args{}

	// Execute
	args.ToCobra(cmd)

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
	test.AssertEqual(t, args.SkipModel, false)
	test.AssertEqual(t, args.SkipTokenizer, false)
}

// TestToPython tests the ArgsProcessForPython.
func TestToPython(t *testing.T) {
	// Init
	args := Args{
		ModelName:         "model",
		ModelModule:       "module",
		ModelClass:        "class",
		DirectoryPath:     "/path/to/download",
		ModelOptions:      []string{"opt1=val1", "opt2=val2"},
		TokenizerClass:    "tokenizer",
		TokenizerOptions:  []string{"tok_opt1=val1"},
		SkipModel:         true,
		SkipTokenizer:     false,
		OnlyConfiguration: true,
		AccessToken:       "token",
	}
	expected := []string{
		TagPrefix + EmfClient, TagPrefix + Overwrite,
		"/path/to/download", "model", "module",
		TagPrefix + ModelClass, "class",
		TagPrefix + ModelOptions, "opt1=val1", "opt2=val2",
		TagPrefix + TokenizerClass, "tokenizer",
		TagPrefix + TokenizerOptions, "tok_opt1=val1",
		TagPrefix + Skip, "model",
		TagPrefix + OnlyConfiguration,
		TagPrefix + AccessToken, "token",
	}

	// Execute
	result := args.ToPython()

	// Assert
	test.AssertEqual(t, len(result), len(expected))
	for i := range expected {
		test.AssertEqual(t, result[i], expected[i])
	}
}
