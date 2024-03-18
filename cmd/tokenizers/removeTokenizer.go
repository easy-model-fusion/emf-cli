package tokenizers

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/model"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/easy-model-fusion/emf-cli/internal/ui"
	"github.com/easy-model-fusion/emf-cli/internal/utils/stringutil"
	"github.com/easy-model-fusion/emf-cli/pkg/huggingface"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

// tokenizerRemoveCmd represents the model remove command
var tokenizerRemoveCmd = &cobra.Command{
	Use:   "remove tokenizers <model_name> [tokenizers..]",
	Short: "Remove one or more tokenizers",
	Long:  "Remove one or more tokenizers",
	Run:   runTokenizerRemove,
}

// runTokenizerRemove runs the model remove command
func runTokenizerRemove(cmd *cobra.Command, argsModel []string, argsTokenizer []string) {
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}
	sdk.SendUpdateSuggestion()

	configModels, err := config.GetModels()
	if err != nil {
		pterm.Error.Println(err.Error())
		return
	}

	// Get models to update
	var selectedModelNames []string
	if len(argsModel) == 0 {
		// No argument provided
		selectedModelNames = selectMode(configModels)
	} else {
		// Remove all the duplicates
		selectedModelNames = stringutil.SliceRemoveDuplicates(argsModel)
	}

	modelsToUse := model.GetModelsByNames(configModels, selectedModelNames)

	var tokenizerNames []string
	if modelsToUse.Module == huggingface.TRANSFORMERS {
		availableNames := modelsToUse.Tokenizers.GetNames()
		if len(availableNames) > 0 && len(argsTokenizer) == 0 {
			message := "Please select the tokenizer(s) to be deleted"
			checkMark := ui.Checkmark{Checked: pterm.Green("+"), Unchecked: pterm.Red("-")}
			tokenizerNames = app.UI().DisplayInteractiveMultiselect(message, availableNames, checkMark, true, true)
			app.UI().DisplaySelectedItems(tokenizerNames)
		} else if len(argsTokenizer) > 0 {
			// Check if selectedTokenizerNames elements exist in tokenizerNames and add them to a new list
			var selectedAndAvailableTokenizerNames []string
			for _, name := range argsTokenizer {
				for _, availableName := range tokenizerNames {
					if name == availableName {
						selectedAndAvailableTokenizerNames = append(selectedAndAvailableTokenizerNames, name)
						break
					}
				}
			}
			tokenizerNames = selectedAndAvailableTokenizerNames
		}
	}

	if len(tokenizerNames) > 0 {

		for _, tokenizerName := range tokenizerNames {
			tokenizer := modelsToUse.Tokenizers[tokenizerName]
			err := os.Remove(tokenizer.Path)
			if err != nil {
				pterm.Error.Println(err.Error())
				return
			}
			model.RemoveTokenizer(tokenizer)
		}
	}
}

func selectMode(currentModels []model.Model) []string {
	// Build a multiselect with each model name
	var modelNames []string
	for _, item := range currentModels {
		modelNames = append(modelNames, item.Name)
	}

	checkMark := ui.Checkmark{Checked: pterm.Red("x"), Unchecked: pterm.Blue("-")}
	message := "Please select a model"
	modelsToChoice := app.UI().DisplayInteractiveMultiselect(message, modelNames, checkMark, false, false)
	app.UI().DisplaySelectedItems(modelsToChoice)
	return modelsToChoice
}

func (m *Model) RemoveTokenizer(tokenizer Tokenizer) {
	for i, t := range m.Tokenizers {
		if t.Path == tokenizer.Path && t.Class == tokenizer.Class {
			m.Tokenizers = append(m.Tokenizers[:i], m.Tokenizers[i+1:]...)
			return
		}
	}
}
