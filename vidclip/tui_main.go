package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Clip represents a single video segment.
type Clip struct {
	Source    string
	StartTime float64
	EndTime   float64
}

// ClipList is the central data structure as requested.
type ClipList []Clip

// Model holds the state of our TUI application.
type model struct {
	clips  ClipList
	cursor int
	size   tea.WindowSizeMsg

	filename string
	err      string

	durationOf func(filename string) (float64, error)
}

// Styles for the TUI components.
var (
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))            // Dark gray/black
	faintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))            // Medium gray for inactive
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))            // Light gray for scrub bg
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))              // Cyan header
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true) // Salmon/Red
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))             // Active Blue
)

func initialModel() model {
	clips := ClipList{
		{Source: "sotu_2024_raw.mp4", StartTime: 120.5, EndTime: 145.6},
		{Source: "sotu_2024_raw.mp4", StartTime: 150.5, EndTime: 165.6},
		{Source: "sotu_2024_raw.mp4", StartTime: 180.5, EndTime: 200.1},
		{Source: "gameplay_capture_01.mkv", StartTime: 340.1, EndTime: 385.3},
		{Source: "outro.mp4", StartTime: 0.0, EndTime: 5.5},
	}

	return model{
		clips:    clips,
		cursor:   1,
		filename: "project.edl*",
		err:      "error: could not open hello.jpg: not a video file",

		durationOf: memoize12(func(filename string) (float64, error) {
			return float64(len(filename) * 50), nil
		}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit

		case "n":
			return newFileModel(m.size, m)

		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case "alt+up", "K":
			if m.cursor > 0 {
				m.clips[m.cursor], m.clips[m.cursor-1] = m.clips[m.cursor-1], m.clips[m.cursor]
				m.cursor--
			}
		case "alt+down", "J":
			if m.cursor < len(m.clips)-1 {
				m.clips[m.cursor], m.clips[m.cursor+1] = m.clips[m.cursor+1], m.clips[m.cursor]
				m.cursor++
			}
		case "d":
			m.clips = slices.Insert(m.clips, m.cursor, m.clips[m.cursor])
		case "del", "x":
			if len(m.clips) <= 1 {
				m.err = "Refusing to delete last clip"
			} else {
				m.clips = slices.Delete(m.clips, m.cursor, m.cursor+1)
			}
		}
	case pathMsg:
		m.clips = append(m.clips, Clip{EndTime: 1, Source: string(msg)})
	case tea.WindowSizeMsg:
		m.size = msg
	}
	m.cursor = min(m.cursor, len(m.clips)-1)
	m.cursor = max(m.cursor, 0)
	return m, nil
}

func formatTime(seconds float64) string {
	minutes := int(seconds) / 60
	secs := math.Mod(seconds, 60)
	return fmt.Sprintf("%02d:%04.1f", minutes, secs)
}

func (m model) renderScrubBar(startPct, endPct float64, width int, isSelected bool) string {
	barWidth := width
	startIdx := int(startPct * float64(barWidth))
	endIdx := int(endPct * float64(barWidth))

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > barWidth {
		endIdx = barWidth
	}
	if endIdx < startIdx {
		endIdx = startIdx
	}

	fgStyle := blueStyle
	if !isSelected {
		fgStyle = faintStyle
	}

	var bar strings.Builder
	for i := range barWidth {
		if i >= startIdx && i < endIdx {
			bar.WriteString(fgStyle.Render("━"))
		} else {
			bar.WriteString(subtleStyle.Render("─"))
		}
	}
	return bar.String()
}

func index[T any](s []T, i int) T {
	if 0 <= i && i < len(s) {
		return s[i]
	}

	var zero T
	return zero
}

func (m model) View() string {
	var s strings.Builder

	// Header Bar
	s.WriteString(cyanStyle.Render(m.filename) + "\n\n")

	// Media Clip List
	for i, clip := range m.clips {
		isSelected := (i == m.cursor)
		sourceDuration, _ := m.durationOf(clip.Source)

		textStyle := defaultStyle
		if isSelected {
			textStyle = textStyle.Bold(true)
		} else {
			textStyle = faintStyle
		}

		clipDuration := clip.EndTime - clip.StartTime
		leftTop := textStyle.Render(fmt.Sprintf("%s - %s  (%0.1fs)", formatTime(clip.StartTime), formatTime(clip.EndTime), clipDuration))

		metadata := clip.Source + "  " + formatTime(sourceDuration)
		if clip.Source == index(m.clips, i-1).Source {
			metadata = ""
		}
		rightTop := textStyle.Render(metadata)

		contentWidth := m.size.Width

		topRow := leftTop + strings.Repeat(" ", max(0, contentWidth-lipgloss.Width(leftTop)-lipgloss.Width(rightTop))) + rightTop

		scrubBar := ""
		if sourceDuration != 0 {
			scrubBar = m.renderScrubBar(clip.StartTime/sourceDuration, clip.EndTime/sourceDuration, contentWidth, isSelected)
		}
		if clip.Source == index(m.clips, i+1).Source {
			scrubBar = strings.ReplaceAll(scrubBar, "─", " ")
		}

		s.WriteString(fmt.Sprintf("%s\n%s\n\n", topRow, scrubBar))
	}

	s.WriteString("\n\n\n")

	// Summary Section
	totalDuration := 0.0
	for _, clip := range m.clips {
		totalDuration += clip.EndTime - clip.StartTime
	}
	s.WriteString(fmt.Sprintf("Total Duration: %0.1fs\n", totalDuration))

	// Assembled Timeline Bar
	timelineWidth := max(m.size.Width-4, 40)

	var timeline strings.Builder
	for i, clip := range m.clips {
		clipDuration := clip.EndTime - clip.StartTime
		segmentWidth := int((clipDuration / totalDuration) * float64(timelineWidth-len(m.clips)))
		if segmentWidth == 0 && clipDuration > 0 {
			segmentWidth = 1
		}

		char := "─"
		style := subtleStyle
		if i == m.cursor {
			style = blueStyle
			char = "━"
		}

		timeline.WriteString(style.Render(strings.Repeat(char, segmentWidth)))
		if i < len(m.clips)-1 {
			timeline.WriteString(" ") // Visual gap between clips
		}
	}
	s.WriteString(timeline.String() + "\n\n")

	// Status/Error Bar
	s.WriteString(errorStyle.Render(m.err))

	return s.String()
}

func tui() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
