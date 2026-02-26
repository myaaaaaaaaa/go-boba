package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestWindowSize(t *testing.T) {
	initial := model{
		width:  10,
		height: 10,
	}
	want := model{
		width:  80,
		height: 24,
	}
	got, _ := initial.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want model \n%+v, got \n%+v", want, got)
	}
}

func TestModelUpdate(t *testing.T) {
	lines := make([]line, 100)
	for i := range lines {
		lines[i] = line{text: "line"}
	}

	tests := []struct {
		name    string
		msg     rune
		initial tea.Model
		want    tea.Model
	}{
		{
			name: "quit",
			msg:  'q',
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			want: finalModel{model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			}},
		},
		{
			name: "cursor down",
			msg:  'j',
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 1,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor down (clamp)",
			msg:  'j',
			initial: model{
				lines:  lines,
				cursor: 99,
				offset: 90,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 99, // clamped to len-1
				offset: 90, // offset = cursor - height + 1 -> 99 - 10 + 1 = 90
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor up",
			msg:  'k',
			initial: model{
				lines:  lines,
				cursor: 5,
				offset: 0,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 4,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor up (clamp)",
			msg:  'k',
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "page down",
			msg:  'f',
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 20, // pageSize = 20 * 3 / 4 = 15
			},
			want: model{
				lines:  lines,
				cursor: 15,
				offset: 0, // offset = cursor - height + 1 -> 15 - 20 + 1 = -4 -> wait, the logic says if cursor >= offset+height (15 >= 0+20 is false, so offset is unchanged)
				width:  10,
				height: 20,
			},
		},
		{
			name: "page up",
			msg:  'w',
			initial: model{
				lines:  lines,
				cursor: 30,
				offset: 20,
				width:  10,
				height: 20, // pageSize = 15
			},
			want: model{
				lines:  lines,
				cursor: 15,
				offset: 15, // cursor < offset (15 < 20) -> offset = cursor
				width:  10,
				height: 20,
			},
		},
		{
			name: "home",
			msg:  'g',
			initial: model{
				lines:  lines,
				cursor: 50,
				offset: 45,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "end",
			msg:  'G',
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			want: model{
				lines:  lines,
				cursor: 99,
				offset: 90, // cursor - height + 1 = 99 - 10 + 1 = 90
				width:  10,
				height: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.msg}}
			got, _ := tt.initial.Update(msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want model \n%+v, got \n%+v", tt.want, got)
			}
		})
	}
}

// TestModelUpdatePanics runs the model through pathological sequences of key presses
// combined with edge-case initial states to ensure the application does not panic.
func TestModelUpdatePanics(t *testing.T) {
	initials := []model{
		{
			lines:  []line{},
			width:  10,
			height: 10,
		},
		{
			width:  10,
			height: 10,
		},
		{
			lines: []line{{text: "a"}, {text: "b"}, {text: "c"}},
		},
		{
			lines:  []line{{text: "a"}, {text: "b"}, {text: "c"}},
			width:  -10,
			height: -10,
		},
		{
			lines:  []line{{text: "a"}, {text: "b"}, {text: "c"}},
			cursor: -5,
			offset: -5,
			width:  10,
			height: 10,
		},
		{
			lines:  []line{{text: "a"}, {text: "b"}, {text: "c"}},
			width:  10,
			height: 1,
		},
	}

	msgs := []string{
		"jjkkGGggffww  ",
		"Ggjkwf ",
		"jjkkffwwGGgg",
		"ffffffwwwwww",
		"q",
		"jjq",
		"GgqqqjjG",
		"  q  ",
	}

	for _, init := range initials {
		for _, msgSeq := range msgs {
			var m tea.Model = init
			m.View()
			for _, r := range msgSeq {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
				m, _ = m.Update(msg)
				m.View()
			}
		}
	}
}

