package command

import (
	"github.com/easy-model-fusion/client/internal/config"
	"github.com/easy-model-fusion/client/internal/sdk"
	"github.com/easy-model-fusion/client/internal/model"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   	"clean",
	Short: 	"Clean project",
	Long: 	"Clean project",
	Run:   	runClean,
}

func runClean(cmd *cobra.Command, args []string) {
	// extensions file removed
    extensions := []string{".exe", ".o", ".obj", ".out"}

    // Get the current working directory
    currentDir, err := os.Getwd()
    if err != nil {
        return err
    }

   	// Walk through all files in the current working directory
    err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            pterm.Error.Printfln("Operation failed.")
			return
        }

        // Check if the file has one of the specified extensions to be removed
        for _, ext := range extensions {
            if filepath.Ext(path) == ext {
                // Remove the file
                err := os.Remove(path)
                if err != nil {
                    pterm.Error.Printfln("Operation failed.")
					return
                }
				pterm.Success.Printfln("Operation succeeded, file %s\n deleted", pat)
                break
            }
        }

        return nil
    })

    if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func runClean(cmd *cobra.Command, args []string)  {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "build" {
			err := deleteDir(path)
			if err != nil {
				pterm.Error.Printfln("Operation failed.")
				return
			}
			else {
				pterm.Success.Printfln("Operation succeeded.")
			}
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
		return nil
	}
	if err == nil {
		pterm.Success.Printfln("Operation succeeded.")
	} else {
		pterm.Error.Printfln("Operation failed.")
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

