package main

import (
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

func TestSplitPct_Property(t *testing.T) {
	r := rand.New(rand.NewSource(4))

	float01 := func() float64 {
		switch r.Intn(10) {
		case 0:
			return 0
		case 1:
			return 1
		default:
			return r.Float64()
		}
	}

	for range 10000 {
		n := r.Intn(6)

		var (
			startPct = float01()
			endPct   = float01()
		)
		left, center, right := splitPct(n, startPct, endPct)

		if left+center+right != n {
			t.Fatalf("FAIL (sum != n): n=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", n, startPct, endPct, left, center, right)
		}
		if left < 0 || center < 0 || right < 0 {
			t.Fatalf("FAIL (negative values): n=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", n, startPct, endPct, left, center, right)
		}
		if n != 0 && center < 1 {
			t.Fatalf("FAIL (center < 1): n=%d, startPct=%f, endPct=%f => left=%d, center=%d, right=%d", n, startPct, endPct, left, center, right)
		}
	}
}
