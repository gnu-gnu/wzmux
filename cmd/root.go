package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.3.0"

var rootCmd = &cobra.Command{
	Use:   "wzmux",
	Short: "WezTerm multiplexer for Claude Code agents",
	Long:  "wzmux manages multiple Claude Code agents via WezTerm tabs with real-time status monitoring.",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print wzmux version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("wzmux %s\n", version)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
