package normalize

import (
	"fmt"
	"os"
	"os/exec"
)

type Options struct {
	MaxColors int
	MaxW      int
	MaxH      int
}

func Normalize(path string, opts Options) error {
	if opts.MaxColors <= 0 {
		opts.MaxColors = 256
	}

	filters := buildFilter(opts)

	tmp := path + ".tmp.png"
	args := []string{"-i", path, "-vf", filters, "-y", tmp}
	cmd := exec.Command("ffmpeg", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("ffmpeg normalize: %w: %s", err, string(out))
	}

	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}

func buildFilter(opts Options) string {
	var pre []string
	if opts.MaxW > 0 && opts.MaxH > 0 {
		pre = append(pre, fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease", opts.MaxW, opts.MaxH))
	}
	palette := fmt.Sprintf("palettegen=max_colors=%d", opts.MaxColors)
	parts := append(pre, "split[s0][s1];[s0]"+palette+"[p];[s1][p]paletteuse")
	var result string
	for i, p := range parts {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}
