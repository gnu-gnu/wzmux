package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:    "hook <event>",
	Short:  "Claude Code hook handler (called by Claude Code, not by users)",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE:   runHook,
}

// hookInput represents the JSON data Claude Code passes via stdin.
type hookInput struct {
	LastAssistantMessage string `json:"last_assistant_message"`
	Message              string `json:"message"`
	Title                string `json:"title"`
}

func runHook(cmd *cobra.Command, args []string) error {
	event := args[0]

	// Guard: only run inside WezTerm
	paneStr := os.Getenv("WEZTERM_PANE")
	if paneStr == "" {
		return nil
	}
	paneID, err := strconv.Atoi(paneStr)
	if err != nil {
		return nil
	}

	cwd, _ := os.Getwd()

	// Read stdin (Claude Code hook JSON)
	stdinData, _ := io.ReadAll(os.Stdin)
	var input hookInput
	json.Unmarshal(stdinData, &input) // best-effort parse

	dirName := filepath.Base(cwd)

	switch event {
	case "PreToolUse":
		status.Write(paneID, "running", cwd, "")
		wezterm.SetTabTitle(paneID, "⚙ "+dirName)

	case "PostToolUse":
		status.Write(paneID, "done", cwd, "")
		wezterm.SetTabTitle(paneID, "✅ "+dirName)

	case "Stop", "Notification":
		msg := sanitizeMsg(input.LastAssistantMessage, 80)
		if msg == "" {
			msg = "Done"
		}
		status.Write(paneID, "waiting", cwd, msg)
		wezterm.SetTabTitle(paneID, "🟡 "+dirName)

	case "Error":
		status.Write(paneID, "error", cwd, input.Message)
		wezterm.SetTabTitle(paneID, "🔴 "+dirName)

	default:
		return fmt.Errorf("unknown hook event: %s", event)
	}

	return nil
}

// sanitizeMsg cleans and truncates a message for display.
func sanitizeMsg(s string, maxLen int) string {
	// Remove newlines
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", "")
	// Truncate to maxLen runes
	if utf8.RuneCountInString(s) > maxLen {
		runes := []rune(s)
		s = string(runes[:maxLen])
	}
	return s
}

func init() {
	rootCmd.AddCommand(hookCmd)
}
