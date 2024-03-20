package config

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/model"
)

// RemoveTokenizersByName removes specified tokenizers
func RemoveTokenizersByName(currentModel model.Model, tokenizersToRemove model.Tokenizers) (failedTokenizers []string, err error) {
	var removedTokenizers model.Tokenizers
	// Trying to remove the tokenizers
	for _, item := range tokenizersToRemove {
		if item.Path != "" {
			// Starting client spinner animation
			spinner := app.UI().StartSpinner(fmt.Sprintf("Removing tokenizer %s...", tokenizersToRemove))
			err := RemoveItemPhysically(item.Path)
			if err != nil {
				spinner.Fail("failed to remove tokenizers")
				failedTokenizers = append(failedTokenizers, item.Class)
				continue
			} else {
				spinner.Success()
			}
		}
		removedTokenizers = append(removedTokenizers, item)
	}
	// update config file
	spinner := app.UI().StartSpinner("Writing tokenizer to configuration file...")
	currentModel.Tokenizers = currentModel.Tokenizers.Difference(removedTokenizers)
	err = AddModels(model.Models{currentModel})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Error while writing the tokenizer to the configuration file: %s", err))
	} else {
		spinner.Success()
	}
	return failedTokenizers, err
}
