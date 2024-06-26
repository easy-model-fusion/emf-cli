package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/spf13/cobra"
	"os"
)

var (
	initUseTorchCuda bool
	initController   controller.InitController
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <project name>",
	Short: "Initialize a EMF project",
	Long:  `Initialize a EMF project.`,
	Args:  fileutil.ValidFileName(1, true),
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	err := initController.Run(args, initUseTorchCuda, "")
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initController = controller.InitController{}
	initCmd.Flags().BoolVarP(&initUseTorchCuda, "cuda", "c", false, "Use torch with cuda")
}
