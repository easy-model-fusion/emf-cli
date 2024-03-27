package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/spf13/cobra"
)

const tokenizerCommandName string = "tokenizer"

// TokenizerCmd represents the tokenizer command
var TokenizerCmd = &cobra.Command{
	Use:   tokenizerCommandName,
	Short: "Palette that contains tokenizer based commands",
	Long:  "Palette that contains tokenizer based commands",
	Run:   runTokenizer,
}

func init() {
	// Preparing to use the hugging face API
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Adding the subcommands
	TokenizerCmd.AddCommand(tokenizerRemoveCmd)
	TokenizerCmd.AddCommand(tokenizerUpdateCmd)
	TokenizerCmd.AddCommand(tokenizerAddCmd)
}

// runTokenizer runs model command
func runTokenizer(cmd *cobra.Command, args []string) {

	// Running command as palette : allowing user to choose subcommand
	err := cobrautil.RunCommandAsPalette(cmd, args, tokenizerCommandName, []string{})
	if err != nil {
		app.UI().Error().Println("Something went wrong :", err)
	}
}
