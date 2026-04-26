package theme

import (
	"strings"
	"testing"
)

func TestGenerateScript(t *testing.T) {
	opts := ScriptOpts{
		Name:   "test",
		Frames: []string{"frame_0001.png", "frame_0002.png"},
		Loop:   true,
	}
	got, err := GenerateScript(opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `new Image("frame_0001.png")`) {
		t.Error("missing frame_0001 ref")
	}
	if !strings.Contains(got, "frameIdx = (frameIdx + 1) % frames.length") {
		t.Error("missing loop logic")
	}
}

func TestValidateScriptForbidden(t *testing.T) {
	bad := `eval("1+1")`
	if err := ValidateScript(bad, nil); err == nil {
		t.Fatal("expected error for forbidden token")
	}
}

func TestValidateScriptAssetRef(t *testing.T) {
	script := `new Image("frame_0001.png")`
	assets := []string{"/tmp/assets/frame_0001.png"}
	if err := ValidateScript(script, assets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateScriptMissingAsset(t *testing.T) {
	script := `new Image("missing.png")`
	assets := []string{"/tmp/assets/frame_0001.png"}
	if err := ValidateScript(script, assets); err == nil {
		t.Fatal("expected error for missing asset ref")
	}
}
