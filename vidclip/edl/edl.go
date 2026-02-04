package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type EditEntry struct {
	Filename string
	Start    float64
	Duration float64
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
			Filename: filename,
			Start:    start,
			Duration: duration,
		})
	}

	return list, nil
}

func (e EditList) Serialize() string {
	lines := []string{"# mpv EDL v0"}

	for _, entry := range e {
		filename := entry.Filename
		if strings.Contains(filename, ",") {
			filename = "\"" + filename + "\""
		}
		lines = append(lines, fmt.Sprintf("%s,%g,%g",
			filename,
			entry.Start,
			entry.Duration,
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
