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
	case tea.MouseMsg:
		const scrollSize = 3
		switch msg.Type {
		case tea.MouseWheelUp:
			m.cursor -= scrollSize
		case tea.MouseWheelDown:
			m.cursor += scrollSize
		}
	}

	m.cursor = min(m.cursor, len(m.lines)-1)
	m.cursor = max(m.cursor, 0)

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

		line := m.at(idx).text
		line = strings.ReplaceAll(line, "\t", "        ")

		if idx == m.cursor {
			style := lipgloss.Style{}.Reverse(true)
			s.WriteString(style.Render(line + strings.Repeat(" ", m.width)))
		} else {
			s.WriteString(line)
		}
	}
	return s.String()
}

func (m model) at(idx int) line {
	if 0 <= idx && idx < len(m.lines) {
		return m.lines[idx]
	}
	return line{}
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

	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))
	m := model{lines: lines}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr), tea.WithAltScreen(), tea.WithMouseCellMotion())
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
	var s strings.Builder
	fmt.Fprintln(&s, "<context>")

	var (
		start = max(m.cursor-10, 0)
		end   = min(m.cursor+10, len(m.lines)-1)
	)
	for i := start; i <= end; i++ {
		line := m.lines[i]
		indent := "  "
		if i == m.cursor {
			indent = ">>"
		}
		fmt.Fprintf(&s, indent+"%d: %s\n", line.num, line.text)
	}
	fmt.Fprintln(&s, "</context>")
	fmt.Fprintln(&s)

	fmt.Fprintln(&s, "<selection>")
	line := m.at(m.cursor)
	fmt.Fprintf(&s, "%s:%d:  %s\n", line.file, line.num, line.text)
	fmt.Fprintln(&s, "</selection>")

	return s.String()
}
