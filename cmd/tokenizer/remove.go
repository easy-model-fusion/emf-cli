package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/tokenizer"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/spf13/cobra"
)

// tokenizerRemoveCmd represents the tokenizer remove command
var tokenizerRemoveCmd = &cobra.Command{
	Use:   "remove <model name> <tokenizer name> [<other tokenizer names>...]",
	Short: "Remove one or more tokenizer",
	Long:  "Remove one or more tokenizer",
	Args:  cobra.MinimumNArgs(1),
	Run:   runTokenizerRemove,
}

// runTokenizerRemove runs the tokenizer remove command
func runTokenizerRemove(cmd *cobra.Command, args []string) {
	sdk.SendUpdateSuggestion()
	tokenizer.TokenizerRemoveCmd(args)
}
