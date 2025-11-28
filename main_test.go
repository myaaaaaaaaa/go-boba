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
	outputWriter := func(p []post) tea.Model {
		for _, p := range p {
			output += p.title
		}
		return nil
	}

	tm := teatest.NewTestModel(t, model{
		posts: generatePosts()[:5:5],
		cb:    outputWriter,
	})

	tm.Send(tea.WindowSizeMsg{Width: 40, Height: 24})

	// Try to go beyond end of list
	for range 10 {
		tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	}

	tm.Type("=====")

	// Try to go beyond beginning of list
	for range 15 {
		tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	}

	// Try to shrink selection size below 1, then grow it to 2
	tm.Type("------=")

	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // selects {1, 2}

	tm.FinalModel(t)

	assertEquals(t, output, `The Art of BakingA Guide to Urban Gardening`)
}
