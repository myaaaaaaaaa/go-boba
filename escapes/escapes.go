package escapes

import (
	"bytes"
	"io"
)

// FindEscapes reads from an io.Reader and splits the stream before each
// occurrence of a newline '\n' or escape '\x1b' character. This is useful
// for parsing terminal output, which often contains ANSI escape codes.
// The function processes the input in a streaming manner, so it doesn't
// need to load the entire input into memory.
func FindEscapes(r io.Reader) []string {
	var result []string
	var current bytes.Buffer
	buf := make([]byte, 64)

	for {
		n, err := r.Read(buf)

		if n > 0 {
			chunk := buf[:n]
			for _, b := range chunk {
				isDelimiter := b == '\n' || b == '\x1b'
				if isDelimiter {
					result = append(result, current.String())
					current.Reset()
				}
				current.WriteByte(b)
			}
		}

		if err == io.EOF {
			// If the input is empty, we should return a slice with an empty string,
			// similar to strings.Split("", ",").
			if len(result) == 0 && current.Len() == 0 {
				return []string{""}
			}
			// Add the last part.
			if current.Len() > 0 {
				result = append(result, current.String())
			}
			break
		}
		if err != nil {
			// Some other error. Add what we have and return.
			if current.Len() > 0 {
				result = append(result, current.String())
			}
			break
		}
	}
	return result
}
