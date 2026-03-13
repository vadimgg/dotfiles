package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ─── Scene analysis ───────────────────────────────────────────────────────────
//
// Pass 1 — FAST scan: ffmpeg decodes only keyframes (-skip_frame noref) and
//           runs the scdet filter, collecting every non-zero score. Keyframe-
//           only decoding is 10-30x faster than full decode on large files and
//           is accurate enough for scoring because scene cuts almost always
//           fall on keyframes.
//
// Pass 2 — Otsu's method separates the "noise" cluster (subtitles, animation,
//           gradual motion — low scores) from the "cut" cluster (hard scene
//           changes — high scores) automatically, with no user input.
//
// Extract — ONE ffmpeg process extracts all frames simultaneously using a
//           select filter with every cut timestamp baked in, and writes them
//           directly to the scenes/ folder. This replaces N sequential ffmpeg
//           calls with a single pass over the file.

const minSceneGapSeconds = 0.5

type scoredFrame struct {
	ts    float64
	score float64
}

// AnalyzeScenes is the main entry point.
func AnalyzeScenes(videoPath, videoDir string) (scenesDir string, count int, err error) {
	if _, err = exec.LookPath("ffmpeg"); err != nil {
		return "", 0, fmt.Errorf("ffmpeg not found — install with: brew install ffmpeg")
	}

	scenesDir = filepath.Join(videoDir, "scenes")
	if err = os.MkdirAll(scenesDir, 0755); err != nil {
		return "", 0, fmt.Errorf("could not create scenes directory: %w", err)
	}

	// ── Pass 1: fast keyframe scan for scene scores ───────────────────────────
	fmt.Printf("%s\n", cBold("🎞️  Scanning for scene cuts (fast keyframe pass)…"))

	frames, err := collectScoredFrames(videoPath)
	if err != nil {
		return scenesDir, 0, fmt.Errorf("scan failed: %w", err)
	}
	if len(frames) == 0 {
		return scenesDir, 0, fmt.Errorf("no scene scores returned — is this a valid video file?")
	}

	// ── Pass 2: Otsu threshold separates noise from real cuts ────────────────
	cutThreshold := otsuThreshold(frames)
	fmt.Printf("  %s scanned %s frames  |  cut boundary: %s\n\n",
		cDim("→"),
		cBold(fmt.Sprintf("%d", len(frames))),
		cBold(fmt.Sprintf("%.2f", cutThreshold)),
	)

	var timestamps []float64
	lastTs := -minSceneGapSeconds - 1.0
	for _, f := range frames {
		if f.score >= cutThreshold && f.ts-lastTs >= minSceneGapSeconds {
			timestamps = append(timestamps, f.ts)
			lastTs = f.ts
		}
	}
	timestamps = prependZero(timestamps)

	fmt.Printf("  %s %s\n\n",
		cDim("→"),
		cBold(fmt.Sprintf("%d scene cuts detected", len(timestamps))),
	)

	// ── Extract all frames in a single ffmpeg pass ────────────────────────────
	fmt.Printf("%s\n", cBold("🖼️  Extracting all scene frames…"))

	// Split into new vs already-existing so we only re-extract what's missing
	var needed []float64
	for _, ts := range timestamps {
		label := timestampLabel(ts)
		if !fileExists(filepath.Join(scenesDir, label+".jpg")) {
			needed = append(needed, ts)
		} else {
			fmt.Printf("  %s %s\n", cDim("skip"), cDim(label+".jpg (exists)"))
			count++
		}
	}

	if len(needed) > 0 {
		extracted, err := extractAllFrames(videoPath, needed, scenesDir)
		if err != nil {
			printWarn("Frame extraction error: %v", err)
		}
		count += extracted
	}

	fmt.Println()
	printSuccess("Scene analysis complete — %d images in: %s", count, cPath(scenesDir))
	return scenesDir, count, nil
}

// collectScoredFrames does a fast keyframe-only decode with scdet so we get
// scene scores without decoding every single frame.
func collectScoredFrames(videoPath string) ([]scoredFrame, error) {
	args := []string{
		"-hide_banner",
		// Skip non-reference (B/P) frames — only decode keyframes (I-frames).
		// This is the main speed-up: 10-30x faster on large files.
		"-skip_frame", "noref",
		"-i", videoPath,
		"-vf", "scdet=threshold=0:sc_pass=0",
		"-an",
		"-f", "null",
		"-",
	}

	cmd := exec.Command("ffmpeg", args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ffmpeg scan failed to start: %w", err)
	}

	reScore := regexp.MustCompile(`lavfi\.scd\.score:\s*([\d.]+),\s*lavfi\.scd\.time:\s*([\d.]+)`)

	var frames []scoredFrame
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		m := reScore.FindStringSubmatch(scanner.Text())
		if m == nil {
			continue
		}
		score, _ := strconv.ParseFloat(m[1], 64)
		ts, err  := strconv.ParseFloat(m[2], 64)
		if err != nil || score == 0 {
			continue
		}
		frames = append(frames, scoredFrame{ts, score})
	}
	_ = cmd.Wait()
	return frames, nil
}

