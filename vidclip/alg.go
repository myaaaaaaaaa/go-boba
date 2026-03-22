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
