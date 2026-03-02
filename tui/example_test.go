package tui_test

import (
	"fmt"

	"github.com/myaaaaaaaaa/go-boba/tui"
)

func ExampleFrom() {
	// Create the buffer.
	// Tabs are used for indentation in the Go source code but are ignored by the function.
	buf := tui.From(`
		+---+
		| A |
		+---+
	`)

	// Update 'A' to 'C' with red style
	buf = buf.MapCell('A', 'C', "\033[31m")

	// Print the buffer's string representation.
	fmt.Print(buf.String())

	// Output:
	// +---+
	// | C |
	// +---+
}

func ExampleFixed_NinePatchScale() {
	// The center 'C' defines the scalable region.
	// We will scale it to 8x5.
	buf := tui.From(`
		+--+
		|CC|
		+--+
	`)

	scalable := buf.NinePatchScale('C')
	scaled := scalable(8, 5)

	// After scaling, we might want to replace the 'C' placeholders with space or another pattern.
	final := scaled.MapCell('C', ' ', "")

	fmt.Print(final.String())

	// Output:
	// +------+
	// |      |
	// |      |
	// |      |
	// +------+
}

func ExampleFixed_MapCell() {
	buf := tui.From(`
		X..X
		....
		X..X
	`)

	// Replace 'X' with '*'
	updated := buf.MapCell('X', '*', "")

	fmt.Print(updated.String())

	// Output:
	// *..*
	// ....
	// *..*
}

func ExampleFixed_MapCell_style() {
	buf := tui.From(`Error: X`)

	styled := buf.MapCell('X', 'X', "\033[31m")

	// Check style of the specific cell (not printed in simple String())
	// but demonstrating API usage.
	// In a real terminal, 'X' would be red.
	// For this test, we just print the content.
	fmt.Print(styled.String())

	// Output:
	// Error: X
}

func ExampleFixed_Fill() {
	target := tui.From(`
		+-------+
		| TTTTT |
		| TTTTT |
		| TTTTT |
		| TTTTT |
		+-------+
	`)

	source := tui.From(`
		+--+
		|..|
		+--+
	`)
	// Create scalable. The center '.' will be tiled.
	scalable := source.NinePatchScale('.')

	// Fill the scalable output into the bounding box of 'T' in target.
	// The bounding box of T is 5x3.
	result := target.Fill('T', scalable)

	fmt.Print(result.String())

	// Output:
	// +-------+
	// | +---+ |
	// | |...| |
	// | |...| |
	// | +---+ |
	// +-------+
}

func ExampleFixed_Fill_with_WrapText() {
	// Target buffer with a region 'T' to fill with text
	target := tui.From(`
		+-------+
		| TTTTT |
		| TTTTT |
		| TTTTT |
		+-------+
	`)

	scalable := tui.WrapText("Hello World This fits")

	// Fill the wrapped text into the 'T' region (5x3)
	result := target.Fill('T', scalable)

	fmt.Print(result.String())

	// Output:
	// +-------+
	// | Hello |
	// | World |
	// | Th... |
	// +-------+
}

func ExampleShader() {
	// Create a target buffer with a region 'S' to be filled by the shader.
	target := tui.From(`
		+-------+
		| SSSSS |
		| SSSSS |
		| SSSSS |
		+-------+
	`)

	// Define a shader function that creates a pattern based on coordinates.
	// This shader creates a split pattern: dots on the left, hashes on the right.
	shaderFunc := func(x, y float64) tui.Cell {
		// x and y are normalized to [0, 1).
		// We'll use a simple threshold to determine the character.
		if x < 0.5 {
			return tui.Cell{Rune: '.'}
		}
		return tui.Cell{Rune: '#'}
	}

	// Create a scalable from the shader function.
	scalable := tui.Shader(shaderFunc)

	// Fill the shader output into the 'S' region of the target buffer.
	result := target.Fill('S', scalable)

	fmt.Print(result.String())

	// Output:
	// +-------+
	// | ...## |
	// | ...## |
	// | ...## |
	// +-------+
}

func ExampleFixed_NearestNeighborScale() {
	target := tui.From(`
		+--------+
		| TTTTTT |
		| TTTTTT |
		| TTTTTT |
		| TTTTTT |
		+--------+
	`)

	// Source is a small pattern we want to upscale
	source := tui.From(`
		AB
		CD
	`)

	// Scale up the source using nearest neighbor.
	scalable := source.NearestNeighborScale()

	result := target.Fill('T', scalable)

	fmt.Print(result.String())

	// Output:
	// +--------+
	// | AAABBB |
	// | AAABBB |
	// | CCCDDD |
	// | CCCDDD |
	// +--------+
}

func ExampleFixed_Draw() {
	// Background with a marker 'O' where we want to draw something
	bg := tui.From(`
		+-------+
		|       |
		|   O   |
		|       |
		+-------+
	`)

	// Source to draw (a small box)
	src := tui.From(`
		+-+
		|X|
		+-+
	`)

	// Draw src centered on 'O'.
	// Anchor (0.5, 0.5) means the center of src is placed at the center of the 'O' marker.
	result := bg.Draw('O', src, 0.5, 0.5)

	fmt.Print(result.String())

	// Output:
	// +-------+
	// |  +-+  |
	// |  |X|  |
	// |  +-+  |
	// +-------+
}

func ExampleFixed_editorLayout() {
	// 1. Define the initial layout topology.
	// S: Sidebar
	// C: Content
	// T: Terminal
	layout := tui.From(`
		SC
		ST
	`)

	// 2. Scale up to 4x4 using Nearest Neighbor.
	// This establishes the regions with sufficient resolution (2x2 each).
	// Intermediate state:
	// SSCC
	// SSCC
	// SSTT
	// SSTT
	nnScalable := layout.NearestNeighborScale()
	intermediate := nnScalable(4, 4)

	// 3. Scale to final dimensions (8x5) using NinePatchScale on 'C'.
	// 'C' (Content) acts as the flexible region.
	// Since 'C' is in the top-right, the 'S' (Sidebar) columns will stay fixed width,
	// and the 'T' (Terminal) rows will stay fixed height relative to the bottom.
	npScalable := intermediate.NinePatchScale('C')
	final := npScalable(8, 5)

	fmt.Print(final.String())

	// Output:
	// SSCCCCCC
	// SSCCCCCC
	// SSCCCCCC
	// SSTTTTTT
	// SSTTTTTT
}
