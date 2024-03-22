package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/spf13/cobra"
)

// tokenizerAddCmd represents the tokenizer add command
var tokenizerAddCmd = &cobra.Command{
	Use:   "add <model name> <tokenizer name> [<other tokenizer names>...]",
	Short: "Add one or more tokenizers",
	Long:  "Add one or more tokenizers",
	Args:  cobra.MinimumNArgs(1),
	Run:   runTokenizerAdd,
}

// runTokenizerAdd runs the tokenizer add command
func runTokenizerAdd(cmd *cobra.Command, args []string) {
	sdk.SendUpdateSuggestion()

}
