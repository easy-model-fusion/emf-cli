package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/tokenizer"
	"github.com/spf13/cobra"
)

var (
	removeTokenizerController tokenizer.RemoveTokenizerController
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
	removeTokenizerController.RunTokenizerRemove(args)
}
