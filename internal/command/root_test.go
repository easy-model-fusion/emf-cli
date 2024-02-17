package command

import (
	"github.com/magiconair/properties/assert"
	"github.com/spf13/cobra"
	"testing"
)

func TestGetAllCommands(t *testing.T) {
	// Init : prepare commands
	rootCmd := &cobra.Command{
		Use:   "root",
		Short: "Root command",
	}
	// Add child commands to the root command
	cmd1 := &cobra.Command{
		Use:   "cmd1",
		Short: "Command 1",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	childCmd1 := &cobra.Command{
		Use:   "childCm1",
		Short: "Child command 1",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	cmd1.AddCommand(childCmd1)
	rootCmd.AddCommand(cmd1)
	cmd2 := &cobra.Command{
		Use:   "cmd2",
		Short: "Command 2",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(cmd2)

	// Execute
	commandsList, commandsMap := getAllCommands(rootCmd, []string{}, map[string]func(*cobra.Command, []string){})

	// Assert
	expectedList := []string{"cmd1", "childCm1", "cmd2"}
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

func TestHideCommands(t *testing.T) {
	// Init
	command1 := "command1"
	command2 := "command2"
	command3 := "command3"
	commandsList := []string{command1, command2, command3}
	commandsMap := map[string]func(*cobra.Command, []string){
		command1: func(cmd *cobra.Command, args []string) {},
		command2: func(cmd *cobra.Command, args []string) {},
		command3: func(cmd *cobra.Command, args []string) {},
	}

	// Execute
	commandsToHide := []string{command1, command3}
	commandsList, commandsMap = hideCommands(commandsList, commandsMap, commandsToHide)

	// Assert
	expectedList := []string{command2}
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
