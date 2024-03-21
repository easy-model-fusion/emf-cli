package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	modelcontroller "github.com/easy-model-fusion/emf-cli/internal/controller/model"
	downloadermodel "github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/spf13/cobra"
)

// addCmd represents the add model by names command
var modelAddCmd = &cobra.Command{
	Use:   "add [<model name>]",
	Short: "Add model by name to your project",
	Long:  `Add model by name to your project`,
	Run:   runAdd,
}

var customArgs downloadermodel.Args
var yes bool

func init() {
	// Bind cobra args to the downloader script args
	customArgs.ToCobra(modelAddCmd)
	customArgs.DirectoryPath = app.DownloadDirectoryPath
	modelAddCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic yes to prompts")
}

// runAddByNames runs the add command to add models by name
func runAdd(cmd *cobra.Command, args []string) {
	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")
	modelcontroller.RunAdd(args, customArgs, yes)
}
