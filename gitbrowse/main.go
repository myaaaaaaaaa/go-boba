package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type line struct {
	text string
	num  int
	file string
}

// model represents the core state of the Bubble Tea terminal application.
// It holds the contents of the files being viewed, the user's current
// selection (cursor), the scrolling view position (offset), and the terminal dimensions.
type model struct {
	lines  []line
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
		case "q", "ctrl+c", "esc", " ", "enter":
			return finalModel{m}, tea.Quit
		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case "pgup", "w":
			m.cursor -= pageSize
		case "pgdown", "f":
			m.cursor += pageSize
		case "home", "g":
			m.cursor = 0
		case "end", "G":
			m.cursor = len(m.lines)
		}

		m.cursor = min(m.cursor, len(m.lines)-1)
		m.cursor = max(m.cursor, 0)
	}

	m.offset = min(m.offset, m.cursor)
	m.offset = max(m.offset, m.cursor-m.height+1)
	m.offset = max(m.offset, 0)

	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	for i := range m.height {
		s.WriteString("\n")
		idx := m.offset + i

		line := ""
		if 0 <= idx && idx < len(m.lines) {
			line = m.lines[idx].text
		}
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

func getFiles() []string {
	if len(os.Args) > 1 {
		return os.Args[1:]
	}

	var files []string

	cmd := exec.Command("git", "ls-files")
	out, _ := cmd.Output()
	for file := range strings.SplitSeq(string(out), "\n") {
		if file != "" {
			files = append(files, file)
		}
	}
	if len(files) > 0 {
		return files
	}

	entries, err := os.ReadDir(".")
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, entry.Name())
			}
		}
	}

	return files
}

func main() {
	var lines []line

	for _, arg := range getFiles() {
		b, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", arg, err)
			continue
		}
		for i, text := range strings.Split(string(b), "\n") {
			lines = append(lines, line{
				text: text,
				num:  i + 1,
				file: arg,
			})
		}
	}

	m := model{lines: lines}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen())
	res, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	// Print the view of the final model (which will be the snippet if we quit normally)
	fmt.Println(res.View())
}

type finalModel struct{ model }

func (m finalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, tea.Quit }
func (m finalModel) View() string {
	if len(m.lines) == 0 {
		return ""
	}

	var s strings.Builder
	start := m.cursor - 10
	if start < 0 {
		start = 0
	}
	end := m.cursor + 10
	if end >= len(m.lines) {
		end = len(m.lines) - 1
	}

	for i := start; i <= end; i++ {
		line := m.lines[i]
		if i == m.cursor {
			style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("204"))
			s.WriteString(style.Render(fmt.Sprintf("> %s:%d: %s", line.file, line.num, line.text)))
		} else {
			s.WriteString(fmt.Sprintf("  %s:%d: %s", line.file, line.num, line.text))
		}
		if i < end {
			s.WriteString("\n")
		}
	}

	return s.String()
}
