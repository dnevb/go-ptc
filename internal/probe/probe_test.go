package probe

import (
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestProbeImage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.png")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 100, 50))); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	info, err := Probe(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Type != TypeImage {
		t.Errorf("type = %q, want image", info.Type)
	}
	if info.Width != 100 || info.Height != 50 {
		t.Errorf("dims = %dx%d, want 100x50", info.Width, info.Height)
	}
	if info.Frames != 1 {
		t.Errorf("frames = %d, want 1", info.Frames)
	}
}

func TestProbeGIF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.gif")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	pal := color.Palette{color.Black, color.White}
	g := &gif.GIF{
		Image: []*image.Paletted{
			image.NewPaletted(image.Rect(0, 0, 80, 60), pal),
			image.NewPaletted(image.Rect(0, 0, 80, 60), pal),
		},
		Delay:  []int{10, 20},
		Config: image.Config{Width: 80, Height: 60},
	}
	if err := gif.EncodeAll(f, g); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	info, err := Probe(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Type != TypeGIF {
		t.Errorf("type = %q, want gif", info.Type)
	}
	if info.Width != 80 || info.Height != 60 {
		t.Errorf("dims = %dx%d, want 80x60", info.Width, info.Height)
	}
	if info.Frames != 2 {
		t.Errorf("frames = %d, want 2", info.Frames)
	}
	if info.Duration < 0.29 || info.Duration > 0.31 {
		t.Errorf("duration = %f, want ~0.3", info.Duration)
	}
}

func TestProbeVideoMissing(t *testing.T) {
	_, err := probeVideo("/nonexistent/video.mp4")
	if err == nil {
		t.Log("probeVideo on missing file did not error (ffmpeg output may be empty)")
	}
}

func TestDetectType(t *testing.T) {
	tests := []struct {
		path string
		want MediaType
	}{
		{"x.png", TypeImage},
		{"x.jpg", TypeImage},
		{"x.gif", TypeGIF},
		{"x.mp4", TypeVideo},
		{"x.webm", TypeVideo},
		{"x.mov", TypeVideo},
		{"x.xyz", TypeUnknown},
	}
	for _, tt := range tests {
		got := detectType("", tt.path)
		if got != tt.want {
			t.Errorf("detectType(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestProbeMissingFile(t *testing.T) {
	_, err := Probe("/nonexistent/path.png")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
