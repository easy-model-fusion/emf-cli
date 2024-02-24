package cmdmodel

import (
	"github.com/easy-model-fusion/emf-cli/internal/config"
	"github.com/easy-model-fusion/emf-cli/internal/sdk"
	"github.com/spf13/cobra"
)

// modelUpdateCmd represents the model update command
var modelUpdateCmd = &cobra.Command{
	Use:   "update <model name> [<other model names>...]",
	Short: "Update one or more models",
	Long:  "Update one or more models",
	Run:   runModelUpdate,
}

// runModelUpdate runs the model update command
func runModelUpdate(cmd *cobra.Command, args []string) {
	if config.GetViperConfig(config.FilePath) != nil {
		return
	}

	sdk.SendUpdateSuggestion()

	// TODO : if args : get config model by names
	// TODO : else : get all models from config

	// TODO : for each config model by name
	// TODO : check if downloaded etc => HOW TO HANDLE THOSE NOT DOWNLOADED?
	// TODO : if source HF : call to HF and see if a new version is available
	// TODO : yes? then offer to overwrite the model
	// TODO : yes again? then download
	// TODO : if any updates : update the configuration file

}
