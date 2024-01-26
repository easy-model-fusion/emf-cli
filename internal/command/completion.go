package command

import (
	"fmt"
	"github.com/easy-model-fusion/client/internal/app"
	"github.com/easy-model-fusion/client/internal/utils"
	"github.com/pterm/pterm"
	"os"

	"github.com/spf13/cobra"
)

const completionUse string = "completion"

var shells = []string{"bash", "zsh", "fish", "powershell"}

var arguments = utils.ArrayStringAsArguments(shells)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   completionUse,
	Short: "Generate completion script",
	Long: fmt.Sprintf(`To load completions:

Bash:

  $ source <(%[1]s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
  # macOS:
  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ %[1]s completion fish | source

  # To load completions for each session, execute once:
  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

PowerShell:

  PS> %[1]s completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> %[1]s completion powershell > %[1]s.ps1
  # and source this file from your PowerShell profile.
`, app.Name),
	DisableFlagsInUseLine: true,
	Run:                   runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) {

	logger := app.L().WithTime(false)

	var selectedShell string

	// No args, asking for a shell input
	if len(args) == 0 {
		selectedShell = askForShell()
	} else {
		selectedShell = args[0]
	}

	// Checks whether the input shell is handled
	if len(selectedShell) == 0 {
		logger.Error(fmt.Sprintf("Please provide a shell. Expected " + arguments))
	} else if utils.ArrayStringContainsItem(shells, selectedShell) {
		switch selectedShell {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	} else {
		logger.Error(fmt.Sprintf("Shell '%s' not recognized. Expected "+arguments, selectedShell))
	}
}

// askForShell asks the user for a shell name and returns it
func askForShell() string {
	// Create an interactive text input with single line input mode
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(false)

	// Show the text input and get the result
	result, _ := textInput.Show("Enter a shell name " + arguments)

	// Print a blank line for better readability
	pterm.Println()

	return result
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
}
