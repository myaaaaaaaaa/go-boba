package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "embed"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	errorStyle  = lipgloss.Style{}.Bold(true).Foreground(lipgloss.Color("#800000"))
	importStyle = lipgloss.Style{}.
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#ff4040"))

	diffAddStyle = lipgloss.Style{}.Bold(true).
			Foreground(lipgloss.Color("#004000")).
			Background(lipgloss.Color("#eeffee"))
	diffDelStyle = lipgloss.Style{}.Bold(true).
			Foreground(lipgloss.Color("#400000")).
			Background(lipgloss.Color("#ffeeee"))
)

type model struct {
	head     bool
	lines    []string
	viewport viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - (headerHeight + footerHeight)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "tab":
			m.head = !m.head
			m.viewport.SetContent(m.content())
		}

	case []string:
		m.lines = msg
		m.viewport.SetContent(m.content())
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m model) content() string {
	var lines []string
	if m.head {
		i := 0
		for _, line := range m.lines {
			if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
				i = 10
			}
			if i > 0 {
				lines = append(lines, line)
			}
			i--
		}
	} else {
		lines = m.lines
	}
	return strings.Join(lines, "\n")
}
func (m model) headerView() string {
	title := titleStyle.Render("git diff")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}
func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func main() {
	p := tea.NewProgram(
		model{viewport: viewport.New(0, 0)},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	go watch(p.Send)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func watch(send func(tea.Msg)) {
	for {
		out, err := exec.Command("git", "diff").CombinedOutput()

		var msg []string
		if err != nil {
			msg = append(msg, errorStyle.Render(
				err.Error(),
			))
		}

		{
			lines := strings.Split(string(out), "\n")
			for i := range lines {
				lines[i] = styleLine(lines[i])
			}
			msg = append(msg, lines...)
		}
		send(msg)

		time.Sleep(time.Second)
	}
}

func styleLine(line string) string {
	if strings.HasPrefix(line, "---") {
	} else if strings.HasPrefix(line, "+++") {
	} else if strings.HasPrefix(line, "-") {
		return diffDelStyle.Render(line)
	} else if strings.HasPrefix(line, "+") {
		return diffAddStyle.Render(line)
	}

	{
		a, b, ok := strings.Cut(line, "import")
		if ok {
			line = importStyle.Render("import")
			line = a + line + b
		}
	}
	return line
}
