package cmdtokenizer

import (
	"github.com/spf13/cobra"
)

// tokenizerAddCmd represents the tokenizer add command
var tokenizerAddCmd = &cobra.Command{
	Use:   "add <model name> [tokenizer names...]",
	Short: "Add one or more tokenizers",
	Long:  "Add one or more tokenizers",
	Run:   runTokenizerAdd,
}

// runTokenizerAdd runs the tokenizer add command
func runTokenizerAdd(cmd *cobra.Command, args []string) {

}
