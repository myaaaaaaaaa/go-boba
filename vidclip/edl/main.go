package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
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
}
