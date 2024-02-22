package cmd

import (
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils/ptermutil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var allFlagDelete bool
var authorizeAllDelete bool

const dirName = "build"

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean project",
	Long:  "Clean project",
	Run:   runClean,
}

func runClean(cmd *cobra.Command, args []string) {
	// Delete all models if flag --all
	if allFlagDelete {
		// Ask for confirmation
		if !authorizeAllDelete {
			yes := ptermutil.AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
			if !yes {
				return
			}
		}
		if config.GetViperConfig(config.FilePath) != nil {
			return
		}
		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		} else {
			pterm.Error.Printfln("Operation failed.")
		}

	}

	// Get the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}
	buildDir := filepath.Join(currentDir, dirName)

	_, err = os.Stat(buildDir)
	if os.IsNotExist(err) {
		pterm.Success.Printfln("Operation succeeded.")
		return
	}
	err = os.RemoveAll(buildDir)
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func init() {
	cleanCmd.Flags().BoolVarP(&allFlagDelete, "all", "a", false, "clean all project")
	cleanCmd.Flags().BoolVarP(&authorizeAllDelete, "yes", "y", false, "authorize all deletions")
	rootCmd.AddCommand(cleanCmd)
}
