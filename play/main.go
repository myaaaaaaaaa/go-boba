package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	var n, c, e time.Duration
	flag.DurationVar(&n, "n", 0, "newline sleep")
	flag.DurationVar(&c, "c", 100*time.Millisecond, "clear sleep")
	flag.DurationVar(&e, "e", 1*time.Second, "ending sleep")
	flag.Parse()

	if n == 0 && c == 0 && e == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Plays back terminal animations.")
		fmt.Fprintln(os.Stderr, "Pauses on newlines and clear screen escape codes.")
		fmt.Fprintln(os.Stderr)
		flag.PrintDefaults()
		os.Exit(1)
	}

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
			time.Sleep(n)
		case '\x1b':
			if isClear(s) {
				time.Sleep(c)
			}
		}
		os.Stdout.Write([]byte(s))
	}
	time.Sleep(e)
}
