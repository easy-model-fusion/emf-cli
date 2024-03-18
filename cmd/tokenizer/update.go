package cmdtokenizer

import (
	"github.com/spf13/cobra"
)

// tokenizerUpdateCmd represents the tokenizer update command
var tokenizerUpdateCmd = &cobra.Command{
	Use:   "update <model name> [tokenizer names...]",
	Short: "Update one or more tokenizers",
	Long:  "Update one or more tokenizers",
	Run:   runTokenizerUpdate,
}

// runTokenizerUpdate runs the tokenizer update command
func runTokenizerUpdate(cmd *cobra.Command, args []string) {

}
