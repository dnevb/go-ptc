package theme

import (
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
