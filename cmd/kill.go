package cmd

import (
	"fmt"
	"strings"

	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <name|--all>",
	Short: "Kill Claude Code agent(s)",
	RunE:  runKill,
}

var killAll bool

func runKill(cmd *cobra.Command, args []string) error {
	agents, err := wezterm.AgentPanes()
	if err != nil {
		return err
	}

	if killAll {
		for _, p := range agents {
			wezterm.KillPane(p.PaneID)
			status.Remove(p.PaneID)
		}
		fmt.Printf("Killed %d agent(s).\n", len(agents))
		return nil
	}

	if len(args) == 0 {
		return fmt.Errorf("specify agent name or use --all")
	}

	target := args[0]
	killed := 0
	for _, p := range agents {
		name := wezterm.AgentName(p)
		if strings.Contains(strings.ToLower(name), strings.ToLower(target)) {
			wezterm.KillPane(p.PaneID)
			status.Remove(p.PaneID)
			fmt.Printf("Killed agent '%s' (pane %d)\n", name, p.PaneID)
			killed++
		}
	}

	if killed == 0 {
		return fmt.Errorf("no agent matching '%s' found", target)
	}
	return nil
}

func init() {
	killCmd.Flags().BoolVarP(&killAll, "all", "a", false, "Kill all agents")
	rootCmd.AddCommand(killCmd)
}
