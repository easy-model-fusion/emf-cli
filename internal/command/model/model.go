package commandmodel

import (
	"fmt"
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/internal/command/model/add"
	"github.com/easy-model-fusion/emf-cli/internal/huggingface"
	"github.com/easy-model-fusion/emf-cli/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const modelCommandName string = "model"

// ModelCmd represents the model command
var ModelCmd = &cobra.Command{
	Use:   modelCommandName,
	Short: "Palette that contains model based commands",
	Long:  "Palette that contains model based commands",
	Run:   runModel,
}

// runModel runs model command
func runModel(cmd *cobra.Command, args []string) {

	// Searching for the currentCmd : when 'cmd' differs from 'addCmd' (i.e. run through parent multiselect)
	currentCmd, found := utils.CobraFindSubCommand(cmd, modelCommandName)
	if !found {
		pterm.Error.Println(fmt.Sprintf("Something went wrong : the '%s' command was not found. Please try again.", modelCommandName))
		return
	}

	// Retrieve all the subcommands of the current command
	commandsList, commandsMap := utils.CobraGetSubCommands(currentCmd, []string{})

	// Users chooses a command and runs it automatically
	utils.CobraSelectCommandToRun(currentCmd, args, commandsList, commandsMap)
}

func init() {
	// Preparing to use the hugging face API
	app.InitHuggingFace(huggingface.BaseUrl, "")

	// Adding the subcommands
	ModelCmd.AddCommand(modelRemoveCmd)
	ModelCmd.AddCommand(modelTidyCmd)
	ModelCmd.AddCommand(add.ModelAddCmd)
}
