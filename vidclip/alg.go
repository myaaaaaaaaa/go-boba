package main

import (
	"math"
	"slices"
)

func memoize12[A1 comparable, R1, R2 any](f func(A1) (R1, R2)) func(A1) (R1, R2) {
	cache := make(map[A1]func() (R1, R2))

	return func(arg A1) (R1, R2) {
		thunk, ok := cache[arg]
		if ok {
			return thunk()
		}

		r1, r2 := f(arg)
		thunk = func() (R1, R2) { return r1, r2 }
		cache[arg] = thunk
		return thunk()
	}
}

func splitScrub(width int, startPct, endPct float64) (left, center, right int) {
	if startPct > endPct {
		startPct, endPct = endPct, startPct
	}
	var (
		startIdx = int(startPct * float64(width))
		endIdx   = int(endPct * float64(width))
	)

	if startIdx == endIdx {
		if startIdx <= (width / 2) {
			endIdx++
		} else {
			startIdx--
		}
	}

	startIdx = max(startIdx, 0)
	endIdx = min(endIdx, width)

	left = startIdx
	center = endIdx - startIdx
	right = width - endIdx
	return
}

func splitTimeline(width int, durations []float64) []int {
	if len(durations) == 0 {
		return nil
	}

	width = width - (len(durations) - 1)
	width = max(width, 0)

	durations = slices.Clone(durations)
	var totalDuration float64
	for d := range durations {
		d := &durations[d]
		*d = max(*d, 0.001)
		totalDuration += *d
	}

	var accumulatedInt int

	var segmentWidths []int
	for _, d := range durations {
		w := int(math.Ceil(float64(width) * d / totalDuration))
		w = max(w, 1)
		segmentWidths = append(segmentWidths, w)
		accumulatedInt += w
	}

	for accumulatedInt > width {
		maxIdx := -1
		maxVal := -1
		for i, w := range segmentWidths {
			if w > maxVal {
				maxVal = w
				maxIdx = i
			}
		}

		segmentWidths[maxIdx]--
		accumulatedInt--
	}

	return segmentWidths
}
