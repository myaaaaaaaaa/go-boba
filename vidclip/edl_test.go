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
		{Source: "f1.mkv", Times: [2]float64{10, 20}},
		{Source: "f2.mkv", Times: [2]float64{10.5, 1.5}},
		{Source: "comma,file.mkv", Times: [2]float64{0, 5}},
		{Source: "video.mp4", Times: [2]float64{100, 50}},
		{Source: "audio.mp3", Times: [2]float64{0, 30.5}},
	}.validate()
}

func TestAbsolute(t *testing.T) {
	list := EditList{
		{Source: "f1.mkv", Times: [2]float64{10, 20}},
		{Source: "/already/abs.mkv", Times: [2]float64{10.5, 1.5}},
		{Source: "dir/f2.mp4", Times: [2]float64{0, 5}},
	}

	got := list.Absolute("/base")

	want := EditList{
		{Source: "/base/f1.mkv", Times: [2]float64{10, 20}},
		{Source: "/already/abs.mkv", Times: [2]float64{10.5, 1.5}},
		{Source: "/base/dir/f2.mp4", Times: [2]float64{0, 5}},
	}

	if len(got) != len(want) {
		t.Fatalf("Got length %d, want %d", len(got), len(want))
	}

	for i := range got {
		if got[i].Source != want[i].Source {
			t.Errorf("At index %d, got source %q, want %q", i, got[i].Source, want[i].Source)
		}
		if got[i].Times != want[i].Times {
			t.Errorf("At index %d, got times %v, want %v", i, got[i].Times, want[i].Times)
		}
	}

	// Verify original list wasn't modified
	if list[0].Source != "f1.mkv" {
		t.Errorf("Original list was modified! list[0].Source = %q", list[0].Source)
	}
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
