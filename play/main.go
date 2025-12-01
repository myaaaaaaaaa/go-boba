package main

import (
	"os"
	"regexp"
	"time"
)

func main() {
	os.Stdout.Write([]byte("\x1b[H\x1b[2J\x1b[3J"))
	os.Stdout.Write([]byte("\x1b[?25l"))
	defer os.Stdout.Write([]byte("\x1b[?25h"))

	isClear := regexp.MustCompile("^\x1b" + `\[\d*J`).MatchString
	for s := range splitBeforeSeq(os.Stdin, "\n\x1b") {
		if len(s) == 0 {
			continue
		}

		switch s[0] {
		case '\n':
			// time.Sleep(time.Millisecond)
		case '\x1b':
			if isClear(s) {
				time.Sleep(time.Millisecond * 100)
			}
		}
		os.Stdout.Write([]byte(s))
	}
	time.Sleep(time.Second)
}
