package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

type EditEntry struct {
	Source string
	Times  [2]float64
}

type EditList []EditEntry

func Parse(s string) (EditList, error) {
	header, s, _ := strings.Cut(s, "\n")
	if header != "# mpv EDL v0" {
		return nil, fmt.Errorf("invalid header: %s", header)
	}

	csvReader := csv.NewReader(bytes.NewBufferString(s))
	csvReader.FieldsPerRecord = 3
	csvReader.TrimLeadingSpace = true

	var list EditList

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		filename := record[0]
		start, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %s", record[1])
		}
		duration, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid duration: %s", record[2])
		}

		list = append(list, EditEntry{
			Source: filename,
			Times:  [2]float64{start, start + duration},
		})
	}

	return list, nil
}

func (e EditList) Serialize() string {
	lines := []string{"# mpv EDL v0"}

	for _, entry := range e {
		filename := entry.Source
		if strings.Contains(filename, ",") {
			filename = "\"" + filename + "\""
		}
		lines = append(lines, fmt.Sprintf("%s,%g,%g",
			filename,
			entry.Times[0],
			entry.Times[1]-entry.Times[0],
		))
	}

	return strings.Join(lines, "\n") + "\n"
}

func (e EditList) validate() {
	str := e.Serialize()
	parsed, err := Parse(str)
	if err != nil {
		panic(err)
	}

	//parsed[0].Duration++

	if fmt.Sprint(e) != fmt.Sprint(parsed) {
		panic(fmt.Sprintf("%v != %v", e, parsed))
	}
}

func quoteBash(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func (e EditList) Absolute(baseDir string) EditList {
	newList := make(EditList, len(e))
	for i, entry := range e {
		if !filepath.IsAbs(entry.Source) {
			entry.Source = filepath.Join(baseDir, entry.Source)
		}
		newList[i] = entry
	}
	return newList
}

func (list EditList) Export() string {
	var w strings.Builder
	fmt.Fprintln(&w, "#!/bin/bash")
	fmt.Fprintln(&w, "set -e")
	fmt.Fprintln(&w, "")
	fmt.Fprintln(&w, "mkdir -p edl_segments")

	concatFile := "edl_concat.txt"
	fmt.Fprintf(&w, "rm -f %s\n", quoteBash(concatFile))

	for i, entry := range list {
		ext := filepath.Ext(entry.Source)
		if ext == "" {
			ext = ".mkv"
		}
		segmentName := fmt.Sprintf("edl_segments/part_%04d%s", i, ext)

		// ffmpeg -ss <start> -to <end> -i <filename> -c copy <i>.<ext>
		fmt.Fprintf(&w, "ffmpeg -y -ss %v -to %v -i %s -c copy %s\n",
			entry.Times[0], entry.Times[1], quoteBash(entry.Source), quoteBash(segmentName))

		// Add to concat file
		// format: file 'path'
		// We use double quotes for echo to allow expansion, but inside we want literal single quotes.
		// We use quoteBash(segmentName) which properly escapes single quotes and wraps the string in single quotes.
		// E.g. segmentName="foo'bar" -> "'foo'\''bar'"
		// echo "file 'foo'\''bar'" prints file 'foo'\''bar' which is what ffmpeg expects.
		fmt.Fprintf(&w, "echo \"file %s\" >> %s\n", quoteBash(segmentName), quoteBash(concatFile))
	}

	fmt.Fprintln(&w, "")
	fmt.Fprintf(&w, "ffmpeg -y -f concat -safe 0 -i %s -c copy \"$1\"\n", quoteBash(concatFile))

	fmt.Fprintln(&w, "")
	fmt.Fprintln(&w, "# echo \"Cleaning up...\"")
	fmt.Fprintf(&w, "# rm -rf edl_segments %s\n", quoteBash(concatFile))

	return w.String()
}
