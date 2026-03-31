package cmd

import (
	"fmt"
	"text/tabwriter"
	"os"

	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List active Claude Code agents",
	RunE:    runLs,
}

func runLs(cmd *cobra.Command, args []string) error {
	agents, err := wezterm.AgentPanes()
	if err != nil {
		return err
	}

	if len(agents) == 0 {
		fmt.Println("No active agents.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "#\tNAME\tSTATUS\tPANE\tCWD")

	for i, p := range agents {
		name := wezterm.AgentName(p)
		cwd := wezterm.NormalizeCWD(p.CWD)

		st := "unknown"
		if s, err := status.Read(p.PaneID); err == nil {
			st = s.Status
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%s\n", i+1, name, st, p.PaneID, cwd)
	}
	w.Flush()
	return nil
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
