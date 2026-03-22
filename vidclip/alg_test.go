package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"testing/quick"
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
	f := func(x int) (int, int) {
		return x * 2, x + 10
	}
	memoized := memoize12(f)

	prop := func(x int) bool {
		f1, f2 := f(x)
		m1, m2 := memoized(x)
		return f1 == m1 && f2 == m2
	}

	if err := quick.Check(prop, &quick.Config{MaxCount: 1000}); err != nil {
		t.Error(err)
	}
}

func TestSplitPct_Property(t *testing.T) {
	rawToFloat := func(raw uint32) float64 {
		switch raw % 10 {
		case 0:
			return 0.0
		case 1:
			return 1.0
		default:
			return float64(raw) / float64(math.MaxUint32)
		}
	}
	prop := func(rawN uint8, rawS, rawE uint32) bool {
		n := int(rawN % 6)

		var (
			startPct = rawToFloat(rawS)
			endPct   = rawToFloat(rawE)
		)

		left, center, right := splitPct(n, startPct, endPct)

		if left+center+right != n {
			fmt.Printf("FAIL (sum != n): n=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d\n", n, startPct, endPct, left, center, right)
			return false
		}
		if n != 0 && center < 1 {
			fmt.Printf("FAIL (center < 1): n=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d\n", n, startPct, endPct, left, center, right)
			return false
		}
		return true
	}

	config := quick.Config{
		MaxCount: 10000,
		Rand:     rand.New(rand.NewSource(4)),
	}
	if err := quick.Check(prop, &config); err != nil {
		t.Error(err)
	}
}
