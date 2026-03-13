package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// videoMeta holds the fields we pull from yt-dlp --dump-json.
type videoMeta struct {
	Title       string `json:"title"`
	Channel     string `json:"channel"`
	Uploader    string `json:"uploader"`
	UploadDate  string `json:"upload_date"` // "YYYYMMDD"
}

// fetchMeta calls yt-dlp --dump-json to get video metadata without downloading.
func fetchMeta(ytDlp, url string, extraArgs []string) (videoMeta, error) {
	args := append([]string{"--dump-json", "--no-playlist"}, extraArgs...)
	args = append(args, url)
	out, err := exec.Command(ytDlp, args...).Output()
	if err != nil {
		return videoMeta{}, err
	}
	var m videoMeta
	if err := json.Unmarshal(out, &m); err != nil {
		return videoMeta{}, err
	}
	return m, nil
}

// buildFolderName returns YYYY_MM_DD_CHANNEL_TITLE from metadata.
// Falls back gracefully when fields are missing.
func buildFolderName(m videoMeta) string {
	// Date: prefer upload_date, fall back to today
	dateStr := m.UploadDate
	var date time.Time
	if len(dateStr) == 8 {
		t, err := time.Parse("20060102", dateStr)
		if err == nil {
			date = t
		}
	}
	if date.IsZero() {
		date = time.Now()
	}
	datePart := date.Format("2006_01_02")

	// Channel: prefer channel, fall back to uploader
	channel := m.Channel
	if channel == "" {
		channel = m.Uploader
	}
	if channel == "" {
		channel = "unknown"
	}

	title := m.Title
	if title == "" {
		title = "untitled"
	}

	return datePart + "_" + slugify(channel, 40) + "_" + slugify(title, 80)
}

// downloadYouTube downloads a YouTube video and returns the file path and output dir.
func downloadYouTube(url, configPath, outputDirOverride, formatOverride string) (filePath, outputDir string) {
	cfg, err := loadConfig(configPath)
	printConfigStatus(configPath, err)

	yt := cfg.YouTube

	// CLI overrides
	if outputDirOverride != "" {
		yt.OutputDir = outputDirOverride
	}
	if formatOverride != "" {
		yt.Format = formatOverride
	}

	// Ensure yt-dlp is available
	ytDlp := cfg.Advanced.YtDlpPath
	if _, err := exec.LookPath(ytDlp); err != nil {
		printError("yt-dlp not found at %q. Install it with:", ytDlp)
		printTip("pip install yt-dlp   or   brew install yt-dlp")
		os.Exit(1)
	}

	// ── Fetch metadata to build the folder name ──────────────────────────────
	fmt.Printf("%s %s\n", cBold("🔍 Fetching metadata:"), cURL(url))

	metaExtraArgs := authArgs(cfg.Advanced)

	meta, err := fetchMeta(ytDlp, url, metaExtraArgs)
	if err != nil {
		printWarn("Could not fetch metadata (%v) — using fallback folder name.", err)
		meta = videoMeta{Title: "untitled", UploadDate: time.Now().Format("20060102")}
	} else {
		fmt.Printf("  %s %s\n", cDim("Channel:"), cBold(meta.Channel))
		fmt.Printf("  %s %s\n\n", cDim("Title:  "), cBold(meta.Title))
	}

	folderName := buildFolderName(meta)

	// ── Prepare base output dir then create per-video subfolder ─────────────
	baseDir, err := prepareOutputDir(yt.OutputDir)
	if err != nil {
		printError("Could not prepare output directory: %v", err)
		os.Exit(1)
	}

	videoDir := filepath.Join(baseDir, folderName)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		printError("Could not create video folder: %v", err)
		os.Exit(1)
	}
	outputDir = videoDir

	// File lives directly in the video folder, named after the video ID
	outputTemplate := filepath.Join(videoDir, "%(id)s.%(ext)s")

	// ── Build format string ──────────────────────────────────────────────────
	format := yt.Format
	if yt.AudioOnly {
		format = "bestaudio/best"
	} else if yt.Resolution != "" && yt.Resolution != "best" {
		format = fmt.Sprintf("bestvideo[height<=%s]+bestaudio/best[height<=%s]", yt.Resolution, yt.Resolution)
	}

	args := []string{
		"--format", format,
		"--output", outputTemplate,
		"--retries", fmt.Sprintf("%d", yt.Retries),
		"--progress",
	}

	if yt.AudioOnly {
		args = append(args,
			"--extract-audio",
			"--audio-format", "mp3",
			"--audio-quality", "0",
		)
		printInfo("Audio-only mode — extracting as mp3.")
	} else {
		args = append(args, "--merge-output-format", yt.MergeFormat)
	}

	if yt.Subtitles {
		args = append(args,
			"--write-subs",
			"--embed-subs",
			"--sub-langs", yt.SubtitleLangs,
		)
	}

	for _, a := range authArgs(cfg.Advanced) {
		args = append(args, a)
	}
	if cfg.Advanced.SkipIfExists {
		args = append(args, "--no-overwrites")
	}
	if cfg.Advanced.WriteMetadata {
		args = append(args, "--write-info-json")
	}
	if cfg.Advanced.WriteThumbnail {
		args = append(args, "--write-thumbnail")
	}
	if cfg.Advanced.RateLimit != "" {
		args = append(args, "--rate-limit", cfg.Advanced.RateLimit)
	}
	if cfg.Advanced.SleepInterval > 0 {
		args = append(args, "--sleep-interval", fmt.Sprintf("%d", cfg.Advanced.SleepInterval))
	}

	args = append(args, url)

	fmt.Printf("%s %s\n",   cBold("⬇️  Downloading:"), cURL(url))
	fmt.Printf("%s %s\n\n", cBold("📁 Saving to:  "), cPath(videoDir))

	cmd := exec.Command(ytDlp, args...)
	cmd.Stdout = os.Stdout
	var stderrBuf strings.Builder
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Run(); err != nil {
		printError("Download failed: %v", err)
		diagnoseFail(stderrBuf.String(), cfg.Advanced)
		os.Exit(1)
	}

	filePath = resolveDownloadedFile(videoDir)
	printSuccess("Done! Saved to: %s", cPath(videoDir))
	return filePath, outputDir
}

