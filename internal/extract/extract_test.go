package extract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseScale(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1920x1080", "scale=1920:1080:force_original_aspect_ratio=decrease"},
		{"", ""},
		{"bad", ""},
		{"0x0", ""},
		{"-1x100", ""},
	}
	for _, tt := range tests {
		got := parseScale(tt.input)
		if got != tt.want {
			t.Errorf("parseScale(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildFilters(t *testing.T) {
	tests := []struct {
		opts Options
		want string
	}{
		{Options{FPS: 30, MaxRes: "1920x1080"}, "fps=30,scale=1920:1080:force_original_aspect_ratio=decrease"},
		{Options{FPS: 30}, "fps=30"},
		{Options{MaxRes: "800x600"}, "scale=800:600:force_original_aspect_ratio=decrease"},
		{Options{}, ""},
	}
	for _, tt := range tests {
		got := buildFilters(tt.opts)
		if got != tt.want {
			t.Errorf("buildFilters(%+v) = %q, want %q", tt.opts, got, tt.want)
		}
	}
}

func TestListFramesSorts(t *testing.T) {
	// Only tests the listing helper without real ffmpeg.
	dir := t.TempDir()
	for _, name := range []string{"frame_0002.png", "frame_0001.png", "frame_0010.png"} {
		f, _ := os.Create(filepath.Join(dir, name))
		f.Close()
	}
	frames, err := listFrames(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(frames) != 3 {
		t.Fatalf("expected 3 frames, got %d", len(frames))
	}
	if !strings.HasSuffix(frames[0], "frame_0001.png") {
		t.Errorf("bad sort: first = %s", frames[0])
	}
	if !strings.HasSuffix(frames[2], "frame_0010.png") {
		t.Errorf("bad sort: last = %s", frames[2])
	}
}
