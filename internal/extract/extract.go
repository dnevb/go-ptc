package extract

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/dnevb/go-ptc/internal/probe"
)

type Options struct {
	FPS    int
	MaxRes string
	OutDir string
}

func ExtractFrames(path string, opts Options) ([]string, error) {
	info, err := probe.Probe(path)
	if err != nil {
		return nil, err
	}

	outDir := opts.OutDir
	if outDir == "" {
		tmp, err := os.MkdirTemp("", "ptc-extract-*")
		if err != nil {
			return nil, fmt.Errorf("mkdir temp: %w", err)
		}
		outDir = tmp
	} else {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return nil, fmt.Errorf("mkdir outdir: %w", err)
		}
	}

	pattern := filepath.Join(outDir, "frame_%04d.png")

	var args []string
	args = append(args, "-i", path)

	if info.Type != probe.TypeImage {
		filters := buildFilters(opts)
		if filters != "" {
			args = append(args, "-vf", filters)
		}
	}

	if info.Type == probe.TypeImage {
		args = append(args, "-frames:v", "1")
	}

	args = append(args, pattern)

	cmd := exec.Command("ffmpeg", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg: %w: %s", err, string(out))
	}

	return listFrames(outDir)
}

func buildFilters(opts Options) string {
	var parts []string
	if opts.FPS > 0 {
		parts = append(parts, fmt.Sprintf("fps=%d", opts.FPS))
	}
	scale := parseScale(opts.MaxRes)
	if scale != "" {
		parts = append(parts, scale)
	}
	return strings.Join(parts, ",")
}

func parseScale(res string) string {
	if res == "" {
		return ""
	}
	parts := strings.Split(res, "x")
	if len(parts) != 2 {
		return ""
	}
	w, wErr := strconv.Atoi(parts[0])
	h, hErr := strconv.Atoi(parts[1])
	if wErr != nil || hErr != nil || w <= 0 || h <= 0 {
		return ""
	}
	return fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", w, h)
}

func listFrames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var frames []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".png" {
			frames = append(frames, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(frames)
	return frames, nil
}
