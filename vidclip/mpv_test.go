package main

import (
	"fmt"
	"testing"
)

func TestParseTimePosEvent(t *testing.T) {
	var events []float64
	cb := func(event float64) {
		events = append(events, event)
	}

	parseTimePosEvent(`{}`, cb)
	parseTimePosEvent(`{"event":"property-change","id":1,"name":"time-pos","data":2.5}`, cb)
	parseTimePosEvent(`{"event":"property-FLARGH","id":1,"name":"time-pos","data":3.5}`, cb)
	parseTimePosEvent(`{"event":"property-change","id":1,"name":"time-FOO","data":4.5}`, cb)
	parseTimePosEvent(`{"event":"property-change","id":1,"name":"time-pos"`, cb)
	parseTimePosEvent(`{"event":"property-change","id":1,"name":"time-pos","data":6.5}`, cb)
	parseTimePosEvent(`{"event":"property-change","id":1,"name":"time-pos","data":"hi"}`, cb)

	got := fmt.Sprint(events)
	want := "[2.5 6.5 0]"
	if got != want {
		t.Error("got", got)
		t.Error("want", want)
	}
}
