package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/tokenizer"
	"github.com/spf13/cobra"
	"os"
)

var (
	updateTokenizerController tokenizer.UpdateTokenizerController
)

// tokenizerUpdateCmd represents the tokenizer update command
var tokenizerUpdateCmd = &cobra.Command{
	Use:   "update <model name> <tokenizer name> [<other tokenizer names>...]",
	Short: "Update one or more tokenizers",
	Long:  "Update one or more tokenizers",
	Run:   runTokenizerUpdate,
}

// runTokenizerUpdate runs the tokenizer update command
func runTokenizerUpdate(cmd *cobra.Command, args []string) {
	err := updateTokenizerController.TokenizerUpdateCmd(args)
	if err != nil {
		os.Exit(1)
	}
}
