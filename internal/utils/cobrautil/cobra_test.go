package cobrautil

import (
	"github.com/easy-model-fusion/emf-cli/internal/app"
	"github.com/easy-model-fusion/emf-cli/test"
	"github.com/easy-model-fusion/emf-cli/test/mock"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/cobra"
	"testing"
)

// prepareSubCmd creates a mock sub-command for testing purposes.
func prepareSubCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Command " + name,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

// prepareRootCmd creates a mock root command for testing purposes.
func prepareRootCmd() *cobra.Command {
	// Init : prepare commands
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	return rootCmd
}

// TestFindSubCommand_NotFound tests the FindSubCommand function when the sub-command is not found.
func TestFindSubCommand_NotFound(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()

	// Execute
	_, found := FindSubCommand(rootCmd, "afnwibqpwifubqwpb")

	// Assert
	test.AssertEqual(t, found, false)
}

// TestFindSubCommand_FromParentSuccess tests the FindSubCommand function when the sub-command is found under the parent command.
func TestFindSubCommand_FromParentSuccess(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd2 := prepareSubCmd("cmd2")
	rootCmd.AddCommand(cmd2)

	// Execute
	resultCmd, found := FindSubCommand(rootCmd, cmd2.Name())

	// Assert
	test.AssertEqual(t, found, true)
	test.AssertEqual(t, resultCmd.Name(), cmd2.Name())
}

// TestFindSubCommand_AsItselfSuccess tests the FindSubCommand function when the sub-command is found as itself.
func TestFindSubCommand_AsItselfSuccess(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd2 := prepareSubCmd("cmd2")
	rootCmd.AddCommand(cmd2)

	// Execute
	resultCmd, found := FindSubCommand(cmd2, cmd2.Name())

	// Assert
	test.AssertEqual(t, found, true)
	test.AssertEqual(t, resultCmd.Name(), cmd2.Name())
}

// TestGetSubCommands_Success tests the GetSubCommands function for retrieving sub-commands.
func TestGetSubCommands_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	cmd2 := prepareSubCmd("cmd2")
	cmd3 := prepareSubCmd("cmd3")
	rootCmd.AddCommand(cmd1)
	rootCmd.AddCommand(cmd2)
	rootCmd.AddCommand(cmd3)

	// Execute
	commandsList, commandsMap := GetSubCommands(rootCmd, []string{cmd2.Use})

	// Assert
	expectedList := []string{cmd1.Use, cmd3.Use}
	assert.Equal(t, len(commandsList), len(expectedList))
	for i, item := range commandsList {
		assert.Equal(t, item, expectedList[i])
	}
	for _, item := range expectedList {
		if _, ok := commandsMap[item]; !ok {
			t.Fail()
		}
	}
}

// TestGetNonProvidedLocalFlags_Success tests the GetNonProvidedLocalFlags function for retrieving non-provided flags.
func TestGetNonProvidedLocalFlags_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Add flags to the command
	cmd1.Flags().Bool("flag1", false, "Test flag 1")
	cmd1.Flags().String("flag2", "", "Test flag 2")
	cmd1.Flags().Bool("help", false, "Help flag")

	// Mark some flags as changed
	err := cmd1.Flags().Set("flag1", "true")
	if err != nil {
		t.Fail()
	}

	// Execute
	flags := GetNonProvidedLocalFlags(cmd1)

	// Assert
	test.AssertEqual(t, len(flags), 1) // skipping help and setting flag1
}

// TestMultiselectSubcommands_Success tests the MultiselectSubcommands function for selecting sub-commands.
func TestMultiselectSubcommands_Success(t *testing.T) {
	// Init
	ui := mock.MockUI{}
	ui.SelectResult = "cmd1"
	app.SetUI(ui)
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	cmd2 := prepareSubCmd("cmd2")
	cmd3 := prepareSubCmd("cmd3")
	rootCmd.AddCommand(cmd1)
	rootCmd.AddCommand(cmd2)
	rootCmd.AddCommand(cmd3)

	// Execute
	MultiselectSubcommands(rootCmd, []string{}, []string{cmd1.Use, cmd2.Use, cmd3.Use}, map[string]func(*cobra.Command, []string){cmd1.Use: cmd1.Run, cmd2.Use: cmd2.Run, cmd3.Use: cmd3.Run})
}

