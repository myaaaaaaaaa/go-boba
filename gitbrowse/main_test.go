package main

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdate(t *testing.T) {
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = "line"
	}

	tests := []struct {
		name     string
		initial  model
		msg      tea.Msg
		expected model
	}{
		{
			name: "window size",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.WindowSizeMsg{Width: 80, Height: 24},
			expected: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  80,
				height: 24,
			},
		},
		{
			name: "quit",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			expected: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor down",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			expected: model{
				lines:  lines,
				cursor: 1,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor down (clamp)",
			initial: model{
				lines:  lines,
				cursor: 99,
				offset: 90,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			expected: model{
				lines:  lines,
				cursor: 99, // clamped to len-1
				offset: 90, // offset = cursor - height + 1 -> 99 - 10 + 1 = 90
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor up",
			initial: model{
				lines:  lines,
				cursor: 5,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			expected: model{
				lines:  lines,
				cursor: 4,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "cursor up (clamp)",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			expected: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "page down",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 20, // pageSize = 20 * 3 / 4 = 15
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}},
			expected: model{
				lines:  lines,
				cursor: 15,
				offset: 0, // offset = cursor - height + 1 -> 15 - 20 + 1 = -4 -> wait, the logic says if cursor >= offset+height (15 >= 0+20 is false, so offset is unchanged)
				width:  10,
				height: 20,
			},
		},
		{
			name: "page up",
			initial: model{
				lines:  lines,
				cursor: 30,
				offset: 20,
				width:  10,
				height: 20, // pageSize = 15
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}},
			expected: model{
				lines:  lines,
				cursor: 15,
				offset: 15, // cursor < offset (15 < 20) -> offset = cursor
				width:  10,
				height: 20,
			},
		},
		{
			name: "home",
			initial: model{
				lines:  lines,
				cursor: 50,
				offset: 45,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}},
			expected: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
		},
		{
			name: "end",
			initial: model{
				lines:  lines,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}},
			expected: model{
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
			m, _ := tt.initial.Update(tt.msg)

			updatedModel := m.(model)
			if !reflect.DeepEqual(updatedModel, tt.expected) {
				t.Errorf("expected model \n%+v, got \n%+v", tt.expected, updatedModel)
			}
		})
	}
}

func TestModelUpdatePanics(t *testing.T) {
	tests := []struct {
		name    string
		initial model
		msg     tea.Msg
	}{
		{
			name: "empty lines",
			initial: model{
				lines:  []string{},
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		},
		{
			name: "nil lines",
			initial: model{
				lines:  nil,
				cursor: 0,
				offset: 0,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}},
		},
		{
			name: "zero height and width",
			initial: model{
				lines:  []string{"a", "b", "c"},
				cursor: 0,
				offset: 0,
				width:  0,
				height: 0,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		},
		{
			name: "negative height and width",
			initial: model{
				lines:  []string{"a", "b", "c"},
				cursor: 0,
				offset: 0,
				width:  -10,
				height: -10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		},
		{
			name: "negative cursor and offset",
			initial: model{
				lines:  []string{"a", "b", "c"},
				cursor: -5,
				offset: -5,
				width:  10,
				height: 10,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		},
		{
			name: "page down with small height",
			initial: model{
				lines:  []string{"a", "b", "c"},
				cursor: 0,
				offset: 0,
				width:  10,
				height: 1,
			},
			msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.View()
			m, _ := tt.initial.Update(tt.msg)
			m.View()
		})
	}
}
