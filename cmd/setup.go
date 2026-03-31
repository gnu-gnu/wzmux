package cmd

import (
	"fmt"

	"github.com/gnu-gnu/wzmux/internal/config"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Register wzmux hooks in Claude Code settings",
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	if err := config.AddHooks(); err != nil {
		return fmt.Errorf("failed to register hooks: %w", err)
	}

	fmt.Println("wzmux hooks registered in", config.SettingsPath())
	fmt.Println()
	fmt.Println("Registered events: PreToolUse, PostToolUse, Stop, Notification")
	fmt.Println()
	fmt.Println("Optional: Add WezTerm status bar integration to your wezterm.lua.")
	fmt.Println("See: https://github.com/gnu-gnu/wzmux#wezterm-integration")

	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
