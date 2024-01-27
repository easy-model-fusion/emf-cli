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

	var selectedShell string

	// No args, asking for a shell input
	if len(args) == 0 {
		selectedShell = utils.AskForUsersInput("Enter a shell name " + arguments)
	} else {
		selectedShell = args[0]
	}

	// Checks whether the input shell is handled
	if len(selectedShell) == 0 {
		pterm.Error.Println(fmt.Sprintf("Please provide a shell. Expected %s", arguments))
	} else if utils.ArrayStringContainsItem(shells, selectedShell) {
		switch selectedShell {
		case "bash":
			err := cmd.Root().GenBashCompletion(os.Stdout)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Error generating script : %s", err))
				return
			}
		case "zsh":
			err := cmd.Root().GenZshCompletion(os.Stdout)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Error generating script : %s", err))
				return
			}
		case "fish":
			err := cmd.Root().GenFishCompletion(os.Stdout, true)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Error generating script : %s", err))
				return
			}
		case "powershell":
			err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Error generating script : %s", err))
				return
			}
		}
	} else {
		pterm.Error.Println(fmt.Sprintf("Shell '%s' not recognized. Expected %s", selectedShell, arguments))
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(completionCmd)
}
