package cli

import (
	"fmt"
	"io"
	"os"

	hullerr "github.com/ebogdum/hull/internal/errors"
	"github.com/ebogdum/hull/internal/policy"
	"github.com/spf13/cobra"
)

// newPolicyCommand exposes `hull policy check <package>` for evaluating
// `policies/*.yaml` rules against the rendered manifest of a package.
func newPolicyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Evaluate package policies against rendered manifests",
	}
	cmd.AddCommand(newPolicyCheckCommand())
	cmd.AddCommand(newPolicyListCommand())
	return cmd
}

func newPolicyCheckCommand() *cobra.Command {
	var manifestFile string
	cmd := &cobra.Command{
		Use:   "check <package-path>",
		Short: "Evaluate the package's policies/ rules against a manifest",
		Long: `Read a rendered manifest (from --manifest <file> or stdin) and evaluate
every policy in <package-path>/policies/ against it. Pipe in 'hull template'
output to check before applying:

    hull template ./pkg | hull policy check ./pkg`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packagePath := args[0]
			rules, err := policy.LoadRules(packagePath)
			if nil != err {
				return err
			}
			if 0 == len(rules) {
				fmt.Fprintln(cmd.OutOrStdout(), "no policies/ rules found")
				return nil
			}
			var manifest []byte
			if "" != manifestFile {
				b, readErr := os.ReadFile(manifestFile)
				if nil != readErr {
					return hullerr.WrapError(hullerr.ErrCLIValidation, "read manifest", readErr)
				}
				manifest = b
			} else {
				b, readErr := io.ReadAll(cmd.InOrStdin())
				if nil != readErr {
					return hullerr.WrapError(hullerr.ErrCLIValidation, "read stdin", readErr)
				}
				manifest = b
			}
			if 0 == len(manifest) {
				return hullerr.NewError(hullerr.ErrCLIValidation,
					"empty manifest (provide --manifest or pipe `hull template` to stdin)")
			}
			violations, err := policy.Evaluate(rules, string(manifest))
			if nil != err {
				return err
			}
			if 0 == len(violations) {
				fmt.Fprintf(cmd.OutOrStdout(), "ok — %d rule(s) passed\n", len(rules))
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), policy.FormatHuman(violations))
			if policy.HasDeny(violations) {
				return hullerr.NewError(hullerr.ErrCLIValidation, "policy violations exist")
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&manifestFile, "manifest", "", "rendered manifest file (default: stdin)")
	return cmd
}

func newPolicyListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <package-path>",
		Short: "List policy rules declared in the package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rules, err := policy.LoadRules(args[0])
			if nil != err {
				return err
			}
			for _, r := range rules {
				fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]\n", r.Name, r.Severity)
			}
			return nil
		},
	}
	return cmd
}
