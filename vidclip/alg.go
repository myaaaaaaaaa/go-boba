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

func splitPct(n int, startPct, endPct float64) (left, center, right int) {
	if startPct > endPct {
		startPct, endPct = endPct, startPct
	}
	var (
		startIdx = int(startPct * float64(n))
		endIdx   = int(endPct * float64(n))
	)

	if startIdx == endIdx {
		if startIdx <= (n / 2) {
			endIdx++
		} else {
			startIdx--
		}
	}

	startIdx = max(startIdx, 0)
	endIdx = min(endIdx, n)

	left = startIdx
	center = endIdx - startIdx
	right = n - endIdx
	return
}
