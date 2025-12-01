package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	newlineSleep := flag.String("n", "0s", "newline sleep")
	clearSleep := flag.String("c", "100ms", "clear sleep")
	endingSleep := flag.String("e", "1s", "ending sleep")
	flag.Parse()

	n, err := time.ParseDuration(*newlineSleep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid newline sleep duration: %v\n", err)
		os.Exit(1)
	}
	c, err := time.ParseDuration(*clearSleep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid clear sleep duration: %v\n", err)
		os.Exit(1)
	}
	e, err := time.ParseDuration(*endingSleep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid ending sleep duration: %v\n", err)
		os.Exit(1)
	}

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
