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
		{Filename: "f1.mkv", Start: 10, Duration: 20},
		{Filename: "f2.mkv", Start: 10.5, Duration: 1.5},
		{Filename: "comma,file.mkv", Start: 0, Duration: 5},
		{Filename: "video.mp4", Start: 100, Duration: 50},
		{Filename: "audio.mp3", Start: 0, Duration: 30.5},
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

	scriptContent := list.Export("final_output.mkv")

	hash := sha256.Sum256([]byte(scriptContent))
	gotHash := hex.EncodeToString(hash[:])

	wantHash := "e3a7faa9cefc9bb29332f544ab703e3b90eed5168b2f64e08c8b8a90bc67ef97"

	if gotHash != wantHash {
		t.Errorf("EditList.Export output hash changed.\nGot:  %s\nWant: %s\nOutput:\n%s", gotHash, wantHash, scriptContent)
	}
}
