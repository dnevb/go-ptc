package probe

import (
	"testing"
)

func TestDetectType(t *testing.T) {
	tests := []struct {
		codec string
		path  string
		want  MediaType
	}{
		{"png", "x.png", TypeImage},
		{"mjpeg", "x.jpg", TypeImage},
		{"gif", "x.gif", TypeGIF},
		{"h264", "x.mp4", TypeVideo},
		{"vp9", "x.webm", TypeVideo},
		{"", "x.mov", TypeVideo},
		{"", "x.png", TypeImage},
		{"unknown", "x.xyz", TypeUnknown},
	}
	for _, tt := range tests {
		got := detectType(tt.codec, tt.path)
		if got != tt.want {
			t.Errorf("detectType(%q, %q) = %q, want %q", tt.codec, tt.path, got, tt.want)
		}
	}
}

func TestParseFrames(t *testing.T) {
	tests := []struct {
		nbFrames string
		fps      string
		duration float64
		want     int
	}{
		{"100", "30/1", 10, 100},
		{"", "30/1", 10, 300},
		{"N/A", "0/0", 0, 1},
		{"", "", 5, 150},
		{"", "", 0, 1},
	}
	for _, tt := range tests {
		got := parseFrames(tt.nbFrames, tt.fps, tt.duration)
		if got != tt.want {
			t.Errorf("parseFrames(%q, %q, %f) = %d, want %d", tt.nbFrames, tt.fps, tt.duration, got, tt.want)
		}
	}
}

func TestProbeMissingFile(t *testing.T) {
	_, err := Probe("/nonexistent/path.png")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
