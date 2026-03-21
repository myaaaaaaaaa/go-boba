package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestViewSmoke(t *testing.T) {
	tests := []struct {
		name  string
		model model
	}{
		{
			name: "empty clips",
			model: model{
				size: tea.WindowSizeMsg{Width: 80, Height: 24},
			},
		},
		{
			name: "single valid clip",
			model: model{
				clips: EditList{
					{Source: "test.mp4", Times: [2]float64{0, 10}},
				},
				size: tea.WindowSizeMsg{Width: 80, Height: 24},
				durationOf: func(filename string) (float64, error) {
					return 100, nil
				},
			},
		},
		{
			name: "end < start (reproducer for crash)",
			model: model{
				clips: EditList{
					{Source: "test.mp4", Times: [2]float64{0, 10}},
					{Source: "test.mp4", Times: [2]float64{10, 5}},
				},
				size: tea.WindowSizeMsg{Width: 80, Height: 24},
				durationOf: func(filename string) (float64, error) {
					return 100, nil
				},
			},
		},
		{
			name: "total duration zero",
			model: model{
				clips: EditList{
					{Source: "test.mp4", Times: [2]float64{5, 5}},
				},
				size: tea.WindowSizeMsg{Width: 80, Height: 24},
				durationOf: func(filename string) (float64, error) {
					return 100, nil
				},
			},
		},
		{
			name: "many clips (negative available width)",
			model: model{
				clips: func() EditList {
					el := EditList{}
					for range 100 {
						el = append(el, EditEntry{Source: "test.mp4", Times: [2]float64{0, 1}})
					}
					return el
				}(),
				size: tea.WindowSizeMsg{Width: 80, Height: 24},
				durationOf: func(filename string) (float64, error) {
					return 100, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.model.View()
		})
	}
}
