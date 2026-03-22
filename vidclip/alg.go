package main

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

	segmentWidths := make([]int, len(durations))
	var totalDuration float64
	for _, d := range durations {
		totalDuration += d
	}

	if totalDuration == 0 {
		for i := range segmentWidths {
			segmentWidths[i] = width / len(durations)
		}
		sum := 0
		for _, w := range segmentWidths {
			sum += w
		}
		segmentWidths[len(segmentWidths)-1] += width - sum
		return segmentWidths
	}

	var accumulatedFloat float64
	var accumulatedInt int

	for i, d := range durations {
		accumulatedFloat += (d / totalDuration) * float64(width)
		expected := int(accumulatedFloat + 0.5)
		w := expected - accumulatedInt
		if w <= 0 && d > 0 {
			w = 1
		}
		w = max(w, 0)
		segmentWidths[i] = w
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
		if maxVal <= 0 {
			break
		}
		segmentWidths[maxIdx]--
		accumulatedInt--
	}

	for accumulatedInt < width {
		maxIdx := -1
		maxVal := -1
		for i, w := range segmentWidths {
			if w > maxVal {
				maxVal = w
				maxIdx = i
			}
		}
		if maxIdx == -1 {
			maxIdx = 0
		}
		segmentWidths[maxIdx]++
		accumulatedInt++
	}

	return segmentWidths
}
