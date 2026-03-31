package cmd

import (
	"fmt"

	"github.com/gnu-gnu/wzmux/internal/config"
	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove wzmux hooks from Claude Code settings and clean up",
	RunE:  runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) error {
	if err := config.RemoveHooks(); err != nil {
		return fmt.Errorf("failed to remove hooks: %w", err)
	}

	status.CleanAll()

	fmt.Println("wzmux hooks removed from", config.SettingsPath())
	fmt.Println("Status files cleaned up.")

	return nil
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
