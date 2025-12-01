package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

func assertEquals(t *testing.T, got any, want any) {
	t.Helper()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Error("got", got, "    want", want)
	}
}

type send func(tea.Msg) string

func (send send) key(t tea.KeyType, repeat int) (rt string) {
	for range repeat {
		rt = send(tea.KeyMsg{Type: t})
	}
	return
}
func (send send) text(s string) (rt string) {
	for _, r := range s {
		rt = send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
		})
	}
	return
}

func teaTester(m tea.Model, out io.Writer) send {
	return func(msg tea.Msg) string {
		m, _ = m.Update(msg)
		rt := ansi.Strip(m.View())
		out.Write([]byte(rt))
		return rt
	}
}

func hash(data []byte) string {
	rt := sha256.Sum256(data)
	return fmt.Sprintf("sha256:%x", rt)
}

func TestApp(t *testing.T) {
	var output string
	outputWriter := func(p []post) tea.Model {
		for _, p := range p {
			output += p.Title
		}
		return nil
	}

	var frames bytes.Buffer
	send := teaTester(model{
		posts: generatePosts()[:5:5],
		cb:    outputWriter,
	}, &frames)

	send(tea.WindowSizeMsg{Width: 40, Height: 24})

	{
		// Try to go beyond end of list
		want := send.key(tea.KeyDown, 10)
		// Try to grow selection, fails
		got := send.text("=====")
		assertEquals(t, got, want)
	}

	// Try to go beyond beginning of list
	send.key(tea.KeyUp, 15)

	// Try to shrink selection size below 1, then grow it to 2
	send.text("------=")

	send.key(tea.KeyDown, 1)
	assertEquals(t, output, "")

	send.key(tea.KeyEnter, 1) // selects {1, 2}
	assertEquals(t, output, `The Art of BakingA Guide to Urban Gardening`)

	// Regression smoke test. Update the checksum if visual change were intentional.
	assertEquals(t, hash(frames.Bytes()), "sha256:6bf5c180ac620fa05ba274944b51a8071eaf051b4c10ad070d7825e84bc63de3")
}
