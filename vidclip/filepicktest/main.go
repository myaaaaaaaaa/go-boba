package main

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type pathMsg string

type fileModel struct {
	filepicker filepicker.Model
	err        error

	callback tea.Model
}

func newFileModel(size tea.WindowSizeMsg, callback tea.Model) (tea.Model, tea.Cmd) {
	fp := filepicker.New()
	//fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.CurrentDirectory = "/tmp"
	fp, _ = fp.Update(size)

	fm := fileModel{
		filepicker: fp,
		callback:   callback,
	}
	return fm, fm.Init()
}

func (m fileModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m fileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m.callback, nil
		}
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)

		if strings.HasPrefix(mimeType, "video/") {
			return m.callback.Update(pathMsg(path))
		}

		m.err = errors.New(path + " is not a video file.")
	}

	return m, cmd
}

func (m fileModel) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString("Pick a video file:")
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

type mainModel struct {
	size   tea.WindowSizeMsg
	files  []string
	cursor int
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case pathMsg:
		if msg == "" {
			break
		}
		m.cursor = min(m.cursor+1, len(m.files))
		m.files = slices.Insert(m.files, m.cursor, string(msg))
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
		case "down", "j":
			m.cursor++
		case "a":
			return newFileModel(m.size, m)
		case "A":
			//For testing convenience
			return m.Update(pathMsg("test.txt"))
		case "d", "r", "backspace", "delete":
			if m.cursor < len(m.files) {
				m.files = slices.Delete(m.files, m.cursor, m.cursor+1)
			}
		}
	}
	m.cursor = min(m.cursor, len(m.files)-1)
	m.cursor = max(m.cursor, 0)
	return m, nil
}

var (
	selectedItemStyle = lipgloss.NewStyle().Bold(true)
	normalItemStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (m mainModel) View() string {
	var s strings.Builder
	s.WriteString("\n  Files:\n")
	if len(m.files) == 0 {
		s.WriteString("  No files added yet.\n")
	} else {
		for i, f := range m.files {
			if m.cursor == i {
				s.WriteString(fmt.Sprintf("  ↳ %s\n", selectedItemStyle.Render(f)))
			} else {
				s.WriteString(fmt.Sprintf("    %s\n", normalItemStyle.Render(f)))
			}
		}
	}
	s.WriteString("\n  [a] Add file  [d] Remove file  [q] Quit\n")
	return s.String()
}

func main() {
	m := mainModel{}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
