package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/gnu-gnu/wzmux/internal/session"
	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List active and inactive Claude Code agents",
	RunE:    runLs,
}

func runLs(cmd *cobra.Command, args []string) error {
	agents, err := wezterm.AgentPanes()
	if err != nil {
		return err
	}

	// Collect active agent names
	activeNames := map[string]bool{}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "#\tNAME\tSTATUS\tPANE\tCWD")

	idx := 0
	for _, p := range agents {
		name := wezterm.AgentName(p)
		activeNames[name] = true
		cwd := wezterm.NormalizeCWD(p.CWD)

		st := "unknown"
		if s, err := status.Read(p.PaneID); err == nil {
			st = s.Status
		}

		idx++
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n", idx, name, st, p.PaneID, cwd)
	}

	// Show inactive (previously recorded) sessions
	sessions, _ := session.List()
	for _, s := range sessions {
		if activeNames[s.Name] {
			continue
		}
		age := formatAge(s.CreatedAt)
		idx++
		fmt.Fprintf(w, "%d\t%s\t%s\t-\t%s\n", idx, s.Name, "exited ("+age+")", s.CWD)
	}

	if idx == 0 {
		fmt.Println("No agents.")
		return nil
	}

	w.Flush()
	return nil
}

func formatAge(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
