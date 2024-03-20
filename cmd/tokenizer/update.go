package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller/tokenizer"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/spf13/cobra"
)

// tokenizerUpdateCmd represents the model update command
var tokenizerUpdateCmd = &cobra.Command{
	Use:   "update <model_name> [tokenizer..]",
	Short: "Update one or more tokenizers",
	Long:  "Update one or more tokenizers",
	Args:  cobra.MinimumNArgs(1),
	Run:   runTokenizerUpdate,
}

// runTokenizerUpdate runs the model remove command
func runTokenizerUpdate(cmd *cobra.Command, args []string) {
	sdk.SendUpdateSuggestion()
	tokenizer.TokenizerUpdateCmd(args)
}
