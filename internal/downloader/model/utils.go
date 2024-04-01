package downloadermodel

import (
	"errors"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/spf13/cobra"
)

// Empty checks if a Model is empty.
func (m *Model) Empty() bool {
	return m.Path == "" && m.Class == ""
}

// Empty checks if a Tokenizer is empty.
func (t *Tokenizer) Empty() bool {
	return t.Path == "" && t.Class == ""
}

// Validate validates the mandatory fields for Args
func (a *Args) Validate() error {

	// Name validity
	if a.ModelName == "" {
		return errors.New("missing name for the model")
	}

	// Module validity
	if a.ModelModule == "" {
		return errors.New("missing module for the model")
	}

	return nil
}

// ToCobra builds the arguments for running the cobra command
func (a *Args) ToCobra(cmd *cobra.Command) {

	// Optional for the model
	cmd.Flags().StringVarP(&a.ModelClass, ModelClass, "c", "", "Python class within the module")
	cmd.Flags().StringSliceVarP(&a.ModelOptions, ModelOptions, "o", []string{}, "List of model options")
	cmd.Flags().StringVarP(&a.ModelModule, ModelModule, "m", "", "Python module used for download")
	cmd.Flags().StringVarP(&a.DirectoryPath, Path, "p", "", "Downloaded Model directory path")

	// Optional for the tokenizer
	cmd.Flags().StringVarP(&a.TokenizerClass, TokenizerClass, "t", "", "Tokenizer class (only for transformers)")
	cmd.Flags().StringArrayVarP(&a.TokenizerOptions, TokenizerOptions, "T", []string{}, "List of tokenizer options (only for transformers)")

	// Situational
	cmd.Flags().BoolVarP(&a.OnlyConfiguration, "only-configuration", "O", false, "Only configure the model without downloading it")
	cmd.Flags().BoolVarP(&a.SkipTokenizer, "skip-tokenizer", "s", false, "Skip tokenizer download")

	// Authorization token
	cmd.Flags().StringVarP(&a.AccessToken, AccessToken, "a", "", "Access token for gated models")
}

// ToCobraTokenizer builds the arguments for running the cobra command
func (a *Args) ToCobraTokenizer(cmd *cobra.Command) {

	// Optional for the tokenizer
	cmd.Flags().StringVarP(&a.TokenizerClass, TokenizerClass, "c", "", "Tokenizer class (only for transformers)")
	cmd.Flags().StringArrayVarP(&a.TokenizerOptions, TokenizerOptions, "o", []string{}, "List of tokenizer options (only for transformers)")

}

// ToPython builds the arguments for running the python script.
// Pre-condition : certain that the user authorized the overwriting when downloading the model.
func (a *Args) ToPython() []string {

	// Mandatory arguments
	cmdArgs := []string{TagPrefix + EmfClient, TagPrefix + Overwrite, a.DirectoryPath, a.ModelName, a.ModelModule}

	// Optional arguments regarding the model
	if a.ModelClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+ModelClass, a.ModelClass)
	}
	if len(a.ModelOptions) != 0 {
		var options []string
		for _, modelOption := range a.ModelOptions {
			options = append(options, stringutil.ParseOptions(modelOption)...)
		}
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + ModelOptions}, options...)...)
	}

	// Optional arguments regarding the model's tokenizer
	if a.TokenizerClass != "" {
		cmdArgs = append(cmdArgs, TagPrefix+TokenizerClass, a.TokenizerClass)
	}
	if len(a.TokenizerOptions) != 0 {
		var options []string
		for _, modelOption := range a.TokenizerOptions {
			options = append(options, stringutil.ParseOptions(modelOption)...)
		}
		cmdArgs = append(cmdArgs, append([]string{TagPrefix + TokenizerOptions}, options...)...)
	}

	// Global tags for the script
	if a.SkipTokenizer {
		cmdArgs = append(cmdArgs, TagPrefix+Skip, SkipValueTokenizer)
	}
	if a.SkipModel {
		cmdArgs = append(cmdArgs, TagPrefix+Skip, SkipValueModel)
	}

	// Only configuration
	if a.OnlyConfiguration {
		cmdArgs = append(cmdArgs, TagPrefix+OnlyConfiguration)
	}

	// Access token
	if a.AccessToken != "" {
		cmdArgs = append(cmdArgs, TagPrefix+AccessToken, a.AccessToken)
	}
	return cmdArgs
}
