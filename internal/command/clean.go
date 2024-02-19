package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var allFlag bool
var authorizeAll bool

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   	"clean",
	Short: 	"Clean project",
	Long: 	"Clean project",
	Run:   	runClean,
}

func runClean(cmd *cobra.Command, args []string)  {
	
	if allFlag {
		if !authorizeAll{
			yes := utils.AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
			if !yes {
				return
			}
		}
		err := config.RemoveAllModels()
		if err == nil {
			pterm.Success.Printfln("Operation succeeded.")
		}
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			pterm.Error.Printfln("Operation failed.")
			return err
		}
		if err := deleteDir(path); err != nil {
            pterm.Error.Printfln("Failed to delete directory: %v", err)
            return err
        }
		return nil
	})
	return
}

func deleteDir(dossier string) error {
	dir, err := os.Open(dossier)
	if err != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(filePath, fileName)

		err = os.RemoveAll(filePath)
		if err != nil {
			pterm.Error.Printfln("Operation failed.")
			return
		}
	}
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func init() {
	cleanCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "clean all project")
	cleanCmd.Flags().BoolVarP(&authorizeAll, "yes", "y", false, "authorize all deletions")
	rootCmd.AddCommand(cleanCmd)
}

