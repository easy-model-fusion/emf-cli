package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/utils/fileutil"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"os"
	"path/filepath"
)

// RemoveTokenizerPhysically only removes the model from the project's downloaded models
func RemoveTokenizerPhysically(tokenizerPath string) error {

	// Check if the tokenizer_path exists
	if exists, err := fileutil.IsExistingPath(tokenizerPath); err != nil {
		// Skipping model : an error occurred
		return err
	} else if exists {

		// Split the path into a slice of strings
		directories := stringutil.SplitPath(tokenizerPath)

		// Removing tokenizer
		err = os.RemoveAll(tokenizerPath)
		if err != nil {
			return err
		}

		// Excluding the tail since it has already been removed
		directories = directories[:len(directories)-1]

		// Cleaning up : removing every empty directory on the way to the model (from tail to head)
		for i := len(directories) - 1; i >= 0; i-- {
			// Build path to parent directory
			path := filepath.Join(directories[:i+1]...)

			// Delete directory if empty
			err = fileutil.DeleteDirectoryIfEmpty(path)
			if err != nil {
			}
		}
	} else {
		// tokenizer path is not in the current project
	}
	return nil
}

// RemoveTokenizersByName removes specified tokenizers
func RemoveTokenizersByName(currentModel model.Model, tokenizersToRemove model.Tokenizers) error {
	var failedTokenizers []string
	var removedTokenizers model.Tokenizers
	// Trying to remove the models
	for _, item := range tokenizersToRemove {
		// Starting client spinner animation
		spinner := app.UI().StartSpinner(fmt.Sprintf("Removing tokenizer %s...", tokenizersToRemove))
		err := RemoveTokenizerPhysically(item.Path)
		if err != nil {
			spinner.Fail("failed to remove tokenizers")
			failedTokenizers = append(failedTokenizers, item.Class)
		} else {
			spinner.Success()
			removedTokenizers = append(removedTokenizers, item)
		}
	}
	// update config file
	spinner := app.UI().StartSpinner("Writing tokenizer to configuration file...")
	currentModel.Tokenizers = currentModel.Tokenizers.Difference(removedTokenizers)
	err := AddModels(model.Models{currentModel})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the tokenizer to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
	return nil
}
