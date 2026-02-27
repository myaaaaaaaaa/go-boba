package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	inputFile := flag.String("i", "", "Input EDL file")
	outputScript := flag.String("o", "", "Output bash script file (default: stdout)")
	outputVideo := flag.String("v", "output.mkv", "Final output video filename")
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file is required.")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	list, err := Parse(string(f))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing EDL: %v\n", err)
		os.Exit(1)
	}

	var out io.Writer
	var scriptFile *os.File
	if *outputScript == "" {
		out = os.Stdout
	} else {
		scriptFile, err = os.Create(*outputScript)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer scriptFile.Close()
		out = scriptFile
	}

	fmt.Fprint(out, list.Export(*outputVideo))

	if scriptFile != nil {
		// Make script executable
		info, err := scriptFile.Stat()
		if err == nil {
			os.Chmod(*outputScript, info.Mode()|0111)
		}
	}
}
