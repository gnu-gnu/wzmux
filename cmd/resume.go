package cmd

import (
	"fmt"
	"os"

	"github.com/gnu-gnu/wzmux/internal/session"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume <name>",
	Short: "Resume a previous Claude Code session by name",
	Args:  cobra.ExactArgs(1),
	RunE:  runResume,
}

func runResume(cmd *cobra.Command, args []string) error {
	name := args[0]
	sessionID := "wzmux-" + name
	cwd, _ := os.Getwd()

	// Use stored CWD from previous session if available
	if entry, err := session.Load(name); err == nil && entry.CWD != "" {
		cwd = entry.CWD
	}

	claudeArgs := []string{"claude", "--resume", sessionID}

	paneID, err := wezterm.Spawn(cwd, claudeArgs...)
	if err != nil {
		return fmt.Errorf("failed to spawn agent: %w", err)
	}

	wezterm.SetTabTitle(paneID, "🤖 "+name)

	self, err := os.Executable()
	if err == nil {
		wezterm.SplitPane(paneID, 25, self, "dashboard", "--watch-pane", fmt.Sprintf("%d", paneID))
	}

	fmt.Printf("Agent '%s' resumed (session: %s, pane %d)\n", name, sessionID, paneID)
	return nil
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
