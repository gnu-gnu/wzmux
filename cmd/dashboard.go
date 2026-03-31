package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gnu-gnu/wzmux/internal/status"
	"github.com/gnu-gnu/wzmux/internal/wezterm"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"dash"},
	Short:   "Live dashboard for monitoring Claude Code agents",
	RunE:    runDashboard,
}

// -- Styles --

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // green
	waitingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // yellow
	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // cyan
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // red
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))  // dim
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// -- Agent data --

type agent struct {
	paneID int
	name   string
	status string
	cwd    string
	msg    string
}

func statusIcon(s string) string {
	switch s {
	case "running":
		return "⚙"
	case "waiting":
		return "🟡"
	case "done":
		return "✅"
	case "error":
		return "🔴"
	default:
		return "🤖"
	}
}

func statusStyle(s string) lipgloss.Style {
	switch s {
	case "running":
		return runningStyle
	case "waiting":
		return waitingStyle
	case "done":
		return doneStyle
	case "error":
		return errorStyle
	default:
		return dimStyle
	}
}

func statusOrder(s string) int {
	switch s {
	case "waiting":
		return 0
	case "running":
		return 1
	case "error":
		return 2
	case "done":
		return 3
	default:
		return 4
	}
}

// -- Bubble Tea model --

type tickMsg time.Time

type model struct {
	agents []agent
	width  int
	height int
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func collectAgents() []agent {
	panes, err := wezterm.AgentPanes()
	if err != nil {
		return nil
	}

	var agents []agent
	for _, p := range panes {
		a := agent{
			paneID: p.PaneID,
			name:   wezterm.AgentName(p),
			cwd:    wezterm.NormalizeCWD(p.CWD),
			status: "unknown",
		}
		if s, err := status.Read(p.PaneID); err == nil {
			a.status = s.Status
			a.msg = s.Msg
		}
		agents = append(agents, a)
	}

	sort.Slice(agents, func(i, j int) bool {
		oi, oj := statusOrder(agents[i].status), statusOrder(agents[j].status)
		if oi != oj {
			return oi < oj
		}
		return agents[i].name < agents[j].name
	})

	return agents
}

func initialModel() model {
	return model{
		agents: collectAgents(),
	}
}

func (m model) Init() tea.Cmd {
	return tickEvery(2 * time.Second)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			// Number keys 1-9: jump to agent
			if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
				idx := int(msg.String()[0]-'0') - 1
				if idx < len(m.agents) {
					wezterm.ActivatePane(m.agents[idx].paneID)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.agents = collectAgents()
		return m, tickEvery(2 * time.Second)
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("  wzmux dashboard"))
	b.WriteString("\n")
	if m.width > 0 {
		b.WriteString(dimStyle.Render(strings.Repeat("─", m.width)))
	} else {
		b.WriteString(dimStyle.Render(strings.Repeat("─", 40)))
	}
	b.WriteString("\n")

	// Summary bar
	counts := map[string]int{}
	for _, a := range m.agents {
		counts[a.status]++
	}
	summary := fmt.Sprintf("  %d agents", len(m.agents))
	for _, s := range []string{"running", "waiting", "done", "error"} {
		if c := counts[s]; c > 0 {
			summary += fmt.Sprintf("  %s %d", statusIcon(s), c)
		}
	}
	b.WriteString(summary)
	b.WriteString("\n\n")

	if len(m.agents) == 0 {
		b.WriteString(dimStyle.Render("  No active agents."))
		b.WriteString("\n")
	}

	// Agent entries
	for i, a := range m.agents {
		icon := statusIcon(a.status)
		style := statusStyle(a.status)

		// Main line: [#] icon name status pane:id
		line := fmt.Sprintf("  [%d] %s %s %s pane:%d",
			i+1, icon, a.name, style.Render(a.status), a.paneID)
		b.WriteString(line)
		b.WriteString("\n")

		// Detail line: cwd + message
		detail := fmt.Sprintf("      %s", dimStyle.Render(a.cwd))
		if a.msg != "" {
			msg := a.msg
			if len(msg) > 55 {
				msg = msg[:55] + "…"
			}
			detail += "  " + dimStyle.Render(msg)
		}
		b.WriteString(detail)
		b.WriteString("\n\n")
	}

	// Footer
	b.WriteString(footerStyle.Render("  [1-9] jump to agent  [q] quit  refresh: 2s"))
	b.WriteString("\n")

	return b.String()
}

func runDashboard(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
