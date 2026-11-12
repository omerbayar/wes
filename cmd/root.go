// Package cmd defines the wes command-line interface.
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/omerbayar/wes/internal/tui"
)

// version is overridden at build time via -ldflags "-X github.com/omerbayar/wes/cmd.version=...".
var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "wes",
	Short:   "Your daily shortcut assistant",
	Long:    "wes — your daily shortcut assistant. Quick, quiet, efficient, like Wes Montgomery.",
	Version: version,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run(".")
	},
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
