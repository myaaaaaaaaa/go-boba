package main

import (
	"os"
	"time"
)

func main() {
	for s := range splitBeforeSeq(os.Stdin, "\n\x1b") {
		time.Sleep(time.Millisecond * 200)
		os.Stdout.Write([]byte(s))
	}
}
