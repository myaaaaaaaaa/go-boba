package tui

import "testing"

func TestStyledString(t *testing.T) {
	// Create a simple buffer
	buf := From(`
		AB
		CD
	`)

	// Apply styles
	buf = buf.MapCell('A', 'A', "<RED>")
	buf = buf.MapCell('D', 'D', "<BLUE>")

	// Call styledWithClose with a dummy close sequence
	got := buf.styledWithClose("</>")

	// Construct expected string
	// A becomes <RED>A</>
	// B is unstyled: B
	// C is unstyled: C
	// D becomes <BLUE>D</>
	// Newlines are preserved
	expected := `<RED>A</>B
C<BLUE>D</>`

	if got != expected {
		t.Errorf("styledWithClose() = %q, want %q", got, expected)
	}

	// Test merged styles
	bufMerged := From("AA")
	bufMerged = bufMerged.MapCell('A', 'A', "<RED>")
	gotMerged := bufMerged.styledWithClose("</>")
	// A and A should be merged into one style block: <RED>AA</>
	expectedMerged := "<RED>AA</>"

	if gotMerged != expectedMerged {
		t.Errorf("styledWithClose() merged = %q, want %q", gotMerged, expectedMerged)
	}
}

func TestStyleInfo(t *testing.T) {
	var style StyleInfo
	testCases := []struct {
		name  string
		style StyleInfo
		text  constString
		want  string
	}{
		{
			"Bold",
			style.Bold(),
			"Bold text",
			"\033[1mBold text\033[0m",
		},
		{
			"Dim",
			style.Dim(),
			"Dim text",
			"\033[2mDim text\033[0m",
		},
		{
			"Italic",
			style.Italic(),
			"Italic text",
			"\033[3mItalic text\033[0m",
		},
		{
			"Underline",
			style.Underline(),
			"Underline text",
			"\033[4mUnderline text\033[0m",
		},
		{
			"Blink",
			style.Blink(),
			"Blink text",
			"\033[5mBlink text\033[0m",
		},
		{
			"Reverse",
			style.Reverse(),
			"Reverse text",
			"\033[7mReverse text\033[0m",
		},
		{
			"Strikethrough",
			style.Strikethrough(),
			"Strikethrough text",
			"\033[9mStrikethrough text\033[0m",
		},
		/*
			{
				"Link",
				style.Link("https://example.com"),
				"Link text",
				"\033]8;;https://example.com\033\\Link text\033[0m",
			},
		*/
		{
			"AnsiColor",
			style.AnsiColor(123),
			"Color text",
			"\033[38;5;123mColor text\033[0m",
		},
		{
			"AnsiBg",
			style.AnsiBg(45),
			"Bg text",
			"\033[48;5;45mBg text\033[0m",
		},
		{
			"Color Hex Short",
			style.Color("#f00"),
			"Red",
			"\033[38;2;255;0;0mRed\033[0m",
		},
		{
			"Color Hex Long",
			style.Color("#00ff00"),
			"Green",
			"\033[38;2;0;255;0mGreen\033[0m",
		},
		{
			"Bg Hex",
			style.Bg("#0000ff"),
			"Blue",
			"\033[48;2;0;0;255mBlue\033[0m",
		},
		{
			"Chained",
			style.Bold().Underline(),
			"Both",
			"\033[1m\033[4mBoth\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := From(tc.text).MapCellFunc(func(c Cell) Cell {
				c.Style = tc.style
				return c
			})
			got := buf.StyledString()
			t.Log(got)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBuffer(t *testing.T) {
	width := 2
	height := 2
	callCount := 0

	// Create a Fixed that counts calls to At
	f := Fixed{
		Width:  width,
		Height: height,
		At: func(x, y int) Cell {
			callCount++
			return Cell{Rune: 'A'}
		},
	}

	// Buffer it
	buffered := f.Buffer()

	// At this point, At should have been called width * height times
	expectedCalls := width * height
	if callCount != expectedCalls {
		t.Errorf("Buffer() calls = %d, want %d", callCount, expectedCalls)
	}

	// Reset call count to ensure no more calls happen
	callCount = 0

	// Access buffered cells multiple times
	for i := 0; i < 3; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				_ = buffered.At(x, y)
			}
		}
	}

	if callCount != 0 {
		t.Errorf("Buffer() subsequent calls = %d, want 0", callCount)
	}

	// Test consistency
	simple := From("AB\nCD")
	bufferedSimple := simple.Buffer()
	if simple.Width != bufferedSimple.Width {
		t.Errorf("Width mismatch: got %d, want %d", bufferedSimple.Width, simple.Width)
	}
	if simple.Height != bufferedSimple.Height {
		t.Errorf("Height mismatch: got %d, want %d", bufferedSimple.Height, simple.Height)
	}
	for y := 0; y < simple.Height; y++ {
		for x := 0; x < simple.Width; x++ {
			if simple.At(x, y) != bufferedSimple.At(x, y) {
				t.Errorf("Content mismatch at %d,%d", x, y)
			}
		}
	}

	// Test 0x0 edge case
	empty := Fixed{}.Buffer()
	if empty.Width != 0 || empty.Height != 0 {
		t.Errorf("Expected empty buffer, got %v", empty)
	}
}
