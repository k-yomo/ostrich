package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const flagNamePath = "path"

var (
	version string

	indexPath = ""
)

// Exec executes command
func Exec() {
	rootCmd := newRootCmd(os.Stdin)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Printf("[error] %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd(out io.Writer) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "ostrich",
		Short:   "CLI for Ostrich",
		Long:    "CLI tool to interact with Ostrich index",
		Version: version,
	}
	rootCmd.SetOut(out)

	rootCmd.PersistentFlags().Bool("help", false, fmt.Sprintf("help for %s", rootCmd.Name()))
	rootCmd.PersistentFlags().StringVarP(&indexPath, flagNamePath, "p", ".", "index directory path")
	rootCmd.AddCommand(
		newSearchCmd(out),
	)
	return rootCmd
}
