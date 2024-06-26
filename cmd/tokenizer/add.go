package cmdtokenizer

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/controller/tokenizer"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/spf13/cobra"
)

var (
	addController tokenizer.AddController
)

// tokenizerAddCmd represents the tokenizer add command
var tokenizerAddCmd = &cobra.Command{
	Use:   "add <model name> <tokenizer name>",
	Short: "Add one or more tokenizers",
	Long:  "Add one or more tokenizers",
	Args:  cobra.MinimumNArgs(2),
	Run:   runTokenizerAdd,
}

var customArgs downloadermodel.Args
var yes bool

func init() {
	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Bind cobra args to the downloader script args
	customArgs.ToCobraTokenizer(tokenizerAddCmd)
	tokenizerAddCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic yes to prompts")
}

// runTokenizerAdd runs the tokenizer add command
func runTokenizerAdd(cmd *cobra.Command, args []string) {
	err := addController.Run(args, customArgs)
	if err != nil {
		return
	}
}
