package main

import (
	"os"
	"regexp"
	"time"
)

func main() {
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
				time.Sleep(time.Millisecond * 300)
			}
		}
		os.Stdout.Write([]byte(s))
	}
}
