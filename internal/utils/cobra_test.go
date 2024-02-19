package utils

import (
	"github.com/easy-model-fusion/client/test"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/cobra"
	"testing"
)

func prepareSubCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Command " + name,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
}

func prepareRootCmd() *cobra.Command {
	// Init : prepare commands
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	return rootCmd
}

func TestCobraFindSubCommand_NotFound(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd2 := prepareSubCmd("cmd2")
	rootCmd.AddCommand(cmd2)

	// Execute
	_, found := CobraFindSubCommand(rootCmd, "afnwibqpwifubqwpb")

	// Assert
	test.AssertEqual(t, found, false)
}

func TestCobraFindSubCommand_Success(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd2 := prepareSubCmd("cmd2")
	rootCmd.AddCommand(cmd2)

	// Execute
	result, found := CobraFindSubCommand(rootCmd, cmd2.Use)

	// Assert
	test.AssertNotEqual(t, result, nil)
	test.AssertEqual(t, result.Use, cmd2.Use)
	test.AssertEqual(t, found, true)
}

func TestCobraGetSubCommands(t *testing.T) {
	// Init
	rootCmd := prepareRootCmd()
	cmd1 := prepareSubCmd("cmd1")
	cmd2 := prepareSubCmd("cmd2")
	cmd3 := prepareSubCmd("cmd3")
	rootCmd.AddCommand(cmd1)
	rootCmd.AddCommand(cmd2)
	rootCmd.AddCommand(cmd3)

	// Execute
	commandsList, commandsMap := CobraGetSubCommands(rootCmd, []string{cmd2.Use})

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

// TestCobraGetNonProvidedFlags tests the CobraGetNonProvidedFlags function.
func TestCobraGetNonProvidedFlags(t *testing.T) {
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
	flags := CobraGetNonProvidedFlags(cmd1)

	// Assert
	test.AssertEqual(t, len(flags), 1) // skipping help and setting flag1
}
