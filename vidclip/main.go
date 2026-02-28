package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func failIf(aProblem bool, format string, a ...any) {
	if aProblem {
		fmt.Fprintf(os.Stderr, format, a...)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

func tryPop(args *[]string, format string, a ...any) string {
	failIf(len(*args) == 0, format, a...)
	val := (*args)[0]
	*args = (*args)[1:]
	return val
}

func main() {
	args := os.Args
	if len(args) < 2 {
		tui()
		return
	}
	exeName := tryPop(&args, "")
	subcmd := tryPop(&args, "")

	switch subcmd {
	case "export":
		f, err := io.ReadAll(os.Stdin)
		failIf(err != nil, "Error reading file: %v", err)

		list, err := Parse(string(f))
		failIf(err != nil, "Error parsing EDL: %v", err)

		fmt.Println(list.Export())
	case "mpv":
		file := tryPop(&args, "Usage: %s mpv <file>", exeName)
		t := 15.0
		mpvChooseTime(file, &t)
		fmt.Println("final time:", t)
	case "probe":
		file := tryPop(&args, "Usage: %s probe <file>", exeName)
		d, err := videoDuration(file)
		failIf(err != nil, "Error finding duration: %v", err)
		fmt.Println("duration:", d)

	default:
		failIf(true, "Unknown subcommand: %s", subcmd)
	}
}

type model struct {
	count     int
	stepSizes []int
	stepIdx   int
}

func initialModel() model {
	return model{
		count:     0,
		stepSizes: []int{1, 5, 10},
		stepIdx:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.count += m.stepSizes[m.stepIdx]
		case "down", "j":
			m.count -= m.stepSizes[m.stepIdx]
		case "tab", " ", "right", "l":
			m.stepIdx = (m.stepIdx + 1) % len(m.stepSizes)
		case "shift+tab", "left", "h":
			m.stepIdx = (m.stepIdx - 1 + len(m.stepSizes)) % len(m.stepSizes)
		}
	}
	return m, nil
}

var (
	titleStyle = lipgloss.Style{}.
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	countStyle = lipgloss.Style{}.
			Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true).
			SetString("Count:")

	selectedStyle = lipgloss.Style{}.
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#04B575")).
			Bold(true)

	dimStyle = lipgloss.Style{}.
			Foreground(lipgloss.Color("#626262"))

	helpStyle = lipgloss.Style{}.
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

func (m model) View() string {
	s := &strings.Builder{}

	fmt.Fprintln(s, titleStyle.Render("VIDCLIP COUNTER"))

	// Count Display
	fmt.Fprintf(s, "%s %d\n\n", countStyle, m.count)

	// Step Sizes
	fmt.Fprintln(s, "Step size: ")
	for i, size := range m.stepSizes {
		str := fmt.Sprintf(" %d ", size)
		style := dimStyle
		if i == m.stepIdx {
			style = selectedStyle
		}
		fmt.Fprint(s, style.Render(str))
	}
	fmt.Fprintln(s)

	// Help
	fmt.Fprintln(s, helpStyle.Render("↑/↓ • k/j: count  │  tab • l/h: step  │  q: quit"))

	return lipgloss.Style{}.Padding(1, 2).Render(s.String())
}

func tui() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
