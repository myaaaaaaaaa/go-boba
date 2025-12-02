package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

//go:embed help.txt
var help string

func main() {
	type flags struct {
		n, c, e time.Duration
	}
	var f flags
	flag.DurationVar(&f.n, "n", 0, "newline sleep")
	flag.DurationVar(&f.c, "c", 0, "clear sleep")
	flag.DurationVar(&f.e, "e", 0, "ending sleep")
	flag.Parse()

	if (f == flags{}) {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file]\n", os.Args[0])
		fmt.Fprint(os.Stderr, help)
		flag.PrintDefaults()
		os.Exit(1)
	}

	os.Stdout.Write([]byte("\x1b[H\x1b[2J\x1b[3J"))
	os.Stdout.Write([]byte("\x1b[?25l"))
	defer os.Stdout.Write([]byte("\x1b[?25h"))

	out := bufio.NewWriterSize(os.Stdout, 65536)
	defer out.Flush()

	sleep := func(t time.Duration) {
		if t == 0 {
			return
		}
		out.Flush()
		time.Sleep(t)
	}

	isClear := regexp.MustCompile("^\x1b" + `\[\d*J`).MatchString
	for s := range splitBeforeSeq(os.Stdin, "\n\x1b") {
		if len(s) == 0 {
			continue
		}

		switch s[0] {
		case '\n':
			sleep(f.n)
		case '\x1b':
			if isClear(s) {
				sleep(f.c)
			}
		}
		out.Write([]byte(s))
	}
	sleep(f.e)
}