// extractAllFrames extracts every timestamp in a SINGLE ffmpeg pass using the
// select filter. Output files are named by timestamp label.
//
// PySceneDetect does the same thing internally — one decode pass, all frames
// written out together — which is why it's much faster than calling ffmpeg
// once per frame.
func extractAllFrames(videoPath string, timestamps []float64, scenesDir string) (int, error) {
	if len(timestamps) == 0 {
		return 0, nil
	}

	// Build a select expression that matches each exact timestamp:
	//   eq(t,12.340)+eq(t,25.680)+...
	// ffmpeg will output exactly one frame per matching timestamp.
	var parts []string
	for _, ts := range timestamps {
		parts = append(parts, fmt.Sprintf("eq(t\\,%.6f)", ts))
	}
	selectExpr := strings.Join(parts, "+")

	// Output pattern: scenesDir/TIMESTAMP.jpg
	// We use the pts_time to derive filenames via a combination of
	// select + vsync + output pattern with frame number, then rename.
	// Simpler: write to a temp pattern and rename after.
	tmpPattern := filepath.Join(scenesDir, "frame_%06d.jpg")

	args := []string{
		"-hide_banner",
		"-i", videoPath,
		"-vf", fmt.Sprintf("select='%s',setpts=N/TB", selectExpr),
		"-vsync", "0",
		"-q:v", "2",
		"-y",
		tmpPattern,
	}

	cmd := exec.Command("ffmpeg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		return 0, fmt.Errorf("ffmpeg extract: %s", lines[len(lines)-1])
	}

	// Rename frame_000001.jpg → 00h03m42s500ms.jpg
	count := 0
	for i, ts := range timestamps {
		src  := filepath.Join(scenesDir, fmt.Sprintf("frame_%06d.jpg", i+1))
		dst  := filepath.Join(scenesDir, timestampLabel(ts)+".jpg")
		if err := os.Rename(src, dst); err != nil {
			// file may not have been written if select missed it
			continue
		}
		fmt.Printf("  %s %s\n",
			cDim(fmt.Sprintf("[%3d/%d]", i+1, len(timestamps))),
			cPath(timestampLabel(ts)+".jpg"),
		)
		count++
	}
	return count, nil
}

// ─── Otsu threshold ───────────────────────────────────────────────────────────

func otsuThreshold(frames []scoredFrame) float64 {
	if len(frames) == 0 {
		return 1.0
	}
	scores := make([]float64, len(frames))
	for i, f := range frames {
		scores[i] = f.score
	}
	sort.Float64s(scores)

	minScore := scores[0]
	maxScore := scores[len(scores)-1]
	if maxScore-minScore < 0.001 {
		return maxScore * 0.9
	}

	const bins = 256
	hist := make([]float64, bins)
	for _, s := range scores {
		bin := int((s - minScore) / (maxScore - minScore) * float64(bins-1))
		if bin >= bins {
			bin = bins - 1
		}
		hist[bin]++
	}

	total  := float64(len(scores))
	sumAll := 0.0
	for i, h := range hist {
		sumAll += float64(i) * h
	}

	bestVariance := -1.0
	bestBin      := 0
	weightBg     := 0.0
	sumBg        := 0.0

	for i, h := range hist {
		weightBg += h
		if weightBg == 0 {
			continue
		}
		weightFg := total - weightBg
		if weightFg == 0 {
			break
		}
		sumBg += float64(i) * h
		meanBg := sumBg / weightBg
		meanFg := (sumAll - sumBg) / weightFg
		v := weightBg * weightFg * math.Pow(meanBg-meanFg, 2)
		if v > bestVariance {
			bestVariance = v
			bestBin = i
		}
	}

	threshold := minScore + float64(bestBin)/float64(bins-1)*(maxScore-minScore)
	if threshold <= minScore+0.001 {
		mean, stddev := meanStddev(scores)
		threshold = mean + 2*stddev
	}
	return threshold
}

func meanStddev(vals []float64) (float64, float64) {
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	mean := sum / float64(len(vals))
	variance := 0.0
	for _, v := range vals {
		d := v - mean
		variance += d * d
	}
	return mean, math.Sqrt(variance / float64(len(vals)))
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func prependZero(ts []float64) []float64 {
	if len(ts) == 0 || ts[0] >= minSceneGapSeconds {
		return append([]float64{0.0}, ts...)
	}
	return ts
}

func timestampLabel(seconds float64) string {
	ms := int(seconds * 1000)
	h  := ms / 3_600_000; ms -= h * 3_600_000
	m  := ms / 60_000;    ms -= m * 60_000
	s  := ms / 1_000;     ms -= s * 1_000
	return fmt.Sprintf("%02dh%02dm%02ds%03dms", h, m, s, ms)
}

// RunSceneAnalysis is called directly from the menu.
func RunSceneAnalysis(videoPath, videoDir string) {
	fmt.Println()
	_, _, err := AnalyzeScenes(videoPath, videoDir)
	if err != nil {
		printError("Scene analysis failed: %v", err)
		printTip("Make sure ffmpeg is installed: brew install ffmpeg")
	}
}
