package main

import (
	"fmt"
	"io"
	"os"
)

func failIf(aProblem bool, format string, a ...any) {
	if aProblem {
		fmt.Fprintf(os.Stderr, format, a...)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

func tryPop(args *[]string, format string, a ...any) string {
	failIf(len(*args) == 0, format, a...)
	val := (*args)[0]
	*args = (*args)[1:]
	return val
}

func main() {
	args := os.Args
	exeName := tryPop(&args, "Error: no program name")
	subcmd := tryPop(&args, "Usage: %s <subcommand> [args...]", exeName)

	switch subcmd {
	case "export":
		f, err := io.ReadAll(os.Stdin)
		failIf(err != nil, "Error reading file: %v", err)

		list, err := Parse(string(f))
		failIf(err != nil, "Error parsing EDL: %v", err)

		fmt.Println(list.Export())
	case "mpv":
		file := tryPop(&args, "Usage: %s mpv <file>", exeName)
		t := 15.0
		mpvChooseTime(file, &t)
		fmt.Println("final time:", t)
	default:
		failIf(true, "Unknown subcommand: %s", subcmd)
	}
}
