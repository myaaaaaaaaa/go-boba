package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	lines  []string
	cursor int
	offset int
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		pageSize := m.height * 3 / 4

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case "pgup", "w":
			m.cursor -= pageSize
		case "pgdown", "f", " ":
			m.cursor += pageSize
		case "home", "g":
			m.cursor = 0
		case "end", "G":
			m.cursor = len(m.lines)
		}

		m.cursor = max(m.cursor, 0)
		m.cursor = min(m.cursor, len(m.lines)-1)
	}

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	for i := range m.height {
		s.WriteString("\n")
		idx := m.offset + i
		if idx >= len(m.lines) {
			break
		}

		line := m.lines[idx]
		line = strings.ReplaceAll(line, "\t", "        ")

		if idx == m.cursor {
			style := lipgloss.NewStyle().Reverse(true)
			s.WriteString(style.Render(line + strings.Repeat(" ", m.width)))
		} else {
			s.WriteString(line)
		}
	}
	return s.String()
}

func main() {
	var lines []string

	for _, arg := range os.Args[1:] {
		b, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", arg, err)
			continue
		}
		lines = append(lines, strings.Split(string(b), "\n")...)
	}

	m := model{lines: lines}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	finalM, ok := finalModel.(model)
	if !ok {
		return
	}

	if len(finalM.lines) > 0 {
		start := finalM.cursor - 10
		if start < 0 {
			start = 0
		}
		end := finalM.cursor + 10
		if end >= len(finalM.lines) {
			end = len(finalM.lines) - 1
		}

		for i := start; i <= end; i++ {
			if i == finalM.cursor {
				style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("204"))
				fmt.Printf("%s\n", style.Render(fmt.Sprintf("> %s", finalM.lines[i])))
			} else {
				fmt.Printf("  %s\n", finalM.lines[i])
			}
		}
	}
}
