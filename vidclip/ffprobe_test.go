package main

import "testing"

func TestParseProbeDuration(t *testing.T) {
	assert := func(s string, want float64) {
		t.Helper()

		got := parseProbeDuration([]byte(s))
		if got != want {
			t.Error("got", got)
			t.Error("want", want)
		}
	}

	assert(`{"format":{"duration":"100"}}`, 100)
	assert(`{"format":{"duration":"3.5"}}`, 3.5)
}
