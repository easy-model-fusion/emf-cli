package command

import (
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
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
	// Delete all models if flag --all
	if allFlagDelete {
		// Ask for confirmation 
		if !authorizeAllDelete{
			yes := utils.AskForUsersConfirmation("Are you sure you want to delete all downloaded models and clean the build files of this project?")
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
		} else{
			pterm.Error.Printfln("Operation failed.")
		}
		
	}

	// Get the current dir
	currentDir, err := os.Getwd()
	if err != nil {
		pterm.Error.Printfln("Operation failed.")
		return
	}

	// Search build dir in the current dir
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// delete all file and dir in build
		if info.IsDir() && info.Name() == "build" {
			err := deleteDir(path)
			if err != nil {
				return err
			}
			// skip search in build and continue
			return filepath.SkipDir
		}
		return nil
	})
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func deleteDir(dossier string) error {
	// Open the dir
	dir, err := os.Open(dossier)
	if err != nil {
		return err
	}
	// Close the dir at the end of function
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	// Delete files and dirs one by one
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
