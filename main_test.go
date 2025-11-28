package main

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func assertEquals(t *testing.T, got any, want any) {
	t.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Error("got", got, "    want", want)
	}
}

func TestApp(t *testing.T) {
	var output string

	m := model{posts: generatePosts()[:5], cb: func(p []post) tea.Model {
		for _, p := range p {
			output += p.title
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

	assertEquals(t, `Exploring the AlpsThe Art of Baking`, output)
}
