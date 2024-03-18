package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
)

// tokenizerRemoveCmd represents the model remove tokenizer command
var tokenizerRemoveCmd = &cobra.Command{
	Use:   "remove tokenizer <model_name> [tokenizer..]",
	Short: "Remove one or more tokenizer",
	Long:  "Remove one or more tokenizer",
	Run:   runTokenizerRemove,
}

// runModelRemove runs the model remove command
func runTokenizerRemove(cmd *cobra.Command, args []string) {
	controller.TokenizerRemoveCmd(args)
}
