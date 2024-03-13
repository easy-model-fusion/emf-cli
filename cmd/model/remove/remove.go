package cmdmodelremove

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const modelRemoveCommandName string = "remove"

// ModelAddCmd represents the add model command
var ModelAddCmd = &cobra.Command{
	Use:   modelRemoveCommandName,
	Short: "Palette that contains add model based commands",
	Long:  "Palette that contains add model based commands",
	Run:   runModelAdd,
}

func init() {
	// Adding subcommands
	//Create new sub command linked to deletion of tokenizers

	//ModelAddCmd.AddCommand(addCustomCmd)
	//add old remove here as subvommand
	//ModelAddCmd.AddCommand(addByNamesCmd)
}

// runModelAdd runs model add command
func runModelAdd(cmd *cobra.Command, args []string) {

	// Running command as palette : allowing user to choose subcommand
	err := cobrautil.RunCommandAsPalette(cmd, args, modelRemoveCommandName, []string{})
	if err != nil {
		pterm.Error.Println("Something went wrong :", err)
	}
}