// resolveDownloadedFile finds the video file inside videoDir by looking for
// known video extensions. This is more reliable than asking yt-dlp to predict
// the filename (which requires a network call and can get the extension wrong).
func resolveDownloadedFile(videoDir string) string {
	exts := map[string]bool{
		".mp4": true, ".mkv": true, ".webm": true,
		".mp3": true, ".m4a": true, ".ogg": true,
	}
	entries, err := os.ReadDir(videoDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if !e.IsDir() && exts[strings.ToLower(filepath.Ext(e.Name()))] {
			return filepath.Join(videoDir, e.Name())
		}
	}
	return ""
}

// diagnoseFail inspects yt-dlp's stderr output and prints targeted advice
// for the most common failure modes so the user knows exactly what to do.
func diagnoseFail(stderr string, adv AdvancedConfig) {
	fmt.Println()
	switch {
	case strings.Contains(stderr, "Sign in to confirm") || strings.Contains(stderr, "not a bot"):
		printError("YouTube is blocking the request — your browser cookies aren't working.")
		fmt.Println()
		printTip("1. Open Safari and make sure you're logged into youtube.com")
		printTip("2. If you're already logged in, try a different browser:")
		printTip("   Edit " + cPath("~/.config/video-downloader/config.yml") + " and set:")
		printTip("   " + cDim("cookies_from_browser: chrome") + "  (or firefox, brave, edge)")
		printTip("3. As a last resort, export cookies manually:")
		printTip("   Install \"Get cookies.txt LOCALLY\" in Chrome → export youtube.com")
		printTip("   Then set: " + cDim("cookies_file: ~/Downloads/youtube-cookies.txt"))

	case strings.Contains(stderr, "n challenge") || strings.Contains(stderr, "JS challenge"):
		printError("YouTube JS challenge failed — deno couldn't solve it.")
		fmt.Println()
		printTip("Make sure deno is installed:  " + cDim("brew install deno"))
		printTip("Then re-run — the solver script downloads automatically.")
		if adv.RemoteComponents == "" {
			printTip("Or add to config:  " + cDim("remote_components: ejs:github"))
		}

	case strings.Contains(stderr, "Private video"):
		printError("This video is private and cannot be downloaded.")

	case strings.Contains(stderr, "This video is not available") || strings.Contains(stderr, "Video unavailable"):
		printError("This video is unavailable in your region or has been removed.")

	case strings.Contains(stderr, "Requested format is not available"):
		printError("The requested video format isn't available.")
		printTip("Try setting " + cDim("resolution: 720") + " or " + cDim("format: best") + " in config.yml")

	default:
		printTip("Keep yt-dlp updated:  pip install -U yt-dlp")
		printTip("Try running directly to see the full error:")
		printTip("  yt-dlp --cookies-from-browser " + adv.CookiesFromBrowser + " <url>")
	}
	fmt.Println()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
