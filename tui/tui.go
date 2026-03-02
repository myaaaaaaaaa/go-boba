package tui

import (
	"fmt"
	"strings"
)

// Package tui provides a simple text-based user interface library.
//
// Pattern: Bounding Box
// Several functions in this package (NinePatchScale, Fill) identify regions of interest
// by finding the bounding box of all cells containing a specific marker rune.
// This allows for flexible layout definitions directly within the string representation
// of the buffer. For example, a 9-patch center region or a target area for drawing
// can be defined by simply filling the desired area with a specific character.
//
// Error Handling:
// This package follows the philosophy that programmer errors should cause panics,
// while user errors (runtime conditions depending on external input) should return errors.
// Functions in this package will panic if called with invalid arguments that indicate
// a logic error in the calling code (e.g., missing 9-patch center markers).

// StyleInfo contains the raw escape codes that will be prepended to the output string.
type StyleInfo string

func (s StyleInfo) Bold() StyleInfo          { return s.appendEscapef("[1m") }
func (s StyleInfo) Dim() StyleInfo           { return s.appendEscapef("[2m") }
func (s StyleInfo) Italic() StyleInfo        { return s.appendEscapef("[3m") }
func (s StyleInfo) Underline() StyleInfo     { return s.appendEscapef("[4m") }
func (s StyleInfo) Blink() StyleInfo         { return s.appendEscapef("[5m") }
func (s StyleInfo) Reverse() StyleInfo       { return s.appendEscapef("[7m") }
func (s StyleInfo) Strikethrough() StyleInfo { return s.appendEscapef("[9m") }

func (s StyleInfo) AnsiColor(n int) StyleInfo  { return s.appendEscapef("[38;5;%dm", n) }
func (s StyleInfo) AnsiBg(n int) StyleInfo     { return s.appendEscapef("[48;5;%dm", n) }
func (s StyleInfo) Color(hex string) StyleInfo { return s.appendTruecolor("38", hex) }
func (s StyleInfo) Bg(hex string) StyleInfo    { return s.appendTruecolor("48", hex) }

// Currently broken
//func (s StyleInfo) Link(url string) StyleInfo { return s.appendEscapef("]8;;%s\033\\", url) }

func (s StyleInfo) appendTruecolor(code, hex string) StyleInfo {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b uint8
	if len(hex) == 3 {
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r, g, b = r*17, g*17, b*17
	} else {
		fmt.Sscanf(hex, "%2x%2x%2x", &r, &g, &b)
	}
	return s.appendEscapef("[%s;2;%d;%d;%dm", code, r, g, b)
}

func (s StyleInfo) appendEscapef(format string, args ...any) StyleInfo {
	return s + StyleInfo(fmt.Sprintf("\033"+format, args...))
}

// Cell is the atomic unit of the grid.
type Cell struct {
	Style StyleInfo
	Rune  rune
}

// Fixed holds the grid of cells.
// At() should not be called with an out-of-bounds coordinate.
type Fixed struct {
	Width  int
	Height int
	At     func(x, y int) Cell
}

// constString is a local type to enforce that From is called with a string constant.
type constString string

// From creates a Fixed from an ASCII string.
// It ignores tab characters to allow for indentation in the source code.
// Newlines separate rows.
func From(art constString) Fixed {
	// Remove tabs
	s := strings.ReplaceAll(string(art), "\t", "")

	// Split into lines
	lines := strings.Split(s, "\n")

	// Determine dimensions.
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	grid := make([][]Cell, height)
	for y, line := range lines {
		row := make([]Cell, width)
		runes := []rune(line)
		for x := 0; x < width; x++ {
			var r rune = ' '
			if x < len(runes) {
				r = runes[x]
			}
			row[x] = Cell{Rune: r}
		}
		grid[y] = row
	}

	return Fixed{
		Width:  width,
		Height: height,
		At:     func(x, y int) Cell { return grid[y][x] },
	}
}

// Contains returns true if the given coordinate is within the buffer's bounds.
func (b Fixed) Contains(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// findBoundingBox returns the min and max coordinates of all cells containing the target rune.
// It panics if the rune is not present in the buffer.
func (b Fixed) findBoundingBox(target rune) (minX, minY, maxX, maxY int) {
	minX, minY = b.Width, b.Height
	maxX, maxY = -1, -1
	found := false

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			if b.At(x, y).Rune == target {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				found = true
			}
		}
	}

	if !found {
		panic(fmt.Sprintf("findBoundingBox: rune %q not found in buffer", target))
	}
	return
}

// Scalable is a function that returns a new Fixed scaled to the given dimensions.
type Scalable func(width, height int) Fixed

// Shader returns a Scalable that generates a buffer using the provided function.
// The function f is called for each cell with x and y coordinates normalized to [0, 1).
func Shader(f func(x, y float64) Cell) Scalable {
	return func(width, height int) Fixed {
		return Fixed{
			Width:  width,
			Height: height,
			At: func(x, y int) Cell {
				u := float64(x) / float64(width)
				v := float64(y) / float64(height)
				return f(u, v)
			},
		}
	}
}

