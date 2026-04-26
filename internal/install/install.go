package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dnevb/go-ptc/internal/plymouth"
)

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func Install(dir, systemDir string) error {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var plymouthPath string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".plymouth" {
			plymouthPath = filepath.Join(dir, e.Name())
			break
		}
	}
	if plymouthPath == "" {
		return fmt.Errorf("no .plymouth file found in %s", dir)
	}

	desc, err := plymouth.Parse(plymouthPath)
	if err != nil {
		return fmt.Errorf("parse .plymouth: %w", err)
	}
	if err := desc.Validate(); err != nil {
		return fmt.Errorf("invalid .plymouth: %w", err)
	}

	scriptPath := filepath.Join(dir, desc.ScriptFile)
	if _, err := os.Stat(scriptPath); err != nil {
		return fmt.Errorf("missing script: %w", err)
	}

	assetsPath := filepath.Join(dir, desc.ImageDir)
	if info, err := os.Stat(assetsPath); err != nil || !info.IsDir() {
		return fmt.Errorf("missing assets dir: %w", err)
	}

	name := desc.Name
	if !nameRe.MatchString(name) {
		return fmt.Errorf("theme name must match [a-zA-Z0-9_-]+")
	}

	// V7: case-insensitive unique check.
	sysEntries, err := os.ReadDir(systemDir)
	if err != nil {
		return fmt.Errorf("read system dir: %w", err)
	}
	lowerName := strings.ToLower(name)
	for _, e := range sysEntries {
		if strings.ToLower(e.Name()) == lowerName {
			return fmt.Errorf("theme %q already installed", name)
		}
	}

	target := filepath.Join(systemDir, name)
	if err := os.MkdirAll(target, 0o755); err != nil {
		return fmt.Errorf("mkdir target: %w", err)
	}

	if err := copyDir(dir, target); err != nil {
		return fmt.Errorf("copy: %w", err)
	}

	return nil
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		srcPath := filepath.Join(src, e.Name())
		dstPath := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := os.MkdirAll(dstPath, 0o755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
