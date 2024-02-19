package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var allFlagDelete bool
var authorizeAllDelete bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   	"clean",
	Short: 	"Clean project",
	Long: 	"Clean project",
	Run:   	runClean,
}

func runClean(cmd *cobra.Command, args []string)  {
	if allFlagDelete {
		if !authorizeAllDelete{
			yes := utils.AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
			if !yes {
				return
			}
		}
		if config.GetViperConfig() != nil {
			return
		}
		sdk.SendUpdateSuggestion()
		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		}
	}

	currentDir, errr := os.Getwd()
	if errr != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}

	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "build" {
			err := deleteDir(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.", err)
	}
	return
}

func deleteDir(dossier string) error {
	dir, err := os.Open(dossier)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		fileToRemove := filepath.Join(dossier, fileName)
		if err := os.Remove(fileToRemove); err != nil {
            return err
        }
    }
    return nil
}

func init() {
	cleanCmd.Flags().BoolVarP(&allFlagDelete, "all", "a", false, "clean all project")
	cleanCmd.Flags().BoolVarP(&authorizeAllDelete, "yes", "y", false, "authorize all deletions")
	rootCmd.AddCommand(cleanCmd)
}

