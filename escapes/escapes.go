package escapes

import (
	"io"
	"iter"
	"slices"
)

// FindEscapes reads from an io.Reader and splits the stream before each
// occurrence of a newline '\n' or escape '\x1b' character. This is useful
// for parsing terminal output, which often contains ANSI escape codes.
// The function processes the input in a streaming manner, so it doesn't
// need to load the entire input into memory.
func FindEscapes(r io.Reader) []string {
	return slices.Collect(splitBeforeSeq(r, "\n\x1b"))
}

func splitBeforeSeq(r io.Reader, splitChars string) iter.Seq[string] {
	var splitSet [256]bool
	for _, c := range []byte(splitChars) {
		splitSet[c] = true
	}

	return func(yield func(string) bool) {
		substr := make([]byte, 0, 1024)
		var readBuffer [1024]byte

		for {
			n, err := r.Read(readBuffer[:])
			readBuffer := readBuffer[:n]

			for _, b := range readBuffer {
				if splitSet[b] {
					if !yield(string(substr)) {
						return
					}
					substr = substr[:0]
				}
				substr = append(substr, b)
			}
			if err != nil {
				break
			}
		}

		// Add the last part.
		if len(substr) > 0 {
			yield(string(substr))
		}
	}
}
