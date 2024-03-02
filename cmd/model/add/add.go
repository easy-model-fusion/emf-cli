package cmdmodeladd

import (
	"github.com/easy-model-fusion/emf-cli/internal/utils/cobrautil"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const modelAddCommandName string = "add"

// ModelAddCmd represents the add model command
var ModelAddCmd = &cobra.Command{
	Use:   modelAddCommandName,
	Short: "Palette that contains add model based commands",
	Long:  "Palette that contains add model based commands",
	Run:   runModelAdd,
}

func init() {
	// Adding subcommands
	ModelAddCmd.AddCommand(addCustomCmd)
	ModelAddCmd.AddCommand(addByNamesCmd)
}

// runModelAdd runs model add command
func runModelAdd(cmd *cobra.Command, args []string) {

	// Running command as palette : allowing user to choose subcommand
	err := cobrautil.RunCommandAsPalette(cmd, args, modelAddCommandName, []string{})
	if err != nil {
		pterm.Error.Println("Something went wrong :", err)
	}
}
