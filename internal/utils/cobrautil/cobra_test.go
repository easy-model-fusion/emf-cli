package cobrautil

import (
	"github.com/easy-model-fusion/emf-cli/test"
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
