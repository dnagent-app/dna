package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Run executes the CLI with the given arguments.
func Run(args []string) error {
	root := newRootCmd()
	root.SetArgs(args[1:])
	return root.Execute()
}

// Fatal prints the error to stderr and exits with code 1.
func Fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dna",
		Short: "Domain Name Agent — manage domains from the command line",
		Long:  "dna is a CLI for searching, listing, and managing domain names via NameCheap.",
		// Don't print usage on RunE errors.
		SilenceUsage: true,
	}
	cmd.AddCommand(domainCmd())
	return cmd
}
