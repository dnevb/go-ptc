package probe

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type MediaType string

const (
	TypeImage MediaType = "image"
	TypeGIF   MediaType = "gif"
	TypeVideo MediaType = "video"
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

type ffprobeStream struct {
	CodecName    string `json:"codec_name"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	NbFrames     string `json:"nb_frames"`
	AvgFrameRate string `json:"avg_frame_rate"`
}

type ffprobeFormat struct {
	Duration string `json:"duration"`
}

type ffprobeOutput struct {
	Streams []ffprobeStream `json:"streams"`
	Format  ffprobeFormat   `json:"format"`
}

func Probe(path string) (*MediaInfo, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("media file unreadable: %w", err)
	}

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_format",
		"-show_streams",
		"-print_format", "json",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("ffprobe failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var raw ffprobeOutput
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("ffprobe parse: %w", err)
	}

	if len(raw.Streams) == 0 {
		return nil, fmt.Errorf("ffprobe: no streams found")
	}

	s := raw.Streams[0]
	info := &MediaInfo{
		Path:  path,
		Width: s.Width,
		Height: s.Height,
	}

	info.Type = detectType(s.CodecName, path)

	if raw.Format.Duration != "" && raw.Format.Duration != "N/A" {
		info.Duration, _ = strconv.ParseFloat(raw.Format.Duration, 64)
	}

	info.Frames = parseFrames(s.NbFrames, s.AvgFrameRate, info.Duration)

	return info, nil
}

func detectType(codec, path string) MediaType {
	switch strings.ToLower(codec) {
	case "gif":
		return TypeGIF
	case "png", "mjpeg", "bmp", "tiff", "webp":
		return TypeImage
	case "h264", "hevc", "vp8", "vp9", "av1", "mpeg4", "theora":
		return TypeVideo
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".webp":
		return TypeImage
	case ".gif":
		return TypeGIF
	case ".mp4", ".webm", ".mov", ".avi", ".mkv":
		return TypeVideo
	}

	return TypeUnknown
}

func parseFrames(nbFrames, avgFrameRate string, duration float64) int {
	if nbFrames != "" && nbFrames != "N/A" {
		if n, err := strconv.Atoi(nbFrames); err == nil && n > 0 {
			return n
		}
	}

	if avgFrameRate != "" && avgFrameRate != "0/0" && avgFrameRate != "N/A" {
		parts := strings.Split(avgFrameRate, "/")
		if len(parts) == 2 {
			num, nErr := strconv.ParseFloat(parts[0], 64)
			den, dErr := strconv.ParseFloat(parts[1], 64)
			if nErr == nil && dErr == nil && den != 0 && duration > 0 {
				return int(num / den * duration)
			}
		}
	}

	if duration == 0 {
		return 1
	}

	return int(duration * 30)
}
