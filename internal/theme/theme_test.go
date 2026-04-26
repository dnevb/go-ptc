package theme

import (
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseRes(t *testing.T) {
	tests := []struct {
		input   string
		wantW   int
		wantH   int
	}{
		{"1920x1080", 1920, 1080},
		{"", 1920, 1080},
		{"bad", 1920, 1080},
		{"0x0", 1920, 1080},
		{"800x600", 800, 600},
	}
	for _, tt := range tests {
		w, h := parseRes(tt.input)
		if w != tt.wantW || h != tt.wantH {
			t.Errorf("parseRes(%q) = %d,%d, want %d,%d", tt.input, w, h, tt.wantW, tt.wantH)
		}
	}
}

func TestCreateInvalidName(t *testing.T) {
	err := Create(CreateOpts{Name: "bad name!", Media: []string{"x.png"}})
	if err == nil {
		t.Fatal("expected error for invalid name")
	}
}

func TestCreateNoMedia(t *testing.T) {
	err := Create(CreateOpts{Name: "test", Media: []string{}})
	if err == nil {
		t.Fatal("expected error for no media")
	}
}

func TestCreateIntegration(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not found")
	}

	dir := t.TempDir()
	imgPath := filepath.Join(dir, "test.png")
	f, err := os.Create(imgPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 100, 100))); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	outDir := t.TempDir()
	opts := CreateOpts{
		Name:      "mytheme",
		Media:     []string{imgPath},
		Res:       "800x600",
		OutputDir: outDir,
	}
	if err := Create(opts); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "mytheme", "mytheme.plymouth")); err != nil {
		t.Error("missing .plymouth")
	}
	if _, err := os.Stat(filepath.Join(outDir, "mytheme", "mytheme.script")); err != nil {
		t.Error("missing .script")
	}
	if _, err := os.Stat(filepath.Join(outDir, "mytheme", "assets", "frame_0001.png")); err != nil {
		t.Error("missing asset")
	}
}
