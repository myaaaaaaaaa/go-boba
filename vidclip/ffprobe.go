package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

func videoDuration(filename string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		filename)
	output, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return 0, fmt.Errorf("ffprobe: %v: %s", err, exitErr.Stderr)
		}
		return 0, err
	}

	return parseProbeDuration(output), nil
}

func parseProbeDuration(b []byte) float64 {
	var entries struct {
		Format struct {
			Duration float64 `json:",string"`
		}
	}

	err := json.Unmarshal([]byte(b), &entries)
	if err != nil || entries.Format.Duration == 0 {
		panic("invalid output from ffprobe: " + string(b) + "\n" + err.Error())
	}

	return entries.Format.Duration
}
