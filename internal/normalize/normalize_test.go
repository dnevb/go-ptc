package normalize

import (
	"strings"
	"testing"
)

func TestBuildFilter(t *testing.T) {
	tests := []struct {
		opts Options
		want string
	}{
		{Options{MaxColors: 256, MaxW: 1920, MaxH: 1080}, "scale=1920:1080:force_original_aspect_ratio=decrease,split[s0][s1];[s0]palettegen=max_colors=256[p];[s1][p]paletteuse"},
		{Options{MaxColors: 128}, "split[s0][s1];[s0]palettegen=max_colors=128[p];[s1][p]paletteuse"},
		{Options{MaxColors: 256, MaxW: 0, MaxH: 0}, "split[s0][s1];[s0]palettegen=max_colors=256[p];[s1][p]paletteuse"},
	}
	for _, tt := range tests {
		got := buildFilter(tt.opts)
		if got != tt.want {
			t.Errorf("buildFilter(%+v) = %q, want %q", tt.opts, got, tt.want)
		}
	}
}

func TestNormalizeMissingFile(t *testing.T) {
	err := Normalize("/nonexistent/file.png", Options{MaxColors: 256})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "ffmpeg") {
		t.Errorf("expected ffmpeg error, got: %v", err)
	}
}
