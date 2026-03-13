package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// compressImages re-encodes every JPEG in dir at the given quality (1–31,
// lower = better quality, larger file; higher = worse quality, smaller file).
// Quality 8 is a good default: visually near-identical to the original at
// roughly 30–50% of the file size.
//
// ffmpeg is used so there are no extra dependencies — it's already required
// for scene detection.
func compressImages(dir string, quality int) (saved int64, count int, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, fmt.Errorf("could not read directory: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.ToLower(filepath.Ext(name)) != ".jpg" {
			continue
		}

		path := filepath.Join(dir, name)

		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		before := info.Size()

		// Write to a temp file first so we never corrupt the original
		tmp := path + ".tmp.jpg"
		args := []string{
			"-hide_banner", "-loglevel", "error",
			"-i", path,
			"-q:v", fmt.Sprintf("%d", quality),
			"-y", tmp,
		}
		if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
			fmt.Printf("  %s %s — %s\n", cDim("skip"), cDim(name), cError(strings.TrimSpace(string(out))))
			_ = os.Remove(tmp)
			continue
		}

		info2, err := os.Stat(tmp)
		if err != nil {
			_ = os.Remove(tmp)
			continue
		}
		after := info2.Size()

		// Only keep the compressed version if it's actually smaller
		if after >= before {
			_ = os.Remove(tmp)
			fmt.Printf("  %s %s %s\n",
				cDim("skip"),
				cDim(name),
				cDim("(already optimal)"),
			)
			continue
		}

		if err := os.Rename(tmp, path); err != nil {
			_ = os.Remove(tmp)
			continue
		}

		reduction := float64(before-after) / float64(before) * 100
		fmt.Printf("  %s %s  %s → %s %s\n",
			cDim(fmt.Sprintf("[%d]", count+1)),
			cPath(name),
			cDim(formatBytes(before)),
			cBold(formatBytes(after)),
			colorize(green, fmt.Sprintf("(−%.0f%%)", reduction)),
		)

		saved += before - after
		count++
	}
	return saved, count, nil
}

func formatBytes(b int64) string {
	switch {
	case b >= 1024*1024:
		return fmt.Sprintf("%.1fMB", float64(b)/1024/1024)
	case b >= 1024:
		return fmt.Sprintf("%.1fKB", float64(b)/1024)
	default:
		return fmt.Sprintf("%dB", b)
	}
}
