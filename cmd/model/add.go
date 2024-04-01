package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/controller/model"
	"github.com/easy-model-fusion/emf-cli/internal/downloader/model"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/spf13/cobra"
	"os"
)

// addCmd represents the add model by names command
var modelAddCmd = &cobra.Command{
	Use:   "add [model name]",
	Short: "Add model by name to your project",
	Long:  `Add model by name to your project`,
	Run:   runAdd,
}

var (
	customArgs    downloadermodel.Args
	addController = modelcontroller.AddController{}
)

func init() {
	// Initialize hugging face api
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Bind cobra args to the downloader script args
	customArgs.ToCobra(modelAddCmd)
	customArgs.DirectoryPath = app.DownloadDirectoryPath
	modelAddCmd.Flags().BoolVarP(&addController.AuthorizeDownload, "yes", "y", false, "Automatic yes to prompts")
	modelAddCmd.Flags().BoolVarP(&addController.SingleFile, "single-file", "S", false, "Use the model as a single file, (usually its a safetensors file)")
}

// runAddByNames runs the add command to add models by name
func runAdd(cmd *cobra.Command, args []string) {
	err := addController.Run(args, customArgs)
	if err != nil {
		app.UI().Error().Println(err.Error())
		os.Exit(1)
	}
}
