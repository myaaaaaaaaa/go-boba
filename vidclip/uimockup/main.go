package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Clip represents a single video segment.
type Clip struct {
	startTime float64
	endTime   float64
	srcVideo  string
	// sourceDuration is the total length of the source file, used for the scrub bar.
	sourceDuration float64
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
}

// Styles for the TUI components.
var (
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // Dark gray/black
	faintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("246")) // Medium gray for inactive
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("251")) // Light gray for scrub bg
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))   // Cyan header
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203")) // Salmon/Red
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))  // Active Blue

	selectedStyle = lipgloss.NewStyle().Padding(0, 1)
	normalStyle   = lipgloss.NewStyle().Padding(0, 1)
)

func initialModel() model {
	clips := ClipList{
		{startTime: 120.5, endTime: 145.6, srcVideo: "sotu_2024_raw.mp4", sourceDuration: 600.0},
		{startTime: 15.0, endTime: 22.8, srcVideo: "interview_b_roll.mov", sourceDuration: 45.0},
		{startTime: 340.1, endTime: 385.3, srcVideo: "gameplay_capture_01.mkv", sourceDuration: 1200.0},
		{startTime: 0.0, endTime: 5.5, srcVideo: "outro_template.mp4", sourceDuration: 10.0},
	}

	return model{
		clips:    clips,
		cursor:   1,
		filename: "project.edl*",
		err:      "error: could not open hello.jpg: not a video file",
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.clips)-1 {
				m.cursor++
			}
		case "n":
			return newFileModel(m.size, m)

		}
	case pathMsg:
		m.clips = append(m.clips, Clip{srcVideo: string(msg)})
	case tea.WindowSizeMsg:
		m.size = msg
	}
	return m, nil
}

func formatTime(seconds float64) string {
	minutes := int(seconds) / 60
	secs := math.Mod(seconds, 60)
	return fmt.Sprintf("%02d:%04.1f", minutes, secs)
}

func (m model) renderScrubBar(clip Clip, width int, isSelected bool) string {
	if clip.sourceDuration == 0 {
		return strings.Repeat(" ", width)
	}

	startPct := clip.startTime / clip.sourceDuration
	endPct := clip.endTime / clip.sourceDuration

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

func (m model) View() string {
	var s strings.Builder

	// Header Bar
	s.WriteString(cyanStyle.Render(m.filename) + "\n\n")

	// Media Clip List
	for i, clip := range m.clips {
		isSelected := (i == m.cursor)

		containerStyle := normalStyle
		textStyle := defaultStyle
		if isSelected {
			textStyle = textStyle.Bold(true)
		} else {
			textStyle = faintStyle
		}

		duration := clip.endTime - clip.startTime

		// Left: Timestamps and duration
		leftTop := textStyle.Render(fmt.Sprintf("%s - %s  (%0.1fs)", formatTime(clip.startTime), formatTime(clip.endTime), duration))

		// Right: Filename and Source Duration
		rightTop := textStyle.Render(clip.srcVideo)
		rightBottom := textStyle.Render(formatTime(clip.sourceDuration))

		contentWidth := max(m.size.Width-4, 40)

		topRow := leftTop + strings.Repeat(" ", max(0, contentWidth-lipgloss.Width(leftTop)-lipgloss.Width(rightTop))) + rightTop
		bottomRow := strings.Repeat(" ", max(0, contentWidth-lipgloss.Width(rightBottom))) + rightBottom

		scrubBar := m.renderScrubBar(clip, contentWidth, isSelected)

		item := fmt.Sprintf("%s\n%s\n%s", topRow, bottomRow, scrubBar)
		s.WriteString(containerStyle.Render(item) + "\n\n")
	}

	// Summary Section
	totalDuration := 0.0
	for _, clip := range m.clips {
		totalDuration += clip.endTime - clip.startTime
	}
	s.WriteString(fmt.Sprintf("Total Duration: %0.1fs\n", totalDuration))

	// Assembled Timeline Bar
	timelineWidth := max(m.size.Width-4, 40)

	var timeline strings.Builder
	for i, clip := range m.clips {
		clipDuration := clip.endTime - clip.startTime
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

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
