package plymouth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.plymouth")
	content := `[Plymouth Theme]
Name=MyTheme
Description=A test theme
ModuleName=script
ImageDir=assets
ScriptFile=mytheme.script
`
	os.WriteFile(path, []byte(content), 0o644)

	d, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if d.Name != "MyTheme" {
		t.Errorf("Name = %q, want MyTheme", d.Name)
	}
	if d.ModuleName != "script" {
		t.Errorf("ModuleName = %q, want script", d.ModuleName)
	}
	if d.ScriptFile != "mytheme.script" {
		t.Errorf("ScriptFile = %q", d.ScriptFile)
	}
}

func TestValidateMissingModuleName(t *testing.T) {
	d := &Descriptor{Name: "x", ImageDir: "assets", ScriptFile: "x.script"}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for missing ModuleName")
	}
}

func TestValidateWrongModuleName(t *testing.T) {
	d := &Descriptor{Name: "x", ModuleName: "two-step", ImageDir: "assets", ScriptFile: "x.script"}
	if err := d.Validate(); err == nil {
		t.Fatal("expected error for wrong ModuleName")
	}
}
