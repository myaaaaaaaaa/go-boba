package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// TestMemoize12_Basic checks that calling the memoized function multiple times
// with the same input only invokes the underlying function once,
// ensuring caching works.
func TestMemoize12_Basic(t *testing.T) {
	var calls int32
	f := func(x int) (int, string) {
		calls++
		return x * 2, "ok"
	}

	memoized := memoize12(f)

	for range 5 {
		r1, r2 := memoized(5)
		if r1 != 10 || r2 != "ok" {
			t.Errorf("Expected 10, 'ok', got %v, %v", r1, r2)
		}
		if calls != 1 {
			t.Errorf("Expected 1 call, got %d", calls)
		}
	}

	for range 5 {
		r1, r2 := memoized(6)
		if r1 != 12 || r2 != "ok" {
			t.Errorf("Expected 12, 'ok', got %v, %v", r1, r2)
		}
		if calls != 2 {
			t.Errorf("Expected 2 calls for new argument, got %d", calls)
		}
	}
}

// TestMemoize12_Property leverages quick check to ensure that
// the memoized function returns the exact same results as the
// original un-memoized function across a wide range of random inputs.
func TestMemoize12_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	f := func(x int) (int, int) {
		return x * 2, x + 10
	}
	memoized := memoize12(f)

	for range 1000 {
		x := r.Int()
		f1, f2 := f(x)
		m1, m2 := memoized(x)
		if f1 != m1 || f2 != m2 {
			t.Fatalf("memoized(%d) = (%d, %d), expected (%d, %d)", x, m1, m2, f1, f2)
		}
	}
}

func float01(r *rand.Rand) float64 {
	switch r.Intn(10) {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return r.Float64()
	}
}

func TestSplitScrub_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))

	for range 10000 {
		var (
			width    = r.Intn(6)
			startPct = float01(r)
			endPct   = float01(r)
		)

		left, center, right := splitScrub(width, startPct, endPct)

		if left+center+right != width {
			t.Fatalf("FAIL (sum != width): width=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", width, startPct, endPct, left, center, right)
		}
		if left < 0 || center < 0 || right < 0 {
			t.Fatalf("FAIL (negative values): width=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", width, startPct, endPct, left, center, right)
		}
		if width != 0 && center < 1 {
			t.Fatalf("FAIL (center < 1): width=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", width, startPct, endPct, left, center, right)
		}
	}
}

func TestSplitTimeline_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))

	for range 10000 {
		width := r.Intn(20)

		var durations []float64
		for range r.Intn(10) + 1 {
			durations = append(durations, float01(r)*10)
		}

		segmentWidths := splitTimeline(width, durations)

		if len(segmentWidths) != len(durations) {
			t.Fatalf("expected %d segments, got %d", len(durations), len(segmentWidths))
		}

		sum := 0
		for i, w := range segmentWidths {
			if w < 0 {
				t.Fatalf("negative segment width: %d at index %d", w, i)
			}
			if w == 0 {
				availableWidth := width - (len(durations) - 1)
				if availableWidth >= len(durations) {
					t.Fatalf("segment width 0 for non-zero duration: %v, widths: %v", durations, segmentWidths)
				}
			}
			sum += w
		}

		expectedSum := width - (len(durations) - 1)
		expectedSum = max(expectedSum, 0)
		if sum != expectedSum {
			t.Fatalf("expected sum %d, got %d. width=%d, numClips=%d", expectedSum, sum, width, len(durations))
		}
	}
}

func parseDurationForTest(str string) float64 {
	var h, m, s float64
	fmt.Sscanf(str, "%f:%f:%f", &h, &m, &s)
	return h*3600 + m*60 + s
}

func TestFormatDuration_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))
	for range 10000 {
		// Test up to 99 hours: 99 hours * 3600 = 356400 seconds
		secs := r.Float64() * 356400

		formatted := formatDuration(secs)
		if len(formatted) != 10 {
			t.Fatalf("expected length 10, got %d (%q) for %f", len(formatted), formatted, secs)
		}

		parsed := parseDurationForTest(formatted)
		if math.Abs(parsed-secs) >= 0.5 {
			t.Fatalf("expected diff < 0.5, got %f for %f (formatted: %q)", math.Abs(parsed-secs), secs, formatted)
		}

		formatted2 := formatDuration(parsed)
		if formatted2 != formatted {
			t.Fatalf("%s != %s", formatted, formatted2)
		}
	}
}

func randomStringOf(r *rand.Rand, chars string) string {
	b := make([]byte, r.Intn(10))
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}

func TestCutZeroes_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))

	for range 10000 {
		leftGen := randomStringOf(r, "0:.")
		rightGen := "1" + randomStringOf(r, "19:.")
		switch r.Intn(10) {
		case 0:
			leftGen = ""
		case 1:
			rightGen = ""
		}

		input := leftGen + rightGen
		left, right := cutZeroes(input)

		if left+right != input {
			t.Fatalf("left+right != input: %q + %q != %q", left, right, input)
		}

		if left != leftGen || right != rightGen {
			t.Fatalf("cutZeroes did not return original strings: expected (%q, %q), got (%q, %q)", leftGen, rightGen, left, right)
		}
	}
}