// NinePatchScale returns a Scalable function that scales the buffer using 9-patch mechanics.
// The center region is defined by the bounding box of all cells containing centerRune.
// The centerRune itself is preserved in the output (it is tiled).
// Use MapCell afterwards to change it if needed.
//
// This function panics if centerRune is not found in the buffer.
func (b Fixed) NinePatchScale(centerRune rune) Scalable {
	// Find bounding box of centerRune
	minX, minY, maxX, maxY := b.findBoundingBox(centerRune)

	// Source regions
	srcLeftW := minX
	srcCenterW := maxX - minX + 1
	srcTopH := minY
	srcCenterH := maxY - minY + 1

	return func(w, h int) Fixed {
		// Destination regions
		destCenterW := max(w-srcLeftW-(b.Width-(maxX+1)), 0)
		destCenterH := max(h-srcTopH-(b.Height-(maxY+1)), 0)

		return Fixed{
			Width:  w,
			Height: h,
			At: func(x, y int) Cell {
				// Determine which region we are in (Source X and Source Y)
				var sx, sy int

				// Vertical mapping
				if y < srcTopH {
					sy = y
				} else if y >= srcTopH+destCenterH {
					sy = srcTopH + srcCenterH + (y - (srcTopH + destCenterH))
					// Clamp to valid source range
					if sy >= b.Height {
						sy = b.Height - 1
					}
				} else {
					// In the center vertical region.
					// Map to source center region (tile/repeat)
					offset := y - srcTopH
					sy = srcTopH + (offset % srcCenterH)
				}

				// Horizontal mapping
				if x < srcLeftW {
					sx = x
				} else if x >= srcLeftW+destCenterW {
					sx = srcLeftW + srcCenterW + (x - (srcLeftW + destCenterW))
					if sx >= b.Width {
						sx = b.Width - 1
					}
				} else {
					// In the center horizontal region.
					offset := x - srcLeftW
					sx = srcLeftW + (offset % srcCenterW)
				}

				return b.At(sx, sy)
			},
		}
	}
}

// NearestNeighborScale returns a Scalable that scales the buffer using nearest-neighbor interpolation.
func (b Fixed) NearestNeighborScale() Scalable {
	return func(width, height int) Fixed {
		if width <= 0 || height <= 0 {
			return Fixed{Width: 0, Height: 0, At: func(x, y int) Cell { return Cell{} }}
		}

		return Fixed{
			Width:  width,
			Height: height,
			At: func(x, y int) Cell {
				// Map destination Y to source Y
				srcY := (y * b.Height) / height
				if srcY >= b.Height {
					srcY = b.Height - 1
				}

				// Map destination X to source X
				srcX := (x * b.Width) / width
				if srcX >= b.Width {
					srcX = b.Width - 1
				}

				return b.At(srcX, srcY)
			},
		}
	}
}

// Fill scales the source buffer using the scalable to match the bounding box of targetBox runes in the receiver buffer.
// The scaled buffer is then drawn on top of the receiver buffer at the top-left corner of the bounding box.
// It returns a new Fixed with the changes.
// It panics if targetBox is not found in the receiver buffer.
func (b Fixed) Fill(targetBox rune, src Scalable) Fixed {
	// Find bounding box of targetBox
	minX, minY, maxX, maxY := b.findBoundingBox(targetBox)

	// Calculate width and height of the target box
	width := maxX - minX + 1
	height := maxY - minY + 1

	// Scale the source buffer
	scaled := src(width, height)

	return Fixed{
		Width:  b.Width,
		Height: b.Height,
		At: func(x, y int) Cell {
			// Check if we are within the target box region
			// and if the relative coordinate is within the scaled buffer
			relX := x - minX
			relY := y - minY

			if scaled.Contains(relX, relY) {
				return scaled.At(relX, relY)
			}
			return b.At(x, y)
		},
	}
}

// Draw overlays the src Fixed onto the receiver Fixed.
// The placement is determined by aligning the point (anchorx, anchory) of the src
// with the center of the bounding box of all targetPoint runes in the receiver.
// anchorx and anchory are normalized coordinates [0.0, 1.0].
// It returns a new Fixed with the changes.
// It panics if targetPoint is not found in the receiver buffer.
func (b Fixed) Draw(targetPoint rune, src Fixed, anchorx, anchory float64) Fixed {
	// Find bounding box of targetPoint
	minX, minY, maxX, maxY := b.findBoundingBox(targetPoint)

	// Calculate center of the target bounding box
	targetCenterX := float64(minX) + float64(maxX-minX+1)/2.0
	targetCenterY := float64(minY) + float64(maxY-minY+1)/2.0

	// Calculate the anchor point in src coordinates
	srcAnchorX := float64(src.Width) * anchorx
	srcAnchorY := float64(src.Height) * anchory

	// Calculate the top-left position of src on the receiver grid
	startX := int(targetCenterX - srcAnchorX)
	startY := int(targetCenterY - srcAnchorY)

	return Fixed{
		Width:  b.Width,
		Height: b.Height,
		At: func(x, y int) Cell {
			// Check if we are within the src region relative to startX, startY
			srcX := x - startX
			srcY := y - startY

			if src.Contains(srcX, srcY) {
				return src.At(srcX, srcY)
			}
			return b.At(x, y)
		},
	}
}

