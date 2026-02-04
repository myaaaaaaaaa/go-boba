package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	err = generateScript(out, list, *outputVideo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating script: %v\n", err)
		os.Exit(1)
	}

	if scriptFile != nil {
		// Make script executable
		info, err := scriptFile.Stat()
		if err == nil {
			os.Chmod(*outputScript, info.Mode()|0111)
		}
	}
}

func quoteBash(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func generateScript(w io.Writer, list EditList, outputVideo string) error {
	fmt.Fprintln(w, "#!/bin/bash")
	fmt.Fprintln(w, "set -e")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "mkdir -p edl_segments")

	concatFile := "edl_concat.txt"
	fmt.Fprintf(w, "rm -f %s\n", quoteBash(concatFile))

	for i, entry := range list {
		ext := filepath.Ext(entry.Filename)
		if ext == "" {
			ext = ".mkv"
		}
		segmentName := fmt.Sprintf("edl_segments/part_%04d%s", i, ext)

		// ffmpeg -ss <start> -i <filename> -t <duration> -c copy <i>.<ext>
		fmt.Fprintf(w, "ffmpeg -y -ss %v -i %s -t %v -c copy %s\n",
			entry.Start, quoteBash(entry.Filename), entry.Duration, quoteBash(segmentName))

		// Add to concat file
		// format: file 'path'
		// We use double quotes for echo to allow expansion, but inside we want literal single quotes.
		// We use quoteBash(segmentName) which properly escapes single quotes and wraps the string in single quotes.
		// E.g. segmentName="foo'bar" -> "'foo'\''bar'"
		// echo "file 'foo'\''bar'" prints file 'foo'\''bar' which is what ffmpeg expects.
		fmt.Fprintf(w, "echo \"file %s\" >> %s\n", quoteBash(segmentName), quoteBash(concatFile))
	}

	fmt.Fprintln(w, "")
	fmt.Fprintf(w, "ffmpeg -y -f concat -safe 0 -i %s -c copy %s\n", quoteBash(concatFile), quoteBash(outputVideo))

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "# echo \"Cleaning up...\"")
	fmt.Fprintf(w, "# rm -rf edl_segments %s\n", quoteBash(concatFile))

	return nil
}
