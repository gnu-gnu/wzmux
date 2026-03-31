package wezterm

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// AgentIcons are emoji prefixes used to identify agent tabs.
const AgentIcons = "⚙🟡✅🔴🤖"

// Pane represents a WezTerm pane from `wezterm cli list --format json`.
type Pane struct {
	PaneID    int    `json:"pane_id"`
	TabID     int    `json:"tab_id"`
	TabTitle  string `json:"tab_title"`
	CWD       string `json:"cwd"`
	IsActive  bool   `json:"is_active"`
	IsFocused bool   `json:"is_zoomed"`
}

// List returns all WezTerm panes.
func List() ([]Pane, error) {
	out, err := exec.Command("wezterm", "cli", "list", "--format", "json").Output()
	if err != nil {
		return nil, fmt.Errorf("wezterm cli list: %w", err)
	}
	var panes []Pane
	if err := json.Unmarshal(out, &panes); err != nil {
		return nil, fmt.Errorf("parse pane list: %w", err)
	}
	return panes, nil
}

// isAgentTab checks if a tab title starts with an agent icon.
func isAgentTab(title string) bool {
	for _, r := range title {
		return strings.ContainsRune(AgentIcons, r)
	}
	return false
}

// AgentPanes returns panes that are agent tabs (filtered and deduplicated by tab_id).
func AgentPanes() ([]Pane, error) {
	all, err := List()
	if err != nil {
		return nil, err
	}
	seen := map[int]bool{}
	var agents []Pane
	for _, p := range all {
		if isAgentTab(p.TabTitle) && !seen[p.TabID] {
			seen[p.TabID] = true
			agents = append(agents, p)
		}
	}
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].TabTitle < agents[j].TabTitle
	})
	return agents, nil
}

// AgentName extracts the agent name by stripping the leading emoji and space.
func AgentName(p Pane) string {
	title := p.TabTitle
	// Strip first rune (emoji) and any following space
	for i, r := range title {
		if i == 0 {
			_ = r
			continue
		}
		rest := title[i:]
		return strings.TrimLeft(rest, " ")
	}
	return title
}

// NormalizeCWD strips file:// scheme and hostname from WezTerm CWD paths.
func NormalizeCWD(cwd string) string {
	if strings.HasPrefix(cwd, "file://") {
		cwd = cwd[7:]
		// strip hostname portion: file://hostname/path → /path
		if idx := strings.Index(cwd, "/"); idx >= 0 {
			cwd = cwd[idx:]
		}
	}
	return cwd
}

// Spawn creates a new WezTerm tab and returns its pane ID.
func Spawn(cwd string, args ...string) (int, error) {
	cmdArgs := []string{"cli", "spawn", "--cwd", cwd}
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, "--")
		cmdArgs = append(cmdArgs, args...)
	}
	out, err := exec.Command("wezterm", cmdArgs...).Output()
	if err != nil {
		return 0, fmt.Errorf("wezterm cli spawn: %w", err)
	}
	var paneID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &paneID); err != nil {
		return 0, fmt.Errorf("parse pane id: %w", err)
	}
	return paneID, nil
}

// SplitPane creates a right split pane with given percentage.
func SplitPane(parentPaneID int, percent int, args ...string) (int, error) {
	cmdArgs := []string{
		"cli", "split-pane",
		"--right",
		"--percent", fmt.Sprintf("%d", percent),
		"--pane-id", fmt.Sprintf("%d", parentPaneID),
	}
	if len(args) > 0 {
		cmdArgs = append(cmdArgs, "--")
		cmdArgs = append(cmdArgs, args...)
	}
	out, err := exec.Command("wezterm", cmdArgs...).Output()
	if err != nil {
		return 0, fmt.Errorf("wezterm cli split-pane: %w", err)
	}
	var paneID int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &paneID); err != nil {
		return 0, fmt.Errorf("parse split pane id: %w", err)
	}
	return paneID, nil
}

// SetTabTitle sets the tab title for a given pane.
func SetTabTitle(paneID int, title string) error {
	return exec.Command("wezterm", "cli", "set-tab-title", title, "--pane-id", fmt.Sprintf("%d", paneID)).Run()
}

// ActivatePane focuses the given pane.
func ActivatePane(paneID int) error {
	return exec.Command("wezterm", "cli", "activate-pane", "--pane-id", fmt.Sprintf("%d", paneID)).Run()
}

// KillPane kills the given pane.
func KillPane(paneID int) error {
	return exec.Command("wezterm", "cli", "kill-pane", "--pane-id", fmt.Sprintf("%d", paneID)).Run()
}