// Buffer returns a new Fixed where the contents of the given Fixed are memoized.
// This is useful when the At function of the input Fixed is expensive to compute.
func (f Fixed) Buffer() Fixed {
	width := f.Width
	height := f.Height

	grid := make([][]Cell, height)
	for y := range height {
		row := make([]Cell, width)
		for x := range width {
			row[x] = f.At(x, y)
		}
		grid[y] = row
	}

	return Fixed{
		Width:  width,
		Height: height,
		At: func(x, y int) Cell {
			return grid[y][x]
		},
	}
}

// MapCellFunc returns a new Fixed where each cell is transformed by the given function.
func (b Fixed) MapCellFunc(f func(Cell) Cell) Fixed {
	return Fixed{
		Width:  b.Width,
		Height: b.Height,
		At: func(x, y int) Cell {
			return f(b.At(x, y))
		},
	}
}

// MapCell returns a new Fixed where all occurrences of targetRune are replaced by newRune.
// If style is not empty, it is applied to the matching cells.
func (b Fixed) MapCell(targetRune, newRune rune, style StyleInfo) Fixed {
	return b.MapCellFunc(func(cell Cell) Cell {
		if cell.Rune == targetRune {
			s := cell.Style
			if style != "" {
				s = style
			}
			return Cell{Rune: newRune, Style: s}
		}
		return cell
	})
}

// String returns a simple string representation of the buffer for debugging/testing.
func (b Fixed) String() string {
	return b.MapCellFunc(func(c Cell) Cell {
		c.Style = ""
		return c
	}).styledWithClose("")
}

func (b Fixed) styledWithClose(closeSequence string) string {
	var sb strings.Builder
	for y := 0; y < b.Height; y++ {
		var currentStyle StyleInfo
		for x := 0; x < b.Width; x++ {
			cell := b.At(x, y)

			if cell.Style != currentStyle {
				if currentStyle != "" {
					sb.WriteString(closeSequence)
				}
				if cell.Style != "" {
					sb.WriteString(string(cell.Style))
				}
				currentStyle = cell.Style
			}

			sb.WriteRune(cell.Rune)
		}
		// Close style at end of row
		if currentStyle != "" {
			sb.WriteString(closeSequence)
		}

		if y < b.Height-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// StyledString returns the string representation of the buffer with ANSI escape codes applied.
func (b Fixed) StyledString() string {
	return b.styledWithClose("\033[0m")
}

// WrapText returns a Scalable that word wraps the given string to fit into the given size,
// truncating to ellipses on overflow.
func WrapText(text string) Scalable {
	return func(width, height int) Fixed {
		if width <= 0 || height <= 0 {
			return Fixed{
				Width:  0,
				Height: 0,
				At:     func(x, y int) Cell { return Cell{} },
			}
		}

		words := strings.Fields(text)
		grid := make([][]Cell, height)
		// Initialize grid with spaces
		for y := range grid {
			row := make([]Cell, width)
			for x := range row {
				row[x] = Cell{Rune: ' '}
			}
			grid[y] = row
		}

		if len(words) == 0 {
			return Fixed{
				Width:  width,
				Height: height,
				At:     func(x, y int) Cell { return grid[y][x] },
			}
		}

		var lines []string
		var currentLine string

		for _, word := range words {
			wLen := len(word)
			lineLen := len(currentLine)
			spaceNeeded := 0
			if lineLen > 0 {
				spaceNeeded = 1
			}

			// If the word fits on the current line
			if lineLen+spaceNeeded+wLen <= width {
				if spaceNeeded > 0 {
					currentLine += " "
				}
				currentLine += word
			} else {
				// Flush current line if not empty
				if currentLine != "" {
					lines = append(lines, currentLine)
				}

				// Handle the word itself
				remaining := word
				// If the word is longer than the width, split it across lines
				for len(remaining) > width {
					lines = append(lines, remaining[:width])
					remaining = remaining[width:]
				}
				currentLine = remaining
			}
		}
		if currentLine != "" {
			lines = append(lines, currentLine)
		}

		// Handle height constraint and ellipses
		if len(lines) > height {
			lines = lines[:height]
			lastLine := lines[height-1]

			dots := "..."
			if width < 3 {
				dots = strings.Repeat(".", width)
			}

			// Truncate to fit ellipses
			if len(lastLine)+len(dots) > width {
				keep := max(width-len(dots), 0)
				lastLine = lastLine[:keep] + dots
			} else {
				lastLine += dots
			}
			lines[height-1] = lastLine
		}

		// Fill the grid
		for y, line := range lines {
			runes := []rune(line)
			for x, r := range runes {
				if x < width {
					grid[y][x] = Cell{Rune: r}
				}
			}
		}

		return Fixed{
			Width:  width,
			Height: height,
			At:     func(x, y int) Cell { return grid[y][x] },
		}
	}
}
