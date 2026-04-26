package probe

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MediaType string

const (
	TypeImage   MediaType = "image"
	TypeGIF     MediaType = "gif"
	TypeVideo   MediaType = "video"
	TypeUnknown MediaType = "unknown"
)

type MediaInfo struct {
	Path     string
	Type     MediaType
	Width    int
	Height   int
	Duration float64
	Frames   int
}

var (
	durationRe    = regexp.MustCompile(`Duration:\s*(\d+):(\d+):(\d+)(?:\.(\d+))?`)
	videoStreamRe = regexp.MustCompile(`Stream.*Video:.*\s(\d+)x(\d+)[\s,]`)
	fpsRe         = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*fps`)
	tbrRe         = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*tbr`)
)

func Probe(path string) (*MediaInfo, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("media file unreadable: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".bmp", ".tiff":
		return probeImage(path)
	case ".gif":
		return probeGIF(path)
	case ".mp4", ".webm", ".mov", ".avi", ".mkv", ".webp":
		return probeVideo(path)
	default:
		if info, err := probeImage(path); err == nil {
			return info, nil
		}
		return probeVideo(path)
	}
}

func probeImage(path string) (*MediaInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, err
	}
	return &MediaInfo{
		Path:   path,
		Type:   TypeImage,
		Width:  cfg.Width,
		Height: cfg.Height,
		Frames: 1,
	}, nil
}

func probeGIF(path string) (*MediaInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g, err := gif.DecodeAll(f)
	if err != nil {
		return nil, err
	}
	var dur time.Duration
	for _, d := range g.Delay {
		dur += time.Duration(d) * 10 * time.Millisecond
	}
	return &MediaInfo{
		Path:     path,
		Type:     TypeGIF,
		Width:    g.Config.Width,
		Height:   g.Config.Height,
		Duration: dur.Seconds(),
		Frames:   len(g.Image),
	}, nil
}

func probeVideo(path string) (*MediaInfo, error) {
	cmd := exec.Command("ffmpeg", "-i", path, "-f", "null", "-")
	out, err := cmd.CombinedOutput()
	output := string(out)
	if output == "" && err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w", err)
	}

	info := &MediaInfo{Path: path, Type: TypeVideo}

	if m := durationRe.FindStringSubmatch(output); len(m) == 5 {
		h, _ := strconv.Atoi(m[1])
		min, _ := strconv.Atoi(m[2])
		sec, _ := strconv.ParseFloat(m[3], 64)
		ms := 0.0
		if m[4] != "" {
			ms, _ = strconv.ParseFloat(m[4], 64)
		}
		info.Duration = float64(h)*3600 + float64(min)*60 + sec + ms/100
	}

	if m := videoStreamRe.FindStringSubmatch(output); len(m) >= 3 {
		info.Width, _ = strconv.Atoi(m[1])
		info.Height, _ = strconv.Atoi(m[2])
	}

	fps := 0.0
	if m := fpsRe.FindStringSubmatch(output); len(m) >= 2 {
		fps, _ = strconv.ParseFloat(m[1], 64)
	} else if m := tbrRe.FindStringSubmatch(output); len(m) >= 2 {
		fps, _ = strconv.ParseFloat(m[1], 64)
	}

	if info.Duration > 0 && fps > 0 {
		info.Frames = int(info.Duration * fps)
	} else if info.Duration == 0 {
		info.Frames = 1
	}

	return info, nil
}

func detectType(_ string, path string) MediaType {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".bmp", ".tiff":
		return TypeImage
	case ".gif":
		return TypeGIF
	case ".mp4", ".webm", ".mov", ".avi", ".mkv", ".webp":
		return TypeVideo
	}
	return TypeUnknown
}
