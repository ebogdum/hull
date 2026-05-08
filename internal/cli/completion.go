package cli

import (
	"fmt"
	"os"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"github.com/spf13/cobra"
)

func newCompletionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for hull.

To load completions:

Bash:
  $ source <(hull completion bash)

Zsh:
  $ source <(hull completion zsh)

Fish:
  $ hull completion fish | source

PowerShell:
  PS> hull completion powershell | Out-String | Invoke-Expression
`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()

			handlers := map[string]func() error{
				"bash":       func() error { return root.GenBashCompletion(os.Stdout) },
				"zsh":        func() error { return root.GenZshCompletion(os.Stdout) },
				"fish":       func() error { return root.GenFishCompletion(os.Stdout, true) },
				"powershell": func() error { return root.GenPowerShellCompletionWithDesc(os.Stdout) },
			}

			handler, ok := handlers[args[0]]
			if !ok {
				return hullerr.CLIError(hullerr.ErrCLIValidation, fmt.Sprintf("unsupported shell: %s", args[0]))
			}
			return handler()
		},
	}
	return cmd
}
