package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/controller"
	"github.com/spf13/cobra"
	"os"
)

var (
	protectedModelsAccessToken string
	installUseTorchCuda        bool
	installController          controller.InstallController
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a EMF project",
	Long:  `Installs an existing EMF project. Basically combining a slightly different init and the tidy commands.`,
	Run:   runInstall,
}

func runInstall(cmd *cobra.Command, args []string) {
	err := installController.Run(args, installUseTorchCuda, protectedModelsAccessToken)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	installCmd.Flags().BoolVarP(&installUseTorchCuda, "cuda", "c", false, "Use torch with cuda")
	installCmd.Flags().StringVarP(&protectedModelsAccessToken, "access-token", "a", "", "Access token for gated models")
}
