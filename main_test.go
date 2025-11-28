package main

import (
	"bytes"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestApp(t *testing.T) {
	// a pipe redirects stdout to a buffer
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// save original stdout and redirect it
	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	m := newModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Simulate window size message
	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Navigate and select some items
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})  // cursor at 2
	tm.Send(tea.KeyMsg{Type: tea.KeySpace}) // selects {2}

	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	tm.Send(tea.KeyMsg{Type: tea.KeyUp}) // cursor at 0

	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'='}}) // range is now 2

	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // selects {0, 1}, now selected is {0, 1, 2}

	// Quit
	tm.Send(tea.KeyMsg{Type: tea.KeyEscape})

	// Run the model by calling FinalModel
	tm.FinalModel(t)

	// close the writer and read the buffer
	w.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatal(err)
	}

	// The output from fmt.Println will have trailing newlines, and map iteration order is not guaranteed.
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	sort.Strings(lines)

	expectedLines := []string{
		"Exploring the Alps",         // index 0
		"The Art of Baking",          // index 1
		"A Guide to Urban Gardening", // index 2
	}
	sort.Strings(expectedLines)

	expected := strings.Join(expectedLines, "\n")
	actual := strings.Join(lines, "\n")

	if actual != expected {
		t.Errorf("expected output:\n\n%q\n\ngot:\n\n%q", expected, actual)
	}
}
