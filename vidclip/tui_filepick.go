package main

import (
	"errors"
	"mime"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
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
