package install

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTheme(t *testing.T) string {
	dir := t.TempDir()
	themeDir := filepath.Join(dir, "mytheme")
	os.MkdirAll(filepath.Join(themeDir, "assets"), 0o755)
	os.WriteFile(filepath.Join(themeDir, "mytheme.plymouth"), []byte(`[Plymouth Theme]
Name=mytheme
ModuleName=script
ImageDir=assets
ScriptFile=mytheme.script
`), 0o644)
	os.WriteFile(filepath.Join(themeDir, "mytheme.script"), []byte(`// script`), 0o644)
	os.WriteFile(filepath.Join(themeDir, "assets", "frame.png"), []byte(`png`), 0o644)
	return themeDir
}

func TestInstallSuccess(t *testing.T) {
	themeDir := setupTheme(t)
	sysDir := t.TempDir()
	if err := Install(themeDir, sysDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(sysDir, "mytheme", "mytheme.plymouth")); err != nil {
		t.Error("plymouth not copied")
	}
}

func TestInstallMissingPlymouth(t *testing.T) {
	dir := t.TempDir()
	if err := Install(dir, t.TempDir()); err == nil {
		t.Fatal("expected error")
	}
}

func TestInstallDuplicateName(t *testing.T) {
	themeDir := setupTheme(t)
	sysDir := t.TempDir()
	os.MkdirAll(filepath.Join(sysDir, "mytheme"), 0o755)
	if err := Install(themeDir, sysDir); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestInstallCaseInsensitiveDuplicate(t *testing.T) {
	themeDir := setupTheme(t)
	sysDir := t.TempDir()
	os.MkdirAll(filepath.Join(sysDir, "MyTheme"), 0o755)
	if err := Install(themeDir, sysDir); err == nil {
		t.Fatal("expected case-insensitive duplicate error")
	}
}
