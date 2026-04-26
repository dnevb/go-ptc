package theme

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dan/plymouth-theme-creator/internal/extract"
	"github.com/dan/plymouth-theme-creator/internal/normalize"
	"github.com/dan/plymouth-theme-creator/internal/probe"
)

type CreateOpts struct {
	Name       string
	Media      []string
	FPS        int
	Res        string
	Loop       bool
	Transition string
	OutputDir  string
}

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func Create(opts CreateOpts) error {
	if !nameRe.MatchString(opts.Name) {
		return fmt.Errorf("theme name must match [a-zA-Z0-9_-]+")
	}
	if len(opts.Media) == 0 {
		return fmt.Errorf("no media files")
	}

	outDir := filepath.Join(opts.OutputDir, opts.Name)
	assetsDir := filepath.Join(outDir, "assets")

	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	// V1: probe all media first (existence + readability checked inside Probe).
	var infos []*probe.MediaInfo
	for _, m := range opts.Media {
		info, err := probe.Probe(m)
		if err != nil {
			return fmt.Errorf("probe %s: %w", m, err)
		}
		infos = append(infos, info)
	}

	// Extract frames for each media file.
	var allFrames []string
	maxW, maxH := parseRes(opts.Res)
	for i, m := range opts.Media {
		info := infos[i]
		var extOpts extract.Options
		if info.Type != probe.TypeImage {
			extOpts.FPS = opts.FPS
		}
		extOpts.MaxRes = opts.Res
		extOpts.OutDir = assetsDir

		frames, err := extract.ExtractFrames(m, extOpts)
		if err != nil {
			return fmt.Errorf("extract %s: %w", m, err)
		}
		allFrames = append(allFrames, frames...)
	}

	// Normalize all extracted frames.
	normOpts := normalize.Options{MaxColors: 256, MaxW: maxW, MaxH: maxH}
	for _, f := range allFrames {
		if err := normalize.Normalize(f, normOpts); err != nil {
			return fmt.Errorf("normalize %s: %w", f, err)
		}
	}

	// V2: generate + validate .plymouth before write.
	plymouthContent, err := GeneratePlymouth(opts.Name)
	if err != nil {
		return err
	}
	if err := ValidatePlymouth(plymouthContent); err != nil {
		return fmt.Errorf("plymouth validation: %w", err)
	}

	// V5: build script with asset basenames.
	var frameNames []string
	for _, f := range allFrames {
		frameNames = append(frameNames, filepath.Base(f))
	}
	scriptContent, err := GenerateScript(ScriptOpts{
		Name:       opts.Name,
		Frames:     frameNames,
		Loop:       opts.Loop,
		Transition: opts.Transition,
	})
	if err != nil {
		return fmt.Errorf("script generation: %w", err)
	}
	if err := ValidateScript(scriptContent, allFrames); err != nil {
		return fmt.Errorf("script validation: %w", err)
	}

	// Write files.
	plymouthPath := filepath.Join(outDir, opts.Name+".plymouth")
	if err := os.WriteFile(plymouthPath, []byte(plymouthContent), 0o644); err != nil {
		return fmt.Errorf("write .plymouth: %w", err)
	}

	scriptPath := filepath.Join(outDir, opts.Name+".script")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o644); err != nil {
		return fmt.Errorf("write .script: %w", err)
	}

	return nil
}

func parseRes(res string) (int, int) {
	parts := strings.Split(res, "x")
	if len(parts) != 2 {
		return 1920, 1080
	}
	w, _ := strconv.Atoi(parts[0])
	h, _ := strconv.Atoi(parts[1])
	if w <= 0 || h <= 0 {
		return 1920, 1080
	}
	return w, h
}
