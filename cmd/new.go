package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gnu-gnu/wzmux/internal/session"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name> [prompt...]",
	Short: "Launch a new Claude Code agent in a WezTerm tab",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runNew,
}

// uniqueName appends -1, -2, ... if the name already exists among active agents.
func uniqueName(base string) string {
	agents, err := wezterm.AgentPanes()
	if err != nil {
		return base
	}
	names := map[string]bool{}
	for _, a := range agents {
		names[wezterm.AgentName(a)] = true
	}
	if !names[base] {
		return base
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s-%d", base, i)
		if !names[candidate] {
			return candidate
		}
	}
}

func runNew(cmd *cobra.Command, args []string) error {
	name := uniqueName(args[0])
	sessionID := "wzmux-" + name
	cwd, _ := os.Getwd()

	// Build claude command
	claudeArgs := []string{"claude", "--name", sessionID}
	if len(args) > 1 {
		prompt := strings.Join(args[1:], " ")
		claudeArgs = append(claudeArgs, "-p", prompt)
	}

	// Spawn new WezTerm tab
	paneID, err := wezterm.Spawn(cwd, claudeArgs...)
	if err != nil {
		return fmt.Errorf("failed to spawn agent: %w", err)
	}

	// Set tab title with robot emoji
	wezterm.SetTabTitle(paneID, "🤖 "+name)

	// Try to spawn dashboard in a right split, watching the agent pane
	self, err := os.Executable()
	if err == nil {
		wezterm.SplitPane(paneID, 25, self, "dashboard", "--watch-pane", fmt.Sprintf("%d", paneID))
	}

	// Record session for later resume
	session.Save(session.Entry{
		Name:      name,
		SessionID: sessionID,
		CWD:       cwd,
		CreatedAt: time.Now(),
	})

	fmt.Printf("Agent '%s' launched (session: %s, pane %d)\n", name, sessionID, paneID)
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
