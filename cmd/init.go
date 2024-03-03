package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/spf13/cobra"
)

var (
	useTorchCuda bool
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
	controller.RunInit(args, useTorchCuda)
}

func init() {
	initCmd.Flags().BoolVarP(&useTorchCuda, "cuda", "c", false, "Use torch with cuda")
}