// TestMultiselectSubcommands_Fail tests the MultiselectSubcommands function for selecting sub-commands when the selected command is not recognized.
func TestMultiselectSubcommands_Fail(t *testing.T) {
	// Init
	ui := mock.MockUI{}
	ui.SelectResult = "cmd4"
	app.SetUI(ui)
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	cmd2 := prepareSubCmd("cmd2")
	cmd3 := prepareSubCmd("cmd3")
	rootCmd.AddCommand(cmd1)
	rootCmd.AddCommand(cmd2)
	rootCmd.AddCommand(cmd3)

	// Execute
	MultiselectSubcommands(rootCmd, []string{}, []string{cmd1.Use, cmd2.Use, cmd3.Use}, map[string]func(*cobra.Command, []string){cmd1.Use: cmd1.Run, cmd2.Use: cmd2.Run, cmd3.Use: cmd3.Run})
}

// TestMultiselectRemainingFlags_Success tests the MultiselectRemainingFlags function for selecting remaining flags.
func TestMultiselectRemainingFlags_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Add flags to the command
	cmd1.Flags().Bool("flag1", false, "Test flag 1")
	cmd1.Flags().String("flag2", "", "Test flag 2")
	cmd1.Flags().Bool("help", false, "Help flag")

	// Execute
	flags, _ := MultiselectRemainingFlags(cmd1)

	// Assert
	test.AssertEqual(t, len(flags), 2) // skipping help
}

// TestMultiselectRemainingFlags_NoFlags tests the MultiselectRemainingFlags function for selecting remaining flags when there are no flags.
func TestMultiselectRemainingFlags_NoFlags(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Execute
	flags, _ := MultiselectRemainingFlags(cmd1)

	// Assert
	test.AssertEqual(t, len(flags), 0)
}

// TestMultiselectRemainingFlags_AllFlags tests the MultiselectRemainingFlags function for selecting remaining flags when all flags are provided.
func TestMultiselectRemainingFlags_AllFlags(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Add flags to the command
	cmd1.Flags().Bool("flag1", false, "Test flag 1")
	cmd1.Flags().String("flag2", "", "Test flag 2")
	cmd1.Flags().Bool("help", false, "Help flag")

	// Mark all flags as changed
	err := cmd1.Flags().Set("flag1", "true")
	if err != nil {
		t.Fail()
	}
	err = cmd1.Flags().Set("flag2", "value")
	if err != nil {
		t.Fail()
	}

	// Execute
	flags, _ := MultiselectRemainingFlags(cmd1)

	// Assert
	test.AssertEqual(t, len(flags), 0)
}

// TestAskFlagInput_Success tests the AskFlagInput function for asking input for a flag.
func TestAskFlagInput_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Add flags to the command
	cmd1.Flags().Bool("flag1", false, "Test flag 1")
	cmd1.Flags().String("flag2", "", "Test flag 2")
	cmd1.Flags().Bool("help", false, "Help flag")

	// Execute
	err := AskFlagInput(cmd1, cmd1.Flags().Lookup("flag1"))

	// Assert
	test.AssertEqual(t, err, nil)
}

// TestRunCommandAsPalette_Success tests the RunCommandAsPalette function when the command is found and successfully run.
func TestRunCommandAsPalette_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	rootCmd.AddCommand(cmd1)

	// Execute
	err := RunCommandAsPalette(rootCmd, []string{}, cmd1.Name(), []string{})

	// Assert

	test.AssertEqual(t, err, nil)
}

// TestRunCommandAsPalette_CommandNotFound tests the RunCommandAsPalette function when the command is not found.
func TestRunCommandAsPalette_CommandNotFound(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()

	// Execute
	err := RunCommandAsPalette(rootCmd, []string{}, "nonexistent", []string{})

	// Assert
	test.AssertNotEqual(t, err, nil)
}