func TestEndToEnd(t *testing.T) {
	var lines []line
	for i := range 50 {
		i := i + 1
		lines = append(lines, line{text: fmt.Sprintf("<%d>", i), num: i, file: "test.txt"})
	}

	var m tea.Model = model{lines: lines}
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 10})

	keys := "jjfkq"
	for _, k := range keys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
		m, _ = m.Update(msg)
	}

	out := m.View()

	// Check if the final output contains expected lines
	if !strings.Contains(out, "<1>") {
		t.Errorf("View should contain '<1>'")
	}
	if !strings.Contains(out, "<19>") {
		t.Errorf("View should contain '<19>'")
	}

	// Check if it doesn't contain unexpected lines
	if strings.Contains(out, "<20>") {
		t.Errorf("View should not contain '<20>'")
	}
	if strings.Contains(out, "<50>") {
		t.Errorf("View should not contain '<50>'")
	}

	// The cursor is at index 8 (which is line 9).
	// start = 8 - 10 = 0 (line 1 to line 19)
	// count lines
	linesOut := strings.Split(out, "\n")
	if len(linesOut) != 19 {
		t.Errorf("View should have 19 lines, got %d", len(linesOut))
	}
}

func TestEndToEndEdgeCases(t *testing.T) {
	tests := []struct {
		name              string
		lineCount         int
		windowWidth       int
		windowHeight      int
		keys              string
		expectedContains  []string
		expectedNotCont   []string
		expectedLineCount int
	}{
		{
			name:              "empty file",
			lineCount:         0,
			windowWidth:       80,
			windowHeight:      24,
			keys:              "q",
			expectedContains:  []string{},
			expectedNotCont:   []string{"<1>"},
			expectedLineCount: 0, // Empty string
		},
		{
			name:              "single line file",
			lineCount:         1,
			windowWidth:       80,
			windowHeight:      24,
			keys:              "q",
			expectedContains:  []string{"<1>"},
			expectedNotCont:   []string{"<2>"},
			expectedLineCount: 1,
		},
		{
			name:              "navigate to very end",
			lineCount:         100,
			windowWidth:       80,
			windowHeight:      24,
			keys:              "Gq",                      // 'G' goes to end
			expectedContains:  []string{"<90>", "<100>"}, // Cursor at 99. View is 89 to 99 (11 lines)
			expectedNotCont:   []string{"<89>", "<101>"},
			expectedLineCount: 11,
		},
		{
			name:              "navigate to start after moving",
			lineCount:         100,
			windowWidth:       80,
			windowHeight:      24,
			keys:              "Ggq",                   // 'G' to end, 'g' to start
			expectedContains:  []string{"<1>", "<11>"}, // Cursor at 0. View is 0 to 10 (11 lines)
			expectedNotCont:   []string{"<12>"},
			expectedLineCount: 11,
		},
		{
			name:              "page down past EOF",
			lineCount:         20,
			windowWidth:       80,
			windowHeight:      10,                       // Page size = 7
			keys:              "fffq",                   // Page down 3 times
			expectedContains:  []string{"<10>", "<20>"}, // Cursor at 19. View is 9 to 19 (11 lines)
			expectedNotCont:   []string{"<9>", "<21>"},
			expectedLineCount: 11,
		},
		{
			name:              "page up past start",
			lineCount:         20,
			windowWidth:       80,
			windowHeight:      10,
			keys:              "fwwwq",                 // Page down once, up three times
			expectedContains:  []string{"<1>", "<11>"}, // Cursor at 0. View is 0 to 10 (11 lines)
			expectedNotCont:   []string{"<12>"},
			expectedLineCount: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var lines []line
			for i := range tt.lineCount {
				i := i + 1
				lines = append(lines, line{text: fmt.Sprintf("<%d>", i), num: i, file: "test.txt"})
			}

			var m tea.Model = model{lines: lines}
			m, _ = m.Update(tea.WindowSizeMsg{Width: tt.windowWidth, Height: tt.windowHeight})

			for _, k := range tt.keys {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
				m, _ = m.Update(msg)
			}

			out := m.View()

			for _, expected := range tt.expectedContains {
				if !strings.Contains(out, expected) {
					t.Errorf("View should contain '%s'", expected)
				}
			}

			for _, unexpected := range tt.expectedNotCont {
				if strings.Contains(out, unexpected) {
					t.Errorf("View should not contain '%s'", unexpected)
				}
			}

			if tt.expectedLineCount == 0 {
				if out != "" {
					t.Errorf("View should be empty, got %q", out)
				}
			} else {
				linesOut := strings.Split(out, "\n")
				if len(linesOut) != tt.expectedLineCount {
					t.Errorf("View should have %d lines, got %d:\n%s", tt.expectedLineCount, len(linesOut), out)
				}
			}
		})
	}
}
