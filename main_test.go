package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestApp(t *testing.T) {
	var buf bytes.Buffer

	m := model{posts: generatePosts()[:5], cb: func(p []post) tea.Model {
		for _, p := range p {
			fmt.Fprintln(&buf, p.title)
		}
		return nil
	}}
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Simulate window size message
	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Navigate and select some items
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyDown}) // cursor at 2

	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp}) // cursor at 0

	tm.Type("---=") // range is now 2

	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // selects {0, 1}

	// Run the model by calling FinalModel
	tm.FinalModel(t)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	expectedLines := []string{
		"Exploring the Alps", // index 0
		"The Art of Baking",  // index 1
	}

	expected := strings.Join(expectedLines, "\n")
	actual := strings.Join(lines, "\n")

	if actual != expected {
		t.Errorf("expected output:\n\n%q\n\ngot:\n\n%q", expected, actual)
	}
}
