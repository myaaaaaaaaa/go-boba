package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"os/exec"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model holds the state of our TUI application.
type model struct {
	saved, clips EditList

	cursor    int
	cursorCol int
	size      tea.WindowSizeMsg

	filename string
	err      string

	durationOf func(filename string) (float64, error)
	changeTime func(video string, t *float64)
}

// Styles for the TUI components.
var (
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))            // Dark gray/black
	faintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))            // Medium gray for inactive
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))            // Light gray for scrub bg
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))              // Cyan header
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true) // Salmon/Red
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))             // Active Blue
)

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
		case "left", "h":
			m.cursorCol = 0
		case "right", "l":
			m.cursorCol = 1
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
		case "ctrl+s":
			data := m.clips.Serialize()
			err := os.WriteFile(m.filename, []byte(data), 0644)
			if err != nil {
				m.err = "save failed: " + err.Error()
			} else {
				m.saved = slices.Clone(m.clips)
			}
		case " ":
			const file = "/tmp/preview.edl"
			data := m.clips.Serialize()
			os.WriteFile(file, []byte(data), 0644)
			go exec.Command("mpv", file).Run()
		case "enter":
			clip := &m.clips[m.cursor]
			m.changeTime(clip.Source, &clip.Times[m.cursorCol])
		}
	case pathMsg:
		m.clips = append(m.clips, EditEntry{Source: string(msg), Times: [2]float64{0, 1}})
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
	{
		style := cyanStyle
		filename := m.filename
		if !slices.Equal(m.saved, m.clips) {
			filename += " *"
			style = style.Bold(true)
		}
		s.WriteString(style.Render(filename) + "\n\n")
	}

	// Media Clip List
	for i, clip := range m.clips {
		isSelected := (i == m.cursor)
		sourceDuration, _ := m.durationOf(clip.Source)

		timeStrs := [2]string{
			formatTime(clip.Times[0]),
			formatTime(clip.Times[1]),
		}
		if isSelected {
			timeStrs[m.cursorCol] = blueStyle.Bold(true).Render(timeStrs[m.cursorCol])
		}

		clipDuration := clip.Times[1] - clip.Times[0]
		leftTop := faintStyle.Render(timeStrs[0]) + faintStyle.Render(" - ") + faintStyle.Render(timeStrs[1]) + defaultStyle.Render(fmt.Sprintf("    (%0.1fs)", clipDuration))

		metadata := clip.Source + "  " + formatTime(sourceDuration)
		if clip.Source == index(m.clips, i-1).Source {
			metadata = ""
		}
		rightTop := defaultStyle.Render(metadata)

		contentWidth := m.size.Width

		topRow := leftTop + strings.Repeat(" ", max(0, contentWidth-lipgloss.Width(leftTop)-lipgloss.Width(rightTop))) + rightTop

		scrubBar := ""
		if sourceDuration != 0 {
			scrubBar = m.renderScrubBar(clip.Times[0]/sourceDuration, clip.Times[1]/sourceDuration, contentWidth, isSelected)
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
		totalDuration += clip.Times[1] - clip.Times[0]
	}
	s.WriteString(fmt.Sprintf("Total Duration: %0.1fs\n", totalDuration))

	// Assembled Timeline Bar
	timelineWidth := max(m.size.Width-4, 40)

	var timeline strings.Builder
	for i, clip := range m.clips {
		clipDuration := clip.Times[1] - clip.Times[0]
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

func tui(file string) {
	m := model{
		filename:   "/tmp/project.edl",
		durationOf: memoize12(videoDuration),
		changeTime: mpvChooseTime,
	}

	if isVideo(file) {
		m.clips = EditList{{Source: file, Times: [2]float64{0, 1}}}
	} else if strings.HasSuffix(file, ".edl") {
		data, err := os.ReadFile(file)
		failIf(err != nil, "failed to read %s: %v", file, err)

		m.filename = file
		m.clips, err = Parse(string(data))
		failIf(err != nil, "invalid .edl file: %s: %v", file, err)
	} else {
		// dummy data
		m.clips = EditList{
			{Source: "sotu_2024_raw.mp4", Times: [2]float64{120.5, 145.6}},
			{Source: "sotu_2024_raw.mp4", Times: [2]float64{150.5, 165.6}},
			{Source: "sotu_2024_raw.mp4", Times: [2]float64{180.5, 200.1}},
			{Source: "gameplay_capture_01.mkv", Times: [2]float64{340.1, 385.3}},
			{Source: "outro.mp4", Times: [2]float64{0.0, 5.5}},
		}

		m.cursor = 1
		m.err = "error: could not open hello.jpg: not a video file"

		m.durationOf = memoize12(func(filename string) (float64, error) {
			return float64(len(filename) * 50), nil
		})
		m.changeTime = func(video string, t *float64) {
			*t += rand.Float64()*4 - 2
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
