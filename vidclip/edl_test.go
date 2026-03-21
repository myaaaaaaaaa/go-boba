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
	list, err := Parse(lines(
		`# mpv EDL v0`,
		`f1.mkv,10,20`,
		`/already/abs.mkv,10.5,1.5`,
		`dir/f2.mp4,0,5`,
	))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	want := lines(
		`# mpv EDL v0`,
		`/base/f1.mkv,10,20`,
		`/already/abs.mkv,10.5,1.5`,
		`/base/dir/f2.mp4,0,5`,
		``,
	)

	for range 3 {
		list.Absolute("/base")
		got := list.Serialize()
		if got != want {
			t.Errorf("Absolute got %v, expected %v", got, want)
		}
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

	wantHash := "ea768d8ded50d558fb4f11f0a42ac2d03b764ff0ec09ce5f7f52ee18059496e7"

	if gotHash != wantHash {
		t.Errorf("EditList.Export output hash changed.\nGot:  %s\nWant: %s\nOutput:\n%s", gotHash, wantHash, scriptContent)
	}
}
