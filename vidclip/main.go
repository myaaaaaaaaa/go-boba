package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "export":
		f, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		list, err := Parse(string(f))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing EDL: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(list.Export())
	case "mpv":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s mpv <file>\n", os.Args[0])
			os.Exit(1)
		}
		t := 15.0
		mpvChooseTime(os.Args[2], &t)
		fmt.Println("final time:", t)
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}
