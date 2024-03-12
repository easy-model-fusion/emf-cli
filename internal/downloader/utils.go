package downloader

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/spf13/cobra"
)

// EmptyModel checks if a Model is empty.
func EmptyModel(model Model) bool {
	return model.Path == "" && model.Class == ""
}

// EmptyTokenizer checks if a Tokenizer is empty.
func EmptyTokenizer(tokenizer Tokenizer) bool {
	return tokenizer.Path == "" && tokenizer.Class == ""
}

// ArgsValidate validates the mandatory fields for Args
func ArgsValidate(args Args) error {

	// Name validity
	if args.ModelName == "" {
		return errors.New("missing name for the model")
	}

	// Module validity
	if args.ModelModule == "" {
		return errors.New("missing module for the model")
	}

	return nil
}

// ArgsGetForCobra builds the arguments for running the cobra command
func ArgsGetForCobra(cmd *cobra.Command, args *Args) {

	// Pseudo mandatory : allowing to customize the calling command
	cmd.Flags().StringVarP(&args.ModelName, ModelName, "n", "", "Name of the model")
	cmd.Flags().StringVarP(&args.ModelModule, ModelModule, "m", "", "Python module used for download")

	// Optional for the model
	cmd.Flags().StringVarP(&args.ModelClass, ModelClass, "c", "", "Python class within the module")
	cmd.Flags().StringSliceVarP(&args.ModelOptions, ModelOptions, "o", []string{}, "List of model options")

	// Optional for the tokenizer
	cmd.Flags().StringVarP(&args.TokenizerClass, TokenizerClass, "t", "", "Tokenizer class (only for transformers)")
	cmd.Flags().StringArrayVarP(&args.TokenizerOptions, TokenizerOptions, "T", []string{}, "List of tokenizer options (only for transformers)")

	// Situational
	cmd.Flags().StringVarP(&args.Skip, Skip, "s", "", "Skip the model or tokenizer download")
}

// ArgsProcessForPython builds the arguments for running the python script.
// Pre-condition : certain that the user authorized the overwriting when downloading the model.
func ArgsProcessForPython(args Args) []string {

	// Mandatory arguments
	cmdArgs := []string{TagPrefix + EmfClient, TagPrefix + Overwrite, app.DownloadDirectoryPath, args.ModelName, args.ModelModule}

	// Optional arguments regarding the model
	if args.ModelClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+ModelClass, args.ModelClass)
	}
	if len(args.ModelOptions) != 0 {
		var options []string
		for _, modelOption := range args.ModelOptions {
			options = append(options, stringutil.ParseOptions(modelOption)...)
		}
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + ModelOptions}, options...)...)
	}

	// Optional arguments regarding the model's tokenizer
	if args.TokenizerClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+TokenizerClass, args.TokenizerClass)
	}
	if len(args.TokenizerOptions) != 0 {
		var options []string
		for _, modelOption := range args.TokenizerOptions {
			options = append(options, stringutil.ParseOptions(modelOption)...)
		}
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + TokenizerOptions}, options...)...)
	}

	// Global tags for the script
	if len(args.Skip) != 0 {
		cmdArgs = append(cmdArgs, TagPrefix+Skip, args.Skip)
	}

	return cmdArgs
}
