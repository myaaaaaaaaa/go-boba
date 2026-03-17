package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func lines(s ...string) string {
	return strings.Join(s, "\n")
}

func TestParse(t *testing.T) {
	parsed, err := Parse(lines(
		`# mpv EDL v0`,
		``,
		``,
		`f1.mkv,10,20`,
		`f2.mkv,10.5,1.5`,
		``,
		`f1.mkv,40,10`,
	))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	parsed.validate()

	got := parsed.Serialize()
	want := lines(
		`# mpv EDL v0`,
		`f1.mkv,10,20`,
		`f2.mkv,10.5,1.5`,
		`f1.mkv,40,10`,
		``,
	)

	if got != want {
		t.Errorf("Parse got %v, expected %v", got, want)
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"BadHeader", "# mpv EDL v1\nf1.mkv,10,20"},
		{"BadFormat", "# mpv EDL v0\nf1.mkv,10"},
		{"BadStart", "# mpv EDL v0\nf1.mkv,foo,20"},
		{"BadDuration", "# mpv EDL v0\nf1.mkv,10,bar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestRoundtrip(t *testing.T) {
	EditList{
		{Source: "f1.mkv", StartTime: 10, EndTime: 20},
		{Source: "f2.mkv", StartTime: 10.5, EndTime: 1.5},
		{Source: "comma,file.mkv", StartTime: 0, EndTime: 5},
		{Source: "video.mp4", StartTime: 100, EndTime: 50},
		{Source: "audio.mp3", StartTime: 0, EndTime: 30.5},
	}.validate()
}

func TestExportGolden(t *testing.T) {
	list, err := Parse(lines(
		`# mpv EDL v0`,
		`f1.mkv,10,20`,
		`f2.mp4,10.5,1.5`,
		`"comma,file.mkv",0,5`,
		`weird'name.mkv,100,50`,
		`no_extension,0,30.5`,
	))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	scriptContent := list.Export()

	hash := sha256.Sum256([]byte(scriptContent))
	gotHash := hex.EncodeToString(hash[:])

	wantHash := "21c17d2fc93de2389cf9e104a747c86dc2c6faa79ac4e84d74ffa2a094d098c5"

	if gotHash != wantHash {
		t.Errorf("EditList.Export output hash changed.\nGot:  %s\nWant: %s\nOutput:\n%s", gotHash, wantHash, scriptContent)
	}
}
