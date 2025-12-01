package escapes

import (
	"io"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestFindEscapes_Unit(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single newline",
			input:    "\n",
			expected: []string{"", "\n"},
		},
		{
			name:     "single escape",
			input:    "\x1b",
			expected: []string{"", "\x1b"},
		},
		{
			name:     "text with newline",
			input:    "hello\nworld",
			expected: []string{"hello", "\nworld"},
		},
		{
			name:     "text with escape",
			input:    "hello\x1bworld",
			expected: []string{"hello", "\x1bworld"},
		},
		{
			name:     "multiple delimiters",
			input:    "a\nb\x1bc",
			expected: []string{"a", "\nb", "\x1bc"},
		},
		{
			name:     "leading delimiter",
			input:    "\nabc",
			expected: []string{"", "\nabc"},
		},
		{
			name:     "trailing delimiter",
			input:    "abc\n",
			expected: []string{"abc", "\n"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			actual := FindEscapes(r)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

// randRead is a custom io.Reader that randomly reads between 1-8 bytes.
type randRead struct {
	s    string
	rand *rand.Rand
}

func (r *randRead) Read(p []byte) (n int, err error) {
	if len(r.s) == 0 {
		return 0, io.EOF
	}
	// Simulate random read sizes.
	readSize := min(r.rand.Intn(8)+1, len(p), len(r.s))
	n = copy(p, r.s[:readSize])
	r.s = r.s[readSize:]
	return n, nil
}

// generateRandomString creates a random string with only the characters
// '\x1b', '\n', ' ', and 'a' to hit corner cases more easily.
func generateRandomString(length int, r *rand.Rand) string {
	chars := []byte{'\x1b', '\n', ' ', 'a'}
	var result []byte
	for range length {
		result = append(result, chars[r.Intn(len(chars))])
	}
	return string(result)
}

// referenceSplitBefore is a simple, non-streaming implementation for testing.
func referenceSplitBefore(s string) []string {
	var result []string
	var current strings.Builder
	for _, r := range s {
		isDelimiter := r == '\n' || r == '\x1b'
		if isDelimiter {
			result = append(result, current.String())
			current.Reset()
		}
		current.WriteRune(r)
	}
	result = append(result, current.String())
	return result
}

func TestFindEscapes_Properties(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := range 10000 {
		origString := generateRandomString(100, r)
		reader := &randRead{s: origString, rand: r}
		result := FindEscapes(reader)
		joined := strings.Join(result, "")

		if joined != origString {
			t.Fatalf("Failed on iteration %d. Original string: %q, joined: %q, result: %q", i, origString, joined, result)
		}

		expected := referenceSplitBefore(origString)
		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("Failed on iteration %d. Split logic does not match reference.\nOriginal: %q\nGot:      %q\nExpected: %q", i, origString, result, expected)
		}
	}
}

func TestFindEscapes_Properties_Structure(t *testing.T) {
	r := rand.New(rand.NewSource(2)) // Use a different fixed seed
	for i := range 10000 {
		origString := generateRandomString(100, r)
		reader := &randRead{s: origString, rand: r}
		result := FindEscapes(reader)

		if len(result) == 0 {
			t.Fatalf("Failed on iteration %d. Result should never be empty. Original: %q", i, origString)
		}

		// Check the first substring.
		if len(result[0]) > 0 {
			// If the original string starts with a delimiter, the first element of the result
			// will be an empty string, and the second will start with the delimiter.
			// Otherwise, the first element should not contain any delimiters.
			startsWithDelimiter := strings.HasPrefix(origString, "\n") || strings.HasPrefix(origString, "\x1b")
			if !startsWithDelimiter && strings.ContainsAny(result[0], "\n\x1b") {
				t.Fatalf("Failed on iteration %d. First substring should not contain delimiters. Original: %q, Got: %q", i, origString, result[0])
			}
		}

		// Check the rest of the substrings.
		for j := 1; j < len(result); j++ {
			s := result[j]
			if len(s) == 0 {
				t.Fatalf("Failed on iteration %d. Subsequent substrings should not be empty. Original: %q, Result: %q", i, origString, result)
			}
			isDelimiter := s[0] == '\n' || s[0] == '\x1b'
			if !isDelimiter {
				t.Fatalf("Failed on iteration %d. Subsequent substring should start with a delimiter. Original: %q, Substring: %q", i, origString, s)
			}
			if strings.ContainsAny(s[1:], "\n\x1b") {
				t.Fatalf("Failed on iteration %d. Subsequent substring should not contain delimiters after the first character. Original: %q, Substring: %q", i, origString, s)
			}
		}
	}
}
